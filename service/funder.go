package service

type funder interface {
	InsertFund(fundCode string) bool
	CheckFundExist(fundCode string) bool
	GetFund(fundCode string) Fund
	GetFundWithStock(fundCode string) FundWithStock
}
