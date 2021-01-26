package main

import (
	"testing"
)

func Benchmark1(b *testing.B) {
	one := One{}

	for n := 0; n < b.N; n++ {
		test1(one)
	}
}

func Benchmark2(b *testing.B) {
	one := One{}

	for n := 0; n < b.N; n++ {
		test2(&one)
	}
}

func Benchmark3(b *testing.B) {
	two := Two{}

	for n := 0; n < b.N; n++ {
		test3(two)
	}
}

func Benchmark4(b *testing.B) {
	two := Two{}

	for n := 0; n < b.N; n++ {
		test4(&two)
	}
}
