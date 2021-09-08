package main

import (
	"net/http"

	se "selFund/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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

func main() {
	router := gin.Default()

	router.POST("/addUser", func(c *gin.Context) {
		var json UserInfo
		if err := c.ShouldBindBodyWith(&json, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, se.AddUser(json.User, json.Desc))
	})
	// 获取用户基金列表
	router.GET("/getUserFund/:username", func(c *gin.Context) {
		username := c.Param("username")
		c.JSON(http.StatusOK, se.GetUserFund(username))
	})

	// 获取用户基金统计
	router.GET("/checkMyRepeatStock/:username", func(c *gin.Context) {
		username := c.Param("username")
		c.JSON(http.StatusOK, se.CheckMyRepeatStock(username))
	})
	// 为一名用户添加基金
	router.POST("/addUserFund", func(c *gin.Context) {
		var json UserFund
		if err := c.ShouldBindBodyWith(&json, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, se.AddUserFund(json.User, json.Fund))
	})
	// 为一名用户删除基金
	router.POST("/deleteUserFund", func(c *gin.Context) {
		var json UserFund
		if err := c.ShouldBindBodyWith(&json, binding.JSON); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, se.DeleteUserFund(json.User, json.Fund))
	})
	// 获取基金详情
	router.GET("/getFundInfo/:fund", func(c *gin.Context) {
		fund := c.Param("fund")
		c.JSON(http.StatusOK, se.GetFundInfo(fund))
	})
	router.Run()
}
