package multicell

import (
	"math"
)

// activation functions

// LeCun-inspired arctan function
func LCatan(omega float64) func(float64) float64 {
	b := omega / Sqrt3
	return func(x float64) float64 {
		return 6.0 * math.Atan(b*x) / math.Pi
	}
}

// LeCun tanh (without multiplying by Sqrt3).
func LCtanh(omega float64) func(float64) float64 {
	b := math.Atanh(1.0/Sqrt3) * omega
	return func(x float64) float64 {
		return math.Tanh(b * x)
	}
}

// Scaled tanh
func SCtanh(omega float64) func(float64) float64 {
	b := math.Atanh(0.99) * omega
	return func(x float64) float64 {
		return math.Tanh(b * x)
	}
}

// degree-1 continuous step function
func CStep1(omega float64) func(float64) float64 {
	return func(x float64) float64 {
		return max(-1.0, min(1.0, omega*x))
	}
}

// smooth degree-3 step function
func SStep3(omega float64) func(float64) float64 {
	return func(x float64) float64 {
		t := omega * x
		if t > 1.0 {
			return 1.0
		} else if t < -1.0 {
			return -1.0
		} else {
			return 0.5 * t * (3.0 - t*t)
		}
	}
}

// smooth degree-5 step function
func SStep5(omega float64) func(float64) float64 {
	return func(x float64) float64 {
		t := omega * x
		if t > 1.0 {
			return 1.0
		} else if t < -1.0 {
			return -1.0
		} else {
			return 0.5 * t * t * t * (5.0 - 3.0*t*t)
		}
	}
}

// abrupt three-step function
func AStep3(omega float64) func(float64) float64 {
	return func(x float64) float64 {
		t := omega * x
		if t >= 1.0 {
			return 1.0
		} else if t <= -1.0 {
			return -1.0
		} else {
			return 0.0
		}
	}
}

// Simple arctan
func Atan(omega float64) func(float64) float64 {
	return func(x float64) float64 {
		return math.Atan(omega * x)
	}
}

func Tanh(omega float64) func(float64) float64 {
	return func(x float64) float64 {
		return math.Tanh(omega * x)
	}
}
