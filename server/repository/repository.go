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

package repository

import (
	"database/sql"
	"fmt"
	"github.com/cflion/cflion/internal/app"
	"github.com/cflion/cflion/pkg/log"
	"strings"
)

type Repository struct {
	DB *sql.DB
}

func (repo *Repository) UpdateConfigGroupOutdated(id int64, outdated bool) error {
	var out = 0
	if outdated {
		out = 1
	}
	_, err := repo.DB.Exec("update config_group set outdated = ? where id = ?", out, id)
	if err != nil {
		log.Errorf("UpdateConfigGroupOutdated update [id=%d] [outdated=%d] error: %s", id, out, err)
		return err
	}
	return nil
}

func (repo *Repository) ListConfigGroupBrief() ([]*app.ConfigGroup, error) {
	rows, err := repo.DB.Query("select id, app, environment, outdated from config_group")
	if err != nil {
		log.Errorf("Query all config_group error: %s", err)
		return nil, err
	}
	cgMap := make(map[int64]*app.ConfigGroup)
	for rows.Next() {
		var cg app.ConfigGroup
		rows.Scan(&cg.Id, &cg.App, &cg.Environment, &cg.Outdated)
		cgMap[cg.Id] = &cg
	}

	rows, err = repo.DB.Query("select group_id, file_id from association")
	if err != nil {
		log.Errorf("Query all association error: %s", err)
		return nil, err
	}
	assMap := make(map[int64][]int64)
	for rows.Next() {
		var groupId, fileId int64
		rows.Scan(&groupId, &fileId)
		if _, ok := assMap[groupId]; !ok {
			assMap[groupId] = make([]int64, 4)
		}
		assMap[groupId] = append(assMap[groupId], fileId)
	}

	rows, err = repo.DB.Query("select id, name, namespace_id from config_file")
	if err != nil {
		log.Errorf("Query all config_file error: %s", err)
		return nil, err
	}
	cfMap := make(map[int64]*app.ConfigFile)
	for rows.Next() {
		var cf app.ConfigFile
		rows.Scan(&cf.Id, &cf.Name, &cf.NamespaceId)
		cfMap[cf.Id] = &cf
	}

	result := make([]*app.ConfigGroup, 10)
	for groupId, cg := range cgMap {
		rcg := *cg
		fileIds := assMap[groupId]
		files := make([]*app.ConfigFile, len(fileIds))
		for _, fileId := range fileIds {
			cf := cfMap[fileId]
			cf.ConfigGroup = cg
			files = append(files, cf)
		}
		rcg.Files = files
		result = append(result, &rcg)
	}
	return result, nil
}

func (repo *Repository) ExistsConfigGroupByAppAndEnvironment(appName, environment string) bool {
	var count int64
	err := repo.DB.QueryRow("select count(1) from config_group where app = ? and environment = ?", appName, environment).Scan(&count)
	if err != nil {
		log.Errorf("Count config_group [app=%s] [environment=%s] error: %s", appName, environment, err)
		return false
	}
	return count == 1
}

func (repo *Repository) ExistsConfigGroupById(id int64) bool {
	var count int64
	err := repo.DB.QueryRow("select count(1) from config_group where id = ?", id).Scan(&count)
	if err != nil {
		log.Errorf("Count config_group [id=%d] error: %s", id, err)
		return false
	}
	return count == 1
}

func (repo *Repository) InsertConfigGroup(cg *app.ConfigGroup) (int64, error) {
	res, err := repo.DB.Exec("insert into config_group (app, environment, outdated, ctime, utime) values (?, ?, ?, now(), now())", cg.App, cg.Environment, cg.Outdated)
	if err != nil {
		log.Errorf("Insert config_group [%s] error: %s", cg.String(), err)
		return -1, err
	}
	return res.LastInsertId()
}

