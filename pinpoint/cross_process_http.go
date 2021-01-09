// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package pinpoint

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/dingyalin/pinpoint-go-agent/internal/cat"
)

// InboundHTTPRequest adds the inbound request metadata to the txnCrossProcess.
func (txp *txnCrossProcess) InboundHTTPRequest(hdr http.Header) error {
	// return txp.handleInboundRequestHeaders(httpHeaderToMetadata(hdr))
	txp.InboundMetadata = httpHeaderToMetadata(hdr)
	return nil
}

// appDataToHTTPHeader encapsulates the given appData value in the correct HTTP
// header.
func appDataToHTTPHeader(appData string) http.Header {
	header := http.Header{}

	if appData != "" {
		header.Add(cat.NewRelicAppDataName, appData)
	}

	return header
}

// httpHeaderToAppData gets the appData value from the correct HTTP header.
func httpHeaderToAppData(header http.Header) string {
	if header == nil {
		return ""
	}

	return header.Get(cat.NewRelicAppDataName)
}

// httpHeaderToMetadata gets the cross process metadata from the relevant HTTP
// headers.
func httpHeaderToMetadata(header http.Header) (metadata crossProcessMetadata) {
	metadata.PinpointPspanid = -1

	if header == nil {
		return
	}

	// pinpointTraceid
	pinpointTraceid := header.Get(cat.PinpointTraceidName)
	traceIDList := strings.Split(pinpointTraceid, "^")
	if len(traceIDList) != 3 {
		return
	}
	agentID, startTimeStr, sequenceIDStr := traceIDList[0], traceIDList[1], traceIDList[2]
	startTime, err := strconv.ParseInt(startTimeStr, 10, 64)
	if err != nil {
		return
	}
	sequenceID, err := strconv.ParseInt(sequenceIDStr, 10, 64)
	if err != nil {
		return
	}
	pinpointTraceidEncoded := encodeTraceID(agentID, startTime, sequenceID)

	// pinpointPapptype
	pAppTypeValueTMP, err := strconv.Atoi(header.Get(cat.PinpointPapptypeName))
	if err != nil {
		return
	}
	pinpointPapptype := int16(pAppTypeValueTMP)

	// pinpointPspanid
	pinpointPspanid, err := strconv.ParseInt(header.Get(cat.PinpointPspanidName), 10, 64)
	if err != nil {
		return
	}

	// pinpointSpanid
	pinpointSpanid, err := strconv.ParseInt(header.Get(cat.PinpointSpanidName), 10, 64)
	if err != nil {
		return
	}

	return crossProcessMetadata{
		//ID:         header.Get(cat.NewRelicIDName),
		//TxnData:    header.Get(cat.NewRelicTxnName),
		//Synthetics: header.Get(cat.NewRelicSyntheticsName),
		PinpointTraceid:        pinpointTraceid,
		PinpointPappname:       header.Get(cat.PinpointPappnameName),
		PinpointPapptype:       pinpointPapptype,
		PinpointPspanid:        pinpointPspanid,
		PinpointSpanid:         pinpointSpanid,
		PinpointFlags:          header.Get(cat.PinpointFlagsName),
		PinpointTraceidEncoded: pinpointTraceidEncoded,
	}
}

// metadataToHTTPHeader creates a set of HTTP headers to represent the given
// cross process metadata.
func metadataToHTTPHeader(metadata crossProcessMetadata) http.Header {
	header := http.Header{}

	/*
		if metadata.ID != "" {
			header.Add(cat.NewRelicIDName, metadata.ID)
		}

		if metadata.TxnData != "" {
			header.Add(cat.NewRelicTxnName, metadata.TxnData)
		}

		if metadata.Synthetics != "" {
			header.Add(cat.NewRelicSyntheticsName, metadata.Synthetics)
		}
	*/

	// pinpoint

	return header
}
