package comic

import "testing"

func BenchmarkSimpleSearch(b *testing.B) {
	testCase := []string{"smack", "name", "big"}
	for i := 0; i < b.N; i++ {
		SimpleSearch(testCase)
	}
}

func BenchmarkIndexSearch(b *testing.B) {
	testCase := []string{"one", "two", "three"}
	for i := 0; i < b.N; i++ {
		IndexSearch(testCase)
	}
}
