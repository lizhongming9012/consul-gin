package routers

import (
	"NULL/consul-gin/middleware/cors"
	"NULL/consul-gin/pkg/export"
	"NULL/consul-gin/pkg/qrcode"
	"NULL/consul-gin/pkg/upload"
	"NULL/consul-gin/routers/api"
	v1 "NULL/consul-gin/routers/api/v1"
	"github.com/gin-gonic/gin"
	"net/http"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.CORSMiddleware())

	r.GET("/health", v1.Health)
	r.GET("/agips", v1.GetAgIPs)
	r.GET("/srvips", v1.GetSrvsIps)
	r.GET("/clearsrv", v1.DeregisterService)
	r.GET("/kvs", v1.GetKVs)
	r.GET("/kv", v1.GetKV)
	r.GET("/putkv", v1.PutKV)
	r.GET("/delkv", v1.DeleteKV)

	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("/update", http.Dir(upload.GetUpdateFullPath()))
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))
	r.Static("/css", "runtime/static/css")
	r.Static("/js", "runtime/static/js")
	r.Static("/img", "runtime/static/img")
	apiv1 := r.Group("/api/v1")
	//apiv1.Use(jwt.JWT())
	{
		//上传文件
		apiv1.POST("/file/upload", api.UploadFile)
		//文件下载
		apiv1.StaticFS("/file/download", http.Dir(upload.GetImageFullPath()))

	}
	return r
}
