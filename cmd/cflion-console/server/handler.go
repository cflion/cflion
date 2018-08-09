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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cflion/cflion/pkg/console/api"
	"github.com/cflion/cflion/pkg/transport/restful"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"strconv"
)

func ListApps(service api.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		data, err := service.ListApps()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, restful.ResponseRet{Data: data})
	}
}

func CreateApp(service api.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var params struct {
			Name string `json:"name" binding:"required"`
			Env  string `json:"env" binding:"required"`
		}
		if err := ctx.ShouldBindWith(&params, binding.JSON); err != nil {
			ctx.JSON(http.StatusBadRequest, restful.ResponseRet{Msg: err.Error()})
			return
		}
		if service.ExistsAppByNameAndEnv(params.Name, params.Env) {
			ctx.JSON(http.StatusBadRequest, restful.ResponseRet{Msg: fmt.Sprintf("App [name=%s] [env=%s] already exists", params.Name, params.Env)})
			return
		}
		managerUrl := getManagerEndpoint(params.Env)
		if len(managerUrl) <= 0 {
			ctx.JSON(http.StatusBadRequest, restful.ResponseRet{Msg: fmt.Sprintf("Can not support [env=%s]", params.Env)})
			return
		}
		// call remote manager
		reqBytes, _ := json.Marshal(map[string]interface{}{"name": params.Name})
		resp, err := http.Post(managerUrl+"/v1/apps", "application/json", bytes.NewBuffer(reqBytes))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			ctx.Status(resp.StatusCode)
			return
		}
		// create local
		_, err = service.CreateApp(params.Name, params.Env)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		ctx.Status(http.StatusCreated)
	}
}

func PublishApp(service api.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var params struct {
			AppId int64 `json:"app_id" binding:"required"`
		}
		if err := ctx.ShouldBindWith(&params, binding.JSON); err != nil {
			ctx.JSON(http.StatusBadRequest, restful.ResponseRet{Msg: err.Error()})
			return
		}
		app, err := service.GetAppById(params.AppId)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, restful.ResponseRet{Msg: err.Error()})
			return
		}
		// call remote manager
		reqBytes, _ := json.Marshal(map[string]string{"name": app.Name})
		client := http.Client{}
		req, err := http.NewRequest("PUT", getManagerEndpoint(app.Env)+"/v1/apps", bytes.NewBuffer(reqBytes))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		resp.Body.Close()
		ctx.Status(resp.StatusCode)
	}
}

func ViewApp(service api.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		appId, err := strconv.ParseInt(ctx.Param("app_id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, restful.ResponseRet{Msg: err.Error()})
			return
		}
		app, err := service.GetAppById(appId)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, restful.ResponseRet{Msg: err.Error()})
			return
		}
		// call remote manager
		resp, err := http.Get(getManagerEndpoint(app.Env) + "/v1/apps/" + app.Name)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			ctx.Status(resp.StatusCode)
			return
		}
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		var result interface{}
		err = json.Unmarshal(respBytes, &result)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}

func UpdateApp(service api.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		appId, err := strconv.ParseInt(ctx.Param("app_id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, restful.ResponseRet{Msg: err.Error()})
			return
		}
		var params struct {
			ConfigFiles []int64 `json:"config_files" binding:"required"`
		}
		if err = ctx.ShouldBindJSON(&params); err != nil {
			ctx.JSON(http.StatusBadRequest, restful.ResponseRet{Msg: err.Error()})
			return
		}
		app, err := service.GetAppById(appId)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, restful.ResponseRet{Msg: err.Error()})
			return
		}
		// call remote manager
		reqBytes, _ := json.Marshal(params)
		client := http.Client{}
		req, err := http.NewRequest("PUT", getManagerEndpoint(app.Env)+"/v1/apps/"+app.Name, bytes.NewBuffer(reqBytes))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		resp.Body.Close()
		ctx.Status(resp.StatusCode)
	}
}

func ListConfigFiles(service api.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// call remote manager
		env := ctx.Query("env")
		managerUrl := getManagerEndpoint(env)
		if len(managerUrl) <= 0 {
			ctx.JSON(http.StatusBadRequest, restful.ResponseRet{Msg: fmt.Sprintf("Can not support [env=%s]", env)})
			return
		}
		resp, err := http.Get(managerUrl + "/v1/config-files")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		defer resp.Body.Close()
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		var result interface{}
		json.Unmarshal(respBytes, &result)
		ctx.JSON(http.StatusOK, result)
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
		app, err := service.GetAppById(params.NamespaceId)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, restful.ResponseRet{Msg: err.Error()})
			return
		}
		// call remote manager
		reqBytes, err := json.Marshal(params)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		resp, err := http.Post(getManagerEndpoint(app.Env)+"/v1/config-files", "application/json", bytes.NewBuffer(reqBytes))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		resp.Body.Close()
		ctx.Status(resp.StatusCode)
	}
}

func ViewConfigFile(service api.Service) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		fileId, err := strconv.ParseInt(ctx.Param("file_id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, restful.ResponseRet{Msg: err.Error()})
			return
		}
		namespaceId, err := strconv.ParseInt(ctx.Query("namespace_id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, restful.ResponseRet{Msg: err.Error()})
			return
		}
		app, err := service.GetAppById(namespaceId)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, restful.ResponseRet{Msg: err.Error()})
			return
		}
		managerUrl := getManagerEndpoint(app.Env)
		// call remote manager
		resp, err := http.Get(fmt.Sprintf("%s/v1/config-files/%d", managerUrl, fileId))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		defer resp.Body.Close()
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		var result interface{}
		json.Unmarshal(respBytes, &result)
		ctx.JSON(resp.StatusCode, result)
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
			NamespaceId int64  `json:"namespace_id" binding:"required"`
			Config      string `json:"config" binding:"required"`
		}
		if err := ctx.ShouldBindWith(&params, binding.JSON); err != nil {
			ctx.JSON(http.StatusBadRequest, restful.ResponseRet{Msg: err.Error()})
			return
		}
		app, err := service.GetAppById(params.NamespaceId)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, restful.ResponseRet{Msg: err.Error()})
			return
		}
		managerUrl := getManagerEndpoint(app.Env)
		// call remote manager
		reqBytes, _ := json.Marshal(params)
		client := http.Client{}
		req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/config-files/%d", managerUrl, fileId), bytes.NewBuffer(reqBytes))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, restful.ResponseRet{Msg: err.Error()})
			return
		}
		resp.Body.Close()
		ctx.Status(resp.StatusCode)
	}
}

func getManagerEndpoint(env string) string {
	return viper.GetString(env + ".manager.endpoint")
}
