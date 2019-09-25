package router

import (
	"net/http"

	"github.com/sankarvj/kubekadal/app/controller"
	"github.com/gorilla/mux"
)

//InitRouter routes traffic to handlers
func InitRouter() *mux.Router {
	router := mux.NewRouter()
	// status
	router.Handle("/app/dialogflow", http.HandlerFunc(controller.DialogFlowHandler)).Methods("POST", "OPTIONS")

	return router
}
