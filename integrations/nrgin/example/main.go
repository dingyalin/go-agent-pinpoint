// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	nrgin "github.com/dingyalin/pinpoint-go-agent/integrations/nrgin"
	"github.com/dingyalin/pinpoint-go-agent/pinpoint"
	"github.com/gin-gonic/gin"
)

func makeGinEndpoint(s string) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Writer.WriteString(s)
	}
}

func v1login(c *gin.Context)  { c.Writer.WriteString("v1 login") }
func v1submit(c *gin.Context) { c.Writer.WriteString("v1 submit") }
func v1read(c *gin.Context)   { c.Writer.WriteString("v1 read") }

func endpoint404(c *gin.Context) {
	c.Writer.WriteHeader(404)
	c.Writer.WriteString("returning 404")
}

func endpoint500(c *gin.Context) {
	a := 1
	b := 0
	_ = a / b
	// panic("sdfjjhdsfjfdsjh")
}

func endpointChangeCode(c *gin.Context) {
	// gin.ResponseWriter buffers the response code so that it can be
	// changed before the first write.
	c.Writer.WriteHeader(404)
	c.Writer.WriteHeader(200)
	c.Writer.WriteString("actually ok!")
}

func endpointResponseHeaders(c *gin.Context) {
	// Since gin.ResponseWriter buffers the response code, response headers
	// can be set afterwards.
	c.Writer.WriteHeader(200)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteString(`{"zip":"zap"}`)
}

func endpointNotFound(c *gin.Context) {
	c.Writer.WriteString("there's no endpoint for that!")
}

func endpointAccessTransaction(c *gin.Context) {
	txn := nrgin.Transaction(c)
	txn.SetName("custom-name")
	c.Writer.WriteString("changed the name of the transaction!")
}

func main() {
	app, err := pinpoint.NewApplication(
		pinpoint.ConfigAppName("GoGinDemo"),
		pinpoint.ConfigAgentID("GoGinDemo"),
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
		fmt.Println(err)
		os.Exit(1)
	}

	router := gin.Default()
	router.Use(nrgin.Middleware(app))

	router.GET("/404", endpoint404)
	router.GET("/500", endpoint500)
	router.GET("/change", endpointChangeCode)
	router.GET("/headers", endpointResponseHeaders)
	router.GET("/txn", endpointAccessTransaction)

	// Since the handler function name is used as the transaction name,
	// anonymous functions do not get usefully named.  We encourage
	// transforming anonymous functions into named functions.
	router.GET("/anon", func(c *gin.Context) {
		c.Writer.WriteString("anonymous function handler")
	})

	v1 := router.Group("/v1")
	v1.GET("/login", v1login)
	v1.GET("/submit", v1submit)
	v1.GET("/read", v1read)

	router.NoRoute(endpointNotFound)

	router.Run(":8000")
}
