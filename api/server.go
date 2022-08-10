package api

import (
	"context"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/ssksameer56/CloudIndexer/cloud"
	"github.com/ssksameer56/CloudIndexer/config"
	"github.com/ssksameer56/CloudIndexer/controllers"
	"github.com/ssksameer56/CloudIndexer/elasticservice"
	"github.com/ssksameer56/CloudIndexer/handlers"
	"github.com/ssksameer56/CloudIndexer/models"
	"github.com/ssksameer56/CloudIndexer/workers"
)

func RunServer() {

	err := config.LoadConfig()
	if err != nil {
		log.Panic().Str("component", "Server").Msg("cant load config")
	}
	_ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := gin.Default()
	router.Use(gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	wg := sync.WaitGroup{}

	accessToken, _ := config.GetAccessToken()
	dropbox := cloud.DropBox{
		AuthKey: accessToken,
		Timeout: time.Minute,
	}
	err = dropbox.Connect(_ctx)
	if err != nil {
		log.Panic().Str("component", "Server").Msg("cant load dropbox client")
	}

	err = dropbox.Ping(_ctx)
	if err != nil {
		log.Panic().Str("component", "Server").Msg("cant connect to dropbox")
	}

	es := elasticservice.ElasticSearchService{}
	err = es.Connect()
	if err != nil {
		log.Panic().Str("component", "Server").Msg("cant load dropbox client")
	}

	IndexerNotificationChannel := make(chan models.CloudWatcherNotification, config.Config.BufferSize)

	cw := workers.CloudWatcher{
		CloudProvider:              &dropbox,
		IndexerNotificationChannel: IndexerNotificationChannel,
	}

	wg.Add(1)
	err = cw.Init(_ctx)
	if err != nil {
		log.Panic().Str("component", "Server").Msg("cant start cloud watcher")
	}
	go cw.Run(&wg)

	esw := workers.ESWorker{
		Service:                    es,
		IndexerNotificationChannel: IndexerNotificationChannel,
	}
	wg.Add(1)
	err = esw.Init(_ctx)
	if err != nil {
		log.Panic().Str("component", "Server").Msg("cant start cloud watcher")
	}
	esw.Run(&wg)

	wg.Add(1)
	go config.AccessTokenLoop(_ctx, &wg)

	sHandler := handlers.SearchHandler{
		CloudProvider:   &dropbox,
		ESSearchService: es,
	}

	sc := controllers.SearchController{
		Handler: sHandler,
	}
	router.GET("/search", sc.Search)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panic().Msgf("listen: %s\n", err)
		}
	}()
	<-_ctx.Done()
	log.Info().Msg("shutting down gracefully. waiting for goroutines now")

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	//wg.Wait()
	srv.Shutdown(_ctx)
	log.Info().Msg("stopped server and all routines")
}
