package fbits

import (
	"math"
	"testing"
)

var usink uint64
var fsink float64
var isink int
var bsink bool

func abs(x float64) float64 {
	if x > 0 {
		return x
	}
	return -x
}
func BenchmarkUlpsBetween(b *testing.B) {
	var u uint64
	f2 := 1.0
	for n := 0; n < b.N; n++ {
		u = UlpsBetween(float64(n), f2)
	}
	usink = u
}

func BenchmarkAdjacent(b *testing.B) {
	var is bool
	f2 := 1.0
	for n := 0; n < b.N; n++ {
		is = Adjacent(float64(n), f2)
	}
	bsink = is
}
func BenchmarkAdjacentFP(b *testing.B) {
	var is bool
	f2 := 1.0
	for n := 0; n < b.N; n++ {
		is = AdjacentFP(float64(n), f2)
	}
	bsink = is
}
func BenchmarkIsPowerOfTwo(b *testing.B) {
	var is bool
	for n := 0; n < b.N; n++ {
		is = IsPowerOfTwo(float64(n))
	}
	bsink = is
}

func BenchmarkIsPowerOfTwoFP(b *testing.B) {
	var is bool
	for n := 0; n < b.N; n++ {
		is = IsPowerOfTwoFP(float64(n))
	}
	bsink = is
}
func BenchmarkIsPowerOfTwoJava(b *testing.B) {
	var is bool
	for n := 0; n < b.N; n++ {
		is = IsPowerOfTwoJava(float64(n))
	}
	bsink = is
}

func BenchmarkUlp(b *testing.B) {
	var y float64
	for n := 0; n < b.N; n++ {
		y = Ulp(float64(n))
	}
	fsink = y
}
func BenchmarkUlpB(b *testing.B) {
	var y float64
	for n := 0; n < b.N; n++ {
		y = UlpB(float64(n))
	}
	fsink = y
}

func BenchmarkLogUlp(b *testing.B) {
	var u int
	for n := 0; n < b.N; n++ {
		u = LogUlp(float64(n))
		// u = Log2(Ulp(float64(n)))
	}
	isink = u
}
func BenchmarkLog2(b *testing.B) {
	var u int
	for n := 0; n < b.N; n++ {
		u = Log2(float64(n))
	}
	isink = u
}

func BenchmarkNextToZero(b *testing.B) {
	var y float64
	for n := 0; n < b.N; n++ {
		y = NextToZero(float64(n))
		// y = math.Abs(float64(n))
	}
	fsink = y
}
func BenchmarkNextToZeroFP(b *testing.B) {
	var y float64
	for n := 0; n < b.N; n++ {
		y = NextToZeroFP(float64(n))
	}
	fsink = y
}

func BenchmarkNextFromZero(b *testing.B) {
	var y float64
	for n := 0; n < b.N; n++ {
		y = NextFromZero(float64(n))
	}
	fsink = y
}
func BenchmarkMathNextafter(b *testing.B) {
	var y float64
	for n := 0; n < b.N; n++ {
		y = math.Nextafter(float64(n), math.MaxFloat64)
	}
	fsink = y
}

func BenchmarkRandomFloat64(b *testing.B) {
	var y float64
	state := uint64(1)
	for n := 0; n < b.N; n++ {
		// y = FiniteFloat64frombits(Splitmix(&state))
		y = RandomFloat64(&state)
	}
	fsink = y
}
// ------------------------------------------------------------- Tests
func TestRandomFloat64(t *testing.T) {
	const rounds int = 1e8*2
	min, max, prop := 999.0, 0.0, 0.0
	state := uint64(1)
	for i := 0; i < rounds; i++ {
		f := RandomFloat64(&state)
		if !IsFinite(f) {
			t.Fatalf("Inf or NaN   %16X", math.Float64bits(f))
		}
		if f < 0 {
			f = -f
		}
		if f < min {
			min = f
		}
		if f > max {
			max = f
		}
		if f >= 0 && f < 1 {
			prop++
		}
	}
	expected := 100.0*1023 / 2047
	t.Logf("Min   %v", min)
	t.Logf("Max   %v", max)
	prop = 100*prop / float64(rounds)
	t.Logf("0-1   %v (%2.6f)", prop, expected)
	if abs(prop - expected) > 0.001 {
		t.Fatalf("Proportion 0-1 failed  %v", prop - expected)
	}
}

