// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package transport

import "golang.org/x/net/context"

// Interceptor defines a transport-level middleware for Inbounds.
//
// Interceptors MAY
//
// - change the context
// - change the request
// - call the ResponseWriter
// - modify the response body by wrapping the ResponseWriter
// - handle the returned error
// - call the given handler zero or more times
//
// Interceptors MUST be thread-safe.
//
// Interceptors are re-used across requests and MAY be called multiple times for
// the same request.
type Interceptor interface {
	Handle(ctx context.Context, req *Request, resw ResponseWriter, h Handler) error
}

// NopInterceptor is a interceptor that does not do anything special. It
// simply calls the underlying Handler.
var NopInterceptor Interceptor = nopInterceptor{}

// ApplyInterceptor applies the given Interceptor to the given Handler.
func ApplyInterceptor(h Handler, i Interceptor) Handler {
	if i == nil {
		return h
	}
	return interceptedHandler{h: h, i: i}
}

// InterceptorFunc adapts a function into an Interceptor.
type InterceptorFunc func(context.Context, *Request, ResponseWriter, Handler) error

// Handle for InterceptorFunc
func (f InterceptorFunc) Handle(ctx context.Context, req *Request, resw ResponseWriter, h Handler) error {
	return f(ctx, req, resw, h)
}

type interceptedHandler struct {
	h Handler
	i Interceptor
}

func (h interceptedHandler) Handle(ctx context.Context, req *Request, resw ResponseWriter) error {
	return h.i.Handle(ctx, req, resw, h.h)
}

type nopInterceptor struct{}

func (nopInterceptor) Handle(ctx context.Context, req *Request, resw ResponseWriter, handler Handler) error {
	return handler.Handle(ctx, req, resw)
}