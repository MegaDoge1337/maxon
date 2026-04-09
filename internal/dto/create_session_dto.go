package dto

type CreateSessionDto struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
