
package fbits

import (
	"math"
)
const maxFloat64 = math.MaxFloat64

// minFbits returns the smaller of x or y.
//
// Special cases are:
//	Min(x, -Inf)   = Min(-Inf, x) = -Inf 
//	Min(x, NaN)    = Min(NaN, x) = NaN
//  Min(-Inf, NaN) = -Inf 
//	Min(-0, ±0)    = Min(±0, -0) = -0
func minFbits(x, y float64) float64 {
	switch  {                             
	case x < y:                           // Min(-Inf, y) 
		return x
	case y < x:                           // Min(x, -Inf) 
		return y
	case x == y:                          // true if x or/and y are not NaNs
		if x == 0 && math.Signbit(x) {    // check negative zero
			return x                      // x = -0
		}
		return y                           // y = -0, +/-Inf or anything but NaN
	case x < -maxFloat64:                  // Here x or/and y are NaNs   
		return x                           // x = -Inf
	case y < -maxFloat64:            
		return y                           // y = -Inf
	case x != x:                           // x != x is true if and only if x is NaN
		return x                           // x = NaN. IEEE 754 prefers to return/propagate 
	}                                      // the same NaN. A NaN has 51 bits for payload
 	return y                               // y = NaN                             
}

// NaN propagation example: https://play.golang.org/p/cRgm6-naFYb
// https://github.com/JuliaLang/julia/issues/7866
// https://github.com/JuliaLang/julia/issues/10729

// Go standard library minGo returns the smaller of x or y.
// https://golang.org/src/math/dim.go
// Compiler: cannot inline minGo: function too complex: cost 138 exceeds budget 80
// case math.IsInf makes this function not inlineable
//
// Special cases are:
//	Min(x, -Inf) = Min(-Inf, x) = -Inf
//	Min(x, NaN) = Min(NaN, x) = NaN
//	Min(-0, ±0) = Min(±0, -0) = -0
func minGo(x, y float64) float64 {
	// special cases
	switch {
	case math.IsInf(x, -1) || math.IsInf(y, -1): // this is here only because of contract  
		return math.Inf(-1)                      // minGo(-Inf, NaN) = -Inf
	case math.IsNaN(x) || math.IsNaN(y):
		return math.NaN()
	case x == 0 && x == y:
		if math.Signbit(x) {
			return x
		}
		return y
	}
	if x < y {
		return x
	}
	return y
}

// https://github.com/golang/go/issues/21913
func minGoProposal(x, y float64) float64 {
	const uvneginf = 0xFFF0000000000000
	const uvnan    = 0x7FF8000000000001  // quiet NaN
	switch {
	case x < y:                           
		return x
	case y < x:                           
		return y
	case x == y:
		// return -0 in preference to +0
		return math.Float64frombits(math.Float64bits(x) | math.Float64bits(y))
	}
	if inf := math.Float64frombits(uvneginf); x == inf || y == inf {
		return inf
	}
	return math.Float64frombits(uvnan)
}

