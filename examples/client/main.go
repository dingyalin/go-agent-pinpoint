// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	pinpoint "github.com/dingyalin/pinpoint-go-agent/pinpoint"
)

func doRequest(txn *pinpoint.Transaction) error {
	req, err := http.NewRequest("GET", "http://localhost:8000/segments", nil)
	if nil != err {
		return err
	}
	client := &http.Client{}
	seg := pinpoint.StartExternalSegment(txn, req)
	defer seg.End()
	resp, err := client.Do(req)
	if nil != err {
		return err
	}
	fmt.Println("response code is", resp.StatusCode)
	return nil
}

func main() {
	app, err := pinpoint.NewApplication(
		pinpoint.ConfigAppName("Client App"),
		pinpoint.ConfigLicense(os.Getenv("PINPOINT_LICENSE_KEY")),
		pinpoint.ConfigDebugLogger(os.Stdout),
		pinpoint.ConfigDistributedTracerEnabled(true),
	)
	if nil != err {
		fmt.Println(err)
		os.Exit(1)
	}

	// Wait for the application to connect.
	if err = app.WaitForConnection(5 * time.Second); nil != err {
		fmt.Println(err)
	}

	txn := app.StartTransaction("client-txn")
	err = doRequest(txn)
	if nil != err {
		txn.NoticeError(err)
	}
	txn.End()

	// Shut down the application to flush data to New Relic.
	app.Shutdown(10 * time.Second)
}
