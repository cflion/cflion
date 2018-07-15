//  Copyright (c) 2018 The cflion Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

// Package restful defines the http server.
package restful

import (
	"context"
	"github.com/cflion/cflion/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"os/signal"
    "github.com/cflion/cflion/internal/app"
)

type Server struct {
	srv *http.Server
}

// NewServer creates a server.
func NewServer(listenAddr string, service app.Service) *Server {
	filePath := viper.GetString("logging.file")
	if len(filePath) > 0 {
		gin.DisableConsoleColor()
		f, _ := os.OpenFile(filePath, os.O_RDWR, 0666)
		gin.DefaultWriter = io.MultiWriter(f)
		gin.DefaultErrorWriter = io.MultiWriter(f)
	}
	if log.IsDebugEnabled() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	initRouter(router, service)

	srv := &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}
	readTimeout := viper.GetDuration("server.readTimeout")
	if readTimeout > 0 {
		srv.ReadTimeout = readTimeout
	}
	writeTimeout := viper.GetDuration("server.writeTimeout")
	if writeTimeout > 0 {
		srv.WriteTimeout = writeTimeout
	}
	idleTimeout := viper.GetDuration("server.idleTimeout")
	if idleTimeout > 0 {
		srv.IdleTimeout = idleTimeout
	}
	return &Server{srv: srv}
}

// Start server.
func (server *Server) Start() {
	go func() {
		if err := server.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("Fatal server listen: %s", err)
			os.Exit(1)
		}
	}()
}

// Stop server.
func (server *Server) Stop() <-chan struct{} {
	ch := make(chan struct{})
	go func(ch chan<- struct{}) {
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt, os.Interrupt)
		<-quit
		log.Info("Shutdown server ...")
		ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("server.quitTimeout"))
		defer cancel()
		if err := server.srv.Shutdown(ctx); err != nil {
			log.Errorf("Shutdown server: %s", err)
			os.Exit(1)
		}
		ch <- struct{}{}
	}(ch)
	return ch
}

func initRouter(router *gin.Engine, service app.Service)  {
    v1 := router.Group("/v1")
    {
        v1.GET("/config-groups", ListConfigGroup(service))
        v1.POST("/config-groups", CreateConfigGroup(service))
        v1.GET("/config-groups/:group_id", ViewConfigGroup(service))
        v1.PUT("/config-groups/:group_id", UpdateConfigGroup(service))

        v1.GET("/config-files", ListConfigFile(service))
        v1.POST("/config-files", CreateConfigFile(service))
        v1.GET("/config-files/:file_id", ViewConfigFile(service))
        v1.PUT("/config-files/:file_id", UpdateConfigFile(service))
    }
}
