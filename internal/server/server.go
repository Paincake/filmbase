package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Paincake/filmbase/internal/auth"
	"github.com/Paincake/filmbase/internal/database"
	"github.com/Paincake/filmbase/internal/dto"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/validator.v2"
	"log/slog"
	"net/http"
	"strings"
)

const (
	Filmbase_authScopes = "filmbase_auth.Scopes"
	DefaultFilmSortKey  = "DESC"
	DefaultFilmSortBy   = "rating"
)

type Response struct {
	Code         int
	Error        error
	ResponseBody any
}

func returnResponse(w http.ResponseWriter, encoder json.Encoder, code int, body any, err error) {
	w.WriteHeader(code)
	encoder.Encode(Response{
		Code:         code,
		Error:        err,
		ResponseBody: body,
	})
}

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
	// CreateActor Create an actor information
	// (POST /actor)
	CreateActor(w http.ResponseWriter, r *http.Request, repository database.FilmbaseRepository, log *slog.Logger)
	// PutActor Change actor information
	// (PUT /actor)
	PutActor(w http.ResponseWriter, r *http.Request, repository database.FilmbaseRepository, log *slog.Logger)
	// GetActorFilms Get an actor's films information
	// (POST /actor/films)
	GetActorFilms(w http.ResponseWriter, r *http.Request, repository database.FilmbaseRepository, log *slog.Logger)
	// DeleteActor Delete actor information
	// (DELETE /actor/{actorId})
	DeleteActor(w http.ResponseWriter, r *http.Request, actorId int64, repository database.FilmbaseRepository, log *slog.Logger)
	// PostActorFilm Add a film information to actor
	// (POST /actor/{actorId}/{filmId})
	PostActorFilm(w http.ResponseWriter, r *http.Request, actorId int64, filmId int64, repository database.FilmbaseRepository, log *slog.Logger)
	// GetFilm Get film information with sorting
	// (GET /film)
	GetFilm(w http.ResponseWriter, r *http.Request, params GetFilmParams, repository database.FilmbaseRepository, log *slog.Logger)
	// CreateFilm Create a film information
	// (POST /film)
	CreateFilm(w http.ResponseWriter, r *http.Request, repository database.FilmbaseRepository, log *slog.Logger)
	// ChangeFilm Change a film information
	// (PUT /film)
	ChangeFilm(w http.ResponseWriter, r *http.Request, repository database.FilmbaseRepository, log *slog.Logger)
	// GetFilmSearch Get film information with with searching by fields
	// (GET /film/search)
	GetFilmSearch(w http.ResponseWriter, r *http.Request, params GetFilmSearchParams, repository database.FilmbaseRepository, log *slog.Logger)
	// DeleteFilm Delete film information
	// (DELETE /film/{filmId})
	DeleteFilm(w http.ResponseWriter, r *http.Request, filmId int64, repository database.FilmbaseRepository, log *slog.Logger)
	// Login logs in the system
	// (POST /login)
	Login(w http.ResponseWriter, r *http.Request, repository database.FilmbaseRepository, log *slog.Logger)
	Signup(w http.ResponseWriter, r *http.Request, repository database.FilmbaseRepository, log *slog.Logger)
}

// BasicServer server implementation that returns http.StatusNotImplemented for each endpoint.
type BasicServer struct{}

