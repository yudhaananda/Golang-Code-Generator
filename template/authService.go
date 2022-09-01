package service

import (
	"[project]/entity"
	"[project]/helper"
	"[project]/input"
	"[project]/repository"
	"errors"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(input input.UserInput) (entity.User, error)
	Login(input input.LoginInput) (entity.User, error)
}

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) *authService {
	return &authService{userRepository}
}

func (s *authService) RegisterUser(input input.UserInput) (entity.User, error) {

	checkUser, err := s.userRepository.FindByUserName(input.UserName)

	if err != nil {
		return entity.User{}, errors.New("error find user")
	}

	if checkUser.UserName != 0 {
		return entity.User{}, errors.New("UserName sudah pernah diinputkan")
	}

	key := rand.Intn(9)
	password, err := bcrypt.GenerateFromPassword([]byte(input.Password), key)
	if err != nil {
		return entity.User{}, errors.New("error encrypt password")
	}

	user := entity.User{
		[registerItem]
	}

	newUser, err := s.userRepository.Save(user)

	if err != nil {
		return newUser, err
	}
	return newUser, nil
}

func (s *authService) Login(input input.LoginInput) (entity.User, error) {

	user, err := s.userRepository.FindByNoPegawai(input.UserName)

	if err != nil {
		return user, err
	}
	if user.Id == 0 {
		input.UserName = helper.FormatNoHp(input.UserName)
		user, err = s.userRepository.FindByNoHp(input.UserName)
		if err != nil {
			return user, err
		}
		if user.Id == 0 {
			return user, errors.New("user with username " + input.UserName + " not found")
		}
	}

	if !user.IsActive {
		return user, errors.New("user is inactived by " + user.NonActivedBy)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

	if err != nil {
		return user, errors.New("wrong password")
	}

	return user, nil
}
