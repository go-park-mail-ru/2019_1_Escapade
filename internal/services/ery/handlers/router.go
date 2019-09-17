package eryhandlers

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"

	// erydb "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/database"

	mi "github.com/go-park-mail-ru/2019_1_Escapade/internal/middleware"
	"github.com/gorilla/mux"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

// параметры пути
const (
	USERID    = "user_id"
	PROJECTID = "project_id"
	SCENEID   = "scene_id"
	OBJECTID  = "object_id"
)

// Router вернуть маршрутизатор путей
func Router(H *Handler, cors config.CORS, cc config.Cookie,
	ca config.Auth, client config.AuthClient) *mux.Router {

	r := mux.NewRouter()

	r.PathPrefix("/swagger").Handler(httpSwagger.Handler(
		httpSwagger.URL("swagger/doc.json"), //The url pointing to API definition"
	))

	r.PathPrefix("/metrics").Handler(promhttp.Handler())

	var noAuth = r.PathPrefix("/ery").Subrouter()
	var withAuth = r.PathPrefix("/ery").Subrouter()

	noAuth.HandleFunc("/users", H.HandleUsers).Methods("GET", "POST")
	noAuth.HandleFunc("/user", H.HandleUser).Methods("OPTIONS", "POST")
	withAuth.HandleFunc("/user", H.HandleUser).Methods("DELETE", "PUT", "GET")
	withAuth.HandleFunc("/user/image", H.HandleUserImage).Methods("POST", "OPTIONS")

	withAuth.HandleFunc("/user/projects", H.HandleProjects).Methods("POST", "GET", "OPTIONS")

	withAuth.HandleFunc("/users/{user_id}", H.HandleUserID).Methods("GET", "OPTIONS")

	noAuth.HandleFunc("/session", H.HandleSession).Methods("GET", "POST", "OPTIONS")
	withAuth.HandleFunc("/session", H.HandleSession).Methods("DELETE", "PUT")

	withAuth.HandleFunc("/projects", H.HandleProjectsSearch).Methods("GET", "OPTIONS")
	withAuth.HandleFunc("/project/{project_id}", H.HandleProjectID).Methods("POST", "DELETE", "GET", "PUT", "OPTIONS")

	withAuth.HandleFunc("/project/{project_id}/members/{user_id}", H.HandleProjectIDMembers).Methods("POST", "DELETE", "PUT", "OPTIONS")

	withAuth.HandleFunc("/project/{project_id}/members/{user_id}/token", H.HandleProjectIDMembersToken).Methods("PUT", "OPTIONS")

	withAuth.HandleFunc("/project/{project_id}/scene", H.HandleScene).Methods("POST", "OPTIONS")

	withAuth.HandleFunc("/project/{project_id}/scene/{scene_id}", H.HandleSceneID).Methods("GET", "PUT", "DELETE", "OPTIONS")

	withAuth.HandleFunc("/project/{project_id}/scene/{scene_id}/erythrocyte", H.HandleSceneErythrocyte).Methods("POST", "OPTIONS")
	withAuth.HandleFunc("/project/{project_id}/scene/{scene_id}/erythrocyte/{object_id}", H.HandleSceneErythrocyteID).Methods("PUT", "DELETE", "OPTIONS")

	withAuth.HandleFunc("/project/{project_id}/scene/{scene_id}/erythrocyte_object", H.HandleSceneErythrocyteObject).Methods("POST", "OPTIONS")
	withAuth.HandleFunc("/project/{project_id}/scene/{scene_id}/erythrocyte_object/{object_id}", H.HandleSceneErythrocyteObjectID).Methods("PUT", "DELETE", "OPTIONS")

	withAuth.HandleFunc("/project/{project_id}/scene/{scene_id}/disease", H.HandleSceneDisease).Methods("POST", "OPTIONS")
	withAuth.HandleFunc("/project/{project_id}/scene/{scene_id}/disease/{object_id}", H.HandleSceneDiseaseID).Methods("PUT", "DELETE", "OPTIONS")

	r.Use(mi.Recover, mi.CORS(cors), mi.Metrics)
	withAuth.Use(mi.Auth(cc, ca, client))

	return r
}
