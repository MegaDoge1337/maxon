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

func (cu *CreateUserCommand) ToDomain() User {
	return User{
		Username: cu.Username,
		Email:    cu.Email,
		Password: cu.Password,
	}
}
