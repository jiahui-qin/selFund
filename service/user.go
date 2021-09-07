package service

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

func AddUser(name, desc string) string {
	db, _ := getConn()
	db.AutoMigrate(&User{})
	user := &User{Name: name, Desc: desc}
	res := db.Debug().Create(user)
	if res.Error == nil {
		var user_in User
		db.Debug().Where(user).First(&user_in)
		fmt.Println(user_in)
		return strconv.Itoa(int(user_in.ID))
	}
	fmt.Println(res.Error)
	return fmt.Sprintf("%d", res.Error)
}

func GetUserFund(username string) []UserFund {
	userid := getUserId(username)
	var userFunds []UserFund
	sqldb, _ := getConn()
	sqldb.Where(&UserFund{UserId: userid}).Find(&userFunds)
	return userFunds
}

func AddUserFund(user string, fund string) string {
	fundid := CheckFundExist(fund)
	userid := getUserId(user)
	if fundid == 0 {
		fundid = SavePosition(fund)
	}
	sqldb, _ := getConn()
	sqldb.Debug().Create(&UserFund{UserId: userid, FundId: fundid})
	return "ok"
}

func getUserId(name string) int {

	var user User
	sqldb, _ := getConn()
	sqldb.Debug().Where(&User{Name: name}).First(&user)
	return int(user.ID)
}

func DeleteUserFund(username string, fundid string) string {
	sqldb, _ := getConn()
	sqldb.Debug().Where("user_id = ? AND fund_id = ?", getUserId(username), getFundId(fundid)).Delete(&UserFund{})
	return "ok"
}

type UserFund struct {
	gorm.Model
	UserId int `gorm:"column:user_id;not null" db:"user_id" json:"user_id" form:"user_id"`
	FundId int `gorm:"column:fund_id;not null " db:"fund_id" json:"fund_id" form:"fund_id"`
}

type User struct {
	gorm.Model
	Name   string `gorm:"column:name; unique" db:"name" json:"name" form:"name"` //user name
	Desc   string `gorm:"column:desc" db:"desc" json:"desc" form:"desc"`         //user desc
	Amount int64  `gorm:"column:amount" db:"amount" json:"amount" form:"amount"`
}
