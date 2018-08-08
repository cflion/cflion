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
	"fmt"
	"github.com/cflion/cflion/pkg/log"
	"github.com/cflion/cflion/pkg/manager/api"
	"strings"
)

type RepositoryImpl struct {
	DB *sql.DB
}

func (repo *RepositoryImpl) ExistsAppById(id int64) bool {
	var count int64
	err := repo.DB.QueryRow("select count(1) from app where id = ?", id).Scan(&count)
	if err != nil {
		log.Errorf("Count app [id=%d] error: %s", id, err)
		return false
	}
	return count == 1
}

func (repo *RepositoryImpl) ExistsAppByName(name string) bool {
	var count int64
	err := repo.DB.QueryRow("select count(1) from app where name = ?", name).Scan(&count)
	if err != nil {
		log.Errorf("Count app [name=%d] error: %s", name, err)
		return false
	}
	return count == 1
}

func (repo *RepositoryImpl) GetAppByName(name string) (*api.App, error) {
	var app api.App
	err := repo.DB.QueryRow("select id, name, outdated from app where name = ?", name).Scan(&app.Id, &app.Name, &app.Outdated)
	if err != nil {
		log.Errorf("Get app [name=%s] error: %s", name, err)
		return nil, err
	}
	return &app, nil
}

func (repo *RepositoryImpl) InsertApp(app *api.App) (int64, error) {
	res, err := repo.DB.Exec("insert into app (name, outdated, ctime, utime) values (?, ?, now(), now())", app.Name, app.Outdated)
	if err != nil {
		log.Errorf("Insert app [%s] error: %s", app.String(), err)
		return -1, err
	}
	return res.LastInsertId()
}

func (repo *RepositoryImpl) RetrieveAppBrief(id int64) (*api.App, error) {
	var app api.App
	err := repo.DB.QueryRow("select id, name, outdated from app where id = ?", id).Scan(&app.Id, &app.Name, &app.Outdated)
	if err != nil {
		log.Errorf("RetrieveAppBrief app [id=%d] when scan app error: %s", id, err)
		return nil, err
	}
	rows, err := repo.DB.Query("select cf.id, cf.name, cf.namespace_id, app.id as app_id, app.name as app_name, app.outdated from association as ass left join config_file as cf on ass.file_id = cf.id left join app on cf.namespace_id = app.id where ass.app_id = ?", id)
	if err != nil {
		log.Errorf("RetrieveAppBrief app [id=%d] when query config_file error: %s", id, err)
		return nil, err
	}
	defer rows.Close()
	cfs := make([]*api.ConfigFile, 0, 8)
	for rows.Next() {
		var cf api.ConfigFile
		cf.App = &api.App{}
		rows.Scan(&cf.Id, &cf.Name, &cf.NamespaceId, &cf.App.Id, &cf.App.Name, &cf.App.Outdated)
		cfs = append(cfs, &cf)
	}
	app.Files = cfs
	return &app, nil
}

func (repo *RepositoryImpl) RetrieveAppDetail(id int64) (*api.App, error) {
	app, err := repo.RetrieveAppBrief(id)
	if err != nil {
		return nil, err
	}
	cfsDetail := make([]*api.ConfigFile, 0, len(app.Files))
	for _, cf := range app.Files {
		cfDetail, err := repo.RetrieveConfigFileDetail(cf.Id)
		if err != nil {
			return nil, err
		}
		cfsDetail = append(cfsDetail, cfDetail)
	}
	app.Files = cfsDetail
	return app, nil
}

func (repo *RepositoryImpl) UpdateAppAssociation(appId int64, addFileIds []int64, delFileIds []int64) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		log.Errorf("UpdateAppAssociation begin transaction error: %s", err)
		return err
	}
	defer tx.Rollback()
	err = insertAppBatchAssociation(tx, appId, addFileIds)
	if err != nil {
		return err
	}
	err = deleteAppBatchAssociation(tx, appId, delFileIds)
	if err != nil {
		return err
	}
	_, err = tx.Exec("update app set outdated = 1 where id = ?", appId)
	if err != nil {
		log.Errorf("UpdateAppAssociation app [id=%d] [outdated=1] error: %s", appId, err)
		return err
	}
	err = tx.Commit()
	return err
}

