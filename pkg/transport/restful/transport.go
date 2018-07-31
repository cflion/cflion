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

package restful

import (
	"context"
	"github.com/cflion/cflion/pkg/log"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type ResponseRet struct {
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type ServerConfig struct {
	ListenAddr      string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	QuitTimeout     time.Duration
	LoggingFilePath string
}

type Server struct {
	srv *http.Server
	cfg *ServerConfig
}

// NewServer creates a server.
func NewServer(cfg *ServerConfig, register func(router *gin.Engine)) *Server {
	filePath := cfg.LoggingFilePath
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
	register(router)
	srv := &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: router,
	}
	if cfg.ReadTimeout > 0 {
		srv.ReadTimeout = cfg.ReadTimeout
	}
	if cfg.WriteTimeout > 0 {
		srv.WriteTimeout = cfg.WriteTimeout
	}
	if cfg.IdleTimeout > 0 {
		srv.IdleTimeout = cfg.IdleTimeout
	}
	return &Server{srv: srv, cfg: cfg}
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
		ctx, cancel := context.WithTimeout(context.Background(), server.cfg.QuitTimeout)
		defer cancel()
		if err := server.srv.Shutdown(ctx); err != nil {
			log.Errorf("Shutdown server: %s", err)
			os.Exit(1)
		}
		ch <- struct{}{}
	}(ch)
	return ch
}
