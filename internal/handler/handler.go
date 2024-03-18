package handler

import (
	"context"
	"github.com/Paincake/filmbase/internal/database"
	"github.com/Paincake/filmbase/internal/errors"
	"github.com/Paincake/filmbase/internal/middleware"
	"github.com/Paincake/filmbase/internal/server"
	"github.com/oapi-codegen/runtime"
	"log/slog"
	"net/http"
)

type ServerInterfaceWrapper struct {
	Logger             *slog.Logger
	Repository         database.FilmbaseRepository
	Handler            server.ServerInterface
	HandlerMiddlewares []middleware.MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

// CreateActor operation middleware
func (siw *ServerInterfaceWrapper) CreateActor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, server.Filmbase_authScopes, []string{"write"})
	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateActor(w, r, siw.Repository, siw.Logger)
	}))
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PutActor operation middleware
func (siw *ServerInterfaceWrapper) PutActor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, server.Filmbase_authScopes, []string{"write"})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PutActor(w, r, siw.Repository, siw.Logger)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetActorFilms operation middleware
func (siw *ServerInterfaceWrapper) GetActorFilms(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, server.Filmbase_authScopes, []string{"write"})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetActorFilms(w, r, siw.Repository, siw.Logger)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// DeleteActor operation middleware
func (siw *ServerInterfaceWrapper) DeleteActor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "actorId" -------------
	var actorId int64

	err = runtime.BindStyledParameterWithOptions("simple", "actorId", r.URL.Query().Get("actorId"), &actorId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &errors.InvalidParamFormatError{ParamName: "actorId", Err: err})
		return
	}

	ctx = context.WithValue(ctx, server.Filmbase_authScopes, []string{"write"})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeleteActor(w, r, actorId, siw.Repository, siw.Logger)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostActorFilm operation middleware
func (siw *ServerInterfaceWrapper) PostActorFilm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "actorId" -------------
	var actorId int64

	err = runtime.BindStyledParameterWithOptions("simple", "actorId", r.URL.Query().Get("actorId"), &actorId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &errors.InvalidParamFormatError{ParamName: "actorId", Err: err})
		return
	}

	// ------------- Path parameter "filmId" -------------
	var filmId int64

	err = runtime.BindStyledParameterWithOptions("simple", "filmId", r.URL.Query().Get("filmId"), &filmId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &errors.InvalidParamFormatError{ParamName: "filmId", Err: err})
		return
	}

	ctx = context.WithValue(ctx, server.Filmbase_authScopes, []string{"write"})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostActorFilm(w, r, actorId, filmId, siw.Repository, siw.Logger)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetFilm operation middleware
func (siw *ServerInterfaceWrapper) GetFilm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	ctx = context.WithValue(ctx, server.Filmbase_authScopes, []string{"read"})

	// Parameter object where we will unmarshal all parameters from the context
	var params server.GetFilmParams

	// ------------- Optional query parameter "sortBy" -------------

	err = runtime.BindQueryParameter("form", true, false, "sortBy", r.URL.Query(), &params.SortBy)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &errors.InvalidParamFormatError{ParamName: "sortBy", Err: err})
		return
	}

	// ------------- Optional query parameter "sortKey" -------------

	err = runtime.BindQueryParameter("form", true, false, "sortKey", r.URL.Query(), &params.SortKey)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &errors.InvalidParamFormatError{ParamName: "sortKey", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetFilm(w, r, params, siw.Repository, siw.Logger)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// CreateFilm operation middleware
func (siw *ServerInterfaceWrapper) CreateFilm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, server.Filmbase_authScopes, []string{"write"})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateFilm(w, r, siw.Repository, siw.Logger)
	}))
	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// ChangeFilm operation middleware
func (siw *ServerInterfaceWrapper) ChangeFilm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, server.Filmbase_authScopes, []string{"write"})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ChangeFilm(w, r, siw.Repository, siw.Logger)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetFilmSearch operation middleware
func (siw *ServerInterfaceWrapper) GetFilmSearch(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	ctx = context.WithValue(ctx, server.Filmbase_authScopes, []string{"read"})

	// Parameter object where we will unmarshal all parameters from the context
	var params server.GetFilmSearchParams

	// ------------- Required query parameter "filmName" -------------

	if paramValue := r.URL.Query().Get("filmName"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &errors.RequiredParamError{ParamName: "filmName"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "filmName", r.URL.Query(), &params.FilmName)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &errors.InvalidParamFormatError{ParamName: "filmName", Err: err})
		return
	}

	// ------------- Required query parameter "actorName" -------------

	if paramValue := r.URL.Query().Get("actorName"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &errors.RequiredParamError{ParamName: "actorName"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "actorName", r.URL.Query(), &params.ActorName)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &errors.InvalidParamFormatError{ParamName: "actorName", Err: err})
		return
	}

	// ------------- Optional query parameter "sortBy" -------------

	err = runtime.BindQueryParameter("form", true, false, "sortBy", r.URL.Query(), &params.SortBy)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &errors.InvalidParamFormatError{ParamName: "sortBy", Err: err})
		return
	}

	// ------------- Optional query parameter "sortKey" -------------

	err = runtime.BindQueryParameter("form", true, false, "sortKey", r.URL.Query(), &params.SortKey)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &errors.InvalidParamFormatError{ParamName: "sortKey", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetFilmSearch(w, r, params, siw.Repository, siw.Logger)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}
	handler.ServeHTTP(w, r.WithContext(ctx))
}

// DeleteFilm operation middleware
func (siw *ServerInterfaceWrapper) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "filmId" -------------
	var filmId int64

	err = runtime.BindStyledParameterWithOptions("simple", "filmId", r.URL.Query().Get("filmId"), &filmId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &errors.InvalidParamFormatError{ParamName: "filmId", Err: err})
		return
	}

	ctx = context.WithValue(ctx, server.Filmbase_authScopes, []string{"write"})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeleteFilm(w, r, filmId, siw.Repository, siw.Logger)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}
func (siw *ServerInterfaceWrapper) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, server.Filmbase_authScopes, []string{"write"})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Login(w, r, siw.Repository, siw.Logger)
	}))

	handler.ServeHTTP(w, r.WithContext(ctx))
}
func (siw *ServerInterfaceWrapper) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, server.Filmbase_authScopes, []string{"write"})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Signup(w, r, siw.Repository, siw.Logger)
	}))

	handler.ServeHTTP(w, r.WithContext(ctx))
}
