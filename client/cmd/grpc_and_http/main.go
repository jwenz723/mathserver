package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jwenz723/mathserver/grpc_and_http/std/pkg/server"
	"github.com/jwenz723/mathserver/pb"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

func main() {
	fs := flag.NewFlagSet("mathcli", flag.ExitOnError)
	var (
		grpcAddr = fs.String("grpc-addr", "", "gRPC address of addsvc")
		httpAddr = fs.String("http-addr", "", "HTTP address of addsvc")
		method   = fs.String("method", "sum", "divide, min, max, multiply, pow, subtract, sum")
	)
	fs.Usage = usageFor(fs, os.Args[0]+" [flags] <a> <b>")
	fs.Parse(os.Args[1:])
	if len(fs.Args()) != 2 {
		fs.Usage()
		os.Exit(1)
	}

	a, _ := strconv.ParseFloat(fs.Args()[0], 10)
	b, _ := strconv.ParseFloat(fs.Args()[1], 10)
	var result float64

	if *grpcAddr != "" {
		conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
		checkErr(err)
		defer conn.Close()

		svc := pb.NewMathClient(conn)
		req := pb.MathOpRequest{A: a, B: b}

		var resp *pb.MathOpReply
		switch *method {
		case "divide":
			resp, err = svc.Divide(context.Background(), &req)
		case "max":
			resp, err = svc.Max(context.Background(), &req)
		case "min":
			resp, err = svc.Min(context.Background(), &req)
		case "multiply":
			resp, err = svc.Multiply(context.Background(), &req)
		case "pow":
			resp, err = svc.Pow(context.Background(), &req)
		case "subtract":
			resp, err = svc.Subtract(context.Background(), &req)
		case "sum":
			resp, err = svc.Sum(context.Background(), &req)
		default:
			fmt.Fprintf(os.Stderr, "error: invalid method %q\n", *method)
			os.Exit(1)
		}
		checkErr(err)

		result = resp.V
	} else if *httpAddr != "" {
		req := server.MathOpRequest{
			A: a,
			B: b,
		}
		b, err := json.Marshal(req)
		if err != nil {
			panic(err)
		}
		resp, err := http.Post(fmt.Sprintf("http://%s/%s", *httpAddr, *method), "application/json", bytes.NewBuffer(b))
		checkErr(err)
		defer resp.Body.Close()

		var m server.MathOpResponse
		err = json.NewDecoder(resp.Body).Decode(&m)
		checkErr(err)

		result = m.V
	} else {
		fmt.Fprintf(os.Stderr, "error: no remote address specified\n")
		os.Exit(1)
	}

	printResult(*method, a, b, result)
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

func printResult(method string, a, b, v float64) {
	var op string
	switch method {
	case "divide":
		op = "/"
	case "max":
		op = "max"
	case "min":
		op = "min"
	case "multiply":
		op = "*"
	case "pow":
		op = "^"
	case "subtract":
		op = "-"
	case "sum":
		op = "+"
	default:
		fmt.Fprintf(os.Stderr, "error: invalid method %q\n", method)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "%f %s %f = %f\n", a, op, b, v)
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}
