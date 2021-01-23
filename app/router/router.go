package router

import (
	"github.com/gin-gonic/gin"
	"powerorder/app/controller/toc/order/addition"
	"powerorder/app/controller/tool"
	"powerorder/app/middlewares"
)

//The routing method is exactly the same as Gin
func RegisterRouter(router *gin.Engine) {

	//can be used to enable CORS with various options
	router.Use(middlewares.CorsMiddleware())

	//健康检查
	Tool := router.Group("/healthcheck")
	Tool.GET("/ping", tool.HealthCheck)

	//order addition
	tocAddition := router.Group("toc/addition", middlewares.VerifySignMiddleware())

	tocAdditionSync := &addition.Sync{}
	tocAddition.POST("/sync", tocAdditionSync.Index)
	tocAdditionBegin := &addition.Begin{}
	tocAddition.POST("/begin", tocAdditionBegin.Index)
	tocAdditionCommit := &addition.Commit{}
	tocAddition.POST("/commit", tocAdditionCommit.Index)
	tocAdditionRollback := &addition.Rollback{}
	tocAddition.POST("/rollback", tocAdditionRollback.Index)

	////order get
	//tocSearchGet := &search.Get{}
	//tocSearch := router.Group("toc/search", middlewares.VerifySignMiddleware())
	//tocSearch.POST("/get", tocSearchGet.Index)
	//tocSearchSearch := &search.Search{}
	//tocSearch.POST("/search", tocSearchSearch.Index)
	//tocSearchQuery := &search.Query{}
	//tocSearch.POST("/query", tocSearchQuery.Index)

	////orderid
	//toOrderId := router.Group("toc/orderid", middlewares.VerifySignMiddleware())
	//tocOrderIdGen := &orderid.Generate{}
	//toOrderId.POST("/generate", tocOrderIdGen.Index)
	//tocOrderIdSearch := &orderid.Search{}
	//toOrderId.POST("/search", tocOrderIdSearch.Index)
	//
	//tocUpdateUpdate := &update.Update{}
	//toUpdate := router.Group("toc/update", middlewares.VerifySignMiddleware())
	//toUpdate.POST("/update", tocUpdateUpdate.Index)

	////order get
	//tobSearchGet := &tobsearch.Get{}
	//tobSearch := router.Group("tob/search", middlewares.VerifySignMiddleware())
	//tobSearch.POST("/get", tobSearchGet.Index)
	//tobSearchSearch := &tobsearch.Search{}
	//tobSearch.POST("/search", tobSearchSearch.Index)
	//tobSearchDSearch := &tobsearch.DSearch{}
	//tobSearch.POST("/dsearch", tobSearchDSearch.Index)

}
