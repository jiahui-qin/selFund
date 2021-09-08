package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	tool "selFund/tool"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func init() {
	db, _ := tool.GetConn()
	db.AutoMigrate(&Fund{})
	db.AutoMigrate(&HoldStock{})
}

func getFundPosition(url string) (Position, error) {
	res, err1 := http.Get(url)
	if err1 != nil {
		return Position{}, err1
	}
	robots, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		return Position{}, err2
	}
	res.Body.Close()
	var tmp Position
	json.Unmarshal([]byte(string(robots)), &tmp)
	fmt.Println("get position: ")
	fmt.Println(tmp)
	return tmp, nil
}

func InsertFund(fundCode string) (FundVO, error) {
	pos, err := getFundPosition("https://api.doctorxiong.club/v1/fund/position?code=" + fundCode)
	db, _ := tool.GetConn()
	if err != nil || pos.Code != 200 {
		return FundVO{}, err
	}
	var fund = Fund{Name: pos.Data.Title, Code: fundCode, Bond: pos.Data.Bond, Cash: pos.Data.Cash, Stock: pos.Data.Stock, Title: pos.Data.Title, Total: pos.Data.Total}
	db.Debug().Create(&fund)
	for _, v := range pos.Data.StockList {
		var holdStock = HoldStock{FundCode: fundCode, StockCode: v[0], Name: v[1], Precent: v[2], Hold: v[3], HoldAmount: v[4]}
		db.Debug().Create(&holdStock)
	}
	return GetFund(fundCode), nil
}

func CheckFundExist(fundCode string) bool {
	db, _ := tool.GetConn()
	var count int64
	db.Debug().Model(&Fund{}).Where(&Fund{Code: fundCode}).Count(&count)
	fmt.Println(count)
	return count != 0
}

func GetFund(fundCode string) FundVO {
	db, _ := tool.GetConn()

	if !CheckFundExist(fundCode) {
		fmt.Println("need insert fund!")
		InsertFund(fundCode)
	}
	var fundVO FundVO
	db.Debug().Model(&Fund{}).Where(&Fund{Code: fundCode}).First(&fundVO)
	fundVO.StockList = GetFundStocks(fundCode)
	return fundVO
}

func GetFundStocks(fundCode string) []StockVO {
	if !CheckFundExist(fundCode) {
		fmt.Println("need insert fund!")
		InsertFund(fundCode)
	}
	db, _ := tool.GetConn()
	var holdStocks []StockVO
	db.Debug().Model(&HoldStock{}).Where(&HoldStock{FundCode: fundCode}).Find(&holdStocks)
	return holdStocks
}

func getFundId(fundCode string) int {
	db, _ := tool.GetConn()
	var id int
	db.Debug().Model(&Fund{}).Where(&Fund{Code: fundCode}).Select("id").Find(&id)
	return id
}

type FundVO struct {
	Name      string    `gorm:"column:name" db:"name" json:"name" form:"name"`
	Code      string    `gorm:"column:code" db:"code" json:"code" form:"code"`
	Bond      string    `gorm:"column:bond" db:"bond" json:"bond" form:"bond"`
	Cash      string    `gorm:"column:cash" db:"cash" json:"cash" form:"cash"`
	Stock     string    `gorm:"column:stock" db:"stock" json:"stock" form:"stock"`
	Title     string    `gorm:"column:title" db:"title" json:"title" form:"title"`
	Total     string    `gorm:"column:total" db:"total" json:"total" form:"total"`
	StockList []StockVO `gorm:"-" json:"stockList"`
}

type StockVO struct {
	StockCode  string `gorm:"column:stock_code" db:"stock_code" json:"stock_code" form:"stock_code"`
	Name       string `gorm:"column:name" db:"name" json:"name" form:"name"`
	Precent    string `gorm:"column:precent" db:"precent" json:"precent" form:"precent"`
	Hold       string `gorm:"column:hold" db:"hold" json:"hold" form:"hold"`
	HoldAmount string `gorm:"column:hold_amount" db:"hold_amount" json:"hold_amount" form:"hold_amount"`
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
