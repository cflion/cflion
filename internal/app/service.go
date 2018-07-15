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

// Package app defines the service's interface of app, config_file.
package app

import (
	"fmt"
	"strings"
)

type Service interface {
	ListConfigGroup() ([]map[string]interface{}, error)
    ExistsConfigGroupByAppAndEnvironment(appName, environment string) bool
    ExistsConfigGroupById(id int64) bool
	CreateConfigGroup(appName, environment string) (int64, error)
	ViewConfigGroup(id int64) (map[string]interface{}, error)
	UpdateConfigGroupAssociation(id int64, fileIds []int64) error

	ListConfigFile() ([]map[string]interface{}, error)
    ExistsConfigFileByNameAndNamespaceId(filename string, namespaceId int64) bool
    ExistsConfigFileById(id int64) bool
    CreateConfigFile(name string, namespaceId int64, content string) (int64, error)
	ViewConfigFile(id int64) (map[string]interface{}, error)
    UpdateConfigFile(id int64, content string) error
}

// ConfigGroup defines the related structure of the config_group table in db.
type ConfigGroup struct {
	Id          int64
	App         string
	Environment string
	Outdated    byte

	Files []*ConfigFile
}

// ConfigFile defines the related structure of the config_file table in db.
type ConfigFile struct {
	Id          int64
	Name        string
	NamespaceId int64

	ConfigGroup *ConfigGroup
	Items       []*ConfigItem
}

// ConfigItem defines the related structure of the config_item table in db.
type ConfigItem struct {
	Id      int64
	FileId  int64
	Name    string
	Value   string
	Comment string
}

func (configGroup *ConfigGroup) String() string {
	return fmt.Sprintf("ConfigGroup {Id=%d | App=%s | Environment=%s | Outdated=%d | Files=%s}", configGroup.Id, configGroup.App, configGroup.Environment, configGroup.Outdated, configGroup.Files)
}

func (configGroup *ConfigGroup) FullName() string {
	return fmt.Sprintf("%s/%s", configGroup.App, configGroup.Environment)
}

func (configGroup *ConfigGroup) Brief() map[string]interface{} {
	configFiles := make([]map[string]interface{}, len(configGroup.Files))
	for _, file := range configGroup.Files {
		configFiles = append(configFiles, file.Brief())
	}
	return map[string]interface{}{
		"id":           configGroup.Id,
		"app":          configGroup.App,
		"environment":  configGroup.Environment,
		"full_name":    configGroup.FullName(),
		"outdated":     configGroup.Outdated,
		"config_files": configFiles,
	}
}

func (configFile *ConfigFile) String() string {
	return fmt.Sprintf("ConfigFile {Id=%d | Name=%s | NamespaceId=%s | ConfigGroup=%s | Items=%s}", configFile.Id, configFile.NamespaceId, configFile.NamespaceId, configFile.ConfigGroup, configFile.Items)
}

func (configFile *ConfigFile) Namespace() string {
	return configFile.ConfigGroup.FullName()
}

func (configFile *ConfigFile) FullName() string {
	return fmt.Sprintf("%s/%s", configFile.Namespace(), configFile.Name)
}

func (configFile *ConfigFile) ConfigFmt() string {
	arr := make([]string, len(configFile.Items))
	for _, item := range configFile.Items {
		arr = append(arr, item.ConfigFmt())
	}
	return strings.Join(arr, "\n")
}

func (configFile *ConfigFile) Brief() map[string]interface{} {
	return map[string]interface{}{
		"id":        configFile.Id,
		"name":      configFile.Name,
		"namespace": configFile.Namespace(),
		"full_name": configFile.FullName(),
	}
}

func (configFile *ConfigFile) Detail() map[string]interface{} {
	detail := configFile.Brief()
	detail["config"] = configFile.ConfigFmt()
	return detail
}

func (configItem *ConfigItem) String() string {
	return fmt.Sprintf("ConfigItem {Id=%d | FileId=%d | Name=%s | Value=%s | Comment=%s}", configItem.Id, configItem.FileId, configItem.Name, configItem.Value, configItem.Comment)
}

func (configItem *ConfigItem) ConfigFmt() string {
	var prefix string
	if len(configItem.Comment) > 0 {
		prefix = fmt.Sprintf("# %s\n", configItem.Comment)
	}
	return fmt.Sprintf("%s%s=%s", prefix, configItem.Name, configItem.Value)
}
