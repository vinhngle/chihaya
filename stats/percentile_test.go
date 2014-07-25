// Copyright 2014 The Chihaya Authors. All rights reserved.
// Use of this source code is governed by the BSD 2-Clause license,
// which can be found in the LICENSE file.

package stats

import (
	"math"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func TestPercentiles(t *testing.T) {
	rand.Seed(time.Now().Unix())

	testSlice(t, uniform(10000, 1), 0.5)
	testSlice(t, uniform(10000, 1), 0.9)
	testSlice(t, uniform(10000, 10000), 0.5)
	testSlice(t, uniform(10000, 10000), 0.9)
}

func TestLogNormPercentiles(t *testing.T) {
	rand.Seed(time.Now().Unix())

	testSlice(t, logNorm(10000, 1), 0.5)
	testSlice(t, logNorm(10000, 1), 0.9)
}

func uniform(n int, scale float64) sort.Float64Slice {
	numbers := make(sort.Float64Slice, n)

	for i := 0; i < n; i++ {
		numbers[i] = rand.Float64() * scale
	}

	return numbers
}

func logNorm(n int, scale float64) sort.Float64Slice {
	numbers := make(sort.Float64Slice, n)

	for i := 0; i < n; i++ {
		numbers[i] = math.Exp(rand.NormFloat64()) * scale
	}

	return numbers
}

func testSlice(t *testing.T, numbers sort.Float64Slice, percentile float64) {
	p := NewPercentile(percentile)

	for i := 0; i < len(numbers); i++ {
		p.AddSample(numbers[i])
	}

	sort.Sort(numbers)
	got := p.Value()
	index := round(float64(len(numbers)) * percentile)

	if got != numbers[index] && got != numbers[index-1] && got != numbers[index+1] {
		t.Errorf("Percentile incorrect\n  actual: %f\nexpected: %f, %f, %f\n", got, numbers[index-1], numbers[index], numbers[index+1])
	}
}

func BenchmarkPercentiles64(b *testing.B) {
	benchmarkSlice(b, uniform(b.N, 1), 64, 0.5)
}

func BenchmarkPercentiles128(b *testing.B) {
	benchmarkSlice(b, uniform(b.N, 1), 128, 0.5)
}

func BenchmarkPercentiles256(b *testing.B) {
	benchmarkSlice(b, uniform(b.N, 1), 256, 0.5)
}

func BenchmarkPercentiles512(b *testing.B) {
	benchmarkSlice(b, uniform(b.N, 1), 512, 0.5)
}

func BenchmarkLNPercentiles128(b *testing.B) {
	benchmarkSlice(b, logNorm(b.N, 1), 128, 0.5)
}

func BenchmarkLNPercentiles256(b *testing.B) {
	benchmarkSlice(b, logNorm(b.N, 1), 258, 0.5)
}

func benchmarkSlice(b *testing.B, numbers sort.Float64Slice, window int, percentile float64) {
	p := NewPercentileWithWindow(percentile, window)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.AddSample(numbers[i])
	}
}
