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

import "fmt"

type Service interface {
	ListApps() ([]map[string]interface{}, error)
	GetAppById(id int64) (*App, error)
	GetAppByName(name string) (*App, error)
	ExistsAppByNameAndEnv(name, env string) bool
	CreateApp(name, env string) (int64, error)
}

type App struct {
	Id       int64
	Name     string
	Env      string
	Outdated byte
}

func (app *App) String() string {
	return fmt.Sprintf("App {Id=%d | Name=%s | Env=%s | Outdated=%d}", app.Id, app.Name, app.Env, app.Outdated)
}
