package server

import (
	"context"
	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"github.com/jwenz723/mathserver/std/pb"
	"github.com/jwenz723/mathserver/std/pkg/mathservice"
)

// compile time assertions to ensure our types are implementing interfaces
var (
	_ pb.MathServer = &server{}
)

type server struct {
	svc mathservice.Service
}

func New(svc mathservice.Service) server {
	return server {
		svc: svc,
	}
}

// GrpcLoggingDecider specifies which methods should have their request/response parameters logged
// by the grpc logging interceptor. Returning false indicates logging should be suppressed.
func (s *server) GrpcLoggingDecider() grpc_logging.ServerPayloadLoggingDecider {
	return func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
		switch fullMethodName {
		default:
			return true
		}
	}
}

// Divide two integers, a/b
func (s *server) Divide(ctx context.Context, req *pb.MathOpRequest) (*pb.MathOpReply, error) {
	v, err := s.svc.Divide(ctx, req.A, req.B)
	return &pb.MathOpReply{
		V:                    v,
		Err:                  err2str(err),
	}, nil
}

// Max two integers, returns the greater value of a and b
func (s *server) Max(ctx context.Context, req *pb.MathOpRequest) (*pb.MathOpReply, error) {
	v, err := s.svc.Max(ctx, req.A, req.B)
	return &pb.MathOpReply{
		V:                    v,
		Err:                  err2str(err),
	}, nil
}

// Min two integers, returns the lesser value of a and b
func (s *server) Min(ctx context.Context, req *pb.MathOpRequest) (*pb.MathOpReply, error) {
	v, err := s.svc.Min(ctx, req.A, req.B)
	return &pb.MathOpReply{
		V:                    v,
		Err:                  err2str(err),
	}, nil
}

// Multiply two integers, a*b
func (s *server) Multiply(ctx context.Context, req *pb.MathOpRequest) (*pb.MathOpReply, error) {
	v, err := s.svc.Multiply(ctx, req.A, req.B)
	return &pb.MathOpReply{
		V:                    v,
		Err:                  err2str(err),
	}, nil
}

// Pow two integers, a^b
func (s *server) Pow(ctx context.Context, req *pb.MathOpRequest) (*pb.MathOpReply, error) {
	v, err := s.svc.Pow(ctx, req.A, req.B)
	return &pb.MathOpReply{
		V:                    v,
		Err:                  err2str(err),
	}, nil
}

// Subtract two integers, a-b
func (s *server) Subtract(ctx context.Context, req *pb.MathOpRequest) (*pb.MathOpReply, error) {
	v, err := s.svc.Subtract(ctx, req.A, req.B)
	return &pb.MathOpReply{
		V:                    v,
		Err:                  err2str(err),
	}, nil
}

// Sums two integers. a+b
func (s *server) Sum(ctx context.Context, req *pb.MathOpRequest) (*pb.MathOpReply, error) {
	v, err := s.svc.Sum(ctx, req.A, req.B)
	return &pb.MathOpReply{
		V:                    v,
		Err:                  err2str(err),
	}, nil
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}