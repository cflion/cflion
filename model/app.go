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

package model

import (
	"fmt"
	"github.com/cflion/cflion/log"
)

type App struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Outdated byte   `json:"outdated"`
}

func (app *App) String() string {
	return fmt.Sprintf("App [Id=%d] [Name=%s] [Outdated=%d]", app.Id, app.Name, app.Outdated)
}

// ExistsApp checks whether the name of the app exists.
func Exists(name string) bool {
	var count int64
	err := db.QueryRow("select count(1) from app where name = ?", name).Scan(&count)
	if err != nil {
		log.Errorf("Query app [%s] error: %s", name, err)
		return false
	}
	if count >= 1 {
		return true
	} else {
		return false
	}
}

func (app *App) Create() (int64, error) {
	stmt, err := db.Prepare("insert into app (name, outdated, ctime, utime) values (?, ?, now(), now())")
	if err != nil {
		log.Error("Prepare error when creates app: ", err)
		return -1, err
	}
	res, err := stmt.Exec(app.Name, app.Outdated)
	if err != nil {
		log.Error("Exec error when creates app: ", err)
		return -1, err
	}
	return res.LastInsertId()
}
