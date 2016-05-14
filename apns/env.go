// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package apns

import (
	"log"

	"github.com/go-ini/ini"
	"github.com/codegangsta/cli"

	"zooinit/cluster"
	"zooinit/config"
)

// This basic discovery service bootstrap env info
type envInfo struct {
	cluster.BaseInfo

	CertPath     string
	CertPassword string
}

func NewEnvInfo(iniobj *ini.File, c *cli.Context) *envInfo {
	env := new(envInfo)

	sec := iniobj.Section(CONFIG_SECTION)
	env.Service = sec.Key("service").String()
	if env.Service == "" {
		log.Fatalln("Config of service section is empty.")
	}

	// parse base info
	env.ParseConfigFile(sec, c)

	keyNow := "cert.path"
	env.CertPath = config.GetValueString(keyNow, sec, c)
	if env.CertPath == "" {
		log.Fatalln("Config of " + keyNow + " is empty.")
	}

	keyNow = "cert.password"
	env.CertPassword = config.GetValueString(keyNow, sec, c)
	if env.CertPassword == "" {
		log.Fatalln("Config of " + keyNow + " is empty.")
	}

	//create uuid
	env.CreateUUID()

	env.GuaranteeSingleRun()

	//register signal watcher
	env.RegisterSignalWatch()

	return env
}