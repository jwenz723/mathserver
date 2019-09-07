package mathendpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/jwenz723/mathserver/gokit/pkg/mathservice"
)

// Set collects all of the endpoints that compose an add service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Set struct {
	DivideEndpoint endpoint.Endpoint
	MaxEndpoint endpoint.Endpoint
	MinEndpoint endpoint.Endpoint
	MultiplyEndpoint endpoint.Endpoint
	PowEndpoint endpoint.Endpoint
	SubtractEndpoint endpoint.Endpoint
	SumEndpoint endpoint.Endpoint
}

// New returns a Set that wraps the provided server, and wires in all of the
// expected endpoint middlewares via the various parameters.
func New(svc mathservice.Service, logger log.Logger, duration metrics.Histogram) Set {
	var divideEndpoint endpoint.Endpoint
	{
		divideEndpoint = MakeDivideEndpoint(svc)
		divideEndpoint = LoggingMiddleware(log.With(logger, "method", "Divide"))(divideEndpoint)
		divideEndpoint = InstrumentingMiddleware(duration.With("method", "Divide"))(divideEndpoint)
	}
	var maxEndpoint endpoint.Endpoint
	{
		maxEndpoint = MakeMaxEndpoint(svc)
		maxEndpoint = LoggingMiddleware(log.With(logger, "method", "Max"))(maxEndpoint)
		maxEndpoint = InstrumentingMiddleware(duration.With("method", "Max"))(maxEndpoint)
	}
	var minEndpoint endpoint.Endpoint
	{
		minEndpoint = MakeMinEndpoint(svc)
		minEndpoint = LoggingMiddleware(log.With(logger, "method", "Min"))(minEndpoint)
		minEndpoint = InstrumentingMiddleware(duration.With("method", "Min"))(minEndpoint)
	}
	var multiplyEndpoint endpoint.Endpoint
	{
		multiplyEndpoint = MakeMultiplyEndpoint(svc)
		multiplyEndpoint = LoggingMiddleware(log.With(logger, "method", "Multiply"))(multiplyEndpoint)
		multiplyEndpoint = InstrumentingMiddleware(duration.With("method", "Multiply"))(multiplyEndpoint)
	}
	var powEndpoint endpoint.Endpoint
	{
		powEndpoint = MakePowEndpoint(svc)
		powEndpoint = LoggingMiddleware(log.With(logger, "method", "Pow"))(powEndpoint)
		powEndpoint = InstrumentingMiddleware(duration.With("method", "Pow"))(powEndpoint)
	}
	var subtractEndpoint endpoint.Endpoint
	{
		subtractEndpoint = MakeSubtractEndpoint(svc)
		subtractEndpoint = LoggingMiddleware(log.With(logger, "method", "Subtract"))(subtractEndpoint)
		subtractEndpoint = InstrumentingMiddleware(duration.With("method", "Subtract"))(subtractEndpoint)
	}
	var sumEndpoint endpoint.Endpoint
	{
		sumEndpoint = MakeSumEndpoint(svc)
		sumEndpoint = LoggingMiddleware(log.With(logger, "method", "Sum"))(sumEndpoint)
		sumEndpoint = InstrumentingMiddleware(duration.With("method", "Sum"))(sumEndpoint)
	}

	return Set{
		DivideEndpoint:   divideEndpoint,
		MaxEndpoint:      maxEndpoint,
		MinEndpoint:      minEndpoint,
		MultiplyEndpoint: multiplyEndpoint,
		PowEndpoint:      powEndpoint,
		SubtractEndpoint: subtractEndpoint,
		SumEndpoint:      sumEndpoint,
	}
}

// Divide implements the service interface, so Set may be used as a service.
// This is primarily useful in the context of a client library.
func (s Set) Divide(ctx context.Context, a, b float64) (float64, error) {
	resp, err := s.DivideEndpoint(ctx, MathOpRequest{A: a, B: b})
	if err != nil {
		return 0, err
	}
	response := resp.(MathOpResponse)
	return response.V, response.Err
}

// MakeDivideEndpoint constructs a Divide endpoint wrapping the service.
func MakeDivideEndpoint(s mathservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(MathOpRequest)
		v, err := s.Divide(ctx, req.A, req.B)
		return MathOpResponse{V: v, Err: err}, nil
	}
}

