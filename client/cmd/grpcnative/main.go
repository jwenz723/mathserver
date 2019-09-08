package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jwenz723/mathserver/pb"
	"google.golang.org/grpc"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

func main() {
	fs := flag.NewFlagSet("mathcli", flag.ExitOnError)
	var (
		grpcAddr = fs.String("grpc-addr", "", "gRPC address of addsvc")
		method   = fs.String("method", "sum", "divide, min, max, multiply, pow, subtract, sum")
	)
	fs.Usage = usageFor(fs, os.Args[0]+" [flags] <a> <b>")
	fs.Parse(os.Args[1:])
	if len(fs.Args()) != 2 {
		fs.Usage()
		os.Exit(1)
	}

	var (
		err error
		op  string
		v   *pb.MathOpReply
	)
	if *grpcAddr == "" {
		fmt.Fprintf(os.Stderr, "error: no remote address specified\n")
		os.Exit(1)
	}
	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	svc := pb.NewMathClient(conn)
	a, _ := strconv.ParseFloat(fs.Args()[0], 10)
	b, _ := strconv.ParseFloat(fs.Args()[1], 10)
	req := pb.MathOpRequest{A: a, B: b}

	switch *method {
	case "divide":
		v, err = svc.Divide(context.Background(), &req)
		op = "/"
	case "max":
		v, err = svc.Max(context.Background(), &req)
		op = "max"
	case "min":
		v, err = svc.Min(context.Background(), &req)
		op = "min"
	case "multiply":
		v, err = svc.Multiply(context.Background(), &req)
		op = "*"
	case "pow":
		v, err = svc.Pow(context.Background(), &req)
		op = "^"
	case "subtract":
		v, err = svc.Subtract(context.Background(), &req)
		op = "-"
	case "sum":
		v, err = svc.Sum(context.Background(), &req)
		op = "+"

	default:
		fmt.Fprintf(os.Stderr, "error: invalid method %q\n", *method)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "%f %s %f = %f\n", a, op, b, v.V)
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
