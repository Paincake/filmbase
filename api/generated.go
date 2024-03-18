// Package filmbase provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
)

const (
	Filmbase_authScopes = "filmbase_auth.Scopes"
)


// GetFilmParamsSortBy defines parameters for GetFilm.
type GetFilmParamsSortBy string

// GetFilmParamsSortKey defines parameters for GetFilm.
type GetFilmParamsSortKey string

type GetFilmParams struct {
	// SortBy Sorting field
	SortBy *GetFilmParamsSortBy `form:"sortBy,omitempty" json:"sortBy,omitempty"`

	// SortKey Sorting key
	SortKey *GetFilmParamsSortKey `form:"sortKey,omitempty" json:"sortKey,omitempty"`
}

type GetFilmSearchParams struct {
	// FilmName Film name fragment
	FilmName string `form:"filmName" json:"filmName"`

	// ActorName Actor name fragment
	ActorName string `form:"actorName" json:"actorName"`

	// SortBy Sorting field
	SortBy *GetFilmSearchParamsSortBy `form:"sortBy,omitempty" json:"sortBy,omitempty"`

	// SortKey Sorting key
	SortKey *GetFilmSearchParamsSortKey `form:"sortKey,omitempty" json:"sortKey,omitempty"`
}
type GetFilmSearchParamsSortBy string

// GetFilmSearchParamsSortKey defines parameters for GetFilmSearch.
type GetFilmSearchParamsSortKey string

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create an actor information
	// (POST /actor)
	CreateActor(w http.ResponseWriter, r *http.Request)
	// Change actor information
	// (PUT /actor)
	PutActor(w http.ResponseWriter, r *http.Request)
	// Get an actor's films information
	// (POST /actor/films)
	GetActorFilms(w http.ResponseWriter, r *http.Request)
	// Delete actor information
	// (DELETE /actor/{actorId})
	DeleteActor(w http.ResponseWriter, r *http.Request, actorId int64)
	// Add a film information to actor
	// (POST /actor/{actorId}/{filmId})
	PostActorFilm(w http.ResponseWriter, r *http.Request, actorId int64, filmId int64)
	// Get film information with sorting
	// (GET /film)
	GetFilm(w http.ResponseWriter, r *http.Request, params GetFilmParams)
	// Create a film information
	// (POST /film)
	CreateFilm(w http.ResponseWriter, r *http.Request)
	// Change a film information
	// (PUT /film)
	ChangeFilm(w http.ResponseWriter, r *http.Request)
	// Get film information with with searching by fields
	// (GET /film/search)
	GetFilmSearch(w http.ResponseWriter, r *http.Request, params GetFilmSearchParams)
	// Delete film information
	// (DELETE /film/{filmId})
	DeleteFilm(w http.ResponseWriter, r *http.Request, filmId int64)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// Create an actor information
