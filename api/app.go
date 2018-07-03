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

func CreateApp(c *gin.Context) {
	var createApp struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindWith(&createApp, binding.JSON); err == nil {
		if model.Exists(createApp.Name) {
			c.JSON(http.StatusUnprocessableEntity, Response{Msg: fmt.Sprintf("App [%s] already exists", createApp.Name)})
			return
		} else {
			app := &model.App{
				Name:     createApp.Name,
				Outdated: 1,
			}
			_, err := app.Create()
			if err != nil {
				c.JSON(http.StatusInternalServerError, Response{Msg: err.Error()})
			} else {
				c.JSON(http.StatusCreated, Response{Msg: fmt.Sprintf("App [%s] creates successfully", createApp.Name)})
			}
		}
	} else {
		c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
	}

}
