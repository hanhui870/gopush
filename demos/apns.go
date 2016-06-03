// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package main

import (
	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
//"github.com/sideshow/apns2/payload"
	"log"
	"github.com/twinj/uuid"
	"fmt"
	"github.com/sideshow/apns2/payload"
)

func main() {
	cert, pemErr := certificate.FromP12File("test.p12", "pass")
	if pemErr != nil {
		log.Fatalln("Cert Error:", pemErr)
	}

	//normal one
	notification := &apns.Notification{}
	notification.DeviceToken = "3523544012e5491b3fe8cf6627eddd123d6aa4191fbebf371191a3ce7d4c02ac"
	notification.ApnsID = uuid.NewV4().String()
	notification.Priority = 10
	notification.Topic = ""
	load := payload.NewPayload()

	load.Badge(1)
	load.AlertTitle("appname")
	load.AlertBody("hello haimi.com test turnvalue2")
	load.Sound("bingbong.aiff")
	load.Custom("TurnType", "OTHER")
	turn := `{"type":"WAP","id":"143573","title":"http://www.haimi.com/","HaimiScheme":"haimi://home"}`
	load.Custom("TurnValue", turn)

	// 终于能跳转了
	//Done push Turn to specific page machanism
	load.Custom("payload", "haimi-541")
	notification.Payload = load // See Payload section below
	// If an encountered value implements the Marshaler interface
	// and is not a nil pointer, Marshal calls its MarshalJSON method
	// to produce JSON.
	jons, _ := notification.MarshalJSON()
	fmt.Println(string(jons))
	//return

	client := apns.NewClient(cert).Production()
	res, err := client.Push(notification)

	if err != nil {
		log.Println("Error:", err)
		return
	}

	log.Println("APNs ID:", res.ApnsID, "Genarated:", notification.ApnsID)
}