func (repo *RepositoryImpl) UpdateAppOutdated(id int64, outdated bool) error {
	var out = 0
	if outdated {
		out = 1
	}
	_, err := repo.DB.Exec("update app set outdated = ? where id = ?", out, id)
	if err != nil {
		log.Errorf("UpdateAppOutdated update [id=%d] [outdated=%d] error: %s", id, out, err)
		return err
	}
	return nil
}

func (repo *RepositoryImpl) ListConfigFilesBrief() ([]*api.ConfigFile, error) {
	rows, err := repo.DB.Query("select cf.id, cf.name, cf.namespace_id, app.id as app_id, app.name as app_name, app.outdated from config_file as cf left join app on cf.namespace_id = app.id")
	if err != nil {
		log.Errorf("ListConfigFileBrief error: %s", err)
		return nil, err
	}
	defer rows.Close()
	cfs := make([]*api.ConfigFile, 0, 8)
	for rows.Next() {
		var cf api.ConfigFile
		cf.App = &api.App{}
		rows.Scan(&cf.Id, &cf.Name, &cf.NamespaceId, &cf.App.Id, &cf.App.Name, &cf.App.Outdated)
		cfs = append(cfs, &cf)
	}
	return cfs, nil
}

func (repo *RepositoryImpl) ExistsConfigFileByNameAndNamespaceId(filename string, namespaceId int64) bool {
	var count int64
	err := repo.DB.QueryRow("select count(1) from config_file where name = ? and namespace_id = ?", filename, namespaceId).Scan(&count)
	if err != nil {
		log.Errorf("Count config_file [name=%s] [namespace_id=%d] error: %s", filename, namespaceId, err)
		return false
	}
	return count == 1
}

func (repo *RepositoryImpl) ExistsConfigFileById(id int64) bool {
	var count int64
	err := repo.DB.QueryRow("select count(1) from config_file where id = ?", id).Scan(&count)
	if err != nil {
		log.Errorf("Count config_file [id=%d] error: %s", id, err)
		return false
	}
	return count == 1
}

func (repo *RepositoryImpl) InsertConfigFileWithItems(cf *api.ConfigFile) (int64, error) {
	tx, err := repo.DB.Begin()
	if err != nil {
		log.Errorf("InsertConfigFileWithItems begin transaction error: %s", err)
		return -1, err
	}
	defer tx.Rollback()
	res, err := tx.Exec("insert into config_file (name, namespace_id, ctime, utime) values (?, ?, now(), now())", cf.Name, cf.NamespaceId)
	if err != nil {
		log.Errorf("Insert config_file [%s] error: %s", cf, err)
		return -1, err
	}
	fileId, err := res.LastInsertId()
	if err != nil {
		log.Errorf("Get config_file insert id error: %s", err)
		return -1, err
	}
	_, err = tx.Exec("insert into association (app_id, file_id, ctime, utime) values (?, ?, now(), now())", cf.NamespaceId, fileId)
	if err != nil {
		log.Errorf("Insert association [app_id=%d] [file_id=%d] error: %s", cf.NamespaceId, fileId, err)
		return -1, err
	}
	// insert config_item
	patterns := make([]string, 0, len(cf.Items))
	params := make([]interface{}, 0, len(cf.Items))
	for _, item := range cf.Items {
		patterns = append(patterns, "(?, ?, ?, ?, now(), now())")
		params = append(params, fileId, item.Name, item.Value, item.Comment)
	}
	query := fmt.Sprintf("insert into config_item (file_id, name, value, comment, ctime, utime) values %s", strings.Join(patterns, ","))
	log.Info("query=", query)
	_, err = tx.Exec(query, params...)
	if err != nil {
		log.Error("Insert batch config_item %s error: %s", cf.Items, err)
		return -1, err
	}
	err = tx.Commit()
	return fileId, err
}

