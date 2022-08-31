package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	VpnSettings       VpnSettings       `yaml:",inline"`
	GlobalNetSettings GlobalNetSettings `yaml:",inline"`
	Security          Security          `yaml:",inline"`
	Flash             Flash             `yaml:",inline"`
	Hardware          Hardware          `yaml:",inline"`
}

type VpnSettings struct {
	PrivateNetwork      string `yaml:"privateNetwork"`
	PingTimesForVpn     int    `yaml:"pingTimesForVpn"`
	RebootOnPrivateFail bool   `yaml:"rebootWhilePrivateFail"`
	VpnRebootTimeout    int    `yaml:"vpnRebootTimeout"`
}

type GlobalNetSettings struct {
	GlobalNetwork       string `yaml:"globalNetwork"`
	GlobalNetWorkPort   int    `yaml:"globalNetworkPort"`
	GlobalRebootTimeout int    `yaml:"globalRebootTimeout"`
	RebootIfFail        bool   `yaml:"rebootWhileGlobalFail"`
}
type Security struct {
	EnableTokenConfigParse bool   `yaml:"enableTokenConfigParse"`
	TokenBotApi            string `yaml:"tokenBotApi"`
}
type Flash struct {
	PathToDev      string `yaml:"pathToDevice"`
	MountPointPath string `yaml:"pathToMountPoint"`
}

type Hardware struct {
	LedIndication bool `yaml:"ledIndication"`
}

func (f *Config) ParseConfig(configFilePath string) *Config {
	yamlFile, yamlParseErr := os.ReadFile(configFilePath)
	if yamlParseErr != nil {
		log.Panicf("an error occured while parsing configFile: %v", yamlParseErr)
	}
	unmarshErr := yaml.Unmarshal(yamlFile, f)
	if unmarshErr != nil {
		log.Panicf("an error occured while unmarshaling configFille: %v", unmarshErr)
	}
	return f
}