// CreateActor Create an actor information
// (POST /actor)
func (_ BasicServer) CreateActor(w http.ResponseWriter, r *http.Request, repository database.FilmbaseRepository, log *slog.Logger) {
	const op = "server.CreateActor POST /actor"
	log.With(op)
	w.Header().Set("Content-Type", "application/json")
	role := r.Context().Value("role")
	encoder := json.NewEncoder(w)
	if role != "admin" {
		log.Info(fmt.Sprintf("Request discarded: forbidden"))
		returnResponse(w, *encoder, http.StatusForbidden, nil, fmt.Errorf("you can't access that resource"))
		return
	}
	decoder := json.NewDecoder(r.Body)
	var actor dto.Actor
	err := decoder.Decode(&actor)
	if err != nil {
		log.Info(fmt.Sprintf("Request discarded: invalid JSON body: %s", err))
		returnResponse(w, *encoder, http.StatusUnprocessableEntity, nil, fmt.Errorf("bad request: %s", err))
		return
	}
	err = validator.Validate(actor)
	if err != nil {
		log.Info(fmt.Sprintf("Request discarded: invalid JSON body: %s", err))
		returnResponse(w, *encoder, http.StatusUnprocessableEntity, nil, fmt.Errorf("bad request: %s", err))
		return
	}

	entityActor := database.Actor{
		Id:        actor.Id,
		Name:      actor.Name,
		Gender:    actor.Gender,
		Birthdate: actor.Birthdate,
	}
	id, err := repository.PostActor(entityActor)
	if err != nil {
		log.Error(fmt.Sprintf("Request discarded: server error: %s", err))
		returnResponse(w, *encoder, http.StatusInternalServerError, nil, fmt.Errorf("internal server error: %s", err))
		return
	}
	returnResponse(w, *encoder, http.StatusOK, id, nil)
}

// PutActor Change actor information
// (PUT /actor)
func (_ BasicServer) PutActor(w http.ResponseWriter, r *http.Request, repository database.FilmbaseRepository, log *slog.Logger) {
	const op = "server.PutActor PUT /actor"
	log.With(op)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	role := r.Context().Value("role")
	if role != "admin" {
		log.Info(fmt.Sprintf("Request discarded: forbidden"))
		returnResponse(w, *encoder, http.StatusForbidden, nil, fmt.Errorf("forbidden"))
		return
	}
	decoder := json.NewDecoder(r.Body)
	var actor dto.Actor
	err := decoder.Decode(&actor)
	if err != nil {
		log.Info(fmt.Sprintf("Request discarded: invalid JSON body: %s", err))
		returnResponse(w, *encoder, http.StatusUnprocessableEntity, nil, fmt.Errorf("bad request: %s", err))
		return
	}
	err = validator.Validate(actor)
	if err != nil {
		log.Info(fmt.Sprintf("Request discarded: invalid JSON body: %s", err))
		returnResponse(w, *encoder, http.StatusUnprocessableEntity, nil, fmt.Errorf("bad request: %s", err))
		return
	}

	entityActor := database.Actor{
		Id:        actor.Id,
		Name:      actor.Name,
		Gender:    actor.Gender,
		Birthdate: actor.Birthdate,
	}
	err = repository.PutActor(entityActor)
	if err != nil {
		log.Error(fmt.Sprintf("Request discarded: server error: %s", err))
		returnResponse(w, *encoder, http.StatusInternalServerError, nil, fmt.Errorf("internal server error: %s", err))
		return
	}
	returnResponse(w, *encoder, http.StatusOK, nil, nil)

}

// GetActorFilms Get an actor's films information
// (GET /actor/films)
func (_ BasicServer) GetActorFilms(w http.ResponseWriter, r *http.Request, repository database.FilmbaseRepository, log *slog.Logger) {
	const op = "server.GetActorFilms GET /actor/films"
	log.With(op)
	encoder := json.NewEncoder(w)
	role := r.Context().Value("role")
	if role != "user" && role != "admin" {
		log.Info(fmt.Sprintf("Request discarded: forbidden"))
		returnResponse(w, *encoder, http.StatusForbidden, nil, fmt.Errorf("forbidden"))
		return
	}
	films, err := repository.GetActorFilms()
	if err != nil {
		log.Error(fmt.Sprintf("Request discarded: server error: %s", err))
		returnResponse(w, *encoder, http.StatusInternalServerError, nil, fmt.Errorf("internal server error: %s", err))
	}
	actorFilmMap := make(map[dto.Actor][]dto.Film)
	for _, e := range films {
		actor := dto.Actor{
			Id:        e.ActorId,
			Name:      e.ActorName,
			Gender:    e.ActorGender,
			Birthdate: e.ActorBirthdate,
		}
		film := dto.Film{
			Id:          e.FilmId,
			Name:        e.FilmName,
			Description: strings.TrimSuffix(e.FilmDescription, " "),
			ReleaseDate: e.FilmReleaseDate,
			Rating:      e.FilmRating,
		}
		if _, ok := actorFilmMap[actor]; ok {
			actorFilmMap[actor] = make([]dto.Film, 0)
		}
		actorFilmMap[actor] = append(actorFilmMap[actor], film)
	}
	var actorFilms []dto.ActorFilm
	for k, v := range actorFilmMap {
		actorFilms = append(actorFilms, dto.ActorFilm{
			Actor: k,
			Films: v,
		})
	}

	returnResponse(w, *encoder, http.StatusOK, actorFilms, nil)
}

