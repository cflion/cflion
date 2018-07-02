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

package main

import (
	"context"
	"flag"
	"github.com/cflion/cflion/api"
	"github.com/cflion/cflion/log"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func init() {
	// parse args
	confPath := flag.String("conf", "conf/app.yml", "path of app.yml")
	flag.Parse()
	// init config
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.readTimeout", 3*time.Second)
	viper.SetDefault("server.writeTimeout", 3*time.Second)
	viper.SetDefault("server.quitTimeout", 5*time.Second)
	viper.SetDefault("logging.level", "INFO")
	viper.SetConfigFile(*confPath)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Errorf("Fatal error config file: %s\n", err)
		os.Exit(1)
	}
	// init logger
	log.SetLevel(viper.GetString("logging.level"))
	filePath := viper.GetString("logging.file")
	if "" != filePath {
		f, _ := os.Create(filePath)
		log.SetOutput(f)
	}
}

func main() {
	srv := api.SetupServer()
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("Fatal server listen: %s\n", err)
			os.Exit(1)
		}
	}()
	// shutdown gracefully
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Interrupt)
	<-quit
	log.Info("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("server.quitTimeout"))
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("Shutdown server: %s", err)
		os.Exit(1)
	}
	log.Info("Server exited")
}
