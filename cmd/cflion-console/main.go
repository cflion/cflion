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
	"github.com/cflion/cflion/cmd/cflion-console/server"
	"github.com/cflion/cflion/cmd/cflion-console/server/repository/mysql"
	"github.com/cflion/cflion/pkg/console/api"
	"github.com/cflion/cflion/pkg/database"
	"github.com/cflion/cflion/pkg/log"
	"github.com/cflion/cflion/pkg/transport/restful"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
	"time"
)

// init setting
func init() {
	confPath := flag.String("conf", "conf/app.yml", "path of app.yml")
	flag.Parse()
	// set default config
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 9090)
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
	dbCfg := &database.DBConfig{
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		Host:     viper.GetString("db.host"),
		Port:     viper.GetInt("db.port"),
		Database: viper.GetString("db.database"),
		MaxIdle:  viper.GetInt("db.maxIdle"),
		MaxOpen:  viper.GetInt("db.maxOpen"),
	}
	db, err := database.ConnectDatabase(dbCfg)
	if err != nil {
		log.Errorf("Fatal error when connect to db: %s", err)
		os.Exit(1)
	}
	var repo server.Repository = &mysql.RepositoryImpl{DB: db}
	var service api.Service = &server.ServiceImpl{Repo: repo}

	srvCfg := &restful.ServerConfig{
		ListenAddr:      fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.port")),
		ReadTimeout:     time.Duration(viper.GetInt("server.readTimeout")) * time.Second,
		WriteTimeout:    time.Duration(viper.GetInt("server.writeTimeout")) * time.Second,
		IdleTimeout:     time.Duration(viper.GetInt("server.idleTimeout")) * time.Second,
		QuitTimeout:     time.Duration(viper.GetInt("server.quitTimeout")) * time.Second,
		LoggingFilePath: viper.GetString("logging.file"),
	}
	srv := restful.NewServer(srvCfg, func(router *gin.Engine) {
		v1 := router.Group("/v1")
		{
			v1.GET("/apps", server.ListApps(service))
			v1.POST("/apps", server.CreateApp(service))
			v1.PUT("/apps", server.PublishApp(service))
			v1.GET("/apps/:app_id", server.ViewApp(service))
			v1.PUT("/apps/:app_id", server.UpdateApp(service))

			v1.GET("/config-files", server.ListConfigFiles(service))
			v1.POST("/config-files", server.CreateConfigFile(service))
			v1.GET("/config-files/:file_id", server.ViewConfigFile(service))
			v1.PUT("/config-files/:file_id", server.UpdateConfigFile(service))
		}
	})
	srv.Start()
	<-srv.Stop()
	log.Info("server exited")
}
