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
	"flag"
	"fmt"
	"github.com/cflion/cflion/pkg/database"
	"github.com/cflion/cflion/pkg/log"
	"github.com/cflion/cflion/server"
	"github.com/cflion/cflion/server/repository"
	"github.com/cflion/cflion/transport/restful"
	"github.com/spf13/viper"
	"os"
)

// init setting
func init() {
	confPath := flag.String("conf", "conf/app.yml", "path of app.yml")
	flag.Parse()
	// set default config
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.readTimeout", 3)
	viper.SetDefault("server.writeTimeout", 3)
	viper.SetDefault("server.quitTimeout", 5)
	viper.SetDefault("logging.level", "INFO")
	viper.SetDefault("db.maxIdle", 20)
	viper.SetDefault("db.maxOpen", 100)
	viper.SetDefault("etcd.requestTimeout", 3)
	viper.SetConfigFile(*confPath)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Errorf("Fatal error config file: %s", err)
		os.Exit(1)
	}
	// init global logger
	log.SetLevel(viper.GetString("logging.level"))
	filePath := viper.GetString("logging.file")
	if len(filePath) > 0 {
		f, _ := os.Create(filePath)
		log.SetOutput(f)
	}
}

func main() {
	dsn := database.FormatDSN(viper.GetString("db.username"), viper.GetString("db.password"), viper.GetString("db.host"), viper.GetInt("db.port"), viper.GetString("database"))
	db, err := database.ConnectDatabase(dsn)
	if err != nil {
		log.Errorf("Fatal error when connect to db: %s", err)
		os.Exit(1)
	}
	repo := &repository.Repository{DB: db}
	service := &server.Service{Repo: repo}
	listenAddr := fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.port"))
	srv := restful.NewServer(listenAddr, service)
	srv.Start()
	<-srv.Stop()
	log.Info("server exited")
}
