package main

import (
	"bytes"
	"encoding/json"
	"github.com/Paincake/filmbase/internal/auth"
	"github.com/Paincake/filmbase/internal/config"
	"github.com/Paincake/filmbase/internal/database"
	"github.com/Paincake/filmbase/internal/database/postgres"
	"github.com/Paincake/filmbase/internal/dto"
	"github.com/Paincake/filmbase/internal/middleware"
	"github.com/Paincake/filmbase/internal/server"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	FilmTableDDl = `
	CREATE TABLE IF NOT EXISTS film (
		id serial PRIMARY KEY,
		name varchar CHECK(char_length(name) > 1),
		description varchar CHECK ( char_length(name) > 1 AND char_length(name) < 1000 ),
		rating smallint CHECK (rating IN (0,1,2,3,4,5,6,7,8,9,10)),
		release_date date
	)
`
	ActorTableDDl = `
	CREATE TABLE IF NOT EXISTS actor (
		id serial PRIMARY KEY,
		name varchar CHECK(char_length(name) > 1),
		gender varchar CHECK(gender IN ('male', 'female')),
		birthdate date
	)
`
	ActorFilmTableDDl = `
	CREATE TABLE IF NOT EXISTS actor_films (
		actorid int REFERENCES actor(id),
		filmid int REFERENCES film(id),
		PRIMARY KEY (actorid, filmid)
	)
`
	UserTableDDl = `
	CREATE TABLE IF NOT EXISTS api_users (
		username varchar PRIMARY KEY,
		password varchar,
		role varchar
	)
`
	ClearTables = `
	TRUNCATE TABLE actor;
	TRUNCATE TABLE film;
	TRUNCATE TABLE actor_films;
	TRUNCATE TABLE api_users;
	ALTER SEQUENCE actor_id_seq RESTART WITH 1;
	ALTER SEQUENCE film_id_seq RESTART WITH 1
`
)

var router http.Handler
var db database.FilmbaseRepository

func TestMain(m *testing.M) {
	setup()
	exitCode := m.Run()
	teardown()
	os.Exit(exitCode)
}

func TestCreateActor_ShouldGet401(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/actor", nil)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, recorder.Code, http.StatusUnauthorized)
}

