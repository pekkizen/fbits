package fbits

import (
	"math"
	"math/bits"
)

const (
	signbit   = 1<<63
	posInf    = 0x7ff0000000000000
	maxUint64 = 1<<64 - 1
)

// UlpsBetween returns the distance between x and y in ulpS.
// 
// The distance in ulps is the number of float64's between x and y - 1.
// Special cases:
// UlpsBetween(+/-Inf, +/-MaxFloat64) = 1
// UlpsBetween(+/-Inf, +/-Inf)        = 0
// UlpsBetween(-Inf, +Inf)            = maxUint64 - 2^53 + 1 (18437736874454810624)
// UlpsBetween(x, NaN)                = maxUint64 
// UlpsBetween(-0, 0)                 = 0
// UlpsBetween(-0, 2^-1074)           = 1
// UlpsBetween(0, -2^-1074)           = 1
// 
func UlpsBetween(x, y float64) (u uint64) {
	k := math.Float64bits(x)
	n := math.Float64bits(y)
	signdiff := k ^ n >= signbit
	k &^= signbit 
	n &^= signbit 
	switch {
	case k > posInf || n > posInf:  // NaNs 
		u = maxUint64	
	case signdiff:
		u = n + k
	case n > k:
		u = n - k
	default:
		u = k - n
	}
	return
}

// Adjacent returns true, if x and y are Adjacent floats.
// 
// Adjacent(x, y) is a faster equivalent to UlpsBetween(x, y) == 1.
// For x > 0 and y > x  Adjacent(x, y) == (math.Nextafter(x, y) == y)
// Special and other cases:
// Adjacent(+Inf, +MaxFloat64) = true
// Adjacent(-Inf, -MaxFloat64) = true
// Adjacent(x, NaN)            = false
// Adjacent(-0, -2^-1074)      = true
// Adjacent(0, 2^-1074)        = true
// Adjacent(-0, 0)             = false
// Adjacent(0, -2^-1074)       = false,  A special case failure and
// Adjacent(-0, 2^-1074)       = false   also this. Only failures.
// 2^-1074 is the smallest nonzero float64.
// 
func Adjacent(x, y float64) bool {
	d := int64(math.Float64bits(x) - math.Float64bits(y))
	return d == 1 || d == -1 
	// or
	// d := math.Float64bits(x) - math.Float64bits(y)
	// return d == 1 || d == maxUint64               
}

// AdjacentFP returns true, if x and y are finite and adjacent floats.
// 
// Only floating-point operations are used.
// This is ~35% (0.3 ns) slower than Adjacent but doesn't fail at zero.
// Special cases different from func Adjacent:
// AdjacentFP(+Inf, +MaxFloat64) = false
// AdjacentFP(-Inf, -MaxFloat64) = false
// AdjacentFP(0, -2^-1074)       = true
// AdjacentFP(-0, 2^-1074)       = true
// 
func AdjacentFP(x, y float64) bool {
	if x == y {
		return false
	}
	mean := x/2 + y/2                // this avoids overflowing x + y to Inf
	if mean != x && mean != y {      // NaNs
		return false
	}
	return -math.MaxFloat64 <= mean && mean <= math.MaxFloat64  // Infs 
}

// Ulp returns the ulp of x as a positive float64. 
// 
// A ulp returned is the distance to the next float64 away from zero.
// If x is a power on two, ulp(x) towards zero is ulp(x)/2 away from zero. 
// All ulps are integer powers of two.
// Special cases:
// Ulp(+/-Inf) = +Inf 
// Ulp(NaN)    = NaN
// 
func Ulp(x float64) float64 {
	u := math.Float64bits(x) &^ signbit
	exp := u >> 52
	switch {
	case exp == 0x7ff:       // Infs and NaNs, returns abs(x)
	case exp > 52:
		u = (exp - 52) << 52
	case exp > 1:
		u = 1 << (exp - 1)
	default:
		u = 1                // x < 2^-2021, Ulp = 2^-1074
	}
	return math.Float64frombits(u)  
}

// LogUlp returns log2(Ulp(x)) as an int, Ulp(x) = 2^LogUlp(x).
// Special cases:
// LogUlp(+/-Inf) = 1024    (2^1024 = +Inf)
// LogUlp(NaN)    = 1024
// 
func LogUlp(x float64) (exp int) {
	exp = int(math.Float64bits(x) &^ signbit >> 52)
	switch {
	case exp == 0x7ff:           // Infs and NaNs
		exp = 1024
	case exp > 0:
		exp -= (1023 + 52)
	default:
		exp = -1074        
	}
	return  
}

