package models

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"net/http"
)

// User defines user format
type User struct {
	Id            int    `json:"id" orm:"column(id);auto"`
	Username      string `json:"username" orm:"column(username);size(128)"`
	Password      string `json:"password" orm:"column(password);size(128)"`
	Salt          string `json:"salt" orm:"column(salt);size(128)"`
}

// LoginRequest defines login request format
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse defines login response
type LoginResponse struct {
	Username    string             `json:"username"`
	UserID      int                `json:"userID"`
	Token       string             `json:"token"`
}

//CreateRequest defines create user request format
type CreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//CreateResponse defines create user response
type CreateResponse struct {
	UserID   int    `json:"userID"`
	Username string `json:"username"`
}

// DoLogin: user login
func DoLogin(lr *LoginRequest) (*LoginResponse, int, error) {
	// get username and password
	username := lr.Username
	password := lr.Password

	//validate username and password if is empty
	if len(username) == 0 || len(password) == 0 {
		return nil, http.StatusBadRequest, errors.New("error: username or password is empty")
	}

	// connect db
	o := orm.NewOrm()

	// check the username if existing
	user := &User{Username: username}
	err := o.Read(user, "username")
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("error: username is not existing")
	}

	// generate the password hash
	hash, err := GeneratePassHash(password, user.Salt)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	if hash != user.Password {
		return nil, http.StatusBadRequest, errors.New("error: password is error")
	}

	// generate token
	tokenString, err := GenerateToken(lr, user.Id, 0)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	return &LoginResponse{
		Username:    user.Username,
		UserID:      user.Id,
		Token:       tokenString,
	}, http.StatusOK, nil
}

// DoCreateUser: create a user
func DoCreateUser(cr *CreateRequest) (*CreateResponse, int, error) {
	// connect db
	o := orm.NewOrm()

	// check username if exist
	userNameCheck := User{Username: cr.Username}
	err := o.Read(&userNameCheck, "username")
	if err == nil {
		return nil, http.StatusBadRequest, errors.New("username has already existed")
	}

	// generate salt
	saltKey, err := GenerateSalt()
	if err != nil {
		logs.Info(err.Error())
		return nil, http.StatusBadRequest, err
	}

	// generate password hash
	hash, err := GeneratePassHash(cr.Password, saltKey)
	if err != nil {
		logs.Info(err.Error())
		return nil, http.StatusBadRequest, err
	}

	// create user
	user := User{}
	user.Username = cr.Username
	user.Password = hash
	user.Salt = saltKey

	_, err = o.Insert(&user)
	if err != nil {
		logs.Info(err.Error())
		return nil, http.StatusBadRequest, err
	}

	return &CreateResponse{
		UserID:   user.Id,
		Username: user.Username,
	}, http.StatusOK, nil

}
