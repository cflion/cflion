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

package server

import (
	"fmt"
	"github.com/cflion/cflion/pkg/manager/api"
	"github.com/cflion/cflion/pkg/transport/restful"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
)

func CreateApp(service api.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var params struct {
			Name string `json:"name" binding:"required"`
		}
		if err := ctx.ShouldBindJSON(&params); err != nil {
			ctx.JSON(http.StatusBadRequest, restful.ResponseRet{Msg: err.Error()})
			return
		}
		if service.ExistsAppByName(params.Name) {
			ctx.JSON(http.StatusUnprocessableEntity, restful.ResponseRet{Msg: fmt.Sprintf("App [name=%s] already exists", params.Name)})
			return
		}
		_, err := service.CreateApp(params.Name)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, restful.ResponseRet{Msg: fmt.Sprintf("App [name=%s] creates successfully", params.Name)})
	}
}

func PublishApp(service api.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var params struct {
			Name string `json:"name" binding:"required"`
		}
		if err := ctx.ShouldBindWith(&params, binding.JSON); err != nil {
			ctx.JSON(http.StatusBadRequest, restful.ResponseRet{Msg: err.Error()})
			return
		}
		app, err := service.GetAppByName(params.Name)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, restful.ResponseRet{Msg: err.Error()})
			return
		}
		err = service.PublishApp(app.Id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		ctx.Status(http.StatusOK)
	}
}

func ViewApp(service api.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		name := ctx.Param("name")
		app, err := service.GetAppByName(name)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, restful.ResponseRet{Msg: err.Error()})
			return
		}
		data, err := service.ViewApp(app.Id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, restful.ResponseRet{Data: data})
	}
}

func UpdateApp(service api.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		name := ctx.Param("name")
		app, err := service.GetAppByName(name)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, restful.ResponseRet{Msg: err.Error()})
			return
		}
		var params struct {
			ConfigFiles []int64 `json:"config_files" binding:"required"`
		}
		if err = ctx.ShouldBindJSON(&params); err != nil {
			ctx.JSON(http.StatusBadRequest, restful.ResponseRet{Msg: err.Error()})
			return
		}
		err = service.UpdateAppAssociation(app.Id, params.ConfigFiles)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		ctx.Status(http.StatusOK)
	}
}

func ListConfigFiles(service api.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		data, err := service.ListConfigFiles()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, restful.ResponseRet{Data: data})
	}
}

func CreateConfigFile(service api.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var params struct {
			NamespaceId int64  `json:"namespace_id" binding:"required"`
			Config      string `json:"config" binding:"required"`
			Filename    string `json:"filename" binding:"required"`
		}
		if err := ctx.ShouldBindWith(&params, binding.JSON); err != nil {
			ctx.JSON(http.StatusBadRequest, restful.ResponseRet{Msg: err.Error()})
			return
		}
		if service.ExistsConfigFileByNameAndNamespaceId(params.Filename, params.NamespaceId) {
			ctx.JSON(http.StatusUnprocessableEntity, restful.ResponseRet{Msg: fmt.Sprintf("Config file [name=%s] [namespace_id=%d] already exists", params.Filename, params.NamespaceId)})
			return
		}
		_, err := service.CreateConfigFile(params.Filename, params.NamespaceId, params.Config)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, restful.ResponseRet{Msg: fmt.Sprintf("Config file [name=%s] creates successfully", params.Filename)})
	}
}

func ViewConfigFile(service api.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fileId, err := strconv.ParseInt(ctx.Param("file_id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, restful.ResponseRet{Msg: err.Error()})
			return
		}
		if !service.ExistsConfigFileById(fileId) {
			ctx.JSON(http.StatusUnprocessableEntity, restful.ResponseRet{Msg: fmt.Sprintf("Config file [id=%d] doesn't exists", fileId)})
			return
		}
		data, err := service.ViewConfigFile(fileId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, restful.ResponseRet{Data: data})
	}
}

func UpdateConfigFile(service api.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fileId, err := strconv.ParseInt(ctx.Param("file_id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, restful.ResponseRet{Msg: err.Error()})
			return
		}
		var params struct {
			Config string `json:"config" binding:"required"`
		}
		if err := ctx.ShouldBindWith(&params, binding.JSON); err != nil {
			ctx.JSON(http.StatusBadRequest, restful.ResponseRet{Msg: err.Error()})
			return
		}
		if !service.ExistsConfigFileById(fileId) {
			ctx.JSON(http.StatusUnprocessableEntity, restful.ResponseRet{Msg: fmt.Sprintf("Config file [id=%d] doesn't exists", fileId)})
			return
		}
		err = service.UpdateConfigFile(fileId, params.Config)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		ctx.Status(http.StatusOK)
	}
}

func QueryWatcher(service api.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var params struct {
			App string `form:"app" binding:"required"`
		}
		if err := ctx.ShouldBindQuery(&params); err != nil {
			ctx.JSON(http.StatusBadRequest, restful.ResponseRet{Msg: err.Error()})
			return
		}
		app := &api.App{Name: params.App}
		endpoints := viper.GetStringSlice("etcd.endpoints")
		ctx.JSON(http.StatusOK, restful.ResponseRet{Data: gin.H{"key": app.Key(), "endpoints": endpoints}})

	}
}
