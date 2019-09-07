package mathservice

import (
	"context"
	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

// LoggingMiddleware takes a logger as a dependency
// and returns a ServiceMiddleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw loggingMiddleware) Divide(ctx context.Context, a, b float64) (v float64, err error) {
	defer func() {
		mw.logger.Log("method", "Divide", "a", a, "b", b, "v", v, "err", err)
	}()
	return mw.next.Sum(ctx, a, b)
}

func (mw loggingMiddleware) Max(ctx context.Context, a, b float64) (v float64, err error) {
	defer func() {
		mw.logger.Log("method", "Max", "a", a, "b", b, "v", v, "err", err)
	}()
	return mw.next.Sum(ctx, a, b)
}

func (mw loggingMiddleware) Min(ctx context.Context, a, b float64) (v float64, err error) {
	defer func() {
		mw.logger.Log("method", "Min", "a", a, "b", b, "v", v, "err", err)
	}()
	return mw.next.Sum(ctx, a, b)
}

func (mw loggingMiddleware) Multiply(ctx context.Context, a, b float64) (v float64, err error) {
	defer func() {
		mw.logger.Log("method", "Multiply", "a", a, "b", b, "v", v, "err", err)
	}()
	return mw.next.Sum(ctx, a, b)
}

func (mw loggingMiddleware) Pow(ctx context.Context, a, b float64) (v float64, err error) {
	defer func() {
		mw.logger.Log("method", "Pow", "a", a, "b", b, "v", v, "err", err)
	}()
	return mw.next.Sum(ctx, a, b)
}

func (mw loggingMiddleware) Subtract(ctx context.Context, a, b float64) (v float64, err error) {
	defer func() {
		mw.logger.Log("method", "Subtract", "a", a, "b", b, "v", v, "err", err)
	}()
	return mw.next.Sum(ctx, a, b)
}

func (mw loggingMiddleware) Sum(ctx context.Context, a, b float64) (v float64, err error) {
	defer func() {
		mw.logger.Log("method", "Sum", "a", a, "b", b, "v", v, "err", err)
	}()
	return mw.next.Sum(ctx, a, b)
}