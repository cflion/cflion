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

package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	Database string
	MaxIdle  int
	MaxOpen  int
}

// ConnectDatabase use dsn and database's configs to connect to db.
func ConnectDatabase(cfg *DBConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", formatDSN(cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database))
	if err != nil {
		return nil, err
	}
	maxIdle := cfg.MaxIdle
	if maxIdle > 0 {
		db.SetMaxIdleConns(maxIdle)
	}
	maxOpen := cfg.MaxOpen
	if maxOpen > 0 {
		db.SetMaxOpenConns(maxOpen)
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// formatDSN generates dataSourceName.
func formatDSN(username string, password string, host string, port int, database string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, database)
}
