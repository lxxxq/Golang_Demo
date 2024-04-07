package gin_common

import (
	"demo/service"
	"github.com/gin-gonic/gin"
)

func DefineRoutes(router *gin.Engine) {
	router.GET("/students/select/:id", service.GetByID)
	router.GET("/students/select", service.GetAll)
	router.DELETE("/students/delete/:id", service.DropById)
	router.PUT("/students/update", service.AlterByIds)
	router.POST("/students/insert", service.BatchInsert)

}
