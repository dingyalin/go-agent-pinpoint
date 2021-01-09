// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dingyalin/pinpoint-go-agent/internal/utilization"
	pinpoint "github.com/dingyalin/pinpoint-go-agent/pinpoint"
)

func main() {
	util := utilization.Gather(utilization.Config{
		DetectAWS:        true,
		DetectAzure:      true,
		DetectDocker:     true,
		DetectPCF:        true,
		DetectGCP:        true,
		DetectKubernetes: true,
	}, pinpoint.NewDebugLogger(os.Stdout))

	js, err := json.MarshalIndent(util, "", "\t")
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		fmt.Printf("%s\n", js)
	}
}
