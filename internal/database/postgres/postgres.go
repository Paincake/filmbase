package postgres

import (
	"fmt"
	"github.com/Paincake/filmbase/internal/database"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type Database struct {
	db *sqlx.DB
}

func New(dbname, username, password, host, port string) (*Database, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?&sslmode=disable",
		username,
		password,
		host,
		port,
		dbname)
	db, err := sqlx.Connect("pgx", connectionString)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return &Database{db: db}, nil
}
func (d *Database) RunMigrations(query ...string) {
	for _, q := range query {
		d.db.Query(q)
	}
}
func (d *Database) PostActor(actor database.Actor) (int64, error) {
	var id int64
	err := d.db.Get(&id, "INSERT INTO actor (name, gender, birthdate) VALUES ($1, $2, $3::date) RETURNING id;", actor.Name, actor.Gender, actor.Birthdate)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (d *Database) PutActor(actor database.Actor) error {
	_, err := d.db.Query("UPDATE actor SET name=$1, gender=$2, birthdate=$3;", actor.Name, actor.Gender, actor.Birthdate)
	if err != nil {
		return err
	}
	return nil
}
func (d *Database) DeleteActorById(actorId int64) error {
	_, err := d.db.Query("DELETE FROM actor WHERE id = $1", actorId)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) GetActorFilms() ([]database.ActorFilm, error) {
	var actors []database.ActorFilm
	err := d.db.Select(&actors,
		"SELECT a.id AS actor_id, a.name AS actor_name, a.gender AS actor_gender, a.birthdate AS actor_birthdate, "+
			"f.id AS film_id, f.name AS film_name, f.description AS film_descr, f.release_date AS film_release_date, f.rating AS film_rating "+
			"FROM actor_films "+
			"JOIN actor a on actor_films.actorId = a.id "+
			"JOIN film f on f.id = actor_films.filmId")
	if err != nil {
		return nil, err
	}
	return actors, nil
}
func (d *Database) PostActorFilm(actorId, filmId int64) error {
	_, err := d.db.Query("INSERT INTO actor_films (actor_id, film_id) VALUES ($1, $2)", actorId, filmId)
	if err != nil {
		return err
	}
	return nil
}
func (d *Database) GetFilmSearch(filmName string, actorName string, sortBy string, sortKey string) ([]database.ActorFilm, error) {
	var actorFilms []database.ActorFilm
	err := d.db.Select(&actorFilms,
		"SELECT a.id AS actor_id, a.name AS actor_name, a.gender AS actor_gender, a.birthdate AS actor_birthdate, "+
			"f.id AS film_id, f.name AS film_name, f.description AS film_descr, f.release_date AS film_release_date, f.rating AS film_rating "+
			"JOIN actor a on actor_films.actorId = a.id "+
			"JOIN film f on f.id = actor_films.filmId "+
			"WHERE actor_name LIKE '%$1%' AND film_name LIKE '%$2%' "+
			"ORDER BY $3 $4", filmName, actorName, sortBy, sortKey)
	if err != nil {
		return nil, err
	}
	return actorFilms, nil
}
func (d *Database) GetFilm(sortBy string, sortKey string) ([]database.Film, error) {
	var films []database.Film
	err := d.db.Select(&films, "SELECT * FROM film ORDER BY $1 $2", sortBy, sortKey)
	if err != nil {
		return nil, err
	}
	return films, nil

}
func (d *Database) PostFilm(film database.Film) (int64, error) {
	var id int64
	err := d.db.Get(&id, "INSERT INTO film (name, description, release_date, rating) VALUES ($1, $2, $3, $4) RETURNING id", film.Name, film.Description, film.ReleaseDate, film.Rating)
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (d *Database) PutFilm(film database.Film) error {
	_, err := d.db.Query("UPDATE film SET name=$1, description=$2, release_date=$3, rating=$4", film.Name, film.Description, film.ReleaseDate, film.Rating)
	if err != nil {
		return err
	}
	return nil
}
func (d *Database) DeleteFilmById(filmId int64) error {
	_, err := d.db.Query("DELETE FROM film WHERE id = $1", filmId)
	if err != nil {
		return err
	}
	return err
}
func (d *Database) Login(username string, password string) (string, error) {
	var user database.User
	err := d.db.Get(&user, "SELECT username, password, role FROM api_users WHERE username = $1", username)
	if err != nil {
		return "", err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}
	return user.Role, nil
}

func (d *Database) Signup(username string, password string) error {
	_, err := d.db.Query("INSERT INTO api_users (username, password) VALUES ($1, $2)", username, password)
	if err != nil {
		return err
	}
	return nil
}
