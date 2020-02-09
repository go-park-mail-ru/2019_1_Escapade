package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	mdb "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database/mocks"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"
	pkgServer "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"
	mserver "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server/mocks"
	madb "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database/mocks"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/proto"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/constants"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/database"
	mgdb "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/database/mocks"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/handlers"
	game "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/service"
	"github.com/gorilla/websocket"
	. "github.com/smartystreets/goconvey/convey"

	mclients "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/clients/mocks"
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

func (conf *ConfigurationArgsFake) Init() *handlers.ConfigurationArgs {
	return &handlers.ConfigurationArgs{
		C:         &config.Configuration{},
		FieldPath: "",
		RoomPath:  "",
	}
}

type A struct{ madb.UserRepositoryI }

type DatabaseArgsFake struct {
}

var InterfacesMock = struct {
	ConsulService *mserver.ConsulServiceI
	ChatService   *mclients.ChatI
}{
	ConsulService: new(mserver.ConsulServiceI),
	ChatService:   new(mclients.ChatI),
}

// system
func TestSystem(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Test borders Init", t, func() {
		var (
			user   = new(madb.UserRepositoryI)
			record = new(madb.RecordRepositoryI)
			game   = new(mgdb.GameRepositoryI)
			chatS  = new(mclients.ChatI)

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
		gameUC.On("Init", game).Return(da.Game)

		da.User.Init(user, record)
		da.Game.Init(game)

		var (
			ca      = new(ConfigurationArgsFake).Init()
			handler = new(handlers.GameHandler)
		)
		db.On("Open", ca.C.DataBase).Return(err)

		userUC.On("Use", db).Return(err)
		gameUC.On("Use", db).Return(err)

		//InterfacesMock.Database = da
		/*
			var args = &handlers.Args{
				Test: true,
				Input: &handlers.Input{
					CMD: &start.CommandLineArgs{
						ConfigurationPath: "game.json",
						PhotoPublicPath:   "./../../../internal/pkg/photo/photo.json",
						PhotoPrivatePath:  "./../../../secret.json",
						MainPort:          "3002",
					},
					FieldPath: "./../../../internal/pkg/constants/field.json",
					RoomPath:  "./../../../internal/pkg/constants/room.json",
				},
			}
			go handlers.Run(args, handlers.InterfacesDefault)
		*/
		testRun(da)

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

		//handler.Handle()
	})
}

func testRun(da *database.Input) {
	args := &pkgServer.Args{
		Input:         generateTestInput(),
		Loader:        generateTestLoader(),
		ConsulService: new(pkgServer.ConsulService),
		Service: &game.Service{
			Chat:     InterfacesMock.ChatService,
			Consul:   InterfacesMock.ConsulService,
			Constant: new(constants.RepositoryFS),
			Database: da,
		},
	}
	pkgServer.Run(args)
}

func generateTestInput() *game.Input {
	var input = new(game.Input)

	input.CallInit = func() {
		input.Data.FieldPath = "./../../../internal/pkg/constants/field.json"
		input.Data.RoomPath = "./../../../internal/pkg/constants/room.json"
		input.Data.MainPort = "3002"
	}
	return input
}

func generateTestLoader() *server.Loader {
	var loader = new(server.Loader)
	loader.Init(new(config.RepositoryFS), "game.json")
	loader.CallExtra = func() error {
		return loader.LoadPhoto(
			"./../../../internal/pkg/photo/photo.json",
			"./../../../secret.json")
	}
	return loader
}
