package domain

type CreateSessionCommand struct {
	Login    string
	Password string
}

type RefreshSessionCommand struct {
	Refresh string
}
