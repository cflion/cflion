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

package server

import (
	"context"
	"github.com/cflion/cflion/internal/app"
	"github.com/cflion/cflion/pkg/common"
	"github.com/cflion/cflion/pkg/log"
	"github.com/coreos/etcd/clientv3"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type Repository interface {
	ListConfigGroupBrief() ([]*app.ConfigGroup, error)
	ExistsConfigGroupByAppAndEnvironment(appName, environment string) bool
	ExistsConfigGroupById(id int64) bool
	InsertConfigGroup(cg *app.ConfigGroup) (int64, error)
	RetrieveConfigGroupBrief(id int64) (*app.ConfigGroup, error)
	RetrieveConfigGroupDetail(id int64) (*app.ConfigGroup, error)
	UpdateConfigGroupAssociation(groupId int64, addFileIds []int64, delFileIds []int64) error
	UpdateConfigGroupOutdated(id int64, outdated bool) error

	ListConfigFileBrief() ([]*app.ConfigFile, error)
	ExistsConfigFileByNameAndNamespaceId(filename string, namespaceId int64) bool
	ExistsConfigFileById(id int64) bool
	InsertConfigFileWithItems(cf *app.ConfigFile) (int64, error)
	RetrieveConfigFileDetail(id int64) (*app.ConfigFile, error)
	UpdateConfigFile(fileId int64, items []*app.ConfigItem) error
}

type Service struct {
	Repo Repository
}

func (service *Service) ListConfigGroup() ([]map[string]interface{}, error) {
	cgs, err := service.Repo.ListConfigGroupBrief()
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, len(cgs))
	for _, cg := range cgs {
		result = append(result, cg.Brief())
	}
	return result, nil
}

func (service *Service) ExistsConfigGroupByAppAndEnvironment(appName, environment string) bool {
	return service.Repo.ExistsConfigGroupByAppAndEnvironment(appName, environment)
}

func (service *Service) ExistsConfigGroupById(id int64) bool {
	return service.Repo.ExistsConfigGroupById(id)
}

func (service *Service) CreateConfigGroup(appName, environment string) (int64, error) {
	cg := &app.ConfigGroup{
		App:         appName,
		Environment: environment,
		Outdated:    1,
	}
	return service.Repo.InsertConfigGroup(cg)
}

func (service *Service) ViewConfigGroup(id int64) (map[string]interface{}, error) {
	cg, err := service.Repo.RetrieveConfigGroupBrief(id)
	if err != nil {
		return nil, err
	}
	return cg.Brief(), nil
}

func (service *Service) UpdateConfigGroupAssociation(id int64, fileIds []int64) error {
	fileIds = common.DistinctInt64Slice(fileIds)
	cg, err := service.Repo.RetrieveConfigGroupBrief(id)
	if err != nil {
		return err
	}
	curFileIds := make([]int64, len(cg.Files))
	for _, cf := range cg.Files {
		curFileIds = append(curFileIds, cf.Id)
	}
	addFileIds, delFileIds := common.DiffTwoInt64Slice(fileIds, curFileIds)
	return service.Repo.UpdateConfigGroupAssociation(id, addFileIds, delFileIds)
}

func (service *Service) PublishConfigGroup(id int64) error {
	cg, err := service.Repo.RetrieveConfigGroupDetail(id)
	if err != nil {
		return err
	}
	value := cg.ConfigFmt()
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   viper.GetStringSlice("etcd." + cg.Environment),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Errorf("Connect to etcd [%s] error: %s", cg.Environment, err)
		return err
	}
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(viper.GetInt("etcd.requestTimeout"))*time.Second)
	_, err = cli.Put(ctx, cg.Key(), value)
	cancel()
	if err != nil {
		log.Errorf("Put [key=%s] [value=%s] into etcd error: %s", cg.Key(), value, err)
		return err
	}
	err = service.Repo.UpdateConfigGroupOutdated(id, false)
	return err
}

func (service *Service) ListConfigFile() ([]map[string]interface{}, error) {
	cfs, err := service.Repo.ListConfigFileBrief()
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, len(cfs))
	for _, cf := range cfs {
		result = append(result, cf.Brief())
	}
	return result, nil
}

func (service *Service) ExistsConfigFileByNameAndNamespaceId(filename string, namespaceId int64) bool {
	return service.Repo.ExistsConfigFileByNameAndNamespaceId(filename, namespaceId)
}

func (service *Service) CreateConfigFile(name string, namespaceId int64, content string) (int64, error) {
	cis := parseContent(content)
	cf := &app.ConfigFile{Name: name, NamespaceId: namespaceId, Items: cis}
	return service.Repo.InsertConfigFileWithItems(cf)
}

func parseContent(content string) []*app.ConfigItem {
	lines := strings.Split(content, "\n")
	current := &app.ConfigItem{}
	items := make([]*app.ConfigItem, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if line[0:1] == "#" {
			current.Comment = strings.TrimSpace(line[1:])
		} else {
			kv := strings.Split(line, "=")
			if len(kv) != 2 {
				log.Errorf("Parse config [%s] failed", line)
				continue
			}
			current.Name, current.Value = kv[0], kv[1]
			items = append(items, current)
			current = &app.ConfigItem{}
		}
	}
	return items
}

func (service *Service) ExistsConfigFileById(id int64) bool {
	return service.Repo.ExistsConfigFileById(id)
}

func (service *Service) ViewConfigFile(id int64) (map[string]interface{}, error) {
	cf, err := service.Repo.RetrieveConfigFileDetail(id)
	if err != nil {
		return nil, err
	}
	return cf.Detail(), nil
}

func (service *Service) UpdateConfigFile(id int64, content string) error {
	cis := parseContent(content)
	return service.Repo.UpdateConfigFile(id, cis)
}
