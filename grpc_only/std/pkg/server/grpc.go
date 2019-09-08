package server

import (
	"context"
	mathservice2 "github.com/jwenz723/mathserver/grpc_only/std/pkg/mathservice"
	"github.com/jwenz723/mathserver/pb"
)

// compile time assertions to ensure our types are implementing interfaces
var (
	_ pb.MathServer = &grpcServer{}
)

type grpcServer struct {
	svc mathservice2.Service
}

func NewGrpcServer(svc mathservice2.Service) grpcServer {
	return grpcServer{
		svc: svc,
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