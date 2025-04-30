package main

import (
	"decleration/dao"
	"decleration/kafka"
	"decleration/lock"
	"decleration/middleware"
	"decleration/router"
	"decleration/service"
	"decleration/utils"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	dao.InitDB1()                                                   // 初始化数据库
	kafka.InitProducer([]string{"localhost:9092"}, "seckill-order") //初始化kafka
	utils.InitLogger()                                              //初始化log
	go kafka.ConsumeOrderMessages(service.CreateOrder)
	dao.AutoMigrate() //创建对应数据库
	dao.InitRedis()   // 初始化 Redis
	lock.InitLock()   // 初始化 Redsync 分布式锁
	h := server.Default()
	h.Use(middleware.RateLimitMiddleware(5)) //限流中间件
	h.Use(middleware.LogMiddleware())        //日志中间件
	router.RegisterRouter(h)
	h.Spin()
}
