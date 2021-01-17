// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dingyalin/pinpoint-go-agent/integrations/nrlogrus"
	pinpoint "github.com/dingyalin/pinpoint-go-agent/pinpoint"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	app, err := pinpoint.NewApplication(
		pinpoint.ConfigAppName("Logrus App"),
		pinpoint.ConfigLicense(os.Getenv("PINPOINT_LICENSE_KEY")),
		nrlogrus.ConfigStandardLogger(),
	)

	if nil != err {
		fmt.Println(err)
		os.Exit(1)
	}

	http.HandleFunc(pinpoint.WrapHandleFunc(app, "/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello world")
	}))

	http.ListenAndServe(":8000", nil)
}
