package domain

const (
	adminRoleName = "admin"
	userRoleName  = "user"
)

// recieved role name => constant role name
var RolesRegistry = map[string]string{
	"admin": adminRoleName,
	"user":  userRoleName,
}

type Role struct {
	ID     int
	UserId int
	Name   string
}