// DeleteActor Delete actor information
// (DELETE /actor/{actorId})
func (_ BasicServer) DeleteActor(w http.ResponseWriter, r *http.Request, actorId int64, repository database.FilmbaseRepository, log *slog.Logger) {
	const op = "server.DeleteActor DELETE /actor"
	log.With(op)
	encoder := json.NewEncoder(w)
	err := repository.DeleteActorById(actorId)
	role := r.Context().Value("role")
	if role != "admin" {
		log.Info(fmt.Sprintf("Request discarded: forbidden"))
		returnResponse(w, *encoder, http.StatusForbidden, nil, fmt.Errorf("forbidden"))
		return
	}
	if err != nil {
		log.Error(fmt.Sprintf("Request discarded: server error: %s", err))
		returnResponse(w, *encoder, http.StatusInternalServerError, nil, fmt.Errorf("internal server error: %s", err))
		return
	}
	returnResponse(w, *encoder, http.StatusOK, nil, nil)
}

// PostActorFilm Add a film information to actor
// (POST /actor/{actorId}/{filmId})
func (_ BasicServer) PostActorFilm(w http.ResponseWriter, r *http.Request, actorId int64, filmId int64, repository database.FilmbaseRepository, log *slog.Logger) {
	const op = "server.PostActorFilm POST /actor"
	log.With(op)
	encoder := json.NewEncoder(w)
	role := r.Context().Value("role")
	if role != "admin" {
		log.Info(fmt.Sprintf("Request discarded: forbidden"))
		returnResponse(w, *encoder, http.StatusForbidden, nil, fmt.Errorf("forbidden"))
		return
	}
	err := repository.PostActorFilm(actorId, filmId)
	if err != nil {
		log.Error(fmt.Sprintf("Request discarded: server error: %s", err))
		returnResponse(w, *encoder, http.StatusInternalServerError, nil, fmt.Errorf("internal server error: %s", err))
		return
	}
	returnResponse(w, *encoder, http.StatusOK, nil, nil)
}

// GetFilm Get film information with sorting
// (GET /film)
func (_ BasicServer) GetFilm(w http.ResponseWriter, r *http.Request, params GetFilmParams, repository database.FilmbaseRepository, log *slog.Logger) {
	const op = "server.GetFilm GET /film"
	log.With(op)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	role := r.Context().Value("role")
	if role != "admin" && role != "user" {
		log.Info(fmt.Sprintf("Request discarded: forbidden"))
		returnResponse(w, *encoder, http.StatusForbidden, nil, fmt.Errorf("forbidden"))
		return
	}
	var sortBy, sortKey string
	if string(*params.SortBy) == "" {
		sortBy = DefaultFilmSortBy
	} else {
		sortBy = string(*params.SortBy)
	}
	if string(*params.SortKey) == "" {
		sortKey = DefaultFilmSortKey
	} else {
		sortKey = string(*params.SortKey)
	}
	films, err := repository.GetFilm(sortBy, sortKey)
	if err != nil {
		log.Error(fmt.Sprintf("Request discarded: server error: %s", err))
		returnResponse(w, *encoder, http.StatusInternalServerError, nil, fmt.Errorf("internal server error: %s", err))
		return
	}
	returnResponse(w, *encoder, http.StatusOK, films, nil)

}

