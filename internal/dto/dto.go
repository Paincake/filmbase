package dto

type Actor struct {
	Id        int64  `json:"id" required:"true"`
	Name      string `json:"name" required:"true" validate:"nonzero"`
	Gender    string `json:"gender" required:"true" validate:"nonzero"`
	Birthdate string `json:"birthdate" required:"true" validate:"nonzero"`
}

type Film struct {
	Id          int64  `json:"id" required:"true"`
	Name        string `json:"name" required:"true" validate:"nonzero, min=1,max=50"`
	Description string `json:"description" required:"true" validate:"nonzero, min=1,max=1000"`
	ReleaseDate string `json:"release-date" required:"true" validate:"nonzero"`
	Rating      int    `json:"rating" required:"true" validate:"nonzero"`
}

type ActorFilm struct {
	Actor Actor  `json:"actor" required:"true" validate:"nonzero"`
	Films []Film `json:"films" required:"true" validate:"nonzero"`
}

type User struct {
	Username string `json:"username" required:"true" validate:"nonzero"`
	Password string `json:"password" required:"true" validate:"nonzero"`
}
