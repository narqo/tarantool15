package main_test

import (
	"testing"

	"github.com/narqo/tarantool15/tntgen/fixtures"
)

func TestPack(t *testing.T) {
	data := fixtures.TestToken.PackSilly()

	tok := fixtures.Token{}
	if err := tok.UnpackTuple(data); err != nil {
		t.Fatal(err)
	}
	if tok.String() != fixtures.TestToken.String() {
		// skip silly compassing for now
		t.Logf("PackSilly: got\n====\n%+v\n====\nwant\n====\n%+v", tok, fixtures.TestToken)
	}

	data = fixtures.TestToken.PackIdx()

	tok = fixtures.Token{}
	if err := tok.UnpackTuple(data); err != nil {
		t.Fatal(err)
	}
	if tok.String() != fixtures.TestToken.String() {
		t.Errorf("PackIdx: got\n====\n%+v\n====\nwant\n====\n%+v", tok, fixtures.TestToken)
	}

	data = fixtures.TestToken.PackTuple()

	tok = fixtures.Token{}
	if err := tok.UnpackTuple(data); err != nil {
		t.Fatal(err)
	}
	if tok.String() != fixtures.TestToken.String() {
		t.Errorf("Pack: got\n====\n%+v\n====\nwant\n====\n%+v", tok, fixtures.TestToken)
	}
}

func BenchmarkPackSilly(b *testing.B) {
	tok := fixtures.TestToken
	for i := 0; i < b.N; i++ {
		ret := tok.PackSilly()
		if ret == nil {
			b.Fatal("unexpected PackSilly result")
		}
	}
}

func BenchmarkPackIdx(b *testing.B) {
	tok := fixtures.TestToken
	for i := 0; i < b.N; i++ {
		ret := tok.PackIdx()
		if ret == nil {
			b.Fatal("unexpected PackIdx result")
		}
	}
}

func BenchmarkPackTuple(b *testing.B) {
	tok := fixtures.TestToken
	for i := 0; i < b.N; i++ {
		ret := tok.PackTuple()
		if ret == nil {
			b.Fatal("unexpected Pack result")
		}
	}
}

func BenchmarkPackTo(b *testing.B) {
	data := make([]byte, 0, 256)
	tok := fixtures.TestToken
	for i := 0; i < b.N; i++ {
		ret := tok.PackTo(data)
		if ret == nil {
			b.Fatal("unexpected Pack result")
		}
		data = data[:0]
	}
}

func BenchmarkUnpackTuple(b *testing.B) {
	data := fixtures.TestToken.PackSilly()
	for i := 0; i < b.N; i++ {
		var t1 fixtures.Token
		err := t1.UnpackTuple(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}
