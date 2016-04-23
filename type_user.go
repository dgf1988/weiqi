package weiqi

import (
	"database/sql"
	"time"
)

const (
	MinLenUsername = 5
	MinLenPassword = 6
)

type User struct {
	Id       int64
	Name     string
	Password string
	Email    string
	Ip       string
	Status   int64
	Register time.Time
}

func (this User) RegisterTime() string {
	return this.Register.Format(ConstStdDatetime)
}

var (
	ErrUserExisted            = NewWeiqiError("用户已经存在")
	ErrUserNotFound           = NewWeiqiError("用户不存在")
	ErrUserNameTooShort       = NewWeiqiError("用户名太短")
	ErrUserPassword           = NewWeiqiError("密码错误")
	ErrUserPasswordTooShort   = NewWeiqiError("密码太短")
	ErrUserPasswordNotTheSame = NewWeiqiError("密码不一致")
)

func getPasswordMd5(password, ip string) string {
	return getMd5(password + ip)
}

//注册用户
func RegisterUser(username, password, password2, email, ip string) (int64, error) {
	if len(username) < MinLenUsername {
		return -1, ErrUserNameTooShort
	}
	if len(password) < MinLenPassword {
		return -1, ErrUserPasswordTooShort
	}
	if password != password2 {
		return -1, ErrUserPasswordNotTheSame
	}

	user := User{}
	err := Users.Get(nil, username).Struct(&user)
	if err == nil {
		return user.Id, ErrUserExisted
	} else if err == sql.ErrNoRows {
		return Users.Add(nil, username, getPasswordMd5(password, ip), email, ip)
	} else {
		return -1, err
	}
}

//验证用户
func loginUser(username, password string) (*User, error) {
	if len(username) < MinLenUsername {
		return nil, ErrUserNameTooShort
	}
	if len(password) < MinLenPassword {
		return nil, ErrUserPasswordTooShort
	}

	var user User
	err := Users.Get(nil, username).Struct(&user)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	if user.Password != getPasswordMd5(password, user.Ip) {
		return nil, ErrUserPassword
	}
	return &user, nil
}
