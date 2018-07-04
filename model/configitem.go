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

type ConfigItem struct {
	Id      uint
	FileId  uint
	Name    string
	Value   string
	Comment string
}

func (configItem *ConfigItem) Create() (int64, error) {
	stmt, err := db.Prepare("insert into config_item (file_id, name, value, comment, ctime, utime) values (?, ?, ?, ?, now(), now())")
	if err != nil {
		return -1, nil
	}
	defer stmt.Close()
	res, err := stmt.Exec(configItem.FileId, configItem.Name, configItem.Value, configItem.Comment)
	if err != nil {
		return -1, nil
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}
