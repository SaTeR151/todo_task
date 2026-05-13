package entity

type User struct {
	ID           string `json:"id"`
	Login        string `json:"login"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
}

type Users []User

type UserCreate struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserUpdate struct {
	ID           string  `json:"id"`
	Login        *string `json:"login"`
	Password     *string `json:"password"`
	RefreshToken *string `json:"refresh_token"`
}

type GetUsersOpts struct {
	ID    string
	Login string
}
