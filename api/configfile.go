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
	"github.com/cflion/cflion/model"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

func CreateConfigFile(c *gin.Context) {
	var params struct {
		Content string `json:"content" binding:"required"`
		AppId   int64  `json:"app_id" binding:"required"`
		Name    string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindWith(&params, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
		return
	}
	if model.ExistsConfigFile(params.AppId, params.Name) {
		c.JSON(http.StatusUnprocessableEntity, Response{Msg: fmt.Sprintf("Config file [%s] already exists", params.Name)})
		return
	}
	configFile := &model.ConfigFile{
		Name:  params.Name,
		AppId: params.AppId,
	}
	err := configFile.CreateWithContent(params.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, Response{Msg: fmt.Sprintf("Config file [%s] creates successfully", params.Name)})
}
