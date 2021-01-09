// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dingyalin/pinpoint-go-agent/integrations/nrmicro"
	proto "github.com/dingyalin/pinpoint-go-agent/integrations/nrmicro/example/proto"
	pinpoint "github.com/dingyalin/pinpoint-go-agent/pinpoint"
	"github.com/micro/go-micro"
)

func subEv(ctx context.Context, msg *proto.HelloRequest) error {
	fmt.Println("Message received from", msg.GetName())
	return nil
}

func publish(s micro.Service, app *pinpoint.Application) {
	c := s.Client()

	for range time.NewTicker(time.Second).C {
		txn := app.StartTransaction("publish")
		msg := c.NewMessage("example.topic.pubsub", &proto.HelloRequest{Name: "Sally"})
		ctx := pinpoint.NewContext(context.Background(), txn)
		fmt.Println("Sending message")
		if err := c.Publish(ctx, msg); nil != err {
			log.Fatal(err)
		}
		txn.End()
	}
}

func main() {
	app, err := pinpoint.NewApplication(
		pinpoint.ConfigAppName("GoMicroPubSubDemo"),
		pinpoint.ConfigAgentID("GoMicroPubSubDemo"),
		//pinpoint.ConfigEnabled(false),
		pinpoint.ConfigCollectorUploaded(false),
		//pinpoint.ConfigCollectorUploadedAgentStat(false),
		pinpoint.ConfigCollectorIP("127.0.0.1"),
		pinpoint.ConfigCollectorTCPPort(9994),
		pinpoint.ConfigCollectorStatPort(9995),
		pinpoint.ConfigCollectorSpanPort(9996),
		pinpoint.ConfigDebugLogger(os.Stdout),
	)
	if nil != err {
		panic(err)
	}
	err = app.WaitForConnection(10 * time.Second)
	if nil != err {
		panic(err)
	}
	defer app.Shutdown(10 * time.Second)

	s := micro.NewService(
		micro.Name("go.micro.srv.pubsub"),
		// Add the New Relic wrapper to the client which will create
		// MessageProducerSegments for each Publish call.
		micro.WrapClient(nrmicro.ClientWrapper()),
		// Add the New Relic wrapper to the subscriber which will start a new
		// transaction for each Subscriber invocation.
		micro.WrapSubscriber(nrmicro.SubscriberWrapper(app)),
	)
	s.Init()

	go publish(s, app)

	micro.RegisterSubscriber("example.topic.pubsub", s.Server(), subEv)

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
