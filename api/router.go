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

package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
)

type Response struct {
	Msg  string `json:"msg,omitempty"`
	Data gin.H  `json:"data,omitempty"`
}

// SetupServer defines router of the server.
func SetupServer() *http.Server {
	// setup logging
	filePath := viper.GetString("logging.file")
	if "" != filePath {
		gin.DisableConsoleColor()
		f, _ := os.OpenFile(filePath, os.O_RDWR, 0666)
		gin.DefaultWriter = io.MultiWriter(f)
	}
	// setup router
	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.POST("/apps", CreateApp)
		v1.GET("/apps/:appId", ViewApp)
		v1.GET("/apps", ListApp)
		v1.POST("/config-files", CreateConfigFile)
	}
	// setup server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.port")),
		Handler:      router,
		ReadTimeout:  viper.GetDuration("server.readTimeout"),
		WriteTimeout: viper.GetDuration("server.writeTimeout"),
	}
	return server
}
