package domain

type RegisterCommand struct {
	Username string
	Email    string
	Password string
}

type LoginCommand struct {
	Username string
	Password string
}

type RefreshCommand struct {
	Refresh string
}
