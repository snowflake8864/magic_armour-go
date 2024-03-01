package main

import (
	"goskeleton/app/global/variable"
	_ "goskeleton/bootstrap"
	"goskeleton/routers"
)

// 这里可以存放后端路由（例如后台管理系统）
func main() {
	router := routers.InitPlatformRouter()
	_ = router.Run(variable.ConfigYml.GetString("HttpServer.Platform.Port"))
}
