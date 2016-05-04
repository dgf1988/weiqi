package weiqi

import (
	"database/sql"
	"time"
)

const (
	c_MinLenUsername = 5
	c_MinLenPassword = 6
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
	return this.Register.Format(constParseDatetimeStd)
}

var (
	ErrUserExisted            = newWeiqiError("用户已经存在")
	ErrUserNotFound           = newWeiqiError("用户不存在")
	ErrUserNameTooShort       = newWeiqiError("用户名太短")
	ErrUserPassword           = newWeiqiError("密码错误")
	ErrUserPasswordTooShort   = newWeiqiError("密码太短")
	ErrUserPasswordNotTheSame = newWeiqiError("密码不一致")
)

func md5Password(password, ip string) string {
	return md5String(password + ip)
}

//注册用户
func registerUser(username, password, password2, email, ip string) (int64, error) {
	if len(username) < c_MinLenUsername {
		return -1, ErrUserNameTooShort
	}
	if len(password) < c_MinLenPassword {
		return -1, ErrUserPasswordTooShort
	}
	if password != password2 {
		return -1, ErrUserPasswordNotTheSame
	}

	user := User{}
	err := Db.User.Get(nil, username).Struct(&user)
	if err == nil {
		return user.Id, ErrUserExisted
	} else if err == sql.ErrNoRows {
		return Db.User.Add(nil, username, md5Password(password, ip), email, ip)
	} else {
		return -1, err
	}
}

//验证用户
func loginUser(username, password string) (*User, error) {
	if len(username) < c_MinLenUsername {
		return nil, ErrUserNameTooShort
	}
	if len(password) < c_MinLenPassword {
		return nil, ErrUserPasswordTooShort
	}

	var user User
	err := Db.User.Get(nil, username).Struct(&user)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	if user.Password != md5Password(password, user.Ip) {
		return nil, ErrUserPassword
	}
	return &user, nil
}
