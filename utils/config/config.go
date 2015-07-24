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

var globalParsed bool = false

var farmConfigs FarmConfigs = FarmConfigs{}
var globalConfig GlobalConfig = GlobalConfig{
	Version:       AppVersion,
	DefaultRegion: "us-east-1",
	Quotas: Quotas{
		MaxInstances: 1,
		MinInstances: 0,
		MaxPrice:     0.07,
		MaxAmount:    1,
	},
}

var Sep string = string(filepath.Separator)
var configDir string = os.Getenv("HOME") + Sep + GlobalDirName + Sep
var farmsDir string = os.Getenv("HOME") + Sep + GlobalDirName + Sep + FarmsDirName + Sep

func init() {
	prepInternals()
}

func GetGlobal() GlobalConfig {
	return globalConfig
}

func GetFarms() FarmConfigs {
	return farmConfigs
}

func GetFarm(name string) *FarmConfig {
	if val, e := farmConfigs[name]; e {
		return &val
	}

	return nil
}

func CreateFarm(name string, fc FarmConfig) FarmConfig {
	if e := os.MkdirAll(farmsDir+name, DirsMask); e != nil {
		panic(e)
	}

	rawConfig, _ := json.Marshal(fc)
	if e := ioutil.WriteFile(farmsDir+name+Sep+FarmConfigName, rawConfig, FilesMask); e != nil {
		panic(e)
	}

	farmConfigs[name] = fc
	return farmConfigs[name]
}

func prepInternals() {
	checkDirs()

	if !globalParsed {
		prepGlobalConfig()
	}

	if len(farmConfigs) == 0 {
		prepFarmsConfig()
	}
}

func prepGlobalConfig() {
	rawConfig, e := ioutil.ReadFile(configDir + GlobalConfigName)
	if e != nil {
		rawConfig, _ := json.Marshal(globalConfig)

		if e = ioutil.WriteFile(configDir+GlobalConfigName, rawConfig, FilesMask); e != nil {
			panic(e)
		}
	} else {
		json.Unmarshal(rawConfig, &globalConfig)
	}

	globalParsed = true
}

func prepFarmsConfig() {
	farmConfigs = FarmConfigs{}
	farms, _ := ioutil.ReadDir(farmsDir)
	for _, descriptor := range farms {
		if !descriptor.Mode().IsDir() {
			continue
		}

		rawFarmConfig, e := ioutil.ReadFile(farmsDir + descriptor.Name() + Sep + FarmConfigName)
		if e != nil {
			continue
		}

		farmConfig := FarmConfig{}
		json.Unmarshal(rawFarmConfig, &farmConfig)
		farmConfigs[descriptor.Name()] = farmConfig
	}
}

func checkDirs() {
	if e := os.MkdirAll(configDir, DirsMask); e != nil {
		panic(e)
	}

	if e := os.MkdirAll(farmsDir, DirsMask); e != nil {
		panic(e)
	}
}
