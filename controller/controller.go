package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"selFund/app"
	se "selFund/service"
)

type UserFund struct {
	User   string `form:"User" json:"User" xml:"User"  binding:"required"`
	Fund   string `form:"Fund" json:"Fund" xml:"Fund" binding:"required"`
	Amount string `form:"Amount" json:"Amount" xml:"Amount" binding:"-"`
}

type UserInfo struct {
	User string `form:"User" json:"User" xml:"User"  binding:"required"`
	Desc string `form:"Desc" json:"Desc" xml:"Desc" binding:"required"`
}

func CheckRepeatStock(c *gin.Context) {
	appG := app.Gin{C: c}
	username := c.Param("username")
	res, err := se.CheckMyRepeatStock(username)
	if err != nil {
		appG.ResponseError("ERROR", err)
		return
	}
	appG.ResponseSuccess(res)
}

func CompareNewFund(c *gin.Context) {
	appG := app.Gin{C: c}
	var json UserFund
	if err := c.ShouldBindBodyWith(&json, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := se.CompareNewFund(json.User, json.Fund)
	if err != nil {
		appG.ResponseError("ERROR", err.Error())
		return
	}
	appG.ResponseSuccess(res)
}
