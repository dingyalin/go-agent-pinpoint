// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	pinpoint "github.com/dingyalin/pinpoint-go-agent/pinpoint"
)

func index(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello world")
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Pinpoint Go Agent Version: "+pinpoint.Version)
}

func noticeError(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "noticing an error")

	txn := pinpoint.FromContext(r.Context())
	txn.NoticeError(errors.New("my error message"))
}

func setName(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "changing the transaction's name")

	txn := pinpoint.FromContext(r.Context())
	txn.SetName("other-name")
}

func ignore(w http.ResponseWriter, r *http.Request) {
	if coinFlip := (0 == rand.Intn(2)); coinFlip {
		txn := pinpoint.FromContext(r.Context())
		txn.Ignore()
		io.WriteString(w, "ignoring the transaction")
	} else {
		io.WriteString(w, "not ignoring the transaction")
	}
}

func segments(w http.ResponseWriter, r *http.Request) {
	txn := pinpoint.FromContext(r.Context())

	func() {
		defer txn.StartSegment("f1").End()

		func() {
			defer txn.StartSegment("f2").End()

			io.WriteString(w, "segments!")
			time.Sleep(50 * time.Millisecond)
		}()

		func() {
			defer txn.StartSegment("f3").End()

			//io.WriteString(w, "segments!")
			time.Sleep(100 * time.Millisecond)
		}()

		time.Sleep(15 * time.Millisecond)
	}()
	time.Sleep(20 * time.Millisecond)
}

func mysql(w http.ResponseWriter, r *http.Request) {
	txn := pinpoint.FromContext(r.Context())
	s := pinpoint.DatastoreSegment{
		StartTime: txn.StartSegmentNow(),
		// Product, Collection, and Operation are the most important
		// fields to populate because they are used in the breakdown
		// metrics.
		Product:    pinpoint.DatastoreMySQL,
		Collection: "users",
		Operation:  "INSERT",

		ParameterizedQuery: "INSERT INTO users (name, age) VALUES ($1, $2)",
		QueryParameters: map[string]interface{}{
			"name": "Dracula",
			"age":  439,
		},
		Host:         "mysql-server-1",
		PortPathOrID: "3306",
		DatabaseName: "my_database",
	}
	defer s.End()

	time.Sleep(20 * time.Millisecond)
	io.WriteString(w, `performing fake query "INSERT * from users"`)
}

func message(w http.ResponseWriter, r *http.Request) {
	txn := pinpoint.FromContext(r.Context())
	s := pinpoint.MessageProducerSegment{
		StartTime:       txn.StartSegmentNow(),
		Library:         "RabbitMQ",
		DestinationType: pinpoint.MessageQueue,
		DestinationName: "myQueue",
	}
	defer s.End()

	time.Sleep(20 * time.Millisecond)
	io.WriteString(w, `producing a message queue message`)
}

func external(w http.ResponseWriter, r *http.Request) {
	txn := pinpoint.FromContext(r.Context())
	req, _ := http.NewRequest("GET", "http://example.com?aa=bb", nil)

	// Using StartExternalSegment is recommended because it does distributed
	// tracing header setup, but if you don't have an *http.Request and
	// instead only have a url string then you can start the external
	// segment like this:
	//
	// es := pinpoint.ExternalSegment{
	// 	StartTime: txn.StartSegmentNow(),
	// 	URL:       urlString,
	// }
	//
	es := pinpoint.StartExternalSegment(txn, req)
	resp, err := http.DefaultClient.Do(req)
	es.End()

	if nil != err {
		io.WriteString(w, err.Error())
		return
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

func roundtripper(w http.ResponseWriter, r *http.Request) {
	// NewRoundTripper allows you to instrument external calls without
	// calling StartExternalSegment by modifying the http.Client's Transport
	// field.  If the Transaction parameter is nil, the RoundTripper
	// returned will look for a Transaction in the request's context (using
	// FromContext). This is recommended because it allows you to reuse the
	// same client for multiple transactions.
	client := &http.Client{}
	client.Transport = pinpoint.NewRoundTripper(client.Transport)

	request, _ := http.NewRequest("GET", "http://example.com", nil)
	// Since the transaction is already added to the inbound request's
	// context by WrapHandleFunc, we just need to copy the context from the
	// inbound request to the external request.
	request = request.WithContext(r.Context())
	// Alternatively, if you don't want to copy entire context, and instead
	// wanted just to add the transaction to the external request's context,
	// you could do that like this:
	//
	//	txn := pinpoint.FromContext(r.Context())
	//	request = pinpoint.RequestWithTransactionContext(request, txn)

	resp, err := client.Do(request)
	if nil != err {
		io.WriteString(w, err.Error())
		return
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

func async(w http.ResponseWriter, r *http.Request) {
	txn := pinpoint.FromContext(r.Context())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(txn *pinpoint.Transaction) {
		defer txn.End()
		defer wg.Done()
		defer txn.StartSegment("async-block").End()
		time.Sleep(100 * time.Millisecond)
	}(txn.NewGoroutine("go1"))

	go func(txn *pinpoint.Transaction) {
		defer txn.End()
		defer txn.StartSegment("async-nonblock").End()
		time.Sleep(500 * time.Millisecond)
	}(txn.NewGoroutine("go2"))

	segment := txn.StartSegment("wg.Wait")
	wg.Wait()

	time.Sleep(150 * time.Millisecond)
	segment.End()
	w.Write([]byte("done!"))
}

func main() {
	app, err := pinpoint.NewApplication(
		pinpoint.ConfigCollectorUploaded(false),
		pinpoint.ConfigFromYaml("./pinpoint.yml"),
		pinpoint.ConfigDebugLogger(os.Stdout),
		pinpoint.ConfigFromEnvironment(),
	)
	if nil != err {
		fmt.Println(err)
		os.Exit(1)
	}

	http.HandleFunc(pinpoint.WrapHandleFunc(app, "/", index))
	http.HandleFunc(pinpoint.WrapHandleFunc(app, "/version", versionHandler))
	http.HandleFunc(pinpoint.WrapHandleFunc(app, "/notice_error", noticeError))
	http.HandleFunc(pinpoint.WrapHandleFunc(app, "/set_name", setName))
	http.HandleFunc(pinpoint.WrapHandleFunc(app, "/ignore", ignore))
	http.HandleFunc(pinpoint.WrapHandleFunc(app, "/segments", segments))
	http.HandleFunc(pinpoint.WrapHandleFunc(app, "/mysql", mysql))
	http.HandleFunc(pinpoint.WrapHandleFunc(app, "/external", external))
	http.HandleFunc(pinpoint.WrapHandleFunc(app, "/roundtripper", roundtripper))
	http.HandleFunc(pinpoint.WrapHandleFunc(app, "/async", async))
	http.HandleFunc(pinpoint.WrapHandleFunc(app, "/message", message))

	http.HandleFunc("/background", func(w http.ResponseWriter, req *http.Request) {
		// Transactions started without an http.Request are classified as
		// background transactions.
		txn := app.StartTransaction("background")
		defer txn.End()

		io.WriteString(w, "background transaction")
		time.Sleep(150 * time.Millisecond)
	})

	http.ListenAndServe(":8000", nil)
}