// UlpFP returns the ulp of x for abs(x) > 0x1p-1022.
// 
// A ulp is calculated as a difference towards zero.
// Special cases:
// UlpFP(+/-Inf) = NaN !!!
// UlpFP(NaN)    = NaN
// For abs(x) <= 0x1p-1022 UlpFP fails and returns 0.
// 
func UlpFP(x float64) float64 {
	y := x - NextToZeroFP(x)         // Inf - Inf = NaN
	if y < 0 { return -y }      
	return y
}

// Log2 returns base 2 logaritm of abs(x) as a rounded towards zero int. 
// For normal floats it is the same as the unbiased IEEE 754 exponent.
// 
// If log2(x) = n,  2^n <= abs(x) < 2^n+1.
// Special cases:
// Log2(+/-Inf)    = 1024
// Log2(NaN)       = 1024
// Log2(-x)        = log2(x)
// 
func Log2(x float64) int {
	u := math.Float64bits(x) &^ signbit
	exp := int(u >> 52)
	if exp == 0 {                        // x is subnormal 
		return bits.Len64(u) - 1075      // Len64(u=2^n) = n + 1, n = 0 - 51
	}
	return exp - 1023                    // x is normal, Inf or NaN
}

// IsPowerOfTwo returns true if float64 x is an integer power of two.
// 
// Cases of interest:
// IsPowerOfTwo(1)      = true
// IsPowerOfTwo(0)      = false
// IsPowerOfTwo(x)      = false for x < 0
// IsPowerOfTwo(+/-Inf) = false
// IsPowerOfTwo(NaN)    = false
// 
func IsPowerOfTwo(x float64) bool {
	s := math.Float64bits(x) 
	e := s >> 52                  // sign bit + 11 exponent bits                                                
	s <<= 12                      // 52 significand bits + zeros 

	return s & (s - 1) == 0 && (s > 0) != (e > 0) && e < 0x7ff
}

// A float64 value x is an integer power of two if and only if the following 
// three conditions are met:
//     s & (s - 1) == 0     -> significand is zero or power of two
//     (s > 0) != (e > 0)   -> significand or exponent is zero, but not both
//     e < 0x7ff            -> x is not +/-Inf, NaN or negative
// 
// Above e > 0 is true for all negative x, but the last condition drops these out.
// s <<= 12 is here faster than masking s &= (1<<52)-1 ? 
// The position of the significand bits is not relevant here.

// IsPowerOfTwoFP returns true if float64 x is an integer power of two.
// 
// https://stackoverflow.com/questions/27566187/code-for-check-if-double-is-a-power-of-2-without-bit-manipulation-in-c
// This is without bit operations and works, but is over 60% slower than IsPowerOfTwo
// Formula x > 0 && 0x1.0p-51/x * x - 0x1.0p-51 == 0 without FMA doesn't work.
// 
func IsPowerOfTwoFP(x float64) bool { 
	return x > 0 && math.FMA(0x1.0p-51/x, x, -0x1.0p-51) == 0 
}

// Java DoubleUtils.IsPowerOfTwo(double x) from com.google.common.math.
// https://www.codota.com/code/java/classes/com.romainpiel.guava.math.DoubleUtils
// This checks twice IsFinite(x).
// 
// public static boolean IsPowerOfTwo(double x) {
//  return x > 0.0 && IsFinite(x) && LongMath.IsPowerOfTwo(getSignificand(x));
// }
// public static boolean IsPowerOfTwo(long x) {
//     return x > 0 & (x & (x - 1)) == 0;
// }
// DoubleUtils.getSignificand(...)
// static long getSignificand(double d) {
//  checkArgument(IsFinite(d), "not a normal value");
//  int exponent = getExponent(d);
//  long bits = doubleToRawLongBits(d);
//  bits &= SIGNIFICAND_MASK;
//  return (exponent == MIN_EXPONENT - 1)
//    ? bits << 1
//    : bits | IMPLICIT_BIT;
// }
	
// IsPowerOfTwoJava implements DoubleUtils.IsPowerOfTwo. 
// The bare cleaned algorithm with the same functionality.
// This small and simple function is still over ~50% slower than IsPowerOfTwo above.
// 
func IsPowerOfTwoJava(x float64) bool {
	bits := math.Float64bits(x)     // bits = doubleToRawLongBits(x)
	exp := bits >> 52               // exponent = getExponent(x) 
	bits &= 1<<52 - 1               // bits &= SIGNIFICAND_MASK
	if exp > 0  {                   // not (exponent == MIN_EXPONENT - 1)
		bits |= 1<<52               // bits | IMPLICIT_BIT (this is the point in the algorithm)
	}
	return bits & (bits - 1) == 0 && bits > 0 && exp < 0x7ff // IsPowerOfTwo(bits) & IsFinite(x) & x > 0.
}

