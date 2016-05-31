// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package apns

import (
	"log"

	"github.com/go-ini/ini"
	"github.com/codegangsta/cli"

	"zooinit/cluster"
	"zooinit/config"
	"gopush/lib"
)

// This basic discovery service bootstrap env info
type EnvInfo struct {
	cluster.BaseInfo

	CertPath     string
	CertPassword string

	PoolConfig   *lib.PoolConfig

	QueueSourceConfig *lib.QueueSourceConfig
}

func NewEnvInfo(iniobj *ini.File, c *cli.Context) *EnvInfo {
	env := new(EnvInfo)

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

	qsConfig:=&lib.QueueSourceConfig{}
	keyNow = "queue.method"
	tmpStr := config.GetValueString(keyNow, sec, c)
	if tmpStr == "" {
		log.Fatalln("Config of " + keyNow + " is empty.")
	}
	if tmpStr!=lib.QUEUE_SOURCE_METHOD_API && tmpStr!=lib.QUEUE_SOURCE_METHOD_FILE && tmpStr!=lib.QUEUE_SOURCE_METHOD_MYSQL {
		log.Fatalln("Config of " + keyNow + " value is not allowed: "+tmpStr)
	}
	qsConfig.Method=tmpStr

	if qsConfig.Method == lib.QUEUE_SOURCE_METHOD_API {
		keyNow = "queue.api.uri"
		tmpStr = config.GetValueString(keyNow, sec, c)
		if tmpStr == "" {
			log.Fatalln("Config of " + keyNow + " is empty.")
		}
		qsConfig.ApiPrefix=tmpStr

		//can be empty
		keyNow = "queue.api.name"
		tmpStr = config.GetValueString(keyNow, sec, c)
		qsConfig.Value=tmpStr
	}else if qsConfig.Method == lib.QUEUE_SOURCE_METHOD_MYSQL {
		keyNow = "queue.mysql.dsn"
		tmpStr = config.GetValueString(keyNow, sec, c)
		if tmpStr == "" {
			log.Fatalln("Config of " + keyNow + " is empty.")
		}
		qsConfig.MysqlDsn=tmpStr

		//can be empty
		keyNow = "queue.mysql.sql"
		tmpStr = config.GetValueString(keyNow, sec, c)
		qsConfig.Value=tmpStr
	}
	//set qsconfig
	env.QueueSourceConfig=qsConfig

	//create uuid
	env.CreateUUID()

	env.GuaranteeSingleRun()

	//register signal watcher
	env.RegisterSignalWatch()

	return env
}

func (p *EnvInfo) CreateWorker() (lib.Worker, error) {
	worker, err := NewWorker(p)
	if err != nil {
		return nil, err
	}

	return worker, nil
}

// TODO destroy
func (p *EnvInfo) DestroyWorker(worker lib.Worker) (error) {
	return nil
}

func (p *EnvInfo) GetPoolConfig() (*lib.PoolConfig) {
	return p.PoolConfig
}

func (p *EnvInfo) GetQueueSourceConfig() (*lib.QueueSourceConfig){
	return p.QueueSourceConfig
}