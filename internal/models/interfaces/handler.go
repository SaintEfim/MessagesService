package interfaces

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler interface {
	ConfigureRoutes(r *mux.Router)
	SendMessage(w http.ResponseWriter, r *http.Request)
}