func TestCreateActor_ShouldGet403(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/actor", nil)
	token, _ := auth.CreateJWT("test", "asdasdasdasd")
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestCreateActor_ShouldGet422(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/actor", nil)
	token, _ := auth.CreateJWT("test", "admin")
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
}
func TestCreateActor_ShouldGet200(t *testing.T) {
	recorder := httptest.NewRecorder()
	body, _ := json.Marshal(dto.Actor{Name: "a", Gender: "male", Birthdate: "2001-01-01"})
	req := httptest.NewRequest("POST", "/actor", bytes.NewBuffer(body))
	token, _ := auth.CreateJWT("test", "admin")
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestPutActor_ShouldGet401(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/actor", nil)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, recorder.Code, http.StatusUnauthorized)
}

func TestPutActor_ShouldGet403(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/actor", nil)
	token, _ := auth.CreateJWT("test", "asdasdasdasd")
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestPutActor_ShouldGet422(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/actor", nil)
	token, _ := auth.CreateJWT("test", "admin")
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
}
func TestPutActor_ShouldGet200(t *testing.T) {
	recorder := httptest.NewRecorder()
	body, _ := json.Marshal(dto.Actor{Name: "a", Gender: "male", Birthdate: "2001-01-01"})
	req := httptest.NewRequest("PUT", "/actor", bytes.NewBuffer(body))
	token, _ := auth.CreateJWT("test", "admin")
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetActorFilms_ShouldGet401(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/actor/films", nil)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, recorder.Code, http.StatusUnauthorized)
}

func TestGetActorFilms_ShouldGet403(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/actor/films", nil)
	token, _ := auth.CreateJWT("test", "asdasdasdasd")
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestGetActorFilms_ShouldGet200(t *testing.T) {
	_, err := db.PostActor(database.Actor{Id: 1, Name: "TEST", Gender: "male", Birthdate: "2001-01-01"})
	if err != nil {
		t.Fail()
	}
	_, err = db.PostFilm(database.Film{Id: 1, Name: "TEST", Description: "TEST", ReleaseDate: "2001-01-01", Rating: 0})
	if err != nil {
		t.Errorf("test failed: %s", err)
	}
	err = db.PostActorFilm(1, 1)
	if err != nil {
		t.Errorf("test failed: %s", err)
	}
	res, err := db.GetActorFilms()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/actor/films", nil)
	token, err := auth.CreateJWT("test", "admin")
	if err != nil {
		t.Errorf("test failed: %s", err)
	}
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)

	var response server.Response
	err = json.NewDecoder(recorder.Body).Decode(&response)
	if err != nil {
		t.Errorf("test failed: %s", err)
	}
	assert.Equal(t, res[0].ActorName, "TEST")
	assert.Equal(t, res[0].FilmName, "TEST")
	assert.NotEqual(t, nil, res)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestDeleteActor_ShouldGet401(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/actor/1", nil)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, recorder.Code, http.StatusUnauthorized)
}

func TestDeleteActor_ShouldGet403(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/actor/1", nil)
	token, _ := auth.CreateJWT("test", "asdasdasdasd")
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestDeleteActor_ShouldGet200(t *testing.T) {
	_, err := db.PostActor(database.Actor{Id: 10, Name: "TESTDELETE", Gender: "female", Birthdate: "2001-01-01"})
	if err != nil {
		t.Fail()
	}
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/actor/1", nil)
	token, err := auth.CreateJWT("test", "admin")
	if err != nil {
		t.Errorf("test failed: %s", err)
	}
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, recorder.Code, http.StatusOK)

}
func TestPostActorFilms_ShouldGet401(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/actor/films/1/1", nil)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, recorder.Code, http.StatusUnauthorized)
}

func TestPostActorFilms_ShouldGet403(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/actor/films/1/1", nil)
	token, _ := auth.CreateJWT("test", "asdasdasdasd")
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestPostActorFilms_ShouldGet200(t *testing.T) {
	_, err := db.PostActor(database.Actor{Id: 20, Name: "TEST", Gender: "male", Birthdate: "2001-01-01"})
	if err != nil {
		t.Fail()
	}
	_, err = db.PostFilm(database.Film{Id: 20, Name: "TEST", Description: "TEST", ReleaseDate: "2001-01-01", Rating: 0})
	if err != nil {
		t.Errorf("test failed: %s", err)
	}
	err = db.PostActorFilm(20, 20)
	if err != nil {
		t.Errorf("test failed: %s", err)
	}
	res, err := db.GetActorFilms()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/actor/films/20/20", nil)
	token, err := auth.CreateJWT("test", "admin")
	if err != nil {
		t.Errorf("test failed: %s", err)
	}
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)

	var response server.Response
	err = json.NewDecoder(recorder.Body).Decode(&response)
	if err != nil {
		t.Errorf("test failed: %s", err)
	}
	assert.NotEqual(t, nil, res)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetFilm_ShouldGet401(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/film", nil)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, recorder.Code, http.StatusUnauthorized)
}

func TestGetFilm_ShouldGet403(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "film", nil)
	token, _ := auth.CreateJWT("test", "asdasdasdasd")
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestGetFilm_ShouldGet200(t *testing.T) {
	_, err := db.PostFilm(database.Film{Id: 30, Name: "TEST", Description: "TEST", ReleaseDate: "2001-01-01", Rating: 0})
	if err != nil {
		t.Errorf("test failed: %s", err)
	}
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/actor/films", nil)
	token, err := auth.CreateJWT("test", "admin")
	if err != nil {
		t.Errorf("test failed: %s", err)
	}
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)

	var response server.Response
	err = json.NewDecoder(recorder.Body).Decode(&response)
	if err != nil {
		t.Errorf("test failed: %s", err)
	}
	assert.NotEqual(t, nil, response.ResponseBody)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestCreateFilm_ShouldGet401(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/film", nil)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, recorder.Code, http.StatusUnauthorized)
}

func TestCreateFilm_ShouldGet403(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/film", nil)
	token, _ := auth.CreateJWT("test", "asdasdasdasd")
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestCreateFilm_ShouldGet422(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/film", nil)
	token, _ := auth.CreateJWT("test", "admin")
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
}

func TestCreateFilm_ShouldGet200(t *testing.T) {
	recorder := httptest.NewRecorder()
	body, _ := json.Marshal(dto.Film{Name: "a", Description: "male", ReleaseDate: "2001-01-01", Rating: 1})
	req := httptest.NewRequest("POST", "/film", bytes.NewBuffer(body))
	token, _ := auth.CreateJWT("test", "admin")
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestPutFilm_ShouldGet401(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/film", nil)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, recorder.Code, http.StatusUnauthorized)
}

func TestPutFilm_ShouldGet403(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/film", nil)
	token, _ := auth.CreateJWT("test", "asdasdasdasd")
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestPutFilm_ShouldGet422(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/film", nil)
	token, _ := auth.CreateJWT("test", "admin")
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
}

func TestPutFilm_ShouldGet200(t *testing.T) {
	recorder := httptest.NewRecorder()
	body, _ := json.Marshal(dto.Film{Name: "a", Description: "male", ReleaseDate: "2001-01-01", Rating: 1})
	req := httptest.NewRequest("POST", "/film", bytes.NewBuffer(body))
	token, _ := auth.CreateJWT("test", "admin")
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}
func TestDeleteFilm_ShouldGet401(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/film/1", nil)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, recorder.Code, http.StatusUnauthorized)
}

func TestDeleteFilm_ShouldGet403(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/film/1", nil)
	token, _ := auth.CreateJWT("test", "asdasdasdasd")
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestDeleteFilm_ShouldGet200(t *testing.T) {
	_, err := db.PostFilm(database.Film{Id: 40, Name: "TEST", Description: "TEST", ReleaseDate: "2001-01-01", Rating: 0})
	if err != nil {
		t.Errorf("test failed: %s", err)
	}
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/film/40", nil)
	token, err := auth.CreateJWT("test", "admin")
	if err != nil {
		t.Errorf("test failed: %s", err)
	}
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, recorder.Code, http.StatusOK)

}

func TestGetFilmSearch_ShouldGet401(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/film/search", nil)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, recorder.Code, http.StatusUnauthorized)
}

func TestGetFilmSearch_ShouldGet403(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/film/search", nil)
	token, _ := auth.CreateJWT("test", "asdasdasdasd")
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestGetFilmSearch_ShouldGet200(t *testing.T) {
	_, err := db.PostActor(database.Actor{Id: 100, Name: "TEST", Gender: "male", Birthdate: "2001-01-01"})
	if err != nil {
		t.Fail()
	}
	_, err = db.PostFilm(database.Film{Id: 100, Name: "TEST", Description: "TEST", ReleaseDate: "2001-01-01", Rating: 0})
	if err != nil {
		t.Errorf("test failed: %s", err)
	}
	err = db.PostActorFilm(1, 1)
	if err != nil {
		t.Errorf("test failed: %s", err)
	}
	films, err := db.GetFilmSearch("TEST", "TEST", "actor_name", "DESC")
	if err != nil {
		t.Errorf("test failed: %s", err)
	}
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/actor/films?actorName=TEST", nil)
	token, err := auth.CreateJWT("test", "admin")
	if err != nil {
		t.Errorf("test failed: %s", err)
	}
	req.Header.Set("Token", token)
	router.ServeHTTP(recorder, req)

	var response server.Response
	err = json.NewDecoder(recorder.Body).Decode(&response)
	if err != nil {
		t.Errorf("test failed: %s", err)
	}

	assert.NotEqual(t, nil, films)
	assert.Equal(t, films[0].FilmName, "TEST")
	assert.Equal(t, films[0].ActorName, "TEST")
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func teardown() {
	db.RunMigrations(ClearTables)
}

func setup() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	configPath := os.Getenv("TEST_CONFIG_PATH")
	if configPath == "" {
		log.Error("TEST_CONFIG_PATH env variable is not set")
	}
	if _, err := os.Stat(configPath); err != nil {
		log.Error("errors opening config file: %s", err)
	}
	var cfg config.Config
	var srv config.HTTPServer
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Error("errors opening config file: %s", err)
	}
	err = cleanenv.ReadConfig(configPath, &srv)
	if err != nil {
		log.Error("errors opening config file: %s", err)
	}
	repository, err := postgres.New(cfg.Name, cfg.User, cfg.Password, cfg.Host, cfg.Port)
	repository.RunMigrations(FilmTableDDl, ActorTableDDl, ActorFilmTableDDl, UserTableDDl, `INSERT INTO api_users VALUES('test', 'test', 'admin')`)
	db = repository
	middlewares := []middleware.MiddlewareFunc{middleware.VerifyJWT}
	opts := HandlerOptions{
		BaseRouter:       *http.NewServeMux(),
		Middlewares:      middlewares,
		ErrorHandlerFunc: nil,
	}
	si := server.BasicServer{}
	router = HandlerWithOptions(si, &opts, repository, log)
}
