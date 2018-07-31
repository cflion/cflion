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

import "github.com/cflion/cflion/pkg/console/api"

type Repository interface {
	QueryAppsBrief() ([]*api.App, error)
	GetAppById(id int64) (*api.App, error)
	GetAppByName(name string) (*api.App, error)
	ExistsAppByNameAndEnv(name, env string) bool
	InsertApp(app *api.App) (int64, error)
}

type ServiceImpl struct {
	Repo Repository
}

func (service *ServiceImpl) ListApps() ([]map[string]interface{}, error) {
	apps, err := service.Repo.QueryAppsBrief()
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, 8)
	for _, app := range apps {
		result = append(result, map[string]interface{}{
			"id":       app.Id,
			"name":     app.Name,
			"env":      app.Env,
		})
	}
	return result, nil
}

func (service *ServiceImpl) GetAppById(id int64) (*api.App, error) {
	return service.Repo.GetAppById(id)
}

func (service *ServiceImpl) GetAppByName(name string) (*api.App, error) {
	return service.Repo.GetAppByName(name)
}

func (service *ServiceImpl) ExistsAppByNameAndEnv(name, env string) bool {
	return service.Repo.ExistsAppByNameAndEnv(name, env)
}

func (service *ServiceImpl) CreateApp(name, env string) (int64, error) {
	app := &api.App{Name: name, Env: env}
	return service.Repo.InsertApp(app)
}