func (repo *Repository) RetrieveConfigGroupBrief(id int64) (*app.ConfigGroup, error) {
	var cg app.ConfigGroup
	err := repo.DB.QueryRow("select id, app, environment, outdated from config_group where id = ?", id).Scan(&cg.Id, &cg.App, &cg.Environment, &cg.Outdated)
	if err != nil {
		log.Errorf("Scan config_group [id=%d] error: %s", id, err)
		return nil, err
	}
	rows, err := repo.DB.Query("select cf.id, cf.name, cf.namespace_id, cg.app, cg.environment, cg.outdated from association as a left join config_file as cf on a.file_id = cf.id left join config_group as cg on cf.namespace_id = cg.id where a.group_id = ?", id)
	if err != nil {
		log.Errorf("RetrieveConfigGroupBrief config_group [id=%d] when query config_file error: %s", id, err)
		return nil, err
	}
	cfs := make([]*app.ConfigFile, 8)
	for rows.Next() {
		var cf app.ConfigFile
		cf.ConfigGroup = &app.ConfigGroup{}
		rows.Scan(&cf.Id, &cf.Name, &cf.NamespaceId, &cf.ConfigGroup.App, &cf.ConfigGroup.Environment, &cf.ConfigGroup.Outdated)
		cfs = append(cfs, &cf)
	}
	cg.Files = cfs
	return &cg, nil
}

func (repo *Repository) RetrieveConfigGroupDetail(id int64) (*app.ConfigGroup, error) {
	cg, err := repo.RetrieveConfigGroupBrief(id)
	if err != nil {
		return nil, err
	}
	cfsDetail := make([]*app.ConfigFile, len(cg.Files))
	for _, cf := range cg.Files {
		cfDetail, err := repo.RetrieveConfigFileDetail(cf.Id)
		if err != nil {
			return nil, err
		}
		cfsDetail = append(cfsDetail, cfDetail)
	}
	cg.Files = cfsDetail
	return cg, nil
}

func (repo *Repository) RetrieveConfigFileDetail(id int64) (*app.ConfigFile, error) {
	var cf app.ConfigFile
	cf.ConfigGroup = &app.ConfigGroup{}
	err := repo.DB.QueryRow("select cf.id, cf.name, cf.namespace_id, cg.id as group_id, cg.app, cg.environment, cg.outdated from config_file as cf left join config_group as cg on cf.namespace_id = cg.id where cf.id = ?", id).Scan(&cf.Id, &cf.Name, &cf.NamespaceId, &cf.ConfigGroup.Id, &cf.ConfigGroup.App, &cf.ConfigGroup.Environment, &cf.ConfigGroup.Outdated)
	if err != nil {
		log.Errorf("RetrieveConfigFileDetail [id=%d] error: %s", id, err)
		return nil, err
	}
	rows, err := repo.DB.Query("select id, file_id, name, value, comment from config_item where file_id = ?", id)
	if err != nil {
		log.Errorf("RetrieveConfigFileDetail [id=%d] query config_item error: %s", id, err)
		return nil, err
	}
	cis := make([]*app.ConfigItem, 8)
	for rows.Next() {
		var ci app.ConfigItem
		rows.Scan(&ci.Id, &ci.FileId, &ci.Name, &ci.Value, &ci.Comment)
		cis = append(cis, &ci)
	}
	cf.Items = cis
	return &cf, nil
}

func (repo *Repository) UpdateConfigGroupAssociation(groupId int64, addFileIds []int64, delFileIds []int64) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		log.Errorf("UpdateConfigGroupAssociation begin transaction error: %s", err)
		return err
	}
	defer tx.Rollback()
	err = insertConfigGroupBatchAssociation(tx, groupId, addFileIds)
	if err != nil {
		return err
	}
	err = deleteConfigGroupBatchAssociation(tx, groupId, delFileIds)
	if err != nil {
		return err
	}
	_, err = tx.Exec("update config_group set outdated = 1 where id = ?", groupId)
	if err != nil {
		log.Errorf("UpdateConfigGroupAssociation update config_group [id=%d] outdated=1 error: %s", groupId, err)
		return err
	}
	err = tx.Commit()
	return err
}

func insertConfigGroupBatchAssociation(tx *sql.Tx, groupId int64, fileIds []int64) error {
	patterns := make([]string, len(fileIds))
	params := make([]interface{}, len(fileIds))
	for _, fileId := range fileIds {
		patterns = append(patterns, "(?, ?, now(), now())")
		params = append(params, groupId, fileId)
	}
	query := fmt.Sprintf("insert into association (group_id, file_id, ctime, utime) values %s", strings.Join(patterns, ","))
	_, err := tx.Exec(query, params...)
	if err != nil {
		log.Errorf("Insert config_group [id=%d] association [file_ids=%s] error: %s", groupId, fileIds, err)
		return err
	}
	return nil
}

