package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jwenz723/mathserver/std/pkg/mathservice"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type httpServer struct {
	logger *zap.Logger
	router *mux.Router
	svc    mathservice.Service
}

func NewHttpRouter(svc mathservice.Service, logger *zap.Logger) *mux.Router {
	s := httpServer{
		logger: logger,
		router: mux.NewRouter(),
		svc:    svc,
	}
	s.routes()
	s.addMiddlewares()
	return s.router
}

func (s *httpServer) routes() {
	s.router.Methods("POST").PathPrefix("/").HandlerFunc(s.mathOpHandlerFunc())
}

func (s *httpServer) addMiddlewares() {
	lmw := loggingMiddleware{
		logger: s.logger,
	}
	s.router.Use(lmw.Middleware)
}

// MathOpRequest collects the request parameters for the math methods.
type MathOpRequest struct {
	A, B float64
}

// MathOpResponse collects the response values for the math methods.
type MathOpResponse struct {
	V   float64   `json:"v"`
	Err error `json:"-"`
}

func (s *httpServer) mathOpHandlerFunc() http.HandlerFunc {
	fmt.Println("setting up math handler")
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var v float64
		switch r.URL.Path {
		case "/divide":
			v, err = s.svc.Divide(context.TODO(), req.A, req.B)
		case "/max":
			v, err = s.svc.Max(context.TODO(), req.A, req.B)
		case "/min":
			v, err = s.svc.Min(context.TODO(), req.A, req.B)
		case "/multiply":
			v, err = s.svc.Multiply(context.TODO(), req.A, req.B)
		case "/pow":
			v, err = s.svc.Pow(context.TODO(), req.A, req.B)
		case "/subtract":
			v, err = s.svc.Subtract(context.TODO(), req.A, req.B)
		case "/sum":
			v, err = s.svc.Sum(context.TODO(), req.A, req.B)
		}

		writeResponse(w, v, err)
	}
}

func decodeRequest(r *http.Request) (MathOpRequest, error) {
	var req MathOpRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func writeResponse(w http.ResponseWriter, v float64, err error) {
	resp := MathOpResponse{
		V:   v,
		Err: err,
	}

	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(js)
}

type loggingMiddleware struct {
	logger *zap.Logger
}

func (lmw *loggingMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(begin time.Time) {
			lmw.logger.Info(r.RequestURI,
				zap.Duration("duration", time.Since(begin)))
		}(time.Now())
		next.ServeHTTP(w, r)
	})
}