package mathtransport

import (
	"context"
	"errors"
	"github.com/jwenz723/mathserver/gokit/pb"
	"github.com/jwenz723/mathserver/gokit/pkg/mathendpoint"
	"github.com/jwenz723/mathserver/gokit/pkg/mathservice"
	"google.golang.org/grpc"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	divide grpctransport.Handler
	max grpctransport.Handler
	min grpctransport.Handler
	multiply grpctransport.Handler
	pow grpctransport.Handler
	subtact grpctransport.Handler
	sum    grpctransport.Handler
}

// NewGRPCServer makes a set of endpoints available as a gRPC MathServer.
func NewGRPCServer(endpoints mathendpoint.Set, logger log.Logger) pb.MathServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &grpcServer{
		divide: grpctransport.NewServer(
			endpoints.DivideEndpoint,
			decodeGRPCMathOpRequest,
			encodeGRPCMathOpResponse,
			options...,
		),
		max:      grpctransport.NewServer(
			endpoints.MaxEndpoint,
			decodeGRPCMathOpRequest,
			encodeGRPCMathOpResponse,
			options...,
		),
		min:      grpctransport.NewServer(
			endpoints.MinEndpoint,
			decodeGRPCMathOpRequest,
			encodeGRPCMathOpResponse,
			options...,
		),
		multiply: grpctransport.NewServer(
			endpoints.MultiplyEndpoint,
			decodeGRPCMathOpRequest,
			encodeGRPCMathOpResponse,
			options...,
		),
		pow:      grpctransport.NewServer(
			endpoints.PowEndpoint,
			decodeGRPCMathOpRequest,
			encodeGRPCMathOpResponse,
			options...,
		),
		subtact:  grpctransport.NewServer(
			endpoints.SubtractEndpoint,
			decodeGRPCMathOpRequest,
			encodeGRPCMathOpResponse,
			options...,
		),
		sum:      grpctransport.NewServer(
			endpoints.SumEndpoint,
			decodeGRPCMathOpRequest,
			encodeGRPCMathOpResponse,
			options...,
		),
	}
}

func (s *grpcServer) Divide(ctx context.Context, req *pb.MathOpRequest) (*pb.MathOpReply, error) {
	_, rep, err := s.divide.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.MathOpReply), nil
}

func (s *grpcServer) Max(ctx context.Context, req *pb.MathOpRequest) (*pb.MathOpReply, error) {
	_, rep, err := s.max.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.MathOpReply), nil
}

func (s *grpcServer) Min(ctx context.Context, req *pb.MathOpRequest) (*pb.MathOpReply, error) {
	_, rep, err := s.min.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.MathOpReply), nil
}

func (s *grpcServer) Multiply(ctx context.Context, req *pb.MathOpRequest) (*pb.MathOpReply, error) {
	_, rep, err := s.multiply.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.MathOpReply), nil
}

func (s *grpcServer) Pow(ctx context.Context, req *pb.MathOpRequest) (*pb.MathOpReply, error) {
	_, rep, err := s.pow.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.MathOpReply), nil
}

func (s *grpcServer) Subtract(ctx context.Context, req *pb.MathOpRequest) (*pb.MathOpReply, error) {
	_, rep, err := s.subtact.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.MathOpReply), nil
}

func (s *grpcServer) Sum(ctx context.Context, req *pb.MathOpRequest) (*pb.MathOpReply, error) {
	_, rep, err := s.sum.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.MathOpReply), nil
}

// NewGRPCClient returns an MathService backed by a gRPC server at the other end
// of the conn. The caller is responsible for constructing the conn, and
// eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
func NewGRPCClient(conn *grpc.ClientConn, logger log.Logger) mathservice.Service {
	var sumEndpoint endpoint.Endpoint
	{
		sumEndpoint = grpctransport.NewClient(
			conn,
			"pb.Math",
			"Sum",
			encodeGRPCMathOpRequest,
			decodeGRPCMathOpResponse,
			pb.MathOpReply{},
		).Endpoint()
	}

	// Returning the endpoint.Set as a service.Service relies on the
	// endpoint.Set implementing the Service methods. That's just a simple bit
	// of glue code.
	return mathendpoint.Set{
		SumEndpoint:    sumEndpoint,
	}
}

// decodeGRPCSumRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC sum request to a user-domain sum request. Primarily useful in a server.
func decodeGRPCMathOpRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.MathOpRequest)
	return mathendpoint.MathOpRequest{A: req.A, B: req.B}, nil
}

// decodeGRPCMathOpResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC MathOp reply to a user-domain MathOp response. Primarily useful in a client.
func decodeGRPCMathOpResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.MathOpReply)
	return mathendpoint.MathOpResponse{V: reply.V, Err: str2err(reply.Err)}, nil
}

// encodeGRPCMathOpResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain MathOp response to a gRPC MathOp reply. Primarily useful in a server.
func encodeGRPCMathOpResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(mathendpoint.MathOpResponse)
	return &pb.MathOpReply{V: resp.V, Err: err2str(resp.Err)}, nil
}

// encodeGRPCMathOpRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain MathOp request to a gRPC MathOp request. Primarily useful in a client.
func encodeGRPCMathOpRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(mathendpoint.MathOpRequest)
	return &pb.MathOpRequest{A: req.A, B: req.B}, nil
}

// These annoying helper functions are required to translate Go error types to
// and from strings, which is the type we use in our IDLs to represent errors.
// There is special casing to treat empty strings as nil errors.

func str2err(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
