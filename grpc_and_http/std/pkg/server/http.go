package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	mathservice2 "github.com/jwenz723/mathserver/grpc_and_http/std/pkg/mathservice"
	"go.uber.org/zap"
	"net/http"
)

type httpServer struct {
	logger *zap.Logger
	router *mux.Router
	svc    mathservice2.Service
}

func NewHttpRouter(svc mathservice2.Service, logger *zap.Logger) *mux.Router {
	s := httpServer{
		logger: logger,
		router: mux.NewRouter(),
		svc:    svc,
	}
	s.routes()
	return s.router
}

func (s *httpServer) routes() {
	s.router.Methods("POST").PathPrefix("/").HandlerFunc(s.mathOpHandlerFunc())
}

// MathOpRequest collects the request parameters for the math methods.
type MathOpRequest struct {
	A, B float64
}

// MathOpResponse collects the response values for the math methods.
type MathOpResponse struct {
	V   float64 `json:"v"`
	Err error   `json:"-"`
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
			v, err = s.svc.Divide(r.Context(), req.A, req.B)
		case "/max":
			v, err = s.svc.Max(r.Context(), req.A, req.B)
		case "/min":
			v, err = s.svc.Min(r.Context(), req.A, req.B)
		case "/multiply":
			v, err = s.svc.Multiply(r.Context(), req.A, req.B)
		case "/pow":
			v, err = s.svc.Pow(r.Context(), req.A, req.B)
		case "/subtract":
			v, err = s.svc.Subtract(r.Context(), req.A, req.B)
		case "/sum":
			v, err = s.svc.Sum(r.Context(), req.A, req.B)
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