func (repo *RepositoryImpl) RetrieveConfigFileDetail(id int64) (*api.ConfigFile, error) {
	var cf api.ConfigFile
	cf.App = &api.App{}
	err := repo.DB.QueryRow("select cf.id, cf.name, cf.namespace_id, app.id as app_id, app.name as app_name, app.outdated from config_file as cf left join app on cf.namespace_id = app.id where cf.id = ?", id).Scan(&cf.Id, &cf.Name, &cf.NamespaceId, &cf.App.Id, &cf.App.Name, &cf.App.Outdated)
	if err != nil {
		log.Errorf("RetrieveConfigFileDetail [id=%d] error: %s", id, err)
		return nil, err
	}
	rows, err := repo.DB.Query("select id, file_id, name, value, comment from config_item where file_id = ?", id)
	if err != nil {
		log.Errorf("RetrieveConfigFileDetail [id=%d] query config_item error: %s", id, err)
		return nil, err
	}
	defer rows.Close()
	cis := make([]*api.ConfigItem, 0, 8)
	for rows.Next() {
		var ci api.ConfigItem
		rows.Scan(&ci.Id, &ci.FileId, &ci.Name, &ci.Value, &ci.Comment)
		cis = append(cis, &ci)
	}
	cf.Items = cis
	return &cf, nil
}

func (repo *RepositoryImpl) UpdateConfigFile(fileId int64, items []*api.ConfigItem) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		log.Errorf("UpdateConfigFile begin transaction error: %s", err)
		return err
	}
	defer tx.Rollback()
	rows, err := tx.Query("select id, name, value, comment from config_item where file_id = ?", fileId)
	if err != nil {
		log.Errorf("UpdateConfigFile config_file [id=%d] query config_item error: %s", fileId, err)
		return err
	}
	oldItems := make(map[string]*api.ConfigItem, len(items))
	for rows.Next() {
		var ci api.ConfigItem
		rows.Scan(&ci.Id, &ci.Name, &ci.Value, &ci.Comment)
		oldItems[ci.Name] = &ci
	}
	outdated := false
	for _, ci := range items {
		if oldItem, ok := oldItems[ci.Name]; ok {
			if ci.Value != oldItem.Value || ci.Comment != oldItem.Comment {
				tx.Exec("update config_item set value = ?, comment = ? where id = ?", ci.Value, ci.Comment, oldItem.Id)
				outdated = true
			}
		} else {
			tx.Exec("insert into config_item (file_id, name, value, comment, ctime, utime) values (?, ?, ?, ?, now(), now())", fileId, ci.Name, ci.Value, ci.Comment)
		}
	}
	if outdated {
		tx.Exec("update app set app.outdated = 1 where app.id in (select ass.app_id from association as ass where ass.file_id = ?)", fileId)
	}
	err = tx.Commit()
	return err
}

func insertAppBatchAssociation(tx *sql.Tx, appId int64, fileIds []int64) error {
	patterns := make([]string, 0, len(fileIds))
	params := make([]interface{}, 0, len(fileIds))
	for _, fileId := range fileIds {
		patterns = append(patterns, "(?, ?, now(), now())")
		params = append(params, appId, fileId)
	}
	query := fmt.Sprintf("insert into association (app_id, file_id, ctime, utime) values %s", strings.Join(patterns, ","))
	_, err := tx.Exec(query, params...)
	if err != nil {
		log.Errorf("Insert app [id=%d] association [file_ids=%s] error: %s", appId, fileIds, err)
		return err
	}
	return nil
}

func deleteAppBatchAssociation(tx *sql.Tx, appId int64, fileIds []int64) error {
	patterns := make([]string, 0, len(fileIds))
	params := make([]interface{}, 0, len(fileIds)+1)
	params = append(params, appId)
	for _, fileId := range fileIds {
		patterns = append(patterns, "?")
		params = append(params, fileId)
	}
	query := fmt.Sprintf("delete form association where app_id = ? and file_id in (%s)", strings.Join(patterns, ","))
	_, err := tx.Exec(query, params...)
	if err != nil {
		log.Errorf("Delete app [id=%d] association [file_ids=%s] error: %s", appId, fileIds, err)
		return err
	}
	return nil
}
