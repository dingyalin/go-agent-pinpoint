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

// Greeter is the server struct
type Greeter struct{}

// Hello is the method on the server being called
func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	name := req.GetName()
	txn := pinpoint.FromContext(ctx)
	txn.AddAttribute("Name", name)
	fmt.Println("Request received from", name)
	rsp.Greeting = "Hello " + name
	return nil
}

func main() {
	app, err := pinpoint.NewApplication(
		pinpoint.ConfigAppName("GoMicroServerDemo"),
		pinpoint.ConfigAgentID("GoMicroServerDemo"),
		//pinpoint.ConfigEnabled(false),
		//pinpoint.ConfigCollectorUploadedAgentStat(false),
		//pinpoint.ConfigCollectorUploaded(false),
		pinpoint.ConfigCollectorIP("127.0.0.1"),
		pinpoint.ConfigCollectorTCPPort(9994),
		pinpoint.ConfigCollectorStatPort(9995),
		pinpoint.ConfigCollectorSpanPort(9996),
		pinpoint.ConfigDebugLogger(os.Stdout),
	)
	if nil != err {
		panic(err)
	}

	/*
		err = app.WaitForConnection(10 * time.Second)
		if nil != err {
			panic(err)
		}
	*/
	defer app.Shutdown(10 * time.Second)

	service := micro.NewService(
		micro.Name("greeter"),
		// Add the New Relic middleware which will start a new transaction for
		// each Handler invocation.
		micro.WrapHandler(nrmicro.HandlerWrapper(app)),
	)

	service.Init()

	proto.RegisterGreeterHandler(service.Server(), new(Greeter))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
