package v1

import (
	"consul-ext/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary sync all consul kv to git repo
// @Tags ConsulKV
// @Produce  json
// @Success 200 {object} resp.Response
// @Router /api/v1/consul-ext/kv [get]
func backupAllKV(ctx *gin.Context) {
	err := service.BackupAllConsulKv()
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, "kv backup success")
}

// @Summary sync git repo to consul
// @Tags ConsulKV
// @Produce  json
// @Success 200 {object} resp.Response
// @Router /api/v1/consul-ext/path/file [put]
// @Param path query string true "path eg:/path/to/file"
func filePut(ctx *gin.Context) {
	path := ctx.Query("path")
	err := service.UpdateKV(path)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, "kv update success")
}
