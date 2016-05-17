// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package apns

import (
	"zooinit/config"

	"github.com/codegangsta/cli"
	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
	"github.com/twinj/uuid"
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

	pool, err := NewPool(Size, Capacity, MiniSpare, MaxSpare)
	if err != nil {
		env.GetLogger().Fatalln("Found error while create push pool:", err)
	}

	notification := &apns.Notification{}
	//bruce
	notification.DeviceToken = "3523544012e5491b3fe8cf6627eddd123d6aa4191fbebf371191a3ce7d4c02ac"
	//jj
	//notification.DeviceToken ="efdd029e3e62ab46bf089bfe7084d3261471b6f9e0e4225f9851b4e5b8e7f57e"
	notification.ApnsID = uuid.NewV1().String()
	notification.Priority = 10
	notification.Topic = ""
	load := payload.NewPayload()

	load.Badge(1)
	load.AlertTitle("appname")
	load.AlertBody("push message")
	//Done push Turn to specific page machanism, addon field
	load.Custom("payload", "haimi-590")

	load.Sound("bingbong.aiff")
	notification.Payload = load

	queue := NewQueueByPool(pool)
	queue.AppendFileDataSource(c.String("queue"))
	env.GetLogger().Println("Using queue file data source: " + c.String("queue"))

	//pool run entry
	pool.Run(queue, notification)
}
