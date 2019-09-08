package server

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jwenz723/mathserver/std/pkg/mathservice"
	"net/http"
)

type httpServer struct {
	router *mux.Router
	svc    mathservice.Service
}

func NewHttpServer(svc mathservice.Service) httpServer {
	s := httpServer{
		router: mux.NewRouter(),
		svc:    svc,
	}
	s.routes()
	return s
}

func (s *httpServer) Router() *mux.Router {
	return s.router
}

func (s *httpServer) routes() {
	s.router.Methods("POST").Path("/divide").HandlerFunc(s.divideHandlerFunc())
	s.router.Methods("POST").Path("/max").HandlerFunc(s.maxHandlerFunc())
	s.router.Methods("POST").Path("/min").HandlerFunc(s.minHandlerFunc())
	s.router.Methods("POST").Path("/multiply").HandlerFunc(s.multiplyHandlerFunc())
	s.router.Methods("POST").Path("/pow").HandlerFunc(s.powHandlerFunc())
	s.router.Methods("POST").Path("/subtract").HandlerFunc(s.subtractHandlerFunc())
	s.router.Methods("POST").Path("/sum").HandlerFunc(s.sumHandlerFunc())
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

func (s *httpServer) divideHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		v, err := s.svc.Divide(context.TODO(), req.A, req.B)
		writeResponse(w, v, err)
	}
}

func (s *httpServer) maxHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		v, err := s.svc.Max(context.TODO(), req.A, req.B)
		writeResponse(w, v, err)
	}
}

func (s *httpServer) minHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		v, err := s.svc.Min(context.TODO(), req.A, req.B)
		writeResponse(w, v, err)
	}
}

func (s *httpServer) multiplyHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		v, err := s.svc.Multiply(context.TODO(), req.A, req.B)
		writeResponse(w, v, err)
	}
}

func (s *httpServer) powHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		v, err := s.svc.Pow(context.TODO(), req.A, req.B)
		writeResponse(w, v, err)
	}
}

func (s *httpServer) subtractHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		v, err := s.svc.Subtract(context.TODO(), req.A, req.B)
		writeResponse(w, v, err)
	}
}

func (s *httpServer) sumHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		v, err := s.svc.Sum(context.TODO(), req.A, req.B)
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