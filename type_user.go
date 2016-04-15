package weiqi

import (
	"database/sql"
	"time"
)

const (
	MinLenUsername = 5
	MinLenPassword = 6
)

type U struct {
	Id       int64
	Name     string
	Email    string
	Ip       string
	Status   int64
	Register time.Time
}

func (this U) RegisterTime() string {
	return this.Register.Format(Time_Def_Format)
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
		return 0, ErrUserNameTooShort
	}
	if len(password) < MinLenPassword {
		return 0, ErrUserPasswordTooShort
	}
	if password != password2 {
		return 0, ErrUserPasswordNotTheSame
	}
	var id int64
	row := db.QueryRow("select user.id from user where uname = ? limit 1", username)
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		res, err := db.Exec("insert into user (uname, upassword, uemail, uip) values(?,?,?, ?)", username, getPasswordMd5(password, ip), email, ip)
		if err != nil {
			return 0, err
		}
		return res.LastInsertId()
	} else if err != nil {
		return 0, err
	} else {
		return id, ErrUserExisted
	}
}

//验证用户
func loginUser(username, password string) (*U, error) {
	if len(username) < MinLenUsername {
		return nil, ErrUserNameTooShort
	}
	if len(password) < MinLenPassword {
		return nil, ErrUserPasswordTooShort
	}

	row := db.QueryRow("select user.id, user.uname, user.upassword, user.uemail, user.uip, user.ustatus, user.uregister from user where user.uname=? limit 1", username)
	var u U
	var upassword string
	err := row.Scan(&u.Id, &u.Name, &upassword, &u.Email, &u.Ip, &u.Status, &u.Register)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	if upassword != getPasswordMd5(password, u.Ip) {
		return nil, ErrUserPassword
	}
	return &u, nil
}
