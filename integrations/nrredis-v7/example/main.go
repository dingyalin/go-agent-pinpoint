// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	nrredis "github.com/dingyalin/pinpoint-go-agent/integrations/nrredis-v7"
	pinpoint "github.com/dingyalin/pinpoint-go-agent/pinpoint"
	redis "github.com/go-redis/redis/v7"
)

func main() {
	app, err := pinpoint.NewApplication(
		pinpoint.ConfigAppName("GoRedisDemo"),
		pinpoint.ConfigAgentID("GoRedisDemo"),
		//pinpoint.ConfigEnabled(false),
		pinpoint.ConfigCollectorUploaded(false),
		pinpoint.ConfigCollectorUploadedAgentStat(false),
		pinpoint.ConfigCollectorIP("127.0.0.1"),
		pinpoint.ConfigCollectorTCPPort(9994),
		pinpoint.ConfigCollectorStatPort(9995),
		pinpoint.ConfigCollectorSpanPort(9996),
		pinpoint.ConfigDebugLogger(os.Stdout),
	)
	if nil != err {
		panic(err)
	}
	app.WaitForConnection(10 * time.Second)
	txn := app.StartTransaction("ping txn")

	opts := &redis.Options{
		Addr: "localhost:6379",
	}
	client := redis.NewClient(opts)

	//
	// Step 1:  Add a nrredis.NewHook() to your redis client.
	//
	client.AddHook(nrredis.NewHook(opts))

	//
	// Step 2: Ensure that all client calls contain a context which includes
	// the transaction.
	//
	ctx := pinpoint.NewContext(context.Background(), txn)
	pipe := client.WithContext(ctx).Pipeline()
	incr := pipe.Incr("pipeline_counter")
	pipe.Expire("pipeline_counter", time.Hour)
	_, err = pipe.Exec()
	fmt.Println(incr.Val(), err)

	txn.End()
	app.Shutdown(5 * time.Second)
}
