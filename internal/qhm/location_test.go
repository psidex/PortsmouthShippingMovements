package qhm

import "testing"

func benchmarkParseAbbreviation(s string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		parseAbbreviation(s)
	}
}

func BenchmarkParseAbbreviationEmpty(b *testing.B) {
	benchmarkParseAbbreviation("", b)
}

func BenchmarkParseAbbreviationNoAbbrv(b *testing.B) {
	benchmarkParseAbbreviation("singleword", b)
}

func BenchmarkParseAbbreviationSimple(b *testing.B) {
	benchmarkParseAbbreviation("OSB to PIP", b)
}

func BenchmarkParseAbbreviationNormal(b *testing.B) {
	benchmarkParseAbbreviation("NCJ ( W ) to SH 6 (DUMMY ANCHORAGES)", b)
}

func BenchmarkParseAbbreviationComplex(b *testing.B) {
	benchmarkParseAbbreviation("NCJ (E) from 4BIII ( W ) to SH 6 ( C ) HBR (S) (FUEL) O/C S'WAY", b)
}
