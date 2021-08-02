package server

import (
	"apitesting/handler"
	"apitesting/middleware"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Run() {
	r := mux.NewRouter()
	userRouter := r.NewRoute().Subrouter()
	adminRouter := r.NewRoute().Subrouter()

	r.Methods("GET").PathPrefix("/testing/").Handler(
		http.StripPrefix(
			"/testing/",
			http.FileServer(http.Dir(viper.GetString("storage.static"))),
		),
	)

	r.Path("/testing/admin/login").Handler(handler.Login)
	userRouter.Path("/testing/admin/adminid").Handler(handler.AdminId)
	userRouter.Path("/testing/admin/api/projects").Handler(handler.ApiProjects)
	userRouter.Path("/testing/admin/testing/projects").Handler(handler.TestingProjects)
	userRouter.Path("/testing/admin/syncapi").Handler(handler.SyncApi)

	userRouter.Path("/testing/admin/env").Handler(handler.SelectEnv)
	adminRouter.Path("/testing/admin/user").Handler(handler.SelectUser)
	userRouter.Path("/testing/admin/api").Handler(handler.SelectApi)
	userRouter.Path("/testing/admin/testing").Handler(handler.SelectTesting)

	userRouter.Path("/testing/admin/env/save").Handler(handler.SaveEnv)
	adminRouter.Path("/testing/admin/user/save").Handler(handler.SaveUser)
	userRouter.Path("/testing/admin/api/save").Handler(handler.SaveApi)
	userRouter.Path("/testing/admin/testing/save").Handler(handler.SaveTesting)

	userRouter.Path("/testing/admin/env/del").Handler(handler.DelEnv)
	userRouter.Path("/testing/admin/user/del").Handler(handler.DelUser)
	userRouter.Path("/testing/admin/api/del").Handler(handler.DelApi)
	userRouter.Path("/testing/admin/testing/del").Handler(handler.DelTesting)

	userRouter.Use(middleware.LogRequest, middleware.UserAuth)
	adminRouter.Use(middleware.LogRequest, middleware.AdminAuth)

	addr := viper.GetString("addr")
	logrus.Infof("ListenAndServe at %s", addr)
	logrus.Fatal(http.ListenAndServe(addr, r))
}
