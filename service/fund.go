package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
)

func GetFundInfo(fundCode string) []Fund {
	dsn := "root:123456@tcp(localhost:3306)/fund?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&Fund{})
	db.AutoMigrate(&HoldStock{})
	var funds []Fund
	var count int64
	db.Debug().Where(&Fund{Code: fundCode}).Find(&funds).Count(&count)
	if count == 0 {
		fundPosition := getFundPosition("https://api.doctorxiong.club/v1/fund/position?code=" + fundCode)
		savePosition(fundPosition, fundCode)
	}
	db.Debug().Where(&Fund{Code: fundCode}).Find(&funds)
	return funds
}

func savePosition(pos Position, fundCode string) {
	dsn := "root:123456@tcp(localhost:3306)/fund?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	var fund = Fund{Name: pos.Data.Title, Code: fundCode, Bond: pos.Data.Bond, Cash: pos.Data.Cash, Stock: pos.Data.Stock, Title: pos.Data.Title, Total: pos.Data.Total}
	db.Debug().Create(&fund)
	for _, v := range pos.Data.StockList {
		var holdStock = HoldStock{FundCode: fundCode, StockCode: v[0], Name: v[1], Precent: v[2], Hold: v[3], HoldAmount: v[4]}
		db.Debug().Create(&holdStock)
	}
}

func getFundPosition(url string) Position {
	res, _ := http.Get(url)
	robots, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	var tmp Position
	json.Unmarshal([]byte(string(robots)), &tmp)
	return tmp
}

func CheckFundExist(fundCode string) int {
	dsn := "root:123456@tcp(localhost:3306)/fund?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	var fund Fund
	db.Debug().Where(&Fund{Code: fundCode}).Find(&fund)
	return int(fund.ID)
}

type Position struct {
	Code int64 `json:"code"`
	Data struct {
		Bond      string     `json:"bond"`
		Cash      string     `json:"cash"`
		Date      string     `json:"date"`
		Stock     string     `json:"stock"`
		StockList [][]string `json:"stockList"`
		Title     string     `json:"title"`
		Total     string     `json:"total"`
	} `json:"data"`
	Message string      `json:"message"`
	Meta    interface{} `json:"meta"`
}

type HoldStock struct {
	gorm.Model
	FundCode   string `gorm:"column:fund_code" db:"fund_code" json:"fund_code" form:"fund_code"`
	StockCode  string `gorm:"column:stock_code" db:"stock_code" json:"stock_code" form:"stock_code"`
	Name       string `gorm:"column:name" db:"name" json:"name" form:"name"`
	Precent    string `gorm:"column:precent" db:"precent" json:"precent" form:"precent"`
	Hold       string `gorm:"column:hold" db:"hold" json:"hold" form:"hold"`
	HoldAmount string `gorm:"column:hold_amount" db:"hold_amount" json:"hold_amount" form:"hold_amount"`
}

type Fund struct {
	gorm.Model
	Name  string `gorm:"column:name" db:"name" json:"name" form:"name"`
	Code  string `gorm:"column:code" db:"code" json:"code" form:"code"`
	Bond  string `gorm:"column:bond" db:"bond" json:"bond" form:"bond"`
	Cash  string `gorm:"column:cash" db:"cash" json:"cash" form:"cash"`
	Stock string `gorm:"column:stock" db:"stock" json:"stock" form:"stock"`
	Title string `gorm:"column:title" db:"title" json:"title" form:"title"`
	Total string `gorm:"column:total" db:"total" json:"total" form:"total"`
}
