package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	mdb "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database/mocks"
	madb "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database/mocks"
	mcc "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/clients/mocks"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/proto"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/constants"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/database"
	mgdb "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/database/mocks"
	"github.com/gorilla/websocket"
	. "github.com/smartystreets/goconvey/convey"
)

type RepositoryFake struct{}

var (
	rc = constants.RoomConfiguration{
		NameMin:          10,
		NameMax:          30,
		TimeToPrepareMin: 0,
		TimeToPrepareMax: 60,
		TimeToPlayMin:    0,
		TimeToPlayMax:    7200,
		PlayersMin:       2,
		PlayersMax:       100,
		ObserversMax:     100,
		Set:              true,
	}
	fc = constants.FieldConfiguration{
		WidthMin:  5,
		WidthMax:  100,
		HeightMin: 5,
		HeightMax: 100,
		Set:       true,
	}
)

// getRoom load Room config(.json file) from FS by its path
func (rfs *RepositoryFake) GetRoom(path string) (constants.RoomConfiguration, error) {

	return rc, nil
}

// getField load field config(.json file) from FS by its path
func (rfs *RepositoryFake) GetField(path string) (constants.FieldConfiguration, error) {
	return fc, nil
}

type ConfigurationArgsFake struct {
}

func (conf *ConfigurationArgsFake) Init() *ConfigurationArgs {
	return &ConfigurationArgs{
		C:         &config.Configuration{},
		FieldPath: "",
		RoomPath:  "",
	}
}

type A struct{ madb.UserRepositoryI }

type DatabaseArgsFake struct {
}

// system
func TestSystem(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Test borders Init", t, func() {
		var (
			user   = new(madb.UserRepositoryI)
			record = new(madb.RecordRepositoryI)
			game   = new(mgdb.GameRepositoryI)
			chatS  = new(mcc.ChatI)

			userUC = new(madb.UserUseCaseI)
			gameUC = new(mgdb.GameUseCaseI)
			db     = new(mdb.DatabaseI)

			da = &database.Input{
				Database: db,

				User: userUC,
				Game: gameUC,
			}
			err    error
			chat   *proto.Chat
			chatID *proto.ChatID
		)

		chatS.On("GetChat", chat).Return(chatID, err)

		userUC.On("Init", user, record).Return(da.User)
		gameUC.On("Init", game, chatS).Return(da.Game)

		da.User.Init(user, record)
		da.Game.Init(game)

		var (
			constRep = new(RepositoryFake)
			ca       = new(ConfigurationArgsFake).Init()
			handler  = new(GameHandler)
		)
		db.On("Open", ca.C.DataBase).Return(err)

		userUC.On("Use", db).Return(err)
		gameUC.On("Use", db).Return(err)

		handler.Init(constRep, chatS, ca, da)

		for i := 0; i < 1; i++ {
			s := httptest.NewServer(http.HandlerFunc(handler.Handle))
			defer s.Close()

			u := "ws" + strings.TrimPrefix(s.URL, "http")

			ws, _, err := websocket.DefaultDialer.Dial(u, nil)
			if err != nil {
				Convey("When websocket dials, the error should be nil", func() {
					So(err, ShouldBeNil)
				})

				return
			}
			defer ws.Close()
		}
	})
}
