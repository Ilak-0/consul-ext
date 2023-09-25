package v1

import (
	"consul-ext/pkg/consul"
	"consul-ext/service"
	"github.com/gin-gonic/gin"

	"net/http"
	"time"
)

// @Summary restore services from mysql to consul
// @Tags backupConsulSvcs
// @Produce  json
// @Success 200 {object} resp.Response
// @Router /api/v1/consul-ext/svc/restore [put]
// @Param svcName query string true "service name,eg:service1,service2	 use 'all' to restore all services"
// @Param backupTime query string true "consul time,eg: 2023-09-01"
// @Param readConsulAddress query string true "read consul agent address"
// @Param writeConsulAddress query string false "write consul agent address"
// @Param deleteCurrentSvcs query bool true "whether delete current svcs"
func restoreSvcs(ctx *gin.Context) {
	params := new(consul.SvcsRestoreParams)
	if err := ctx.ShouldBindJSON(params); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if params.WriteConsulAddress == "" {
		params.WriteConsulAddress = params.ReadConsulAddress
	}
	if params.BackupTime.IsZero() {
		params.BackupTime = time.Now()
	}
	err := service.ConsulSvcsRestore(params)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, "svcs restore success")
}
