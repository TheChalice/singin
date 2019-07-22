package main

import (
	"fmt"
	"github.com/TheChalice/singin/handler"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {

	router := handle()
	s := &http.Server{
		Addr:           ":10086",
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 0,
	}
	fmt.Println("START TOKEN SERVER")
	s.ListenAndServe()

}

func handle() (router *gin.Engine) {
	//设置全局环境：1.开发环境（datafoundry_docker.DebugMode） 2.线上环境（datafoundry_docker.ReleaseMode）
	gin.SetMode(gin.DebugMode)
	//获取路由实例
	router = gin.Default()

	router.POST("/signin", handler.GetToken)


	//未知调用方式
	router.NoMethod(func(context *gin.Context) {
		context.String(404, "Not method")
	})

	return
}
