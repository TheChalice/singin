package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/TheChalice/singin/handler"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "10.1.236.50:10086",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := client.Ping().Result()

	fmt.Println(pong, err)
}

func main() {
	handler.Setkey(client, "a", "b")
	value := handler.Get(client, "a")
	fmt.Println("a", value)
	router := handle()
	s := &http.Server{
		Addr:           ":10080",
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

	router.GET("/cluster", handler.Cluster)

	//未知调用方式
	router.NoMethod(func(context *gin.Context) {
		context.String(404, "Not method")
	})

	return
}
