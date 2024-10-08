package main

import (
	"shortly/internal/app/api/routers"
	"shortly/internal/app/config"
)

func main() {
	appConfig := config.Init()
	routers.Run(appConfig)
}