// IsInf returns true if x is +/-Inf.
func IsInf(x float64) bool {
	return math.Float64bits(x) &^ signbit == posInf 
}

// IsFinite returns true if x is not +/-Inf of NaN.
func IsFinite(x float64) bool {
	return math.Float64bits(x) &^ signbit < posInf 
}

// NextToZero returns the next float64 after x towards zero.
// 
// NextToZero(x) is equivalent to math.Nextafter(x, 0)
// In a benchmark loop it is faster than math.Abs.
// Special cases:
// NextToZero(+/-Inf)   = +/-MaxFloat64
// NextToZero(NaN)      = NaN
// NextToZero(0)        = 0
// NextToZero(-0)       = -0  
// NextToZero(2^1074)   = 0
// NextToZero(-2^-1074) = -0
// 
func NextToZero(x float64) float64 {
	if y := NextToZeroFP(x); y != x {    // Nan != NaN is true            
		return y                               
	}
	if x == 0 { return x } 
	return math.Float64frombits(math.Float64bits(x) - 1)
}

// NextToZeroFP is equivalent to NextToZero for abs(x) > 2^-1022. 
// In (0, 2^-1022] NextToZeroFP(x) fails and returns x.
// The constant 1 - 0x1p-53 converts exactly to the next float64 from 1 towards zero. 
// Float64bits(1):           3FF0000000000000.
// Float64bits(1 - 0x1p-53): 3FEFFFFFFFFFFFFF.
// Float64bits(1 + 0x1p-53): 3FF0000000000000.
// Float64bits(1 + 0x1p-52): 3FF0000000000001.
// 
func NextToZeroFP(x float64) float64 {
	return x * (1 - 0x1p-53)
}

// NextFromZero returns the next float64 after x away from zero.
// 
// NextFromZero(+/-abs(x) is equivalent to math.Nextafter(+/-abs(x), math.Inf(1/-1)).
// Special cases:
// NextFromZero(NaN)           = NaN
// NextFromZero(+/-Inf)        = +/-Inf 
// NextFromZero(+/-MaxFloat64) = +/-Inf  
// NextFromZero(0)             = 2^-1074 
// NextFromZero(-0)            = -2^-1074
// 
func NextFromZero(x float64) float64 {
	if y := NextFromZeroFP(x); x != y {    // Nan != NaN is true     
		return y                            
	}
	if IsInf(x) { return x }
	return math.Float64frombits(math.Float64bits(x) + 1)
}

// NextFromZeroFP is equivalent to NextFromZero for abs(x) >= 2^-1019. 
// For abs(x) < 2^-1019 NextFromZeroFP(x) fails and returns x.
// The evident formula x * (1 + 0x1p-52) doesn't work.
// 
func NextFromZeroFP(x float64) float64 {
	return x + x * 0x1.25p-53            // Inf + Inf = Inf
}

// RandomFloat64 returns a random float64's from [-MaxFloat64, MaxFloat64].
// Every float has an equal probability 1 / (2^64 - 2^53). 
// 
func RandomFloat64(state *uint64) float64 {
	return FiniteFloat64frombits(Splitmix(state))
}

// FiniteFloat64frombits returns math.Float64frombits(u), except for Infs
// and NaNs the exponent (0x7ff) is replaced by u mod 0x7ff (0 - 2046).
// 
func FiniteFloat64frombits(u uint64) float64 {
	if u &^ signbit >= posInf {  
		u = u &^ posInf | (u % 0x7ff) << 52
	}
	return math.Float64frombits(u)
}

// RandomFloat64RS uses resampling in the case of Inf or Nan.
// This gives a provable unbiased distribution of floats assuming that the
// random  number generator Splitmix gives unbiased uniform distribution 
// of Uint64's. The modulo operation u % 0x7ff above has a _very_ small bias 
// for small numbers.
// The Go compiler can inline functions with goto loops, but not with for loops.
func RandomFloat64RS(state *uint64) float64 {
	again:
	u := Splitmix(state)
	if u & posInf == posInf {  
		goto again                         // resample, 1/2048 of cases
	}
	return math.Float64frombits(u)
}

// Splitmix is a 64-bit state SplitMix64 pseudo-random number generator
// from http://prng.di.unimi.it/splitmix64.c .
// The pointer parameter state is used as the random state.
// Uint64 *state can be initialized to any value.
// 
func Splitmix(state *uint64) uint64 {
	*state += 0x9e3779b97f4a7c15 
	z := *state
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}