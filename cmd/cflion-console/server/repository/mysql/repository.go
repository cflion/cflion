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

package mysql

import (
	"database/sql"
	"github.com/cflion/cflion/pkg/console/api"
	"github.com/cflion/cflion/pkg/log"
)

type RepositoryImpl struct {
	DB *sql.DB
}

func (repo *RepositoryImpl) QueryAppsBrief() ([]*api.App, error) {
	rows, err := repo.DB.Query("select id, name, env from app")
	if err != nil {
		log.Errorf("Query all apps error: %s", err)
		return nil, err
	}
	apps := make([]*api.App, 0, 8)
	for rows.Next() {
		var app api.App
		rows.Scan(&app.Id, &app.Name, &app.Env)
		apps = append(apps, &app)
	}
	return apps, nil
}

func (repo *RepositoryImpl) GetAppById(id int64) (*api.App, error) {
	var app api.App
	err := repo.DB.QueryRow("select id, name, env from app where id = ?", id).Scan(&app.Id, &app.Name, &app.Env)
	if err != nil {
		log.Errorf("Get app info [id=%d] error: %s", id, err)
		return nil, err
	}
	return &app, nil
}

func (repo *RepositoryImpl) GetAppByName(name string) (*api.App, error) {
	var app api.App
	err := repo.DB.QueryRow("select id, name, env from app where name = ?", name).Scan(&app.Id, &app.Name, &app.Env)
	if err != nil {
		log.Errorf("Get app info [name=%s] error: %s", name, err)
		return nil, err
	}
	return &app, nil
}

func (repo *RepositoryImpl) ExistsAppByNameAndEnv(name, env string) bool {
	var count int64
	err := repo.DB.QueryRow("select count(1) from app where name = ? and env = ?", name, env).Scan(&count)
	if err != nil {
		log.Errorf("Count app [name=%s] [env=%s] error: %s", name, env, err)
		return false
	}
	return count == 1
}

func (repo *RepositoryImpl) InsertApp(app *api.App) (int64, error) {
	res, err := repo.DB.Exec("insert into app (name, env, ctime, utime) values (?, ?, now(), now())", app.Name, app.Env)
	if err != nil {
		log.Errorf("Create app [name=%s] [env=%s] error: %s", app.Name, app.Env, err)
		return -1, err
	}
	return res.LastInsertId()
}
