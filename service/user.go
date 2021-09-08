package service

import (
	"fmt"
	tool "selFund/tool"
	"strconv"

	"gorm.io/gorm"
)

func init() {
	db, _ := tool.GetConn()
	db.AutoMigrate(&User{})
	db.AutoMigrate(&UserFund{})
}
func AddUser(name, desc string) string {
	db, _ := tool.GetConn()
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

func GetUserFund(username string) []Fund {
	userid := getUserId(username)
	var userFunds []UserFund
	sqldb, _ := tool.GetConn()
	sqldb.Debug().Where(&UserFund{UserId: userid}).Find(&userFunds)
	var fundIds []int
	for _, uf := range userFunds {
		fundIds = append(fundIds, uf.FundId)
	}
	var funds []Fund
	if fundIds != nil {
		sqldb.Debug().Find(&funds, fundIds)
		return funds
	}
	return nil
}

func CheckMyRepeatStock(user string) map[string]*UserStock {
	funds := GetUserFund(user)
	sqldb, _ := tool.GetConn()
	holdStockMap := make(map[string]*UserStock)

	for _, fund := range funds {
		var holdStocks []HoldStock
		sqldb.Debug().Where(&HoldStock{FundCode: fund.Code}).Find(&holdStocks)
		for _, holdStock := range holdStocks {
			userStock, exist := holdStockMap[holdStock.Name]
			if exist {
				holdStockMap[holdStock.Name].HoldFundCount = userStock.HoldFundCount + 1
			} else {
				holdStockMap[holdStock.Name] = &UserStock{HoldStock: holdStock, HoldFundCount: 1}
			}
		}
	}
	return holdStockMap
}

type UserStock struct {
	HoldStock
	HoldFundCount int `json:"holdFundcCount"`
}

func AddUserFund(user string, fund string) string {
	fundExist := CheckFundExist(fund)
	if !fundExist {
		InsertFund(fund)
	}
	userid := getUserId(user)
	fundid := getFundId(fund)
	sqldb, _ := tool.GetConn()
	var userCount int64
	sqldb.Debug().Find(&UserFund{UserId: userid, FundId: fundid}).Count(&userCount)
	if userCount == 0 {
		sqldb.Debug().Create(&UserFund{UserId: userid, FundId: fundid})
	}
	return "ok"
}

func getUserId(name string) int {
	var user User
	sqldb, _ := tool.GetConn()
	sqldb.Debug().Where(&User{Name: name}).First(&user)
	return int(user.ID)
}

func DeleteUserFund(username string, fundid string) string {
	sqldb, _ := tool.GetConn()
	sqldb.Debug().Where("user_id = ? AND fund_id = ?", getUserId(username), getFundId(fundid)).Delete(&UserFund{})
	return "ok"
}

type UserFund struct {
	gorm.Model
	UserId int `gorm:"column:user_id;not null;primaryKey" db:"user_id" json:"user_id" form:"user_id"`
	FundId int `gorm:"column:fund_id;not null;primaryKey" db:"fund_id" json:"fund_id" form:"fund_id"`
}

type User struct {
	gorm.Model
	Name   string `gorm:"column:name; unique" db:"name" json:"name" form:"name"` //user name
	Desc   string `gorm:"column:desc" db:"desc" json:"desc" form:"desc"`         //user desc
	Amount int64  `gorm:"column:amount" db:"amount" json:"amount" form:"amount"`
}
