// deu/pkg/router/router.go (Updated)
package router

import (
	"net/http"

	"deu/internal/users"
	"deu/internal/places"
)

type Config struct {
	UserHandler *users.Handler
	PlaceHandler *places.Handler
}

func NewRouter(cfg Config) http.Handler {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /", fs)

	mux.HandleFunc("GET /users", cfg.UserHandler.GetAll)
	mux.HandleFunc("POST /users", cfg.UserHandler.Create)
	mux.HandleFunc("DELETE /users", cfg.UserHandler.DeleteAll)

	mux.HandleFunc("GET /users/{id}", cfg.UserHandler.GetById)
	mux.HandleFunc("PATCH /users/{id}", cfg.UserHandler.Update)
	mux.HandleFunc("DELETE /users/{id}", cfg.UserHandler.DeleteById)

	mux.HandleFunc("POST /users/{id}/places/{place_id}", cfg.UserHandler.AddVisitedPlace)
	mux.HandleFunc("GET /users/{id}/places/{place_id}", cfg.UserHandler.CheckIfVisited)
	mux.HandleFunc("DELETE /users/{id}/places/{place_id}", cfg.UserHandler.RemoveVisitedPlace)

	mux.HandleFunc("GET /places", cfg.PlaceHandler.GetAll)
	mux.HandleFunc("POST /places", cfg.PlaceHandler.Create)
	mux.HandleFunc("DELETE /places", cfg.PlaceHandler.DeleteAll)

	mux.HandleFunc("GET /places/{id}", cfg.PlaceHandler.GetById)
	mux.HandleFunc("PATCH /places/{id}", cfg.PlaceHandler.Update)
	mux.HandleFunc("DELETE /places/{id}", cfg.PlaceHandler.DeleteById)

	return mux
}