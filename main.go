package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

var (
	tr         *http.Transport
	client     *http.Client
	logger     *logrus.Logger
	db         *gorm.DB
	conf       Config
	err        error
	signalChan chan os.Signal
)

func init() {
	logger = initLogger("message-pusher.log", logrus.DebugLevel)
	loadConfig()
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{
		Transport: tr,
		Timeout:   5 * time.Second,
	}

}

func main() {
	if db, err = initConnect(); err != nil {
		logger.Errorln(err.Error())
		return
	}
	initDatabase()

	router := newRouter()

	httpServer := &http.Server{
		Addr:           fmt.Sprintf(":%d", conf.Server.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("HTTP server listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	signalChan = make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)
	sig := <-signalChan
	logger.Println("Get Signal:", sig)

	logger.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	logger.Println("Server exiting")
}
