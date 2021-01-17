// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dingyalin/pinpoint-go-agent/pinpoint"
)

func main() {
	app, err := pinpoint.NewApplication(
		pinpoint.ConfigAppName("Short Lived App"),
		pinpoint.ConfigLicense(os.Getenv("PINPOINT_LICENSE_KEY")),
		pinpoint.ConfigDebugLogger(os.Stdout),
	)
	if nil != err {
		fmt.Println(err)
		os.Exit(1)
	}

	// Wait for the application to connect.
	if err := app.WaitForConnection(5 * time.Second); nil != err {
		fmt.Println(err)
	}

	// Do the tasks at hand.  Perhaps record them using transactions and/or
	// custom events.
	tasks := []string{"white", "black", "red", "blue", "green", "yellow"}
	for _, task := range tasks {
		txn := app.StartTransaction("task")
		time.Sleep(10 * time.Millisecond)
		txn.End()
		app.RecordCustomEvent("task", map[string]interface{}{
			"color": task,
		})
	}

	// Shut down the application to flush data to New Relic.
	app.Shutdown(10 * time.Second)
}