func TestIsInf(t *testing.T) {
	t.Logf("+Inf           %v", IsInf(math.Inf(1)))
	t.Logf("-Inf           %v", IsInf(math.Inf(-1)))
	t.Logf("NaN            %v", IsInf(math.NaN()))
	t.Logf("MaxFloat64     %v", IsInf(math.MaxFloat64))
}
func TestIsFinite(t *testing.T) {
	t.Logf("+Inf           %v", IsFinite(math.Inf(1)))
	t.Logf("-Inf           %v", IsFinite(math.Inf(-1)))
	t.Logf("NaN            %v", IsFinite(math.NaN()))
	t.Logf("MaxFloat64     %v", IsFinite(math.MaxFloat64))
	t.Logf("-MaxFloat64    %v", IsFinite(-math.MaxFloat64))
}
func TestInfNaN(t *testing.T) {
	f1 := math.Inf(1)
	f2 := f1 + 1
	f3 := math.Float64frombits(math.Float64bits(f1) + 1)
	f4 := f3 + 1
	t.Logf("f1  %X %v", math.Float64bits(f1), f1)
	t.Logf("f2  %X %v", math.Float64bits(f2), f2)
	t.Logf("f3  %X %v !!!", math.Float64bits(f3), f3)
	t.Logf("-f3 %X %v\n", math.Float64bits(-f3), -f3)
	t.Logf("f4  %X %v", math.Float64bits(f4), f4)
	t.Logf("f3 != f3 %v", f3 != f3)
}

func TestNextToZero(t *testing.T) {
    const rounds int = 1e8*3
	state := uint64(1)
	zero, inf, nan, min := 0.0, math.Inf(1), math.NaN(), 0x1p-1074

	t.Logf("zero     %v", NextToZero(zero))
	t.Logf("-zero    %v", NextToZero(-zero))
	t.Logf("min      %v", NextToZero(min))
	t.Logf("-min     %v", NextToZero(-min))
	t.Logf("+inf     %v", NextToZero(inf))
	t.Logf("-inf     %v", NextToZero(-inf))
	t.Logf("NaN      %v", NextToZero(nan))
	for i := 0; i < rounds; i++ {
		f1 := RandomFloat64(&state) 
		f2 := NextToZero(f1)
		f3 := math.Nextafter(f1, 0)
		if f2 != f3 {
			t.Logf("Nextafter i   %d", i)
			t.Logf("F1  %v", f1)
			t.Logf("F2  %v", f2)
			t.Fatalf("F3  %v", f3)
		}
		Ulps := UlpsBetween(f1, f2)
		if Ulps != 1 && f1 > 0 && !IsInf(f1) {	
			t.Logf("Ulps %v", Ulps)
			t.Logf("i    %d", i)
			t.Logf("F1   %v", f1)
			t.Fatalf("F2   %v", f2)
		}
	}
}

func TestNextToZeroFP(t *testing.T) {
    const rounds int = 1e8*3
	state := uint64(1)
	for i := 0; i < rounds; i++ {
		f1 := RandomFloat64(&state) 
		f2 := NextToZeroFP(f1)
		f3 := math.Nextafter(f1, 0)
		if f2 != f3 && f1 > 0x1p-1022 {
			t.Logf("Nextafter i=   %d", i)
			t.Logf("F1  %v", f1)
			t.Logf("F2  %v", f2)
			t.Fatalf("F3  %v", f3)
		}
	}
}

func TestNextFromZero(t *testing.T) {
    const rounds int = 1e8
	zero, max, inf, nan := 0.0, math.MaxFloat64, math.Inf(1), math.NaN()
	t.Logf("zero     %v", NextFromZero(zero))
	t.Logf("-zero    %v", NextFromZero(-zero))
	t.Logf("max      %v", NextFromZero(max))
	t.Logf("-max     %v", NextFromZero(-max))
	t.Logf("+inf     %v", NextFromZero(inf))
	t.Logf("-inf     %v", NextFromZero(-inf))
	f := NextFromZero(nan)
	t.Logf("NaN      %v %16X", f, math.Float64bits(f))
	state := uint64(1)
	for i := 0; i < rounds; i++ {
		f1 := RandomFloat64(&state) 
		f2 := NextFromZero(f1)
		d := inf
		if f1 < 0 { d = -d }
		f3 := math.Nextafter(f1, d)
		if f2 != f3 {
			t.Logf("i   %d", i)
			t.Logf("F1  %v", f1)
			t.Logf("F2  %v", f2)
			t.Logf("F3  %v", f3)
			
			t.Logf("F1  %X" , math.Float64bits(f1))
			t.Logf("F2  %X" , math.Float64bits(f2))
			t.Fatalf("F3  %X" , math.Float64bits(f3))
		
		}
	}
}

