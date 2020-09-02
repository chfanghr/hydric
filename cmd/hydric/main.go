package main

import (
	"flag"
	"github.com/chfanghr/hydric/core/auth"
	"github.com/chfanghr/hydric/core/config"
	"github.com/chfanghr/hydric/core/greeting"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var (
	configFile = flag.String("config", "config.json", "specify config file")
	debug      = flag.Bool("debug", true, "is in debug env or not")
)

func main() {
	flag.Parse()

	if !*debug {
		gin.SetMode(gin.ReleaseMode)
	}

	conf := config.MustFromFile(*configFile).MustParse()
	authService := auth.MustFromConfig(conf)

	router := gin.Default()
	router.Use(conf.MakeCheckMaintenanceStatusMiddleware())

	api := router.Group("/api")

	authService.SetupDefaultAuthAPI(api)
	api.GET("/greet/:name", authService.MakeIsAuthMiddleware(), greeting.GreetHandler)

	if err := http.Serve(conf.Listener, router); err != nil {
		log.Fatalln(err)
	}
}
