package oauth2manager

import (
	err "errors"

	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/store"

	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx"
	pg "github.com/vgarvardt/go-oauth2-pg"
	"github.com/vgarvardt/go-pg-adapter/pgxadapter"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	gconfiguration "github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure/configuration"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/infrastructure/configuration"
)

type TokenManager struct {
	manager *manage.Manager
	store   *pg.TokenStore
}

func New(
	config configuration.ConfigurationRepository,
	DBConfig gconfiguration.ConnectionStringRepository,
	logger infrastructure.Logger,
) (*TokenManager, error) {
	// check token configuration repository given
	if config == nil {
		return nil, err.New(ErrNoConfiguration)
	}
	c := config.Get()

	// check database configuration repository given
	if DBConfig == nil {
		return nil, err.New(ErrNoDBConfiguration)
	}

	//overriding the nil value of Logger
	if logger == nil {
		logger = new(infrastructure.LoggerNil)
	}

	manager := manage.NewDefaultManager()
	cfg := &manage.Config{
		AccessTokenExp:    c.Token.AccessExpire,
		RefreshTokenExp:   c.Token.RefreshExpire,
		IsGenerateRefresh: c.Token.IsGenerateRefresh,
	}

	manager.SetPasswordTokenCfg(cfg)
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	manager.MapAccessGenerate(
		generates.NewJWTAccessGenerate(
			[]byte(c.JWT),
			jwt.SigningMethodHS512),
	)

	pgxConnConfig, err := pgx.ParseURI(
		DBConfig.Get().ToString(DRIVER),
	)
	if err != nil {
		return nil, err
	}

	pgxConn, err := pgx.Connect(pgxConnConfig)
	if err != nil {
		return nil, err
	}

	adapter := pgxadapter.NewConn(pgxConn)
	tokenStore, err := pg.NewTokenStore(
		adapter,
		pg.WithTokenStoreGCInterval(c.GCInterval),
		pg.WithTokenStoreLogger(logger),
	)
	if err != nil {
		return nil, err
	}

	manager.MapTokenStorage(tokenStore)

	clientStore, err := pg.NewClientStore(
		adapter,
		pg.WithClientStoreLogger(logger),
	)
	if err != nil {
		tokenStore.Close()
		return nil, err
	}

	addClients(logger, clientStore, c.WhiteList...)

	manager.MapClientStorage(clientStore)

	return &TokenManager{
		manager: manager,
		store:   tokenStore,
	}, nil
}

func addClients(
	logger infrastructure.Logger,
	store *pg.ClientStore,
	clients ...*models.Client,
) {
	for _, client := range clients {
		err := store.Create(client)
		if err != nil {
			logger.Println("Warning:", err.Error())
		}
	}
}

func (tm *TokenManager) Manager() *manage.Manager {
	return tm.manager
}
func (tm *TokenManager) Store() *pg.TokenStore {
	return tm.store
}

func (tm *TokenManager) Close() error {
	return tm.store.Close()
}
