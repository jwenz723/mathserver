package mathendpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	mathservice2 "github.com/jwenz723/mathserver/grpc_and_http/gokit/pkg/mathservice"
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
func New(svc mathservice2.Service, logger log.Logger) Set {
	return Set{
		DivideEndpoint:   MakeDivideEndpoint(svc),
		MaxEndpoint:      MakeMaxEndpoint(svc),
		MinEndpoint:      MakeMinEndpoint(svc),
		MultiplyEndpoint: MakeMultiplyEndpoint(svc),
		PowEndpoint:      MakePowEndpoint(svc),
		SubtractEndpoint: MakeSubtractEndpoint(svc),
		SumEndpoint:      MakeSumEndpoint(svc),
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
func MakeDivideEndpoint(s mathservice2.Service) endpoint.Endpoint {
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
func MakeMaxEndpoint(s mathservice2.Service) endpoint.Endpoint {
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
func MakeMinEndpoint(s mathservice2.Service) endpoint.Endpoint {
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
func MakeMultiplyEndpoint(s mathservice2.Service) endpoint.Endpoint {
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
func MakePowEndpoint(s mathservice2.Service) endpoint.Endpoint {
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
func MakeSubtractEndpoint(s mathservice2.Service) endpoint.Endpoint {
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
func MakeSumEndpoint(s mathservice2.Service) endpoint.Endpoint {
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