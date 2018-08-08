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
	"github.com/cflion/cflion/pkg/common"
	"github.com/cflion/cflion/pkg/log"
	"github.com/cflion/cflion/pkg/manager/api"
	"github.com/coreos/etcd/clientv3"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type Repository interface {
	ExistsAppById(id int64) bool
	ExistsAppByName(name string) bool
	GetAppByName(name string) (*api.App, error)
	InsertApp(app *api.App) (int64, error)
	RetrieveAppBrief(id int64) (*api.App, error)
	RetrieveAppDetail(id int64) (*api.App, error)
	UpdateAppAssociation(appId int64, addFileIds []int64, delFileIds []int64) error
	UpdateAppOutdated(id int64, outdated bool) error

	ListConfigFilesBrief() ([]*api.ConfigFile, error)
	ExistsConfigFileByNameAndNamespaceId(filename string, namespaceId int64) bool
	ExistsConfigFileById(id int64) bool
	InsertConfigFileWithItems(cf *api.ConfigFile) (int64, error)
	RetrieveConfigFileDetail(id int64) (*api.ConfigFile, error)
	UpdateConfigFile(fileId int64, items []*api.ConfigItem) error
}

type ServiceImpl struct {
	Repo Repository
}

func (service *ServiceImpl) ListApps() ([]map[string]interface{}, error) {
	panic("implement me")
}

func (service *ServiceImpl) ExistsAppById(id int64) bool {
	return service.Repo.ExistsAppById(id)
}

func (service *ServiceImpl) ExistsAppByName(name string) bool {
	return service.Repo.ExistsAppByName(name)
}

func (service *ServiceImpl) GetAppByName(name string) (*api.App, error) {
	return service.Repo.GetAppByName(name)
}

func (service *ServiceImpl) CreateApp(name string) (int64, error) {
	app := &api.App{Name: name, Outdated: 1}
	return service.Repo.InsertApp(app)
}

func (service *ServiceImpl) ViewApp(id int64) (map[string]interface{}, error) {
	app, err := service.Repo.RetrieveAppBrief(id)
	if err != nil {
		return nil, err
	}
	return app.Brief(), nil
}

func (service *ServiceImpl) UpdateAppAssociation(id int64, fileIds []int64) error {
	fileIds = common.DistinctInt64Slice(fileIds)
	cg, err := service.Repo.RetrieveAppBrief(id)
	if err != nil {
		return err
	}
	curFileIds := make([]int64, 0, len(cg.Files))
	for _, cf := range cg.Files {
		curFileIds = append(curFileIds, cf.Id)
	}
	addFileIds, delFileIds := common.DiffTwoInt64Slice(fileIds, curFileIds)
	return service.Repo.UpdateAppAssociation(id, addFileIds, delFileIds)
}

func (service *ServiceImpl) PublishApp(id int64) error {
	app, err := service.Repo.RetrieveAppDetail(id)
	if err != nil {
		return err
	}
	value := app.ConfigFmt()
	etcdEndpoints := viper.GetStringSlice("etcd.endpoints")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Errorf("Connect to etcd [%s] error: %s", etcdEndpoints, err)
		return err
	}
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(viper.GetInt("etcd.requestTimeout"))*time.Second)
	_, err = cli.Put(ctx, app.Key(), value)
	cancel()
	if err != nil {
		log.Errorf("Put [key=%s] [value=%s] into etcd error: %s", app.Key(), value, err)
		return err
	}
	err = service.Repo.UpdateAppOutdated(id, false)
	return err
}

func (service *ServiceImpl) ListConfigFiles() ([]map[string]interface{}, error) {
	cfs, err := service.Repo.ListConfigFilesBrief()
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, 0, len(cfs))
	for _, cf := range cfs {
		result = append(result, cf.Brief())
	}
	return result, nil
}

func (service *ServiceImpl) ExistsConfigFileByNameAndNamespaceId(filename string, namespaceId int64) bool {
	return service.Repo.ExistsConfigFileByNameAndNamespaceId(filename, namespaceId)
}

func (service *ServiceImpl) ExistsConfigFileById(id int64) bool {
	return service.Repo.ExistsConfigFileById(id)
}

func (service *ServiceImpl) CreateConfigFile(name string, namespaceId int64, content string) (int64, error) {
	cis := parseContent(content)
	cf := &api.ConfigFile{Name: name, NamespaceId: namespaceId, Items: cis}
	return service.Repo.InsertConfigFileWithItems(cf)
}

func (service *ServiceImpl) ViewConfigFile(id int64) (map[string]interface{}, error) {
	cf, err := service.Repo.RetrieveConfigFileDetail(id)
	if err != nil {
		return nil, err
	}
	return cf.Detail(), nil
}

func (service *ServiceImpl) UpdateConfigFile(id int64, content string) error {
	cis := parseContent(content)
	return service.Repo.UpdateConfigFile(id, cis)
}

func parseContent(content string) []*api.ConfigItem {
	lines := strings.Split(content, "\n")
	current := &api.ConfigItem{}
	items := make([]*api.ConfigItem, 0, len(lines))
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
            current = &api.ConfigItem{}
		}
	}
	return items
}
