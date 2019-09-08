package mathtransport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	mathendpoint2 "github.com/jwenz723/mathserver/grpc_and_http/gokit/pkg/mathendpoint"
	mathservice2 "github.com/jwenz723/mathserver/grpc_and_http/gokit/pkg/mathservice"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// NewHTTPHandler returns an HTTP handler that makes a set of endpoints
// available on predefined paths.
func NewHTTPHandler(endpoints mathendpoint2.Set, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	m := http.NewServeMux()
	m.Handle("/divide", httptransport.NewServer(
		endpoints.DivideEndpoint,
		decodeHTTPMathOpRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	m.Handle("/max", httptransport.NewServer(
		endpoints.MaxEndpoint,
		decodeHTTPMathOpRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	m.Handle("/min", httptransport.NewServer(
		endpoints.MinEndpoint,
		decodeHTTPMathOpRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	m.Handle("/multiply", httptransport.NewServer(
		endpoints.MultiplyEndpoint,
		decodeHTTPMathOpRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	m.Handle("/pow", httptransport.NewServer(
		endpoints.PowEndpoint,
		decodeHTTPMathOpRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	m.Handle("/subtract", httptransport.NewServer(
		endpoints.SubtractEndpoint,
		decodeHTTPMathOpRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	m.Handle("/sum", httptransport.NewServer(
		endpoints.SumEndpoint,
		decodeHTTPMathOpRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	return m
}

// NewHTTPClient returns an MathService backed by an HTTP server living at the
// remote instance. We expect instance to come from a service discovery system,
// so likely of the form "host:port". We bake-in certain middlewares,
// implementing the client library pattern.
func NewHTTPClient(instance string, logger log.Logger) (mathservice2.Service, error) {
	// Quickly sanitize the instance string.
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	var divideEndpoint endpoint.Endpoint
	{
		divideEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/divide"),
			encodeHTTPGenericRequest,
			decodeHTTPMathOpResponse,
		).Endpoint()
	}
	var maxEndpoint endpoint.Endpoint
	{
		maxEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/max"),
			encodeHTTPGenericRequest,
			decodeHTTPMathOpResponse,
		).Endpoint()
	}
	var minEndpoint endpoint.Endpoint
	{
		minEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/min"),
			encodeHTTPGenericRequest,
			decodeHTTPMathOpResponse,
		).Endpoint()
	}
	var multiplyEndpoint endpoint.Endpoint
	{
		multiplyEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/multiply"),
			encodeHTTPGenericRequest,
			decodeHTTPMathOpResponse,
		).Endpoint()
	}
	var powEndpoint endpoint.Endpoint
	{
		powEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/pow"),
			encodeHTTPGenericRequest,
			decodeHTTPMathOpResponse,
		).Endpoint()
	}
	var subtractEndpoint endpoint.Endpoint
	{
		subtractEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/subtract"),
			encodeHTTPGenericRequest,
			decodeHTTPMathOpResponse,
		).Endpoint()
	}
	var sumEndpoint endpoint.Endpoint
	{
		sumEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/sum"),
			encodeHTTPGenericRequest,
			decodeHTTPMathOpResponse,
		).Endpoint()
	}

	// Returning the endpoint.Set as a service.Service relies on the
	// endpoint.Set implementing the Service methods. That's just a simple bit
	// of glue code.
	return mathendpoint2.Set{
		DivideEndpoint:   divideEndpoint,
		MaxEndpoint:      maxEndpoint,
		MinEndpoint:      minEndpoint,
		MultiplyEndpoint: multiplyEndpoint,
		PowEndpoint:      powEndpoint,
		SubtractEndpoint: subtractEndpoint,
		SumEndpoint:      sumEndpoint,
	}, nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func err2code(err error) int {
	switch err {
	case mathservice2.ErrDivideByZero, mathservice2.ErrNoMax, mathservice2.ErrNoMin:
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

func errorDecoder(r *http.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

type errorWrapper struct {
	Error string `json:"error"`
}

// decodeHTTPMathOpRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded MathOp request from the HTTP request body. Primarily useful in a
// server.
func decodeHTTPMathOpRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req mathendpoint2.MathOpRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// decodeHTTPMathOpResponse is a transport/http.DecodeResponseFunc that decodes a
// JSON-encoded MathOp response from the HTTP response body. If the response has a
// non-200 status code, we will interpret that as an error and attempt to decode
// the specific error message from the response body. Primarily useful in a
// client.
func decodeHTTPMathOpResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp mathendpoint2.MathOpResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

// encodeHTTPGenericRequest is a transport/http.EncodeRequestFunc that
// JSON-encodes any request to the request body. Primarily useful in a client.
func encodeHTTPGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// encodeHTTPGenericResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
