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

type mathServer interface {
	Divide(ctx context.Context, a, b float64) (float64, error)
	Max(ctx context.Context, a, b float64) (float64, error)
	Min(ctx context.Context, a, b float64) (float64, error)
	Multiply(ctx context.Context, a, b float64) (float64, error)
	Pow(ctx context.Context, a, b float64) (float64, error)
	Subtract(ctx context.Context, a, b float64) (float64, error)
	Sum(ctx context.Context, a, b float64) (float64, error)
}

type httpMathServer struct {
	addr string
}

func (h httpMathServer) handleMathRequest(method string, ctx context.Context, a, b float64) (float64, error) {
	req := server.MathOpRequest{
		A: a,
		B: b,
	}
	j, err := json.Marshal(req)
	if err != nil {
		return 0, err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/%s", h.addr, method), "application/json", bytes.NewBuffer(j))
	checkErr(err)
	defer resp.Body.Close()

	var m server.MathOpResponse
	err = json.NewDecoder(resp.Body).Decode(&m)
	return m.V, err
}

func (h httpMathServer) Divide(ctx context.Context, a, b float64) (float64, error) {
	return h.handleMathRequest("divide", ctx, a, b)
}

func (h httpMathServer) Max(ctx context.Context, a, b float64) (float64, error) {
	return h.handleMathRequest("max", ctx, a, b)
}

func (h httpMathServer) Min(ctx context.Context, a, b float64) (float64, error) {
	return h.handleMathRequest("min", ctx, a, b)
}

func (h httpMathServer) Multiply(ctx context.Context, a, b float64) (float64, error) {
	return h.handleMathRequest("multiply", ctx, a, b)
}

func (h httpMathServer) Pow(ctx context.Context, a, b float64) (float64, error) {
	return h.handleMathRequest("pow", ctx, a, b)
}

func (h httpMathServer) Subtract(ctx context.Context, a, b float64) (float64, error) {
	return h.handleMathRequest("subtract", ctx, a, b)
}

func (h httpMathServer) Sum(ctx context.Context, a, b float64) (float64, error) {
	return h.handleMathRequest("sum", ctx, a, b)
}

type grpcMathServer struct {
	g pb.MathClient
}

func (g grpcMathServer) Divide(ctx context.Context, a, b float64) (float64, error) {
	r, err := g.g.Divide(ctx, &pb.MathOpRequest{
		A: a,
		B: b,
	})
	return r.V, err
}

func (g grpcMathServer) Max(ctx context.Context, a, b float64) (float64, error) {
	r, err := g.g.Max(ctx, &pb.MathOpRequest{
		A: a,
		B: b,
	})
	return r.V, err
}

func (g grpcMathServer) Min(ctx context.Context, a, b float64) (float64, error) {
	r, err := g.g.Min(ctx, &pb.MathOpRequest{
		A: a,
		B: b,
	})
	return r.V, err
}

func (g grpcMathServer) Multiply(ctx context.Context, a, b float64) (float64, error) {
	r, err := g.g.Multiply(ctx, &pb.MathOpRequest{
		A: a,
		B: b,
	})
	return r.V, err
}

func (g grpcMathServer) Pow(ctx context.Context, a, b float64) (float64, error) {
	r, err := g.g.Pow(ctx, &pb.MathOpRequest{
		A: a,
		B: b,
	})
	return r.V, err
}

func (g grpcMathServer) Subtract(ctx context.Context, a, b float64) (float64, error) {
	r, err := g.g.Subtract(ctx, &pb.MathOpRequest{
		A: a,
		B: b,
	})
	return r.V, err
}

func (g grpcMathServer) Sum(ctx context.Context, a, b float64) (float64, error) {
	r, err := g.g.Sum(ctx, &pb.MathOpRequest{
		A: a,
		B: b,
	})
	return r.V, err
}

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

	var m mathServer
	if *grpcAddr != "" {
		conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
		checkErr(err)
		defer conn.Close()

		svc := pb.NewMathClient(conn)
		m = grpcMathServer{g: svc}
	} else if *httpAddr != "" {
		m = httpMathServer{addr: *httpAddr}
	}

	var (
		err  error
		op   string
		resp float64
	)
	switch *method {
	case "divide":
		resp, err = m.Divide(context.Background(), a, b)
		op = "/"
	case "max":
		resp, err = m.Max(context.Background(), a, b)
		op = "max"
	case "min":
		resp, err = m.Min(context.Background(), a, b)
		op = "min"
	case "multiply":
		resp, err = m.Multiply(context.Background(), a, b)
		op = "*"
	case "pow":
		resp, err = m.Pow(context.Background(), a, b)
		op = "^"
	case "subtract":
		resp, err = m.Subtract(context.Background(), a, b)
		op = "-"
	case "sum":
		resp, err = m.Sum(context.Background(), a, b)
		op = "+"
	default:
		fmt.Fprintf(os.Stderr, "error: invalid method %q\n", *method)
		os.Exit(1)
	}
	checkErr(err)
	fmt.Fprintf(os.Stdout, "%f %s %f = %f\n", a, op, b, resp)
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

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}
