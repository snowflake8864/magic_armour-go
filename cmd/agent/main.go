package main

import (
	"goskeleton/app/global/variable"
	_ "goskeleton/bootstrap"
	"goskeleton/routers"
)

// 这里可以存放门户类网站入口
func main() {
	router := routers.InitAgentRouter()
	_ = router.Run(variable.ConfigYml.GetString("HttpServer.Agent.Port"))
}