func TestNextFromZeroFP(t *testing.T) {
    const rounds int = 1e8
	state := uint64(1)
	inf :=  math.Inf(1)
	for i := 0; i < rounds; i++ {
		f1 := RandomFloat64(&state) 
		if f1 < 0 { f1 = -f1 }
		f2 := NextFromZeroFP(f1)
		f3 := math.Nextafter(f1, inf)
		if f2 != f3 && f1 >= 0x1p-1019 {
			t.Logf("i   %d", i)
			t.Logf("F1  %v", f1)
			t.Logf("F2  %v", f2)
			t.Logf("F3  %v", f3)
			
			t.Logf("F1  %X" , math.Float64bits(f1))
			t.Logf("F2  %X" , math.Float64bits(f2))
			t.Fatalf("F3  %X" , math.Float64bits(f3))
		
		}
	}
}

func TestAdjacent(t *testing.T) {
    const rounds int = 1e8
    state := uint64(1)
    // Adjacent := AdjacentFP
	zero, min, max, inf, nan := 0.0, 0x1p-1074, math.MaxFloat64, math.Inf(1), math.NaN()

	t.Logf("zero min    %v", Adjacent(zero, min))
	t.Logf("-zero -min  %v", Adjacent(-zero, -min))
	t.Logf("-zero min   %v !!!", Adjacent(-zero, min))
	t.Logf("zero -min   %v !!!", Adjacent(zero, -min))
	t.Logf("zero -zero  %v", Adjacent(zero, -zero))
	t.Logf("max inf     %v", Adjacent(max, inf))
	t.Logf("-max -inf   %v", Adjacent(-max, -inf))
	t.Logf("NaN inf     %v", Adjacent(nan, inf))
	t.Logf("-max zero   %v", Adjacent(-max, 0))
	for i := 0; i < rounds; i++ {
		f1 := RandomFloat64(&state)
		// f2 := NextFromZero(f1)
		f2 := NextToZero(f1)
		if i & 15 == 0 {
			f2 *= 2
		}
		Ulps := UlpsBetween(f1, f2)
		if Adjacent(f1, f2) != (Ulps == 1) {
			t.Logf("Ulps %v", Ulps)
			t.Logf("i    %d", i)
			t.Logf("F1   %v", f1)
			t.Logf("F2   %v", f2)
			t.Fatalf("     %v %v", Adjacent(f1, f2), AdjacentFP(f1, f2))
		}
	}
}
func TestAdjacentFP(t *testing.T) {
    const rounds int = 1e8
    state := uint64(1)
    Adjacent := AdjacentFP
	zero, min, max, inf, nan := 0.0, 0x1p-1074, math.MaxFloat64, math.Inf(1), math.NaN()

	t.Logf("zero min    %v", Adjacent(zero, min))
	t.Logf("-zero -min  %v", Adjacent(-zero, -min))
	t.Logf("-zero min   %v", Adjacent(-zero, min))
	t.Logf("zero -min   %v", Adjacent(zero, -min))
	t.Logf("zero -zero  %v", Adjacent(zero, -zero))
	t.Logf("max inf     %v", Adjacent(max, inf))
	t.Logf("-max -inf   %v", Adjacent(-max, -inf))
	t.Logf("NaN inf     %v", Adjacent(nan, inf))
	t.Logf("-max zero   %v", Adjacent(-max, 0))
	for i := 0; i < rounds; i++ {
		f1 := RandomFloat64(&state)
		// f2 := NextFromZero(f1)
		f2 := NextToZero(f1)
		if i & 15 == 0 {
			f2 *= 2
		}
		Ulps := UlpsBetween(f1, f2)
		if Adjacent(f1, f2) != (Ulps == 1) {
			t.Logf("Ulps %v", Ulps)
			t.Logf("i    %d", i)
			t.Logf("F1   %v", f1)
			t.Logf("F2   %v", f2)
			t.Fatalf("     %v %v", Adjacent(f1, f2), AdjacentFP(f1, f2))
		}
	}
}

