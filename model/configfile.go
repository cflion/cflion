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
	"strings"
)

type ConfigFile struct {
	Id    int64
	Name  string
	AppId int64
}

func parseContent(content string) []ConfigItem {
	lines := strings.Split(content, "\n")
	current := ConfigItem{}
	items := make([]ConfigItem, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line[0:1] == "#" {
			current.Comment = line[1:]
		} else {
			kv := strings.Split(line, "=")
			if len(kv) != 2 {
				log.Errorf("Parse config [%s] failed", line)
				continue
			}
			current.Name, current.Value = kv[0], kv[1]
			items = append(items, current)
			current = ConfigItem{}
		}
	}
	return items
}

func ExistsConfigFile(appId int64, name string) bool {
	var count int64
	err := db.QueryRow("select count(1) from config_file where app_id = ? and name = ?", appId, name).Scan(&count)
	if err != nil {
		log.Errorf("Query config_file [%s] error: %s", name, err)
		return false
	}
	if count >= 1 {
		return true
	} else {
		return false
	}
}

func (configFile *ConfigFile) Create() (int64, error) {
	stmt, err := db.Prepare("insert into config_file (name, app_id, ctime, utime) values (?, ?, now(), now())")
	if err != nil {
		log.Error("Prepare error when creates config_file: ", err)
		return -1, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(configFile.Name, configFile.AppId)
	if err != nil {
		log.Error("Exec error when creates config_file: ", err)
		return -1, err
	}
	return res.LastInsertId()
}

func (configFile *ConfigFile) CreateWithContent(content string) error {
	items := parseContent(content)
	tx, err := db.Begin()
	defer tx.Rollback()
	// insert config_file
	res, err := tx.Exec("insert into config_file (name, app_id, ctime, utime) values (?, ?, now(), now())", configFile.Name, configFile.AppId)
	if err != nil {
		log.Error("Exec error when creates config_file: ", err)
		return err
	}
	fileId, err := res.LastInsertId()
	if err != nil {
		log.Error("LastInsertId error of config_file: ", err)
		return err
	}
	// insert association
	_, err = tx.Exec("insert into association (app_id, file_id, ctime, utime) values (?, ?, now(), now())", configFile.AppId, fileId)
	if err != nil {
		log.Error("Exec error when insert association: ", err)
		return err
	}
	// insert config_item
	patterns := make([]string, len(items))
	params := make([]interface{}, len(items))
	for _, item := range items {
		patterns = append(patterns, "(?, ?, ?, ?, now(), now())")
		params = append(params, fileId, item.Name, item.Value, item.Comment)
	}
	sql := fmt.Sprintf("insert into config_item (file_id, name, value, comment, ctime, utime) values %s", strings.Join(patterns, ","))
	_, err = tx.Exec(sql, params...)
	if err != nil {
		log.Error("Exec error when insert config_item: ", err)
		return err
	}

	err = tx.Commit()
	return err
}