// Max implements the service interface, so Set may be used as a service.
// This is primarily useful in the context of a client library.
func (s Set) Max(ctx context.Context, a, b float64) (float64, error) {
	resp, err := s.MaxEndpoint(ctx, MathOpRequest{A: a, B: b})
	if err != nil {
		return 0, err
	}
	response := resp.(MathOpResponse)
	return response.V, response.Err
}

// MakeMaxEndpoint constructs a Max endpoint wrapping the service.
func MakeMaxEndpoint(s mathservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(MathOpRequest)
		v, err := s.Max(ctx, req.A, req.B)
		return MathOpResponse{V: v, Err: err}, nil
	}
}

// Min implements the service interface, so Set may be used as a service.
// This is primarily useful in the context of a client library.
func (s Set) Min(ctx context.Context, a, b float64) (float64, error) {
	resp, err := s.MinEndpoint(ctx, MathOpRequest{A: a, B: b})
	if err != nil {
		return 0, err
	}
	response := resp.(MathOpResponse)
	return response.V, response.Err
}

// MakeMinEndpoint constructs a Min endpoint wrapping the service.
func MakeMinEndpoint(s mathservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(MathOpRequest)
		v, err := s.Min(ctx, req.A, req.B)
		return MathOpResponse{V: v, Err: err}, nil
	}
}

// Multiply implements the service interface, so Set may be used as a service.
// This is primarily useful in the context of a client library.
func (s Set) Multiply(ctx context.Context, a, b float64) (float64, error) {
	resp, err := s.MultiplyEndpoint(ctx, MathOpRequest{A: a, B: b})
	if err != nil {
		return 0, err
	}
	response := resp.(MathOpResponse)
	return response.V, response.Err
}

// MakeMultiplyEndpoint constructs a Multiply endpoint wrapping the service.
func MakeMultiplyEndpoint(s mathservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(MathOpRequest)
		v, err := s.Multiply(ctx, req.A, req.B)
		return MathOpResponse{V: v, Err: err}, nil
	}
}

// Pow implements the service interface, so Set may be used as a service.
// This is primarily useful in the context of a client library.
func (s Set) Pow(ctx context.Context, a, b float64) (float64, error) {
	resp, err := s.PowEndpoint(ctx, MathOpRequest{A: a, B: b})
	if err != nil {
		return 0, err
	}
	response := resp.(MathOpResponse)
	return response.V, response.Err
}

// MakePowEndpoint constructs a Pow endpoint wrapping the service.
func MakePowEndpoint(s mathservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(MathOpRequest)
		v, err := s.Pow(ctx, req.A, req.B)
		return MathOpResponse{V: v, Err: err}, nil
	}
}

// Subtract implements the service interface, so Set may be used as a service.
// This is primarily useful in the context of a client library.
func (s Set) Subtract(ctx context.Context, a, b float64) (float64, error) {
	resp, err := s.SubtractEndpoint(ctx, MathOpRequest{A: a, B: b})
	if err != nil {
		return 0, err
	}
	response := resp.(MathOpResponse)
	return response.V, response.Err
}

// MakeSubtractEndpoint constructs a Subtract endpoint wrapping the service.
func MakeSubtractEndpoint(s mathservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(MathOpRequest)
		v, err := s.Subtract(ctx, req.A, req.B)
		return MathOpResponse{V: v, Err: err}, nil
	}
}

// Sum implements the service interface, so Set may be used as a service.
// This is primarily useful in the context of a client library.
func (s Set) Sum(ctx context.Context, a, b float64) (float64, error) {
	resp, err := s.SumEndpoint(ctx, MathOpRequest{A: a, B: b})
	if err != nil {
		return 0, err
	}
	response := resp.(MathOpResponse)
	return response.V, response.Err
}

// MakeSumEndpoint constructs a Sum endpoint wrapping the service.
func MakeSumEndpoint(s mathservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(MathOpRequest)
		v, err := s.Sum(ctx, req.A, req.B)
		return MathOpResponse{V: v, Err: err}, nil
	}
}

// compile time assertions for our response types implementing endpoint.Failer.
var (
	_ endpoint.Failer = MathOpResponse{}
)

// MathOpRequest collects the request parameters for the math methods.
type MathOpRequest struct {
	A, B float64
}

// MathOpResponse collects the response values for the math methods.
type MathOpResponse struct {
	V   float64   `json:"v"`
	Err error `json:"-"` // should be intercepted by Failed/errorEncoder
}

// Failed implements endpoint.Failer.
func (r MathOpResponse) Failed() error { return r.Err }