package mathservice

import (
	"context"
	"errors"
	"math"
)

// Service describes a service that adds things together.
type Service interface {
	// Divide two integers, a/b
	Divide(ctx context.Context, a, b float64) (float64, error)
	// Max two integers, returns the greater value of a and b
	Max(ctx context.Context, a, b float64) (float64, error)
	// Min two integers, returns the lesser value of a and b
	Min(ctx context.Context, a, b float64) (float64, error)
	// Multiply two integers, a*b
	Multiply(ctx context.Context, a, b float64) (float64, error)
	// Pow two integers, a^b
	Pow(ctx context.Context, a, b float64) (float64, error)
	// Subtract two integers, a-b
	Subtract(ctx context.Context, a, b float64) (float64, error)
	// Sums two integers. a+b
	Sum(ctx context.Context, a, b float64) (float64, error)
}

var (
	ErrDivideByZero = errors.New("can't divide by zero")
	ErrNoMax = errors.New("no maximum value, a and b are the same")
	ErrNoMin = errors.New("no minimum value, a and b are the same")
)

// NewBasicService returns a naïve, stateless implementation of Service.
func NewBasicService() basicService {
	return basicService{}
}

type basicService struct{}

func (s basicService) Divide(ctx context.Context, a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivideByZero
	}
	return a/b, nil
}

func (s basicService) Max(ctx context.Context, a, b float64) (float64, error) {
	if a == b {
		return 0, ErrNoMax
	}
	return math.Max(a, b), nil
}

func (s basicService) Min(ctx context.Context, a, b float64) (float64, error) {
	if a == b {
		return 0, ErrNoMin
	}
	return math.Min(a, b), nil
}

func (s basicService) Multiply(ctx context.Context, a, b float64) (float64, error) {
	return a*b, nil
}

func (s basicService) Pow(ctx context.Context, a, b float64) (float64, error) {
	return math.Pow(a, b), nil
}

func (s basicService) Subtract(ctx context.Context, a, b float64) (float64, error) {
	return a-b, nil
}

func (s basicService) Sum(ctx context.Context, a, b float64) (float64, error) {
	return a + b, nil
}
