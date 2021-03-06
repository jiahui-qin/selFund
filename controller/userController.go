package controller

import (
	"net/http"
	"selFund/app"
	se "selFund/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func AddUser(c *gin.Context) {
	appG := app.Gin{C: c}
	var json UserInfo
	if err := c.ShouldBindBodyWith(&json, binding.JSON); err != nil {
		appG.ResponseError("PARA_ANALYSIS_ERROR", err)
		return
	}
	msg, err := se.AddUser(json.User, json.Desc)
	if err != nil {
		appG.ResponseError("DATA_INSERT_ERROR", err)
		return
	}
	appG.ResponseSuccess(map[string]interface{}{
		"msg": msg,
	})
}

func GetUserFund(c *gin.Context) {
	appG := app.Gin{C: c}
	username := c.Param("username")
	res, err := se.GetUserFund(username)
	if err != nil {
		appG.ResponseError("ERROR", err.Error())
		return
	}
	appG.ResponseSuccess(res)
}

func AddUserFund(c *gin.Context) {
	appG := app.Gin{C: c}
	var json UserFund
	if err := c.ShouldBindBodyWith(&json, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := se.AddUserFund(json.User, json.Fund)
	if err != nil {
		appG.ResponseError("ERROR", err)
		return
	}
	appG.ResponseSuccess(res)
}

func DeleteUserFund(c *gin.Context) {
	appG := app.Gin{C: c}
	var json UserFund
	if err := c.ShouldBindBodyWith(&json, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := se.DeleteUserFund(json.User, json.Fund)
	if err != nil {
		appG.ResponseError("ERROR", err)
		return
	}
	appG.ResponseSuccess(res)
}
