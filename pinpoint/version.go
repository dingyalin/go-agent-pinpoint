// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package pinpoint

import (
	"runtime"

	"github.com/dingyalin/pinpoint-go-agent/internal"
)

const (
	// Version is the full string version of this Go Agent. [3.9.0]
	Version = "1.0.0"
)

var (
	goVersionSimple = minorVersion(runtime.Version())
)

func init() {
	internal.TrackUsage("Go", "Version", Version)
	internal.TrackUsage("Go", "Runtime", "Version", goVersionSimple)
	internal.TrackUsage("Go", "gRPC", "Version", grpcVersion)
}
