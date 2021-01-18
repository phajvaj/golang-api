package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"nan_api_main/controller"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

func getName(msg string) {
	fmt.Println("Debug: " + msg)
}

func main() {
	viper.SetConfigName("config") // ชื่อ config file
	viper.AddConfigPath(".")      // ระบุ path ของ config file
	viper.AutomaticEnv()          // อ่าน value จาก ENV variable
	// แปลง _ underscore ใน env เป็น . dot notation ใน viper
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// อ่าน config
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s \n", err))
	}

	srv := &http.Server{
		Addr:              viper.GetString("app.host"),
		Handler:           controller.SetRoutes().Path(),
		TLSConfig:         &tls.Config{},
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      map[string]func(*http.Server, *tls.Conn, http.Handler){},
		ConnState: func(net.Conn, http.ConnState) {
		},
		ErrorLog: &log.Logger{},
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	//srv.ListenAndServe()
	//controller.SetRoutes().Run(viper.GetString("app.host"))
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shuting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
