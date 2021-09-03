package service

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetUserFund(userId int) []UserFund {
	dsn := "root:123456@tcp(localhost:3306)/fund?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&UserFund{})
	var userFunds []UserFund
	db.Debug().Where(&UserFund{UserId: userId}).Find(&userFunds)
	return userFunds
}

func AddUserFund(user int, fund string) string {
	dsn := "root:123456@tcp(localhost:3306)/fund?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	fundid := CheckFundExist(fund)
	if fundid == nil {
		//插入fundid！
	}
	db.Debug().Create(&UserFund{UserId: user, FundId: fund})
}

func DeleteUserFund(user, fund string) string {
	return "user is " + user
}

type UserFund struct {
	gorm.Model
	UserId int `gorm:"column:user_id" db:"user_id" json:"user_id" form:"user_id"`
	FundId int `gorm:"column:fund_id" db:"fund_id" json:"fund_id" form:"fund_id"`
}
