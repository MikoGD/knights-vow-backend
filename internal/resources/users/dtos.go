package users

type userDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type createUserResponse struct {
	Message string `json:"name"`
	Token   string `json:"token"`
}

type loginUserResponse = createUserResponse

type authStatusResponse struct {
	IsAuthenticated bool `json:"isAuthenticated"`
}
