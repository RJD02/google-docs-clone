package routes

import (
	"github.com/RJD02/google-docs-clone/handlers"
	"github.com/gorilla/mux"
)

var ApiRouter = mux.NewRouter()

func Run() {
	ApiRouter.HandleFunc("/", handlers.HandleWebSocket)
}
