package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

// Response setting gin.JSON
func (g *Gin) ResponseError(errorMsg string, data interface{}) {
	g.C.JSON(http.StatusInternalServerError, Response{
		Code: 500,
		Msg:  GetMsg(errorMsg),
		Data: data,
	})
}

func (g *Gin) ResponseSuccess(data interface{}) {
	g.C.JSON(http.StatusOK, Response{
		Code: 200,
		Data: data,
	})
}
