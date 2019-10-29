package main

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/constants"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	api "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/photo"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"net/http"
	"os"
)

func main() {
	var (
		configuration *config.Configuration
		handler       *api.Handler
		db            *database.DataBase
		err           error
	)

	if len(os.Args) < 6 {
		utils.Debug(false, "Game service need 5 command line arguments. But only",
			len(os.Args)-1, "get.")
		return
	}

	var (
		configurationPath = os.Args[1]
		photoPublicPath   = os.Args[2]
		photoPrivatePath  = os.Args[3]
		fieldPath         = os.Args[4]
		roomPath          = os.Args[5]
	)

	configuration, err = config.Init(configurationPath)
	if err != nil {
		utils.Debug(false, "Initialization error with main configuration:", err.Error())
		return
	}

	err = photo.Init(photoPublicPath, photoPrivatePath)
	if err != nil {
		utils.Debug(false, "Initialization error with photo configuration:", err.Error())
		return
	}

	err = constants.InitField(fieldPath)
	if err != nil {
		utils.Debug(false, "Initialization error with field constants:", err.Error())
		return
	}

	err = constants.InitRoom(roomPath)
	if err != nil {
		utils.Debug(false, "Initialization error with room constants:", err.Error())
		return
	}
	db, err = database.Init(configuration.DataBase)
	if err != nil {
		utils.Debug(false, "Initialization error with database:", err.Error())
		return
	}
	handler = api.Init(db, configuration)

	metrics.InitApi()

	var (
		r    = router.HistoryRouter(handler, configuration.Cors)
		port = router.Port(configuration)
	)

	if err = http.ListenAndServe(port, r); err != nil {
		utils.Debug(false, "Serving error:", err.Error())
	}
}
