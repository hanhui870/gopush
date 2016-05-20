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

	env.GetLogger().Println("Push Worker Size:", Size)
	env.GetLogger().Println("Push Worker Capacity:", Capacity)
	env.GetLogger().Println("Push Worker MiniSpare:", MiniSpare)
	env.GetLogger().Println("Push Worker MaxSpare:", MaxSpare)

	if Size <= 0 || Capacity <= 0 || MiniSpare <= 0 || MaxSpare <= 0 {
		env.GetLogger().Fatalln("All Size, Capacity, MiniSpare, MaxSpare parameters must all >0")
	}
	if Size < MiniSpare {
		Size = MiniSpare
		env.GetLogger().Println("Size<MiniSpare, will change to equal to MiniSpare")
	}
	if Size > Capacity {
		Size = Capacity
		env.GetLogger().Println("Size>Capacity, will change to equal to Capacity")
	}
	if MiniSpare > MaxSpare {
		env.GetLogger().Fatalln("MiniSpare must <= MaxSpare")
	}
	if Size > Capacity || MiniSpare > Capacity || MaxSpare >= Capacity {
		env.GetLogger().Fatalln("Capacity must be the greatest parameter within Size, Capacity, MiniSpare, MaxSpare")
	}

	pool, err := lib.NewPool(Size, Capacity, MiniSpare, MaxSpare, env)
	if err != nil {
		env.GetLogger().Fatalln("Found error while create push pool:", err)
	}

	// no need next
	server := api.NewApiV1Server(env)
	server.SetPool(pool)
	err = server.Start()

	if err != nil {
		env.GetLogger().Fatalln("Found error while server.Start():", err)
	}

	return
}
