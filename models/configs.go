package models

import (
	"errors"
	// "fmt"
	"io/ioutil"
	"log"
	// "reflect"
	// "regexp"
	// "sort"
	// "strconv"
	// "strings"
	// "time"

	"github.com/BurntSushi/toml"
)

var (
	configErrors = []string{}
)

type Config struct {
	VCenterURL		string
	VCenterUsername	string
	VCenterPassword	string
	Datacenter		string
	TemplateResourcePool	string
	StartingPortGroup	int
	EndingPortGroup		int
	Https           bool
	Port            int
	Cert            string
	Key             string
	TargetResourcePool		string
	Domain			string
	WanPortGroup	string
	JwtSecret		string
}

/*
	Load config settings into given config object
*/
func ReadConfig(conf *Config, configPath string) {
	fileContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalln("Configuration file ("+configPath+") not found:", err)
	}
	if md, err := toml.Decode(string(fileContent), &conf); err != nil {
		log.Fatalln(err)
	} else {
		for _, undecoded := range md.Undecoded() {
			errMsg := "[WARN] Undecoded scoring configuration key \"" + undecoded.String() + "\" will not be used."
			configErrors = append(configErrors, errMsg)
			log.Println(errMsg)
		}
	}
}

/*
	Check for config errors and set defaults
*/
func CheckConfig(conf *Config) error {

	if conf.VCenterURL == "" {
		return errors.New("illegal config: vCenterURL must be defined")
	}
	if conf.VCenterUsername == "" {
		return errors.New("illegal config: vCenterUsername must be defined")
	}
	if conf.VCenterPassword == "" {
		return errors.New("illegal config: vCenterPassword must be defined")
	}
	if conf.Datacenter == "" {
		return errors.New("illegal config: Datacenter must be defined")
	}
	if conf.TemplateResourcePool == "" {
		return errors.New("illegal config: TemplateResourcePool must be defined")
	}

	if conf.StartingPortGroup == 0 || conf.EndingPortGroup == 0 {
		return errors.New("illegal config: StartingPortGroup and EndingPortGroup must be defined")
	}
	if conf.Port == 0 {
		if conf.Https {
			conf.Port = 443
		} else {
			conf.Port = 80
		}
	}

	if conf.Https {
		if conf.Cert == "" || conf.Key == "" {
			return errors.New("illegal config: https requires a cert and key pair")
		}
	}

	if conf.JwtSecret == "" {
		return errors.New("illegal config: JwtSecret must be defined")
	}

	return nil
}