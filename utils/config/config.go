package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	AppVersion       = "0.0.1"
	DirsMask         = 0770
	FilesMask        = 0600
	GlobalDirName    = ".gofarmer"
	FarmsDirName     = "farms"
	GlobalConfigName = "config.json"
	FarmConfigName   = "farm.json"
)

var farmConfigs map[string]Config
var globalConfig *GlobalConfig
var configDir = os.Getenv("HOME") + string(filepath.Separator) + GlobalDirName + string(filepath.Separator)
var farmsDir = os.Getenv("HOME") + string(filepath.Separator) + GlobalDirName + string(filepath.Separator) + FarmsDirName + string(filepath.Separator)

func GetGlobal() GlobalConfig {
	prepInternals()
	return *globalConfig
}

func Get(name string) Config {
	prepInternals()
	if val, e := farmConfigs[name]; e {
		return val
	}
	panic("Undefined farm")
}

func prepInternals() {
	/*
		1. Check and prepare internal dirs
		2. Check and read global config
		3. Find all farms subdirs and read configs of each in farmConfigs
	*/
	checkDirs()

	if globalConfig == nil {
		prepGlobalConfig()
	}

	if farmConfigs == nil {
		prepFarmsConfig()
	}
}

func prepGlobalConfig() {
	rawConfig, e := ioutil.ReadFile(configDir + GlobalConfigName)
	if e != nil {
		gc := new(GlobalConfig)
		gc.Version = AppVersion
		//Nasty
		globalConfig = gc
		rawConfig, _ := json.Marshal(globalConfig)

		e = ioutil.WriteFile(configDir+GlobalConfigName, rawConfig, 0600)
		if e != nil {
			panic(e)
		}
	} else {
		json.Unmarshal(rawConfig, &globalConfig)
	}
}

func prepFarmsConfig() {
	farms, _ := ioutil.ReadDir(farmsDir)
	for _, descriptor := range farms {
		if !descriptor.Mode().IsDir() {
			continue
		}

		rawFarmConfig, e := ioutil.ReadFile(farmsDir + descriptor.Name() + string(filepath.Separator) + FarmConfigName)
		if e != nil {
			continue
		}

		farmConfig := new(Config)
		json.Unmarshal(rawFarmConfig, &farmConfig)
		farmConfigs[descriptor.Name()] = *farmConfig
	}
}

func checkDirs() {
	e := os.MkdirAll(configDir, DirsMask)
	if e != nil {
		panic(e)
	}

	e = os.MkdirAll(farmsDir, DirsMask)
	if e != nil {
		panic(e)
	}
}