func TestUlpsBetween(t *testing.T) {
	const rounds int = 1e8
	log2 := math.Log2
	zero, max, inf, nan, min := 0.0, math.MaxFloat64, math.Inf(1), math.NaN(), 0x1p-1074

	t.Logf("+Inf 0     %v", UlpsBetween(inf, 0))
	t.Logf("-Inf 0     %v", UlpsBetween(-inf, 0))
	t.Logf("NaN 0      %v", UlpsBetween(nan, 0))
	t.Logf("+Inf +Inf  %v", UlpsBetween(inf, inf))
	t.Logf("-Inf +Inf  %v", UlpsBetween(-inf, inf))
	t.Logf("-max max   %v (log2)", log2(float64(UlpsBetween(-max, max))))
	t.Logf("-min min   %v", UlpsBetween(-min, min))
	t.Logf("-zero zero %v", float64(UlpsBetween(-zero, zero)))
	t.Logf("-zero min  %v", float64(UlpsBetween(-zero, min)))
	t.Logf("zero min   %v", UlpsBetween(zero, min))
	t.Logf("zero -min  %v", UlpsBetween(zero, -min))
	t.Logf("-zero -min %v", UlpsBetween(-zero, -min))
	t.Logf("0.5 - 1    %v (log2)", log2(float64(UlpsBetween(0.5, 1.0))))
	t.Logf("0 - 1      %v (log2)", log2(float64(UlpsBetween(0, 1.0))))
	t.Logf("1 - Inf    %v (log2)", log2(float64(UlpsBetween(1.0, inf))))
 	t.Logf("subNorm    %v (log2)", log2(float64(UlpsBetween(0, 0x1p-1022))))
	
	state := uint64(1)
	for i := 1; i < rounds; i++ {
		dist := 1.0 + float64(Splitmix(&state) & ((1<<32) - 1))
		f1 := RandomFloat64(&state)
        f2 := f1
        u1 := Ulp(f1)
 		if f2 < 0 {
			f2 -= dist * u1
		} else {
			f2 += dist * u1
        }
        if u1 != Ulp(f2) {
            continue
        }
		Ulps := UlpsBetween(f1, f2)
		if Ulps != uint64(dist) {
			t.Logf("Ulps %v", Ulps)
			t.Logf("i    %d", i)
			t.Logf("F1   %v", f1)
			t.Logf("F2   %v", f2)
			t.Logf("F1   %X" , math.Float64bits(f1))
			t.Fatalf("F2   %X" , math.Float64bits(f2))
		}
	}
}


func TestUlp(t *testing.T) {
	const rounds int = 1e8*2
	log2 := math.Log2
	// Ulp := UlpB

	t.Logf("Value        Log2(Ulp)")
	t.Logf("0            %v", log2(Ulp(0)))
	t.Logf("0.5          %v", log2(Ulp(0.5)))
	t.Logf("1            %v", log2(Ulp(1)))
	t.Logf("2            %v", log2(Ulp(2)))
	t.Logf("2^51         %v", log2(Ulp(0x1p+51)))
	t.Logf("2^52         %v", log2(Ulp(0x1p+52)))
	t.Logf("2^53         %v", log2(Ulp(0x1p+53)))
	t.Logf("0x1p-1074    %v", log2(Ulp(0x1p-1074)))
	t.Logf("0x1p-1025    %v", log2(Ulp(0x1p-1025)))
	t.Logf("0x1p-1021    %v", log2(Ulp(0x1p-1021)))
	t.Logf("0x1p-1021-   %v", log2(Ulp(0x1p-1021 - 0x1p-1050)))
	t.Logf("+/-Inf       %v", log2(Ulp(math.Inf(1))))
	t.Logf("NaN          %v", log2(Ulp(math.NaN())))
	t.Logf("MaxFloat64   %v", log2(Ulp(math.MaxFloat64)))
	state := uint64(1)
	for i := 0; i < rounds; i++ {
		f1 := RandomFloat64(&state)
		f3 := Ulp(f1)
		f2 := f1 + f3
		if !Adjacent(f1, f2) || !IsPowerOfTwo(f3) {
			t.Logf("Ulps %v", UlpsBetween(f1, f2))
			t.Logf("i    %d", i)
			t.Logf("F1   %v", f1)
			t.Logf("F2   %v", f2)
			t.Logf("F3   %v", f3)
			t.Logf("F1   %X" , math.Float64bits(f1))
			t.Logf("F2   %X" , math.Float64bits(f2))
			t.Fatalf("F3   %X" , math.Float64bits(f3))
		}
	}
}

