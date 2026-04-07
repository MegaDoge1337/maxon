package domain

type CreateUserCommand struct {
	Username string
	Email    string
	Password string
}

type UpdateUserCommand struct {
	Username    string
	Email       string
	OldPassword string
	NewPassword string
}