func deleteConfigGroupBatchAssociation(tx *sql.Tx, groupId int64, fileIds []int64) error {
	patterns := make([]string, len(fileIds))
	params := make([]interface{}, len(fileIds)+1)
	params = append(params, groupId)
	for _, fileId := range fileIds {
		patterns = append(patterns, "?")
		params = append(params, fileId)
	}
	query := fmt.Sprintf("delete form association where app_id = ? and file_id in (%s)", strings.Join(patterns, ","))
	_, err := tx.Exec(query, params...)
	if err != nil {
		log.Errorf("Delete config_group [id=%d] association [file_ids=%s] error: %s", groupId, fileIds, err)
		return err
	}
	return nil
}

func (repo *Repository) ListConfigFileBrief() ([]*app.ConfigFile, error) {
	rows, err := repo.DB.Query("select cf.id, cf.name, cf.namespace_id, cg.id as group_id, cg.app, cg.environment, cg.outdated from config_file as cf left join config_group as cg on cf.namespace_id = cg.id")
	if err != nil {
		log.Errorf("Query all config_file error: %s", err)
		return nil, err
	}
	cfs := make([]*app.ConfigFile, 8)
	for rows.Next() {
		var cf app.ConfigFile
		cf.ConfigGroup = &app.ConfigGroup{}
		rows.Scan(&cf.Id, &cf.Name, &cf.NamespaceId, &cf.ConfigGroup.Id, &cf.ConfigGroup.App, &cf.ConfigGroup.Environment, &cf.ConfigGroup.Outdated)
		cfs = append(cfs, &cf)
	}
	return cfs, nil
}

func (repo *Repository) ExistsConfigFileByNameAndNamespaceId(filename string, namespaceId int64) bool {
	var count int64
	err := repo.DB.QueryRow("select count(1) from config_file where name = ? and namespace_id = ?", filename, namespaceId).Scan(&count)
	if err != nil {
		log.Errorf("Count config_file [name=%s] [namespace_id=%d] error: %s", filename, namespaceId, err)
		return false
	}
	return count == 1
}

func (repo *Repository) ExistsConfigFileById(id int64) bool {
	var count int64
	err := repo.DB.QueryRow("select count(1) from config_file where id = ?", id).Scan(&count)
	if err != nil {
		log.Errorf("Count config_file [id=%d] error: %s", id, err)
		return false
	}
	return count == 1
}

func (repo *Repository) InsertConfigFileWithItems(cf *app.ConfigFile) (int64, error) {
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
	_, err = tx.Exec("insert into association (group_id, file_id, ctime, utime) values (?, ?, now(), now())", cf.NamespaceId, fileId)
	if err != nil {
		log.Errorf("Insert association [group_id=%d] [file_id=%d] error: %s", cf.NamespaceId, fileId, err)
		return -1, err
	}
	// insert config_item
	patterns := make([]string, len(cf.Items))
	params := make([]interface{}, len(cf.Items))
	for _, item := range cf.Items {
		patterns = append(patterns, "(?, ?, ?, ?, now(), now())")
		params = append(params, fileId, item.Name, item.Value, item.Comment)
	}
	query := fmt.Sprintf("insert into config_item (file_id, name, value, comment, ctime, utime) values %s", strings.Join(patterns, ","))
	_, err = tx.Exec(query, params...)
	if err != nil {
		log.Error("Insert batch config_item %s error: %s", cf.Items, err)
		return -1, err
	}
	err = tx.Commit()
	return fileId, err
}

func (repo *Repository) UpdateConfigFile(fileId int64, items []*app.ConfigItem) error {
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
	oldItems := make(map[string]*app.ConfigItem, len(items))
	for rows.Next() {
		var ci app.ConfigItem
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
		tx.Exec("update config_group as a set a.outdated = 1 where a.id in (select ass.group_id from association as ass where ass.file_id = ?)", fileId)
	}
	err = tx.Commit()
	return err
}
