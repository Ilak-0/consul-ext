package v1

import "consul-ext/router/app"

func init() {
	apiResource := app.Get().Group("/api/v1/consul-ext")
	{
		apiResource.PUT("/svc/restore", restoreSvcs) // restore services
		apiResource.GET("/kv", backupAllKV)
		apiResource.PUT("/path/file", filePut)
		apiResource.POST("/:repoType/webhook", webhook)
	}
}