// (POST /actor)
func (_ Unimplemented) CreateActor(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Change actor information
// (PUT /actor)
func (_ Unimplemented) PutActor(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get an actor's films information
// (POST /actor/films)
func (_ Unimplemented) GetActorFilms(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Delete actor information
// (DELETE /actor/{actorId})
func (_ Unimplemented) DeleteActor(w http.ResponseWriter, r *http.Request, actorId int64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Add a film information to actor
// (POST /actor/{actorId}/{filmId})
func (_ Unimplemented) PostActorFilm(w http.ResponseWriter, r *http.Request, actorId int64, filmId int64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get film information with sorting
// (GET /film)
func (_ Unimplemented) GetFilm(w http.ResponseWriter, r *http.Request, params GetFilmParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Create a film information
// (POST /film)
func (_ Unimplemented) CreateFilm(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Change a film information
// (PUT /film)
func (_ Unimplemented) ChangeFilm(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get film information with with searching by fields
// (GET /film/search)
func (_ Unimplemented) GetFilmSearch(w http.ResponseWriter, r *http.Request, params GetFilmSearchParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Delete film information
// (DELETE /film/{filmId})
func (_ Unimplemented) DeleteFilm(w http.ResponseWriter, r *http.Request, filmId int64) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// CreateActor operation middleware
func (siw *ServerInterfaceWrapper) CreateActor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, Filmbase_authScopes, []string{"write"})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateActor(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PutActor operation middleware
func (siw *ServerInterfaceWrapper) PutActor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, Filmbase_authScopes, []string{"write"})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PutActor(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetActorFilms operation middleware
func (siw *ServerInterfaceWrapper) GetActorFilms(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, Filmbase_authScopes, []string{"write"})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetActorFilms(w, r)
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

	err = runtime.BindStyledParameterWithOptions("simple", "actorId", chi.URLParam(r, "actorId"), &actorId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "actorId", Err: err})
		return
	}

	ctx = context.WithValue(ctx, Filmbase_authScopes, []string{"write"})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeleteActor(w, r, actorId)
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

	err = runtime.BindStyledParameterWithOptions("simple", "actorId", chi.URLParam(r, "actorId"), &actorId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "actorId", Err: err})
		return
	}

	// ------------- Path parameter "filmId" -------------
	var filmId int64

	err = runtime.BindStyledParameterWithOptions("simple", "filmId", chi.URLParam(r, "filmId"), &filmId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "filmId", Err: err})
		return
	}

	ctx = context.WithValue(ctx, Filmbase_authScopes, []string{"write"})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostActorFilm(w, r, actorId, filmId)
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

	ctx = context.WithValue(ctx, Filmbase_authScopes, []string{"read"})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetFilmParams

	// ------------- Optional query parameter "sortBy" -------------

	err = runtime.BindQueryParameter("form", true, false, "sortBy", r.URL.Query(), &params.SortBy)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "sortBy", Err: err})
		return
	}

	// ------------- Optional query parameter "sortKey" -------------

	err = runtime.BindQueryParameter("form", true, false, "sortKey", r.URL.Query(), &params.SortKey)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "sortKey", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetFilm(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// CreateFilm operation middleware
func (siw *ServerInterfaceWrapper) CreateFilm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, Filmbase_authScopes, []string{"write"})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateFilm(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// ChangeFilm operation middleware
func (siw *ServerInterfaceWrapper) ChangeFilm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, Filmbase_authScopes, []string{"write"})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ChangeFilm(w, r)
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

	ctx = context.WithValue(ctx, Filmbase_authScopes, []string{"read"})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetFilmSearchParams

	// ------------- Required query parameter "filmName" -------------

	if paramValue := r.URL.Query().Get("filmName"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "filmName"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "filmName", r.URL.Query(), &params.FilmName)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "filmName", Err: err})
		return
	}

	// ------------- Required query parameter "actorName" -------------

	if paramValue := r.URL.Query().Get("actorName"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "actorName"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "actorName", r.URL.Query(), &params.ActorName)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "actorName", Err: err})
		return
	}

	// ------------- Optional query parameter "sortBy" -------------

	err = runtime.BindQueryParameter("form", true, false, "sortBy", r.URL.Query(), &params.SortBy)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "sortBy", Err: err})
		return
	}

	// ------------- Optional query parameter "sortKey" -------------

	err = runtime.BindQueryParameter("form", true, false, "sortKey", r.URL.Query(), &params.SortKey)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "sortKey", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetFilmSearch(w, r, params)
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

	err = runtime.BindStyledParameterWithOptions("simple", "filmId", chi.URLParam(r, "filmId"), &filmId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "filmId", Err: err})
		return
	}

	ctx = context.WithValue(ctx, Filmbase_authScopes, []string{"write"})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeleteFilm(w, r, filmId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("errors unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// TODO CHANGE CHI TO net/http
// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ServerOptions{})
}

type ServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/actor", wrapper.CreateActor)
	})
	r.Group(func(r chi.Router) {
		r.Put(options.BaseURL+"/actor", wrapper.PutActor)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/actor/films", wrapper.GetActorFilms)
	})
	r.Group(func(r chi.Router) {
		r.Delete(options.BaseURL+"/actor/{actorId}", wrapper.DeleteActor)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/actor/{actorId}/{filmId}", wrapper.PostActorFilm)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/film", wrapper.GetFilm)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/film", wrapper.CreateFilm)
	})
	r.Group(func(r chi.Router) {
		r.Put(options.BaseURL+"/film", wrapper.ChangeFilm)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/film/search", wrapper.GetFilmSearch)
	})
	r.Group(func(r chi.Router) {
		r.Delete(options.BaseURL+"/film/{filmId}", wrapper.DeleteFilm)
	})

	return r
}