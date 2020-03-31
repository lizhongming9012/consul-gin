package v1

import (
	"NULL/consul-gin/pkg/app"
	"NULL/consul-gin/pkg/consul"
	"NULL/consul-gin/pkg/e"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Health(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, e.SUCCESS, "service running healthily...")
}
func GetAgIPs(c *gin.Context) {
	appG := app.Gin{C: c}
	agIPs, err := consul.GetAgentIPs()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, err)
	}
	appG.Response(http.StatusOK, e.SUCCESS, agIPs)
}
func GetSrvsIps(c *gin.Context) {
	appG := app.Gin{C: c}
	srvIps := make([]string, 0)
	srvs := consul.GetHealthServices()
	for _, srv := range srvs {
		ips := consul.CatalogService(srv)
		srvIps = append(srvIps, ips...)
	}
	appG.Response(http.StatusOK, e.SUCCESS, srvIps)
}
func DeregisterService(c *gin.Context) {
	appG := app.Gin{C: c}
	srvID := c.Query("srvID")
	if err := consul.DeregisterService(srvID); err != nil {
		appG.Response(http.StatusOK, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
func GetKVs(c *gin.Context) {
	appG := app.Gin{C: c}
	kvs, err := consul.GetKVs()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, err)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, kvs)
}
func GetKV(c *gin.Context) {
	appG := app.Gin{C: c}
	key := c.Query("key")
	val := consul.GetKV(key)
	appG.Response(http.StatusOK, e.SUCCESS, val)
}
func PutKV(c *gin.Context) {
	appG := app.Gin{C: c}
	key := c.Query("key")
	val := c.Query("val")
	appG.Response(http.StatusOK, e.SUCCESS, consul.PutKV(key, val))
}
func DeleteKV(c *gin.Context) {
	appG := app.Gin{C: c}
	key := c.Query("key")
	appG.Response(http.StatusOK, e.SUCCESS, consul.DeleteKV(key))
}
