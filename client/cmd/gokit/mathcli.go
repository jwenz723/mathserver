package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jwenz723/mathserver/gokit/pkg/mathservice"
	"github.com/jwenz723/mathserver/gokit/pkg/mathtransport"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
	"google.golang.org/grpc"
	"github.com/go-kit/kit/log"
)

func main() {
	fs := flag.NewFlagSet("mathcli", flag.ExitOnError)
	var (
		httpAddr       = fs.String("http-addr", "", "HTTP address of addsvc")
		grpcAddr       = fs.String("grpc-addr", "", "gRPC address of addsvc")
		method         = fs.String("method", "sum", "divide, min, max, multiply, pow, subtract, sum")
	)
	fs.Usage = usageFor(fs, os.Args[0]+" [flags] <a> <b>")
	fs.Parse(os.Args[1:])
	if len(fs.Args()) != 2 {
		fs.Usage()
		os.Exit(1)
	}

	// This is a demonstration client, which supports multiple transports.
	// Your clients will probably just define and stick with 1 transport.
	var (
		svc mathservice.Service
		err error
		op string
		v float64
	)
	if *httpAddr != "" {
		svc, err = mathtransport.NewHTTPClient(*httpAddr, log.NewNopLogger())
	} else if *grpcAddr != "" {
		conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			os.Exit(1)
		}
		defer conn.Close()
		svc = mathtransport.NewGRPCClient(conn, log.NewNopLogger())
	} else {
		fmt.Fprintf(os.Stderr, "error: no remote address specified\n")
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	a, _ := strconv.ParseFloat(fs.Args()[0], 10)
	b, _ := strconv.ParseFloat(fs.Args()[1], 10)

	switch *method {
	case "divide":
		v, err = svc.Divide(context.Background(), a, b)
		op = "/"
	case "max":
		v, err = svc.Max(context.Background(), a, b)
		op = "max"
	case "min":
		v, err = svc.Min(context.Background(), a, b)
		op = "min"
	case "multiply":
		v, err = svc.Multiply(context.Background(), a, b)
		op = "*"
	case "pow":
		v, err = svc.Pow(context.Background(), a, b)
		op = "^"
	case "subtract":
		v, err = svc.Subtract(context.Background(), a, b)
		op = "-"
	case "sum":
		v, err = svc.Sum(context.Background(), a, b)
		op = "+"

	default:
		fmt.Fprintf(os.Stderr, "error: invalid method %q\n", *method)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "%f %s %f = %f\n", a, op, b, v)
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
