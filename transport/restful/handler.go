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
	"fmt"
	"github.com/cflion/cflion/internal/app"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
)

type ResponseRet struct {
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func ListConfigGroup(service app.Service) func(*gin.Context) {
	return func(ctx *gin.Context) {
		data, err := service.ListConfigGroup()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ResponseRet{Msg: err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, ResponseRet{Data: data})
	}
}

func CreateConfigGroup(service app.Service) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var params struct {
			App         string `json:"app" binding:"required"`
			Environment string `json:"environment" binding:"required"`
		}
		if err := ctx.ShouldBindWith(&params, binding.JSON); err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseRet{Msg: err.Error()})
			return
		}
		if service.ExistsConfigGroupByAppAndEnvironment(params.App, params.Environment) {
			ctx.JSON(http.StatusUnprocessableEntity, ResponseRet{Msg: fmt.Sprintf("App [name=%s] [environment=%s] already exists", params.App, params.Environment)})
			return
		}
		_, err := service.CreateConfigGroup(params.App, params.Environment)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ResponseRet{Msg: err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, ResponseRet{Msg: fmt.Sprintf("App [name=%s] creates successfully", params.App)})
	}
}

func ViewConfigGroup(service app.Service) func(*gin.Context) {
	return func(ctx *gin.Context) {
		groupId, err := strconv.ParseInt(ctx.Param("group_id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseRet{Msg: err.Error()})
			return
		}
		if !service.ExistsConfigGroupById(groupId) {
			ctx.JSON(http.StatusUnprocessableEntity, ResponseRet{Msg: fmt.Sprintf("App [id=%d] doesn't exists", groupId)})
			return
		}
		data, err := service.ViewConfigGroup(groupId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ResponseRet{Msg: err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, ResponseRet{Data: data})
	}
}

func UpdateConfigGroup(service app.Service) func(*gin.Context) {
	return func(ctx *gin.Context) {
		groupId, err := strconv.ParseInt(ctx.Param("group_id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseRet{Msg: err.Error()})
			return
		}
		if !service.ExistsConfigGroupById(groupId) {
			ctx.JSON(http.StatusUnprocessableEntity, ResponseRet{Msg: fmt.Sprintf("App [id=%d] doesn't exists", groupId)})
			return
		}
		var params struct {
			ConfigFiles []int64 `json:"config_files" binding:"required"`
		}
		if err = ctx.ShouldBindWith(&params, binding.JSON); err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseRet{Msg: err.Error()})
			return
		}
		err = service.UpdateConfigGroupAssociation(groupId, params.ConfigFiles)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ResponseRet{Msg: err.Error()})
			return
		}
		ctx.Status(http.StatusOK)
	}
}

func ListConfigFile(service app.Service) func(*gin.Context) {
	return func(ctx *gin.Context) {
		data, err := service.ListConfigFile()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ResponseRet{Msg: err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, ResponseRet{Data: data})
	}
}

func CreateConfigFile(service app.Service) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var params struct {
			NamespaceId int64  `json:"namespace_id" binding:"required"`
			Config      string `json:"config" binding:"required"`
			Filename    string `json:"filename" binding:"required"`
		}
		if err := ctx.ShouldBindWith(&params, binding.JSON); err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseRet{Msg: err.Error()})
			return
		}
		if service.ExistsConfigFileByNameAndNamespaceId(params.Filename, params.NamespaceId) {
			ctx.JSON(http.StatusUnprocessableEntity, ResponseRet{Msg: fmt.Sprintf("Config file [name=%s] [namespace_id=%d] already exists", params.Filename, params.NamespaceId)})
			return
		}
		_, err := service.CreateConfigFile(params.Filename, params.NamespaceId, params.Config)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ResponseRet{Msg: err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, ResponseRet{Msg: fmt.Sprintf("Config file [name=%s] creates successfully", params.Filename)})
	}
}

func ViewConfigFile(service app.Service) func(*gin.Context) {
	return func(ctx *gin.Context) {
		fileId, err := strconv.ParseInt(ctx.Param("file_id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseRet{Msg: err.Error()})
			return
		}
		if !service.ExistsConfigFileById(fileId) {
			ctx.JSON(http.StatusUnprocessableEntity, ResponseRet{Msg: fmt.Sprintf("Config file [id=%d] doesn't exists", fileId)})
			return
		}
		data, err := service.ViewConfigFile(fileId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ResponseRet{Msg: err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, ResponseRet{Data: data})
	}
}

func UpdateConfigFile(service app.Service) func(*gin.Context) {
	return func(ctx *gin.Context) {
		fileId, err := strconv.ParseInt(ctx.Param("file_id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseRet{Msg: err.Error()})
			return
		}
		var params struct {
			Config string `json:"config" binding:"required"`
		}
		if err := ctx.ShouldBindWith(&params, binding.JSON); err != nil {
			ctx.JSON(http.StatusBadRequest, ResponseRet{Msg: err.Error()})
			return
		}
		if !service.ExistsConfigFileById(fileId) {
			ctx.JSON(http.StatusUnprocessableEntity, ResponseRet{Msg: fmt.Sprintf("Config file [id=%d] doesn't exists", fileId)})
			return
		}
		err = service.UpdateConfigFile(fileId, params.Config)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ResponseRet{Msg: err.Error()})
			return
		}
		ctx.Status(http.StatusOK)
	}
}