// CreateFilm Create a film information
// (POST /film)
func (_ BasicServer) CreateFilm(w http.ResponseWriter, r *http.Request, repository database.FilmbaseRepository, log *slog.Logger) {
	const op = "server.CreateFilm POST /film"
	log.With(op)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	role := r.Context().Value("role")
	if role != "admin" {
		log.Info(fmt.Sprintf("Request discarded: forbidden"))
		returnResponse(w, *encoder, http.StatusForbidden, nil, fmt.Errorf("forbidden"))
		return
	}
	decoder := json.NewDecoder(r.Body)
	var film dto.Film
	err := decoder.Decode(&film)
	if err != nil {
		log.Info(fmt.Sprintf("Request discarded: invalid JSON body: %s", err))
		returnResponse(w, *encoder, http.StatusUnprocessableEntity, nil, fmt.Errorf("bad request: %s", err))
		return
	}
	err = validator.Validate(film)
	if err != nil {
		log.Info(fmt.Sprintf("Request discarded: invalid JSON body: %s", err))
		returnResponse(w, *encoder, http.StatusUnprocessableEntity, nil, fmt.Errorf("bad request: %s", err))
		return
	}
	entityFilm := database.Film{
		Id:          film.Id,
		Name:        film.Name,
		Description: film.Description,
		ReleaseDate: film.ReleaseDate,
		Rating:      film.Rating,
	}
	id, err := repository.PostFilm(entityFilm)
	if err != nil {
		log.Error(fmt.Sprintf("Request discarded: server error: %s", err))
		returnResponse(w, *encoder, http.StatusInternalServerError, nil, fmt.Errorf("internal server error: %s", err))
		return
	}
	returnResponse(w, *encoder, http.StatusOK, id, nil)
}

// ChangeFilm Change a film information
// (PUT /film)
func (_ BasicServer) ChangeFilm(w http.ResponseWriter, r *http.Request, repository database.FilmbaseRepository, log *slog.Logger) {
	const op = "server.ChangeFilm PUT /film"
	log.With(op)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	role := r.Context().Value("role")
	if role != "admin" {
		log.Info(fmt.Sprintf("Request discarded: forbidden"))
		returnResponse(w, *encoder, http.StatusForbidden, nil, fmt.Errorf("forbidden"))
		return
	}
	decoder := json.NewDecoder(r.Body)
	var film dto.Film
	err := decoder.Decode(&film)
	if err != nil {
		log.Info(fmt.Sprintf("Request discarded: invalid JSON body: %s", err))
		returnResponse(w, *encoder, http.StatusUnprocessableEntity, nil, fmt.Errorf("bad request: %s", err))
		return
	}
	err = validator.Validate(film)
	if err != nil {
		log.Info(fmt.Sprintf("Request discarded: invalid JSON body: %s", err))
		returnResponse(w, *encoder, http.StatusUnprocessableEntity, nil, fmt.Errorf("bad request: %s", err))
		return
	}
	entityFilm := database.Film{
		Id:          film.Id,
		Name:        film.Name,
		Description: film.Description,
		ReleaseDate: film.ReleaseDate,
		Rating:      film.Rating,
	}
	err = repository.PutFilm(entityFilm)
	if err != nil {
		log.Error(fmt.Sprintf("Request discarded: server error: %s", err))
		returnResponse(w, *encoder, http.StatusInternalServerError, nil, fmt.Errorf("internal server error: %s", err))
		return
	}
	returnResponse(w, *encoder, http.StatusOK, nil, nil)
}

// GetFilmSearch Get film information with searching by fields
// (GET /film/search)
func (_ BasicServer) GetFilmSearch(w http.ResponseWriter, r *http.Request, params GetFilmSearchParams, repository database.FilmbaseRepository, log *slog.Logger) {
	const op = "server.GetFilmSearch GET /film/search"
	log.With(op)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	role := r.Context().Value("role")
	if role != "admin" && role != "user" {
		log.Info(fmt.Sprintf("Request discarded: forbidden"))
		returnResponse(w, *encoder, http.StatusForbidden, nil, fmt.Errorf("forbidden"))
		return
	}
	var sortBy, sortKey string
	if string(*params.SortBy) == "" {
		sortBy = DefaultFilmSortBy
	} else {
		sortBy = string(*params.SortBy)
	}
	if string(*params.SortKey) == "" {
		sortKey = DefaultFilmSortKey
	} else {
		sortKey = string(*params.SortKey)
	}
	films, err := repository.GetFilmSearch(params.FilmName, params.ActorName, sortBy, sortKey)
	if err != nil {
		log.Error(fmt.Sprintf("Request discarded: server error: %s", err))
		returnResponse(w, *encoder, http.StatusInternalServerError, nil, fmt.Errorf("internal server error: %s", err))
		return
	}
	returnResponse(w, *encoder, http.StatusOK, films, nil)
}

