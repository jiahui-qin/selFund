package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
)

func GetFundInfo(fundCode string) Fund {
	// jsons, _ := json.Marshal(fundPosition)
	db, _ := sql.Open("mysql", "root:123456@tcp(localhost:3306)/fund")
	var fund Fund
	db.QueryRow("select * from fund where code = ?", fundCode).Scan(&fund.Id, &fund.Name, &fund.Code, &fund.Bond, &fund.Cash, &fund.Stock, &fund.Title, &fund.Total)

	rows, err := db.Query("select * from fund where code = ?", fundCode)
	for rows.Next() {
		var fund1 Fund
		GetData(rows, &fund1)
		fmt.Println(fund1)
	}

	defer db.Close()
	if err == sql.ErrNoRows {
		fmt.Println(fundCode)
		fundPosition := getFundPosition("https://api.doctorxiong.club/v1/fund/position?code=" + fundCode)
		savePosition(fundPosition, fundCode)
	}
	return fund
}

type Fund struct {
	Id    int    `gorm:"column:id" db:"id" json:"id" form:"id"`
	Name  string `gorm:"column:name" db:"name" json:"name" form:"name"`
	Code  string `gorm:"column:code" db:"code" json:"code" form:"code"`
	Bond  string `gorm:"column:bond" db:"bond" json:"bond" form:"bond"`
	Cash  string `gorm:"column:cash" db:"cash" json:"cash" form:"cash"`
	Stock string `gorm:"column:stock" db:"stock" json:"stock" form:"stock"`
	Title string `gorm:"column:title" db:"title" json:"title" form:"title"`
	Total string `gorm:"column:total" db:"total" json:"total" form:"total"`
}

func savePosition(fund Position, fundCode string) {
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/fund")
	if err != nil {
		panic(err.Error())
	}
	for _, v := range fund.Data.StockList {
		res, _ := db.Exec("INSERT INTO `fund`.`hold_stock` (`fund_code`, `stock_code`, `name`, `precent`, `hold`, `hold_amount`) VALUES(?,?,?,?,?,?)", fundCode, v[0], v[1], v[2], v[3], v[4])
		fmt.Println(res)
	}
	db.Exec("INSERT INTO `fund`.`fund` (`name`, `code`, `bond`, `cash`, `stock`, `title`, `total`) VALUES( ?, ?,?,?,?,?,?)", fund.Data.Title, fundCode, fund.Data.Bond, fund.Data.Cash, fund.Data.Stock, fund.Data.Title, fund.Data.Total)
	defer db.Close()
}

func getFundPosition(url string) Position {
	res, _ := http.Get(url)
	robots, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	var tmp Position
	json.Unmarshal([]byte(string(robots)), &tmp)
	return tmp
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

type Stock struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	Precent    string `json:"precent"`
	Hold       string `json:"hold"`
	HoldAmount string `json:"holdAmount"`
}

//
// func newStock(code, name, precent, hold, holdAmount string) *Stock {
// 	return &Stock{
// 		Code:       code,
// 		Name:       name,
// 		Precent:    precent,
// 		Hold:       hold,
// 		HoldAmount: holdAmount,
// 	}
// }
//
func GetData(rows *sql.Rows, dest interface{}) error {
	// 取得資料的每一列的名稱

	col_names, err := rows.Columns()
	if err != nil {
		return err
	}
	// 取得變數對象的值跟類型資訊
	v := reflect.ValueOf(dest)
	if v.Elem().Type().Kind() != reflect.Struct {
		return errors.New("give me  a struct")
	}
	// 宣告一個interface{}的slice
	scan_dest := []interface{}{}
	// 建立一個string, interface{}的map
	addr_by_col_name := map[string]interface{}{}

	for i := 0; i < v.Elem().NumField(); i++ {
		propertyName := v.Elem().Field(i)
		col_name := v.Elem().Type().Field(i).Tag.Get("db")
		if col_name == "" {
			if v.Elem().Field(i).CanInterface() == false {
				continue
			}
			col_name = propertyName.Type().Name()
		}
		// Addr() 返回該屬性的記憶體位置的指針
		// Interface() 返回該屬性真正的值, 這裡還是存著位置
		addr_by_col_name[col_name] = propertyName.Addr().Interface()
	}
	// 把實際各成員屬性的位置, 給加到scan_dest中
	for _, col_name := range col_names {
		scan_dest = append(scan_dest, addr_by_col_name[col_name])
	}
	// 執行Scan
	return rows.Scan(scan_dest...)
}
