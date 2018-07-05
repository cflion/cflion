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
	"strconv"
)

func CreateApp(c *gin.Context) {
	var params struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindWith(&params, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
		return
	}
	if model.ExistsApp(params.Name) {
		c.JSON(http.StatusUnprocessableEntity, Response{Msg: fmt.Sprintf("App [%s] already exists", params.Name)})
		return
	}
	app := &model.App{
		Name:     params.Name,
		Outdated: 1,
	}
	_, err := app.Create()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Msg: err.Error()})
	} else {
		c.JSON(http.StatusCreated, Response{Msg: fmt.Sprintf("App [%s] creates successfully", params.Name)})
	}
}

func ViewApp(c *gin.Context) {
	appId, err := strconv.ParseInt(c.Param("appId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
		return
	}
	_, err = model.RetrieveApp(appId)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, Response{Msg: err.Error()})
		return
	}
	appBrief, err := model.ViewApp(appId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, Response{Data: appBrief})
}

func ListApp(c *gin.Context) {
	appsBrief, err := model.ListApp()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, Response{Data: gin.H{"apps": appsBrief}})
}

// UpdateApp updates app relationship.
func UpdateApp(c *gin.Context) {
	appId, err := strconv.ParseInt(c.Param("appId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
		return
	}
	app, err := model.RetrieveApp(appId)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, Response{Msg: err.Error()})
		return
	}
	var params struct {
		FileIds []int64 `json:"fileIds" binding:"required"`
	}
	if err := c.ShouldBindWith(&params, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, Response{Msg: err.Error()})
		return
	}
	err = app.UpdateAssociation(params.FileIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Msg: err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
