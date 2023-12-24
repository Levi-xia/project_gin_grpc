package dao

import (
	"errors"
	"com.levi/project-common/base"
)

type userDao struct {
}

var UserDao = new(userDao)


type User struct {
    base.ID
    Name string `json:"name" gorm:"not null;comment:用户名称"`
    Mobile string `json:"mobile" gorm:"not null;index;comment:用户手机号"`
    Password string `json:"password" gorm:"not null;default:'';comment:用户密码"`
    base.Timestamps
    base.SoftDeletes
}

func (userDao *userDao) GetUserInfo(id int) (user User, err error) {
	err = base.Mysql.Find(&user, id).Error
    if err != nil {
        err = errors.New("DB Query Failed")
    }
	return
}