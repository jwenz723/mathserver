package server

import (
	"context"
	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	"github.com/jwenz723/mathserver/pb"
	"github.com/jwenz723/mathserver/std/pkg/mathservice"
)

// compile time assertions to ensure our types are implementing interfaces
var (
	_ pb.MathServer = &grpcServer{}
)

type grpcServer struct {
	svc mathservice.Service
}

func NewGrpcServer(svc mathservice.Service) grpcServer {
	return grpcServer{
		svc: svc,
	}
}

// GrpcLoggingDecider specifies which methods should have their request/response parameters logged
// by the grpc logging interceptor. Returning false indicates logging should be suppressed.
func (s *grpcServer) GrpcLoggingDecider() grpc_logging.ServerPayloadLoggingDecider {
	return func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
		switch fullMethodName {
		default:
			return true
		}
	}
}

// Divide two integers, a/b
func (s *grpcServer) Divide(ctx context.Context, req *pb.MathOpRequest) (*pb.MathOpReply, error) {
	v, err := s.svc.Divide(ctx, req.A, req.B)
	return &pb.MathOpReply{
		V:                    v,
		Err:                  err2str(err),
	}, nil
}

// Max two integers, returns the greater value of a and b
func (s *grpcServer) Max(ctx context.Context, req *pb.MathOpRequest) (*pb.MathOpReply, error) {
	v, err := s.svc.Max(ctx, req.A, req.B)
	return &pb.MathOpReply{
		V:                    v,
		Err:                  err2str(err),
	}, nil
}

// Min two integers, returns the lesser value of a and b
func (s *grpcServer) Min(ctx context.Context, req *pb.MathOpRequest) (*pb.MathOpReply, error) {
	v, err := s.svc.Min(ctx, req.A, req.B)
	return &pb.MathOpReply{
		V:                    v,
		Err:                  err2str(err),
	}, nil
}

// Multiply two integers, a*b
func (s *grpcServer) Multiply(ctx context.Context, req *pb.MathOpRequest) (*pb.MathOpReply, error) {
	v, err := s.svc.Multiply(ctx, req.A, req.B)
	return &pb.MathOpReply{
		V:                    v,
		Err:                  err2str(err),
	}, nil
}

// Pow two integers, a^b
func (s *grpcServer) Pow(ctx context.Context, req *pb.MathOpRequest) (*pb.MathOpReply, error) {
	v, err := s.svc.Pow(ctx, req.A, req.B)
	return &pb.MathOpReply{
		V:                    v,
		Err:                  err2str(err),
	}, nil
}

// Subtract two integers, a-b
func (s *grpcServer) Subtract(ctx context.Context, req *pb.MathOpRequest) (*pb.MathOpReply, error) {
	v, err := s.svc.Subtract(ctx, req.A, req.B)
	return &pb.MathOpReply{
		V:                    v,
		Err:                  err2str(err),
	}, nil
}

// Sums two integers. a+b
func (s *grpcServer) Sum(ctx context.Context, req *pb.MathOpRequest) (*pb.MathOpReply, error) {
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