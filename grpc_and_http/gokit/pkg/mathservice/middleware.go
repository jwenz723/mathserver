package mathservice

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"time"
)

type Middleware func(Service) Service

// ObservabilityMiddleware implements both logging and prometheus metrics for each Service method
func ObservabilityMiddleware(duration metrics.Histogram, logger log.Logger) Middleware {
	return func(next Service) Service {
		return observabilityMiddleware{duration,logger, next}
	}
}

type observabilityMiddleware struct {
	duration metrics.Histogram
	logger   log.Logger
	next     Service
}

func (mw observabilityMiddleware) observeMethodExecution(ctx context.Context, method string, a, b, v float64, begin time.Time, err error) {
	duration := time.Since(begin)

	mw.logger.Log("msg", "method executed",
		"method", method,
		"a", a,
		"b", b,
		"v", v,
		"duration", duration,
		"err", err)
	mw.duration.With("method", method, "success", fmt.Sprint(err == nil)).Observe(duration.Seconds())
}

func (mw observabilityMiddleware) Divide(ctx context.Context, a, b float64) (v float64, err error) {
	defer func(begin time.Time) {
		m := "Divide"
		mw.observeMethodExecution(ctx, m, a, b, v, begin, err)
	}(time.Now())
	return mw.next.Divide(ctx, a, b)
}

func (mw observabilityMiddleware) Max(ctx context.Context, a, b float64) (v float64, err error) {
	defer func(begin time.Time) {
		m := "Max"
		mw.observeMethodExecution(ctx, m, a, b, v, begin, err)
	}(time.Now())
	return mw.next.Max(ctx, a, b)
}

func (mw observabilityMiddleware) Min(ctx context.Context, a, b float64) (v float64, err error) {
	defer func(begin time.Time) {
		m := "Min"
		mw.observeMethodExecution(ctx, m, a, b, v, begin, err)
	}(time.Now())
	return mw.next.Min(ctx, a, b)
}

func (mw observabilityMiddleware) Multiply(ctx context.Context, a, b float64) (v float64, err error) {
	defer func(begin time.Time) {
		m := "Multiply"
		mw.observeMethodExecution(ctx, m, a, b, v, begin, err)
	}(time.Now())
	return mw.next.Multiply(ctx, a, b)
}

func (mw observabilityMiddleware) Pow(ctx context.Context, a, b float64) (v float64, err error) {
	defer func(begin time.Time) {
		m := "Pow"
		mw.observeMethodExecution(ctx, m, a, b, v, begin, err)
	}(time.Now())
	return mw.next.Pow(ctx, a, b)
}

func (mw observabilityMiddleware) Subtract(ctx context.Context, a, b float64) (v float64, err error) {
	defer func(begin time.Time) {
		m := "Subtract"
		mw.observeMethodExecution(ctx, m, a, b, v, begin, err)
	}(time.Now())
	return mw.next.Subtract(ctx, a, b)
}

func (mw observabilityMiddleware) Sum(ctx context.Context, a, b float64) (v float64, err error) {
	defer func(begin time.Time) {
		m := "Sum"
		mw.observeMethodExecution(ctx, m, a, b, v, begin, err)
	}(time.Now())
	return mw.next.Sum(ctx, a, b)
}