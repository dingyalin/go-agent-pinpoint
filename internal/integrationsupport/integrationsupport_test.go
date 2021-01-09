// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package integrationsupport

import (
	"sync"
	"testing"

	"github.com/dingyalin/pinpoint-go-agent/internal"
	pinpoint "github.com/dingyalin/pinpoint-go-agent/pinpoint"
)

func TestNilTransaction(t *testing.T) {
	var txn *pinpoint.Transaction

	AddAgentAttribute(txn, pinpoint.AttributeHostDisplayName, "hostname", nil)
	AddAgentSpanAttribute(txn, pinpoint.SpanAttributeAWSOperation, "operation")
}

func TestEmptyTransaction(t *testing.T) {
	txn := &pinpoint.Transaction{}

	AddAgentAttribute(txn, pinpoint.AttributeHostDisplayName, "hostname", nil)
	AddAgentSpanAttribute(txn, pinpoint.SpanAttributeAWSOperation, "operation")
}

func testApp(t *testing.T) *pinpoint.Application {
	app, err := pinpoint.NewApplication(
		pinpoint.ConfigAppName("appname"),
		pinpoint.ConfigLicense("0123456789012345678901234567890123456789"),
		pinpoint.ConfigEnabled(false),
		pinpoint.ConfigDistributedTracerEnabled(true),
	)
	if nil != err {
		t.Fatal(err)
	}
	replyfn := func(reply *internal.ConnectReply) {
		reply.SetSampleEverything()
	}
	internal.HarvestTesting(app.Private, replyfn)
	return app
}

func TestSuccess(t *testing.T) {
	app := testApp(t)
	txn := app.StartTransaction("hello")
	AddAgentAttribute(txn, pinpoint.AttributeHostDisplayName, "hostname", nil)
	segment := txn.StartSegment("mySegment")
	AddAgentSpanAttribute(txn, pinpoint.SpanAttributeAWSOperation, "operation")
	segment.End()
	txn.End()

	app.Private.(internal.Expect).ExpectTxnEvents(t, []internal.WantEvent{
		{
			AgentAttributes: map[string]interface{}{
				pinpoint.AttributeHostDisplayName: "hostname",
			},
		},
	})
	app.Private.(internal.Expect).ExpectSpanEvents(t, []internal.WantEvent{
		{
			Intrinsics: map[string]interface{}{
				"name":     "Custom/mySegment",
				"parentId": internal.MatchAnything,
				"category": "generic",
			},
			AgentAttributes: map[string]interface{}{
				pinpoint.SpanAttributeAWSOperation: "operation",
			},
		},
		{
			Intrinsics: map[string]interface{}{
				"name":             "OtherTransaction/Go/hello",
				"transaction.name": "OtherTransaction/Go/hello",
				"category":         "generic",
				"nr.entryPoint":    true,
			},
			AgentAttributes: map[string]interface{}{
				"host.displayName": "hostname",
			},
		},
	})
}

func TestConcurrentCalls(t *testing.T) {
	// This test will fail with a data race if the txn is not properly locked
	app := testApp(t)
	txn := app.StartTransaction("hello")
	defer txn.End()
	defer txn.StartSegment("mySegment").End()

	var wg sync.WaitGroup
	addAttr := func() {
		AddAgentSpanAttribute(txn, pinpoint.SpanAttributeAWSOperation, "operation")
		wg.Done()
	}

	wg.Add(1)
	go addAttr()
	wg.Wait()
}
