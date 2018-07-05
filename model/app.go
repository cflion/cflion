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
	"github.com/cflion/cflion/util"
	"strings"
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
func ExistsApp(name string) bool {
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

func RetrieveApp(id int64) (*App, error) {
	var app App
	err := db.QueryRow("select id, name, outdated from app where id = ?", id).Scan(&app.Id, &app.Name, &app.Outdated)
	if err != nil {
		log.Errorf("Query app_id [%s] error: %s", id, err)
		return nil, err
	}
	return &app, nil
}

func (app *App) Create() (int64, error) {
	stmt, err := db.Prepare("insert into app (name, outdated, ctime, utime) values (?, ?, now(), now())")
	if err != nil {
		log.Error("Prepare error when creates app: ", err)
		return -1, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(app.Name, app.Outdated)
	if err != nil {
		log.Error("Exec error when creates app: ", err)
		return -1, err
	}
	return res.LastInsertId()
}

func (app *App) Brief() (map[string]interface{}, error) {
	result := make(map[string]interface{}, 4)
	result["id"] = app.Id
	result["name"] = app.Name
	result["outdated"] = app.Outdated
	rows, err := db.Query("select cf.id, cf.name from config_file as cf where cf.id in (select a.file_id from association as a where a.app_id = ?)", app.Id)
	if err != nil {
		log.Errorf("Query config files by app [%s] error: %s", app.Name, err)
		return nil, err
	}
	configFiles := make([]map[string]interface{}, 10)
	for rows.Next() {
		var fileId int64
		var filename string
		err = rows.Scan(&fileId, &filename)
		if err != nil {
			log.Error("Scan config file error: ", err)
			return nil, err
		}
		configFile := map[string]interface{}{
			"id":       fileId,
			"name":     filename,
			"app_name": app.Name,
		}
		configFiles = append(configFiles, configFile)
	}
	result["config_files"] = configFiles
	return result, nil
}

func (app *App) UpdateAssociation(fileIds []int64) error {
	// query association file
	rows, err := db.Query("select file_id from association where app_id = ?", app.Id)
	if err != nil {
		log.Errorf("Query association with app_id [%s] error: %s", app.Id, err)
		return err
	}
	currentFileIds := make([]interface{}, 10)
	for rows.Next() {
		var fileId int64
		err = rows.Scan(&fileId)
		if err != nil {
			log.Error("Scan association file error: ", err)
			return err
		}
		currentFileIds = append(currentFileIds, fileId)
	}
	// convert to set
	updateFileIds := util.Int64SliceToInterfaceSlice(fileIds)
	addSet, delSet := util.DiffSet(util.NewSetBySlice(updateFileIds), util.NewSetBySlice(currentFileIds))
	// add association file
	addFileIds := addSet.ConvertToSlice()
	patterns := make([]string, len(addFileIds))
	params := make([]interface{}, len(addFileIds))
	for _, fileId := range addFileIds {
		patterns = append(patterns, "(?, ?, now(), now())")
		params = append(params, app.Id, fileId)
	}
	sql := fmt.Sprintf("insert into association (app_id, file_id, ctime, utime) values %s", strings.Join(patterns, ","))
	_, err = db.Exec(sql, params...)
	if err != nil {
		log.Errorf("Exec add association of app_id [%s] error: %s", app.Id, err)
		return err
	}
	// delete association file
	delFileIds := delSet.ConvertToSlice()
	patterns = make([]string, len(delFileIds))
	params = make([]interface{}, len(delFileIds)+1)
	params = append(params, app.Id)
	for _, fileId := range delFileIds {
		patterns = append(patterns, "?")
		params = append(params, fileId)
	}
	sql = fmt.Sprintf("delete form association where app_id = ? and file_id in (%s)", strings.Join(patterns, ","))
	_, err = db.Exec(sql, params...)
	if err != nil {
		log.Errorf("Exec delete association of app_id [%s] error: %s", app.Id, err)
		return err
	}
	// update app outdated
	_, err = db.Exec("update app set outdated = 1 where id = ?", app.Id)
	if err != nil {
		log.Errorf("Update app_id [%s] outdated error: %s", app.Id, err)
		return err
	}
	return nil
}

func ListApp() ([]map[string]interface{}, error) {
	rows, err := db.Query("select id, name, outdated from app")
	if err != nil {
		log.Error("Query all app error: ", err)
		return nil, err
	}
	result := make([]map[string]interface{}, 10)
	for rows.Next() {
		app := App{}
		err = rows.Scan(&app.Id, &app.Name, &app.Outdated)
		if err != nil {
			log.Error("Scan app error: ", err)
			return nil, err
		}
		appBrief, err := app.Brief()
		if err != nil {
			log.Errorf("Get app [%s] brief error: %s", app.Name, err)
			return nil, err
		}
		result = append(result, appBrief)
	}
	return result, nil
}

func ViewApp(appId int64) (map[string]interface{}, error) {
	var app App
	err := db.QueryRow("select id, name, outdated from app where id = ?", appId).Scan(&app.Id, &app.Name, &app.Outdated)
	if err != nil {
		log.Errorf("Query app_id [%s] error: %s", appId, err)
		return nil, err
	}
	return app.Brief()
}
