// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package nrgrpc

import (
	"context"
	"net/http"
	"strings"

	pinpoint "github.com/dingyalin/pinpoint-go-agent/pinpoint"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func startTransaction(ctx context.Context, app *pinpoint.Application, fullMethod string) *pinpoint.Transaction {
	method := strings.TrimPrefix(fullMethod, "/")

	var hdrs http.Header
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		hdrs = make(http.Header, len(md))
		for k, vs := range md {
			for _, v := range vs {
				hdrs.Add(k, v)
			}
		}
	}

	target := hdrs.Get(":authority")
	url := getURL(method, target)

	webReq := pinpoint.WebRequest{
		Header:    hdrs,
		URL:       url,
		Method:    method,
		Transport: pinpoint.TransportHTTP,
	}
	txn := app.StartTransaction(method)
	txn.SetWebRequest(webReq)

	return txn
}

// UnaryServerInterceptor instruments server unary RPCs.
//
// Use this function with grpc.UnaryInterceptor and a pinpoint.Application to
// create a grpc.ServerOption to pass to grpc.NewServer.  This interceptor
// records each unary call with a transaction.  You must use both
// UnaryServerInterceptor and StreamServerInterceptor to instrument unary and
// streaming calls.
//
// Example:
//
//	app, _ := pinpoint.NewApplication(
//		pinpoint.ConfigAppName("gRPC Server"),
//		pinpoint.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
//		pinpoint.ConfigDebugLogger(os.Stdout),
//	)
//	server := grpc.NewServer(
//		grpc.UnaryInterceptor(nrgrpc.UnaryServerInterceptor(app)),
//		grpc.StreamInterceptor(nrgrpc.StreamServerInterceptor(app)),
//	)
//
// These interceptors add the transaction to the call context so it may be
// accessed in your method handlers using pinpoint.FromContext.
//
// Full example:
// https://github.com/pinpoint/go-agent/blob/master/v3/integrations/nrgrpc/example/server/server.go
//
func UnaryServerInterceptor(app *pinpoint.Application) grpc.UnaryServerInterceptor {
	if nil == app {
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			return handler(ctx, req)
		}
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		txn := startTransaction(ctx, app, info.FullMethod)
		defer txn.End()

		ctx = pinpoint.NewContext(ctx, txn)
		resp, err = handler(ctx, req)
		txn.SetWebResponse(nil).WriteHeader(int(status.Code(err)))
		return
	}
}

type wrappedServerStream struct {
	grpc.ServerStream
	txn *pinpoint.Transaction
}

func (s wrappedServerStream) Context() context.Context {
	ctx := s.ServerStream.Context()
	return pinpoint.NewContext(ctx, s.txn)
}

func newWrappedServerStream(stream grpc.ServerStream, txn *pinpoint.Transaction) grpc.ServerStream {
	return wrappedServerStream{
		ServerStream: stream,
		txn:          txn,
	}
}

// StreamServerInterceptor instruments server streaming RPCs.
//
// Use this function with grpc.StreamInterceptor and a pinpoint.Application to
// create a grpc.ServerOption to pass to grpc.NewServer.  This interceptor
// records each streaming call with a transaction.  You must use both
// UnaryServerInterceptor and StreamServerInterceptor to instrument unary and
// streaming calls.
//
// Example:
//
//	app, _ := pinpoint.NewApplication(
//		pinpoint.ConfigAppName("gRPC Server"),
//		pinpoint.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
//		pinpoint.ConfigDebugLogger(os.Stdout),
//	)
//	server := grpc.NewServer(
//		grpc.UnaryInterceptor(nrgrpc.UnaryServerInterceptor(app)),
//		grpc.StreamInterceptor(nrgrpc.StreamServerInterceptor(app)),
//	)
//
// These interceptors add the transaction to the call context so it may be
// accessed in your method handlers using pinpoint.FromContext.
//
// Full example:
// https://github.com/pinpoint/go-agent/blob/master/v3/integrations/nrgrpc/example/server/server.go
//
func StreamServerInterceptor(app *pinpoint.Application) grpc.StreamServerInterceptor {
	if nil == app {
		return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			return handler(srv, ss)
		}
	}

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		txn := startTransaction(ss.Context(), app, info.FullMethod)
		defer txn.End()

		err := handler(srv, newWrappedServerStream(ss, txn))
		txn.SetWebResponse(nil).WriteHeader(int(status.Code(err)))
		return err
	}
}
