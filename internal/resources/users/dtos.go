package users

type userDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type createUserResponse struct {
	message string
	token   string
}
