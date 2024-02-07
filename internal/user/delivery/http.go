package delivery

import "study-planner/internal/user"

type UserController struct {
	userRepository user.Repository
}

func NewUserController(userRepository user.Repository) *UserController {
	return &UserController{userRepository: userRepository}
}
