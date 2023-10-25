package server

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/i2dou/sponge/pkg/gofile"
	"github.com/i2dou/sponge/pkg/logger"
	"github.com/i2dou/sponge/pkg/utils"
)

var (
	dataFile = saveDir + "/data.json"
	rcd      *record
)

type parameters struct {
	ServerName    string `json:"serverName"`
	ProjectName   string `json:"projectName"`
	ModuleName    string `json:"moduleName"`
	RepoAddr      string `json:"repoAddr"`
	ProtobufFile  string `json:"-"`
	YamlFile      string `json:"-"`
	Dsn           string `json:"dsn"`
	TableName     string `json:"tableName"`
	Embed         bool   `json:"embed"`
	IncludeInitDB bool   `json:"includeInitDB"`
	UpdateAt      string `json:"updateAt"`
}

type record struct {
	mux        *sync.Mutex
	HostRecord map[string]*parameters
}

func initRecord() {
	rcd = &record{
		mux:        new(sync.Mutex),
		HostRecord: make(map[string]*parameters),
	}

	data, err := os.ReadFile(dataFile)
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, &rcd.HostRecord)
}

func recordObj() *record {
	return rcd
}

func (r *record) set(ip string, commandType string, params *parameters) {
	utils.SafeRunWithTimeout(time.Second*3, func(cancel context.CancelFunc) {
		r.mux.Lock()
		defer func() {
			r.mux.Unlock()
			cancel()
		}()

		key := getKey(ip, commandType)
		r.HostRecord[key] = params
		data, err := json.Marshal(r.HostRecord)
		if err != nil {
			logger.Warn("json marshal error", logger.Err(err))
			return
		}

		var file = dataFile
		if gofile.IsWindows() {
			file = strings.ReplaceAll(dataFile, "/", "\\")
		}
		dir := gofile.GetFileDir(file)
		_ = gofile.CreateDir(dir)
		err = os.WriteFile(file, data, 0666)
		if err != nil {
			logger.Warn("WriteFile error", logger.Err(err))
			return
		}
	})
}

func getKey(ip string, commandType string) string {
	if ip == "::1" {
		ip = "127.0.0.1"
	}
	return ip + "-" + commandType
}

func (r *record) get(ip string, commandType string) *parameters {
	r.mux.Lock()
	defer r.mux.Unlock()
	key := getKey(ip, commandType)
	return r.HostRecord[key]
}

func parseCommandArgs(args []string) *parameters {
	var params = &parameters{UpdateAt: time.Now().Format("20060102T150405")}
	for _, v := range args {
		ss := strings.SplitN(v, "=", 2)
		if len(ss) == 1 {
			if ss[0] == "--embed" {
				params.Embed = true
			}
			if ss[0] == "--include-init-db" {
				params.IncludeInitDB = true
			}
		} else {
			val := ss[1]
			switch ss[0] {
			case "--db-dsn":
				params.Dsn = val
			case "--db-table":
				params.TableName = val
			case "--embed":
				if val == "true" {
					params.Embed = true
				} else {
					params.Embed = false
				}
			case "--include-init-db":
				if val == "true" {
					params.IncludeInitDB = true
				} else {
					params.IncludeInitDB = false
				}
			case "--module-name":
				params.ModuleName = val
			case "--project-name":
				params.ProjectName = val
			case "--server-name":
				if val != "" {
					val = strings.ReplaceAll(val, "-", "_")
				}
				params.ServerName = val
			case "--repo-addr":
				params.RepoAddr = val
			case "--protobuf-file":
				params.ProtobufFile = val
			case "--yaml-file":
				params.YamlFile = val
			}
		}
	}

	return params
}
