package database

type FilmbaseRepository interface {
	RunMigrations(query ...string)
	PostActor(actor Actor) (int64, error)
	PutActor(actor Actor) error
	DeleteActorById(actorId int64) error
	GetActorFilms() ([]ActorFilm, error)
	PostActorFilm(actorId, filmId int64) error
	GetFilmSearch(filmName string, actorName string, sortBy string, sortKey string) ([]ActorFilm, error)
	GetFilm(sortBy string, sortKey string) ([]Film, error)
	PostFilm(film Film) (int64, error)
	PutFilm(film Film) error
	DeleteFilmById(filmId int64) error
	Login(username string, password string) (string, error)
	Signup(username string, password string) error
}

type Actor struct {
	Id        int64  `db:"id" required:"true"`
	Name      string `db:"name" required:"true"`
	Gender    string `db:"gender" required:"true"`
	Birthdate string `db:"birthdate" required:"true"`
}

type Film struct {
	Id          int64  `db:"id" required:"true"`
	Name        string `db:"name" required:"true"`
	Description string `db:"description" required:"true"`
	ReleaseDate string `db:"release_date" required:"true"`
	Rating      int    `db:"rating" required:"true"`
}

type ActorFilm struct {
	ActorId         int64  `db:"actor_id" required:"true"`
	ActorName       string `db:"actor_name" required:"true"`
	ActorGender     string `db:"actor_gender" required:"true"`
	ActorBirthdate  string `db:"actor_birthdate" required:"true"`
	FilmId          int64  `db:"film_id" required:"true"`
	FilmName        string `db:"film_name" required:"true"`
	FilmDescription string `db:"film_descr" required:"true"`
	FilmReleaseDate string `db:"film_release_date" required:"true"`
	FilmRating      int    `db:"film_rating" required:"true"`
}

type User struct {
	Username string `db:"username" required:"true"`
	Password string `db:"password" required:"true"`
	Role     string `db:"role" required:"true"`
}
