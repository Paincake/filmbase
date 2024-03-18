package main

import (
	"fmt"
	"github.com/Paincake/filmbase/internal/config"
	"github.com/Paincake/filmbase/internal/database"
	"github.com/Paincake/filmbase/internal/database/postgres"
	"github.com/Paincake/filmbase/internal/handler"
	"github.com/Paincake/filmbase/internal/middleware"
	"github.com/Paincake/filmbase/internal/server"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, srv := config.MustLoad()
	si := server.BasicServer{}
	repository, err := postgres.New(cfg.Name, cfg.User, cfg.Password, cfg.Host, cfg.Port)
	if err != nil {
		logger.Error(fmt.Sprintf("Error loading config file:%s", err))
	}
	logger.Debug("Config loaded")

	middlewares := []middleware.MiddlewareFunc{middleware.VerifyJWT}

	opts := HandlerOptions{
		BaseRouter:       *http.NewServeMux(),
		Middlewares:      middlewares,
		ErrorHandlerFunc: nil,
	}

	router := HandlerWithOptions(si, &opts, repository, logger)
	logger.Info(fmt.Sprintf("Starting server at address %s", srv.Address))
	err = http.ListenAndServe("localhost:8080", router)
}

type HandlerOptions struct {
	BaseRouter       http.ServeMux
	Middlewares      []middleware.MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si server.ServerInterface, options *HandlerOptions, repository database.FilmbaseRepository, logger *slog.Logger) http.Handler {
	r := &options.BaseRouter

	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := handler.ServerInterfaceWrapper{
		Logger:             logger,
		Repository:         repository,
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.HandleFunc("POST "+"/actor", wrapper.CreateActor)
	r.HandleFunc("GET "+"/actor/films", wrapper.GetActorFilms)
	r.HandleFunc("PUT "+"/actor", wrapper.PutActor)
	r.HandleFunc("DELETE "+"/actor/{actorId}", wrapper.DeleteActor)
	r.HandleFunc("POST "+"/actor/{actorId}/{filmId}", wrapper.PostActorFilm)
	r.HandleFunc("GET "+"/film", wrapper.GetFilm)
	r.HandleFunc("POST "+"/film", wrapper.CreateFilm)
	r.HandleFunc("PUT "+"/film", wrapper.ChangeFilm)
	r.HandleFunc("GET "+"/film/search", wrapper.GetFilmSearch)
	r.HandleFunc("DELETE "+"/film/{filmId}", wrapper.DeleteFilm)
	r.HandleFunc("POST "+"/login", wrapper.Login)
	r.HandleFunc("POST "+"/sign", wrapper.Signup)

	return r
}
