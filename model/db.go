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
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var db *sql.DB

// SetupDB init a *sql.DB object.
func SetupDB() error {
	host := viper.GetString("db.host")
	port := viper.GetInt("db.port")
	username := viper.GetString("db.username")
	password := viper.GetString("db.password")
	database := viper.GetString("db.database")
	maxIdle := viper.GetInt("db.maxIdle")
	maxOpen := viper.GetInt("db.maxOpen")

	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, database)
	tempDb, err := sql.Open("mysql", url)
	if err != nil {
		return err
	}
	tempDb.SetMaxIdleConns(maxIdle)
	tempDb.SetMaxOpenConns(maxOpen)
	err = tempDb.Ping()
	if err != nil {
		return err
	}
	db = tempDb
	return nil
}
