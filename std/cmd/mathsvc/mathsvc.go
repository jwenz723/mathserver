package main

import (
	"flag"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/jwenz723/mathserver/std/pb"
	"github.com/jwenz723/mathserver/std/pkg/mathservice"
	"github.com/jwenz723/mathserver/std/pkg/server"
	"github.com/oklog/oklog/pkg/group"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"
)

func main() {
	// Define our flags. Your service probably won't need to bind listeners for
	// *all* supported transports, or support both Zipkin and LightStep, and so
	// on, but we do it here for demonstration purposes.
	fs := flag.NewFlagSet("mathsvc", flag.ExitOnError)
	var (
		debugAddr      = fs.String("debug.addr", ":8080", "Debug and metrics listen address")
		httpAddr       = fs.String("http-addr", ":8081", "HTTP listen address")
		grpcAddr       = fs.String("grpc-addr", ":8082", "gRPC listen address")
	)
	fs.Usage = usageFor(fs, os.Args[0]+" [flags]")
	fs.Parse(os.Args[1:])

	logger, _ := zap.NewProduction()

	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())

	var (
		service = mathservice.NewBasicService()
		grpcSvc = server.NewGrpcServer(service)
		httpSvc = server.NewHttpServer(service)
	)

	var g group.Group
	{
		// The debug listener mounts the http.DefaultServeMux, and serves up
		// stuff like the Prometheus metrics route, the Go debug and profiling
		// routes, and so on.
		debugListener, err := net.Listen("tcp", *debugAddr)
		if err != nil {
			logger.Error("failed to start listener",
				zap.String("transport", "debug/HTTP"),
				zap.String("during", "Listen"),
				zap.Error(err))
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Info("starting listener",
				zap.String("transport", "debug/HTTP"),
				zap.String("addr", *debugAddr))
			return http.Serve(debugListener, http.DefaultServeMux)
		}, func(error) {
			debugListener.Close()
		})
	}
	{
		// The HTTP listener mounts the Go kit HTTP handler we created.
		httpListener, err := net.Listen("tcp", *httpAddr)
		if err != nil {
			logger.Error("failed to start httpSvc listener",
				zap.Error(err))
			os.Exit(1)
		}

		g.Add(func() error {
			logger.Info("starting httpSvc listener",
				zap.String("addr", *httpAddr))
			return http.Serve(httpListener, httpSvc.Router())
		}, func(error) {
			httpListener.Close()
		})
	}
	{
		logger.Info("starting grpcSvc listener",
			zap.String("addr", *grpcAddr))
		lis, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			logger.Error("failed to start grpcSvc listener", zap.Error(err))
		}

		g.Add(func() error {
			grpcServer := grpc.NewServer(
				grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
					grpc_prometheus.UnaryServerInterceptor,
					grpc_zap.UnaryServerInterceptor(logger),
					grpc_zap.PayloadUnaryServerInterceptor(logger, grpcSvc.GrpcLoggingDecider()),
				)),
			)

			pb.RegisterMathServer(grpcServer, &grpcSvc)
			return grpcServer.Serve(lis)
		}, func(error) {
			lis.Close()
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	logger.Error("exit", zap.Error(g.Run()))
}

func usageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t-%s %s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}
