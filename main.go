package main

import (
	"selFund/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// 新增用户
	router.POST("/addUser", controller.AddUser)
	// 获取用户基金列表
	router.GET("/getUserFund/:username", controller.GetUserFund)
	// 获取用户基金统计
	router.GET("/checkMyRepeatStock/:username", controller.CheckRepeatStock)
	// 统计新基金与持有基金的股票重合数
	router.POST("/compareNewFund", controller.CompareNewFund)
	// 为一名用户添加基金
	router.POST("/addUserFund", controller.AddUserFund)
	// 为一名用户删除基金
	router.POST("/deleteUserFund", controller.DeleteUserFund)
	// 获取基金详情
	router.GET("/getFund/:fund", controller.GetFund)
	router.Run()
}
