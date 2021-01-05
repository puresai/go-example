package service

import (
	"context"
	"demo7-docker-compose/model"
	"errors"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

type UserInfoDTO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

var (
	ErrUserExisted = errors.New("user is existed")
	ErrPassword    = errors.New("email and password are not match")
	ErrRegistering = errors.New("email is registering")
)

type RegisterUser struct {
	Username string
	Password string
	Email    string
}

type UserService interface {
	Login(ctx context.Context, email, pass string) (*UserInfoDTO, error)
	Register(ctx context.Context, user *RegisterUser) (*UserInfoDTO, error)
}

type UserServiceImpl struct {
	userDao model.UserDao
}

func MakeUserServiceImpl(userDao model.UserDao) UserService {
	return &UserServiceImpl{
		userDao,
	}
}

func (userService *UserServiceImpl) Login(ctx context.Context, email, password string) (*UserInfoDTO, error) {
	user, err := userService.userDao.SelectByEmail(email)
	if err == nil {
		if user.Password == password {
			return &UserInfoDTO{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
			}, nil
		} else {
			return nil, ErrPassword
		}
	} else {
		log.Printf("err : %s", err)
	}
	return nil, err
}

func (userService UserServiceImpl) Register(ctx context.Context, user *RegisterUser) (*UserInfoDTO, error) {
	ret := model.RedisClient.SetNX(user.Email, 1, time.Duration(5)*time.Second)
	if ret.Val() == false {
		return nil, ErrRegistering
	}
	defer model.RedisClient.Del(user.Email)

	existUser, err := userService.userDao.SelectByEmail(user.Email)

	if (err == nil && existUser == nil) || err == gorm.ErrRecordNotFound {
		newUser := &model.UserEntity{
			Username: user.Username,
			Password: user.Password,
			Email:    user.Email,
		}
		err = userService.userDao.Save(newUser)
		if err == nil {
			return &UserInfoDTO{
				ID:       newUser.ID,
				Username: newUser.Username,
				Email:    newUser.Email,
			}, nil
		}
	}
	if err == nil {
		err = ErrUserExisted
	}
	return nil, err

}