// DeleteFilm Delete film information
// (DELETE /film/{filmId})
func (_ BasicServer) DeleteFilm(w http.ResponseWriter, r *http.Request, filmId int64, repository database.FilmbaseRepository, log *slog.Logger) {
	encoder := json.NewEncoder(w)
	role := r.Context().Value("role")
	if role != "admin" {
		log.Info(fmt.Sprintf("Request discarded: forbidden"))
		returnResponse(w, *encoder, http.StatusForbidden, nil, fmt.Errorf("forbidden"))
		return
	}
	err := repository.DeleteFilmById(filmId)
	if err != nil {
		log.Error(fmt.Sprintf("Request discarded: server error: %s", err))
		returnResponse(w, *encoder, http.StatusInternalServerError, nil, fmt.Errorf("internal server error: %s", err))
		return
	}
	returnResponse(w, *encoder, http.StatusOK, nil, nil)

}

func (_ BasicServer) Login(w http.ResponseWriter, r *http.Request, repository database.FilmbaseRepository, log *slog.Logger) {
	encoder := json.NewEncoder(w)
	creds := r.Header.Get("Authorization")
	if creds == "" || len(strings.Split(creds, " ")) < 2 {
		returnResponse(w, *encoder, http.StatusUnauthorized, nil, nil)
		return
	}
	creds = strings.Split(creds, " ")[1]
	raw, err := base64.StdEncoding.DecodeString(creds)
	if err != nil {
		returnResponse(w, *encoder, http.StatusUnauthorized, err, nil)
		return
	}
	decodedCreds := strings.Split(string(raw), ":")
	role, err := repository.Login(decodedCreds[0], decodedCreds[1])
	if err != nil {
		returnResponse(w, *encoder, http.StatusUnauthorized, err, nil)
		return
	}
	token, err := auth.CreateJWT(decodedCreds[0], role)
	if err != nil {

		log.Error(fmt.Sprintf("Request discarded: server error: %s", err))

		returnResponse(w, *encoder, http.StatusInternalServerError, nil, fmt.Errorf("internal server error: %s", err))
		return
	}
	returnResponse(w, *encoder, http.StatusOK, token, nil)

}
func (_ BasicServer) Signup(w http.ResponseWriter, r *http.Request, repository database.FilmbaseRepository, log *slog.Logger) {
	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)
	var user dto.User
	err := decoder.Decode(&user)
	if err != nil {
		log.Info(fmt.Sprintf("Request discarded: invalid JSON body: %s", err))
		returnResponse(w, *encoder, http.StatusUnprocessableEntity, nil, fmt.Errorf("bad request: %s", err))
	}
	err = validator.Validate(user)
	if err != nil {
		log.Info(fmt.Sprintf("Request discarded: invalid JSON body: %s", err))
		returnResponse(w, *encoder, http.StatusUnprocessableEntity, nil, fmt.Errorf("bad request: %s", err))
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		log.Error(fmt.Sprintf("Request discarded: server error: %s", err))
		returnResponse(w, *encoder, http.StatusInternalServerError, nil, fmt.Errorf("internal server error: %s", err))
		return
	}
	err = repository.Signup(user.Username, string(hashedPassword))
	if err != nil {
		log.Info(fmt.Sprintf("Request discarded: forbidden"))
		returnResponse(w, *encoder, http.StatusForbidden, nil, fmt.Errorf("forbidden"))
		return
	}
	returnResponse(w, *encoder, http.StatusCreated, nil, nil)

}
