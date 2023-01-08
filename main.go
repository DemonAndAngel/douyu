package main

import (
	"context"
	"douyu/middlewares/common"
	"douyu/models"
	"douyu/router"
	"douyu/timer"
	"douyu/utils/config"
	_ "douyu/utils/config"
	_ "douyu/utils/helpers"
	"douyu/utils/snowflake"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone
}

var configFile = flag.String("c", "./config/dev.json", "the config file")

func main() {
	// 初始化配置文件
	flag.Parse()
	config.InitConfig(configFile)
	models.Migrate()
	timer.Timer()
	fmt.Println(snowflake.GenId())

	port := viper.GetInt("port")
	if port <= 0 {
		port = 8764
	}

	engine := gin.Default()
	if viper.GetString("appEnv") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine.Use(common.StrictHeader())

	router.InitRouter(engine)

	srv := &http.Server{
		Addr:         fmt.Sprintf("127.0.0.1:%d", port),
		Handler:      engine,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println(fmt.Sprintf("listen: %d", port))

	// 监听停止服务信号
	quit, stop := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt, os.Kill)
	defer stop()
	<-quit.Done()
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

}
