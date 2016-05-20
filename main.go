// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package main

import (
	"os"
	"gopush/apns"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Author = "bruce"
	app.Email = "bruce@haimi.com"
	app.Copyright = "haimi.com All rights reseverd."
	app.Name = "gopush"
	app.Usage = "Go push service for app."
	app.Version = "0.0.1"

	cfgFlag := &cli.StringFlag{
		Name:  "config, f",
		Value: "runtime/config/config.ini",
		Usage: "Configuration file of gopush.",
	}

	queueFlag := &cli.StringFlag{
		Name:  "queue",
		Value: "runtime/data/queue.txt",
		Usage: "Default text queue file of push.",
	}

	logChannel := &cli.StringFlag{
		Name:  "log.channel",
		Value: "",
		Usage: "Configuration of runtime log channel: file, write to file; stdout, write to stdout; multi, write both.",
	}

	logPath := &cli.StringFlag{
		Name:  "log.path, log",
		Value: "",
		Usage: "Configuration of runtime log path.",
	}

	pidPath := &cli.StringFlag{
		Name:  "pid.path, pid",
		Value: "",
		Usage: "Configuration of runtime log path.",
	}

	size := &cli.IntFlag{
		Name:  "size",
		Value: 20,
		Usage: "Worker pool size, init size.",
	}

	capacity := &cli.IntFlag{
		Name:  "capacity",
		Value: 200,
		Usage: "Worker pool capacity.",
	}

	miniSpare := &cli.IntFlag{
		Name:  "spare.mini",
		Value: 10,
		Usage: "Worker pool miniSpare worker.",
	}

	maxSpare := &cli.IntFlag{
		Name:  "spare.max",
		Value: 50,
		Usage: "Worker pool maxSpare worker.",
	}

	app.Commands = []cli.Command{
		{
			Name:    "apns",
			Usage:   "Usage: " + os.Args[0] + " apns -f config.ini \nSend apns push notificatioins.",
			Action:  apns.Bootstrap,
			Flags: []cli.Flag{
				cfgFlag, logChannel, logPath, pidPath, size, capacity, miniSpare, maxSpare, queueFlag,
			},
		},
	}
	app.Run(os.Args)
}