func TestLogUlp(t *testing.T) {
	const rounds int = 1e8
	t.Logf("MaxFloat64   %d", LogUlp(math.MaxFloat64))
	t.Logf("+/-Inf       %d", LogUlp(math.Inf(1)))
	t.Logf("NaN          %d", LogUlp(math.NaN()))
	state := uint64(1)
	for i := 0; i < rounds; i++ {
		f := RandomFloat64(&state)
		if i == 0 {
			f = -0x1p-1074
		}
		u := Ulp(f)
		l1 := LogUlp(f)
		l2 := Log2(u)
		if u != math.Ldexp(1, l1) || l1 != l2 {
			t.Logf("i    %d", i)
			t.Logf("F    %v", f)
			t.Fatalf("F    %X" , math.Float64bits(f))
		}
	}
}

func TestUlpWithNextafter(t *testing.T) {
    const rounds int = 1e8*2
	state := uint64(1)
	for i := 0; i < rounds; i++ {
		f1 := RandomFloat64(&state)
		if f1 < 0 {
			f1 = -f1
		}
		u1 := Ulp(f1)
		u2 := math.Nextafter(f1, math.MaxFloat64) - f1
		if u1 != u2 {
			t.Logf("i    %d", i)
			t.Logf("F1   %v", f1)
			t.Fatalf("F1   %X" , math.Float64bits(f1))
		}
	}
}
func TestLog2(t *testing.T) {
    const rounds int = 1e8
	state := uint64(1)
	for i := 0; i < rounds; i++ {
		f1 := math.Abs(RandomFloat64(&state))
		if i == 0 {
            // f1 = 0.5
            // f1 = 0x1p-1074
            f1 = 0x1p+1023
		}
		d := math.Log2(f1) - float64(Log2(f1))
        if 0 <= d && d < 1 {
            continue
        }
		t.Logf("i    %d", i)
        t.Logf("F1   %v", f1)
        t.Fatalf("F1   %X" , math.Float64bits(f1))
	}
}

func TestIsPowerOfTwo(t *testing.T) {
	const rounds int = 1e8
	// IsPowerOfTwo := IsPowerOfTwoFP
	// IsPowerOfTwo := IsPowerOfTwoJava
	zero := 0.0
	t.Logf("0            %v", IsPowerOfTwo(zero))
	t.Logf("-0           %v", IsPowerOfTwo(-zero))
	t.Logf("1            %v", IsPowerOfTwo(1))
	t.Logf("-1           %v", IsPowerOfTwo(-1))
	t.Logf("2^50         %v", IsPowerOfTwo(0x1p+50))
	t.Logf("2^-50        %v", IsPowerOfTwo(0x1p-50))
	t.Logf("2^50 - 1     %v", IsPowerOfTwo(0x1p-50 - 1))
	t.Logf("2^-1074      %v", IsPowerOfTwo(0x1p-1074))
	t.Logf("-2^-1074     %v", IsPowerOfTwo(-0x1p-1074))
	t.Logf("2^-1022      %v", IsPowerOfTwo(0x1p-1022))
	t.Logf("+Inf         %v", IsPowerOfTwo(math.Inf(1)))
	t.Logf("NaN          %v", IsPowerOfTwo(math.NaN()))

	state := uint64(1)
	for i := 0; i < rounds; i++ {
		f1 := Ulp(RandomFloat64(&state))                      // Ulp is power of two 
		f2 := math.Float64frombits(math.Float64bits(f1) + 5)  // not power of two
		if !IsPowerOfTwo(f1) || IsPowerOfTwo(f2) {
			t.Logf("i       %d", i)
			t.Logf("f1      %v", f1)
			t.Logf("f2      %v", f2)
			t.Logf("f1      %16X" , math.Float64bits(f1))
			t.Fatalf("f2      %16X" , math.Float64bits(f2))
		}
	}
}
