package delivery

import "github.com/ruscalworld/study-planner/internal/user"

type UserController struct {
	userRepository user.Repository
}

func NewUserController(userRepository user.Repository) *UserController {
	return &UserController{userRepository: userRepository}
}
