package controller

import (
	"selFund/app"
	se "selFund/service"

	"github.com/gin-gonic/gin"
)

func GetFund(c *gin.Context) {
	appG := app.Gin{C: c}
	fund := c.Param("fund")
	res, err := se.GetFund(fund)
	if err != nil {
		appG.ResponseError("ERROR", err)
		return
	}
	appG.ResponseSuccess(res)
}
