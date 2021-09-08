package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

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

func AddUser(c *gin.Context) {
	var json UserInfo
	if err := c.ShouldBindBodyWith(&json, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, se.AddUser(json.User, json.Desc))
}

func GetUserFund(c *gin.Context) {
	username := c.Param("username")
	c.JSON(http.StatusOK, se.CheckMyRepeatStock(username))
}

func CheckRepeatStock(c *gin.Context) {
	username := c.Param("username")
	c.JSON(http.StatusOK, se.CheckMyRepeatStock(username))
}

func AddUserFund(c *gin.Context) {
	var json UserFund
	if err := c.ShouldBindBodyWith(&json, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, se.AddUserFund(json.User, json.Fund))
}

func DeleteUserFund(c *gin.Context) {
	var json UserFund
	if err := c.ShouldBindBodyWith(&json, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, se.DeleteUserFund(json.User, json.Fund))
}

func GetFund(c *gin.Context) {
	fund := c.Param("fund")
	c.JSON(http.StatusOK, se.GetFund(fund))
}
