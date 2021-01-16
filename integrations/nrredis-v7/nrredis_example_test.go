// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package nrredis_test

import (
	"context"
	"fmt"

	nrredis "github.com/dingyalin/pinpoint-go-agent/integrations/nrredis-v7"
	pinpoint "github.com/dingyalin/pinpoint-go-agent/pinpoint"
	redis "github.com/go-redis/redis/v7"
)

func getTransaction() *pinpoint.Transaction { return nil }

func Example_client() {
	opts := &redis.Options{Addr: "localhost:6379"}
	client := redis.NewClient(opts)

	//
	// Step 1:  Add a nrredis.NewHook() to your redis client.
	//
	client.AddHook(nrredis.NewHook(opts))

	//
	// Step 2: Ensure that all client calls contain a context with includes
	// the transaction.
	//
	txn := getTransaction()
	ctx := pinpoint.NewContext(context.Background(), txn)
	pong, err := client.WithContext(ctx).Ping().Result()
	fmt.Println(pong, err)
}

func Example_clusterClient() {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},
	})

	//
	// Step 1:  Add a nrredis.NewHook() to your redis cluster client.
	//
	client.AddHook(nrredis.NewHook(nil))

	//
	// Step 2: Ensure that all client calls contain a context with includes
	// the transaction.
	//
	txn := getTransaction()
	ctx := pinpoint.NewContext(context.Background(), txn)
	pong, err := client.WithContext(ctx).Ping().Result()
	fmt.Println(pong, err)
}
