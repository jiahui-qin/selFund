package service

import (
	"fmt"
	tool "selFund/tool"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func init() {
	db, _ := tool.GetConn()
	db.AutoMigrate(&User{})
	db.AutoMigrate(&UserFund{})
}
func CompareNewFund(user, newFundCode string) (map[string]*StockWithFunds, error) {
	funds, err := GetUserFund(user)
	if err != nil {
		return nil, errors.Wrap(err, "get user fund error")
	}
	newFund, err := GetFund(newFundCode)
	if err != nil {
		return nil, errors.Wrap(err, "get new fund error")
	}
	for _, f := range funds {
		if newFundCode == f.Code {
			return nil, errors.New("already have this fund!")
		}
	}
	stockMap, err := getStockFromUserFundList(funds)
	if err != nil {
		return nil, errors.Wrap(err, "repeat fund")
	}
	newStockMap := make(map[string]*StockWithFunds)
	for _, v := range newFund.StockList {
		if stockMap[v.StockCode] != nil {
			newStockMap[v.StockCode] = stockMap[v.StockCode]
		}
	}
	return newStockMap, nil
}

func getStockFromUserFundList(funds []Fund) (map[string]*StockWithFunds, error) {
	stockMap := make(map[string]*StockWithFunds)
	for _, f := range funds {
		stocks, err := GetFundStocks(f.Code)
		if err != nil {
			return nil, errors.Wrap(err, "get fund detail "+f.Code+" open errors")
		}
		for _, s := range stocks {
			if stockMap[s.StockCode] == nil {
				stockdetail := &StockWithFunds{
					s.StockCode,
					s.Name,
					[]FundCodeName{{f.Name, f.Code}},
				}
				stockMap[s.StockCode] = stockdetail
			} else {
				fundCodeName := FundCodeName{f.Name, f.Code}
				stockMap[s.StockCode].RealtedFund = append(stockMap[s.StockCode].RealtedFund, fundCodeName)
			}
		}
	}
	return stockMap, nil
}

type StockWithFunds struct {
	StockCode   string         `gorm:"column:stock_code" db:"stock_code" json:"stock_code" form:"stock_code"`
	Name        string         `gorm:"column:name" db:"name" json:"name" form:"name"`
	RealtedFund []FundCodeName `json:"fund"`
}
type FundCodeName struct {
	Name string `gorm:"column:name" db:"name" json:"name" form:"name"`
	Code string `gorm:"column:code" db:"code" json:"code" form:"code"`
}

func AddUser(name, desc string) (User, error) {
	db, _ := tool.GetConn()
	user := &User{Name: name, Desc: desc}
	res := db.Debug().Create(user)
	if res.Error == nil {
		var user_in User
		db.Debug().Where(user).First(&user_in)
		fmt.Println(user_in)
		return user_in, nil
	}
	fmt.Println(res.Error)
	return User{}, res.Error
}

func GetUserFund(username string) ([]Fund, error) {
	userid, err := getUserId(username)
	if err != nil {
		return nil, err
	}
	var userFunds []UserFund
	sqldb, err := tool.GetConn()
	if err != nil {
		return nil, err
	}
	sqldb.Debug().Where(&UserFund{UserId: userid}).Find(&userFunds)
	var fundIds []int
	for _, uf := range userFunds {
		fundIds = append(fundIds, uf.FundId)
	}
	var funds []Fund
	if fundIds != nil {
		sqldb.Debug().Find(&funds, fundIds)
		return funds, nil
	}
	return nil, nil
}

func CheckMyRepeatStock(user string) (map[string]*UserStock, error) {
	funds, err := GetUserFund(user)
	if err != nil {
		return nil, err
	}
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
	for id, holdStock := range holdStockMap {
		if holdStock.HoldFundCount < 2 {
			delete(holdStockMap, id)
		}
	}
	return holdStockMap, nil
}

type UserStock struct {
	HoldStock
	HoldFundCount int `json:"holdFundcCount"`
}

func AddUserFund(user string, fund string) (string, error) {
	fundExist := CheckFundExist(fund)
	if !fundExist {
		InsertFund(fund)
	}
	userid, err := getUserId(user)
	if err != nil {
		return "", err
	}
	fundid := getFundId(fund)
	sqldb, _ := tool.GetConn()
	var userCount int64
	sqldb.Debug().Find(&UserFund{UserId: userid, FundId: fundid}).Count(&userCount)
	if userCount == 0 {
		sqldb.Debug().Create(&UserFund{UserId: userid, FundId: fundid})
	}
	return "add fund success", nil
}

func getUserId(name string) (int, error) {
	var user User
	sqldb, _ := tool.GetConn()
	res := sqldb.Debug().Where(&User{Name: name}).First(&user)
	if res.Error != nil {
		fmt.Print("not find user")
		return -1, errors.New("not find this user: " + name)
	}
	return int(user.ID), nil
}

func DeleteUserFund(username string, fundid string) (string, error) {
	userid, err := getUserId(username)
	if err != nil {
		return "", err
	}
	sqldb, _ := tool.GetConn()
	res := sqldb.Debug().Where("user_id = ? AND fund_id = ?", userid, getFundId(fundid)).Delete(&UserFund{})
	if res.Error != nil {
		return "", res.Error
	}
	return "delete success", nil
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
