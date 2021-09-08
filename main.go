package main

import (
	"selFund/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/addUser", controller.AddUser)

	// 获取用户基金列表
	router.GET("/getUserFund/:username", controller.GetUserFund)

	// 获取用户基金统计
	router.GET("/checkMyRepeatStock/:username", controller.CheckRepeatStock)
	// 为一名用户添加基金
	router.POST("/addUserFund", controller.AddUserFund)
	// 为一名用户删除基金
	router.POST("/deleteUserFund", controller.DeleteUserFund)
	// 获取基金详情
	router.GET("/getFund/:fund", controller.GetFund)
	router.Run()
}
