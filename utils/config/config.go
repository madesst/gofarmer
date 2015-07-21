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
	GlobalDir        = ".gofarmer"
	FarmsDir         = "farms"
	GlobalConfigName = "config.json"
	FarmConfigName   = "farm.json"
)

var farmConfigs map[string]Config = nil
var globalConfig *GlobalConfig = nil
var configDir = os.Getenv("HOME") + string(filepath.Separator) + GlobalDir + string(filepath.Separator)

func GetGlobal() *GlobalConfig {
	/*
		1. Check and prepare internal dirs
		2. Check and read global config
		3. Find all farms subdirs and read configs of each in farmConfigs
	*/
	checkDirs()

	if globalConfig == nil {

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

	return globalConfig
}

func Get() map[string]Config {
	return farmConfigs
}

func checkDirs() {
	e := os.MkdirAll(configDir, DirsMask)
	if e != nil {
		panic(e)
	}

	e = os.MkdirAll(configDir+FarmsDir, DirsMask)
	if e != nil {
		panic(e)
	}
}
