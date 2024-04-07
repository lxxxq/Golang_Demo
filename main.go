package main

import (
	"demo/gin_common"
	"demo/mysql_common/mysql"
	"github.com/gin-gonic/gin"
)

func main() {

	//初始化mysql
	mysql.InitMysqlFile()

	//初始化Gin引擎
	router := gin.Default()

	//调用DefineRoutes来注册路由
	gin_common.DefineRoutes(router)

	//启动Gin服务器
	router.Run(":8080")

}
