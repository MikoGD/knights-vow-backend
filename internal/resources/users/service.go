package users

import (
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	myjwt "knights-vow/pkg/jwt"
)

type UserService interface {
	CreateUser(user *userDTO) (int, error)
	AuthenticateUser(user *userDTO) (string, error)
	CheckUserAuthStatus(token string, userID int) (bool, error)
}

type userService struct {
	repository UserRepository
}

func createNewUserService(repository UserRepository) UserService {
	return &userService{repository}
}

func (service *userService) CreateUser(userDTO *userDTO) (int, error) {
	user, err := service.repository.GetUserByUsername(userDTO.Username)
	if err != nil {
		return -1, err
	}

	if user != nil {
		return -1, &UserExistsError{user.Username}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return -1, err
	}

	userID, err := service.repository.SaveUser(userDTO.Username, string(hashedPassword))
	if err != nil {
		return -1, err
	}

	return userID, nil
}

func (service *userService) AuthenticateUser(userDTO *userDTO) (string, error) {
	user, err := service.repository.GetUserByUsername(userDTO.Username)

	if err != nil {
		return "", err
	}

	if user == nil {
		return "", &InvalidLoginError{}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDTO.Password))

	if err != nil {
		return "", &InvalidLoginError{}
	}

	return myjwt.CreateJWT(user.ID), nil
}

func (service *userService) CheckUserAuthStatus(tokenString string, userID int) (bool, error) {
	token := myjwt.ParseJWT(tokenString)

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if claims["client_id"] != float64(userID) {
			return false, &InvalidLoginError{}
		}
	}

	if token.Valid {
		return true, nil
	}

	return false, nil
}
