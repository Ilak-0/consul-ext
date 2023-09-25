package v1

import (
	"consul-ext/pkg/git_repo"
	"consul-ext/pkg/models"
	"consul-ext/service"
	"github.com/gin-gonic/gin"
	gitlabSDK "github.com/xanzy/go-gitlab"
	"log"
	"net/http"
)

// @Summary webhook sync git repo change to consul
// @Tags ConsulKV
// @Produce  json
// @Success 200 {object} resp.Response
// @Router /api/v1/consul-ext/{repoType}/webhook [post]
// @Param repoType path string true "repo type eg:gitlab,gitea only use one of them"
func webhook(ctx *gin.Context) {
	repo := ctx.Param("repoType")
	var err error
	event := git_repo.APIFunc().NewEvent()
	if err := ctx.ShouldBindJSON(event); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	switch repo {
	case "gitea":
		err = service.UpdateFromGiteaRepo(event.(*models.GiteaPayloadCommit))
	case "gitlab":
		err = service.UpdateFromGitLabRepo(event.(*gitlabSDK.PushEvent))
	}
	if err != nil {
		log.Println("webhook update kv failed", err)
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, "kv update success")
}
