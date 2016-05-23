// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package apns

import (
	"zooinit/config"

	"github.com/codegangsta/cli"

	"gopush/api"
	"gopush/lib"
)

const (
	CONFIG_SECTION = "system.apns"
)

var (
	env *EnvInfo
)

func Bootstrap(c *cli.Context) {
	fname := config.GetConfigFileName(c.String("config"))
	iniobj := config.GetConfigInstance(fname)

	env = NewEnvInfo(iniobj, c)

	//flush last log info
	defer env.Logger.Sync()

	var Size, Capacity, MiniSpare, MaxSpare int
	Size = c.Int("size")
	Capacity = c.Int("capacity")
	MiniSpare = c.Int("spare.mini")
	MaxSpare = c.Int("spare.max")

	poolCfg, err := lib.NewPoolConfig(Size, Capacity, MiniSpare, MaxSpare)
	if err != nil {
		env.GetLogger().Fatalln("Found error while create PoolConfig:", err)
	}

	env.GetLogger().Println("Push Worker Size:", poolCfg.Size)
	env.GetLogger().Println("Push Worker Capacity:", poolCfg.Capacity)
	env.GetLogger().Println("Push Worker MiniSpare:", poolCfg.MiniSpare)
	env.GetLogger().Println("Push Worker MaxSpare:", poolCfg.MaxSpare)

	env.PoolConfig = poolCfg

	// no need next
	server := api.NewApiV1Server(env)
	err = server.Start()

	if err != nil {
		env.GetLogger().Fatalln("Found error while server.Start():", err)
	}

	return
}
