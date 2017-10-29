package tarantool15

import (
	"encoding/hex"
	"fmt"
	"math"
	"testing"
)

var testCases = []struct {
	in     interface{}
	outhex string
}{
	{nil, ""},
	{int8(1), "01"},
	{int16(3), "03000000"}, // int16 packed as int32
	{int32(5), "05000000"},
	{int64(7), "0700000000000000"},
	{uint8(2), "02"},
	{uint16(4), "04000000"},
	{uint32(6), "06000000"},
	{uint64(8), "0800000000000000"},
	{uint64(0), "0000000000000000"},
	{byte('a'), "61"},
	{"testestest testestest", "746573746573746573742074657374657374657374"},
	{"", ""},
}

func TestPack(t *testing.T) {
	for _, tc := range testCases {
		name := fmt.Sprintf("%T", tc.in)
		t.Run(name, func(t *testing.T) {
			b, err := Pack(nil, tc.in)
			if err != nil {
				t.Fatal(err)
			}
			gothex := hex.EncodeToString(b)
			if gothex != tc.outhex {
				t.Fatalf("got %q, want %q", gothex, tc.outhex)
			}
		})
	}
}

func TestPackTupleField(t *testing.T) {
	cases := []struct {
		name   string
		packfn func() []byte
		outhex string
	}{
		{
			name: "PackTupleFieldUint32",
			packfn: func() []byte {
				return PackTupleFieldUint32(nil, uint32(0))
			},
			outhex: "0400000000",
		},
		{
			name: "PackTupleFieldUint32(MaxUint32)",
			packfn: func() []byte {
				return PackTupleFieldUint32(nil, math.MaxUint32)
			},
			outhex: "04ffffffff",
		},
		{
			name: "PackTupleFieldUint64",
			packfn: func() []byte {
				return PackTupleFieldUint64(nil, uint64(156))
			},
			outhex: "089c00000000000000",
		},
		{
			name: "PackTupleFieldString",
			packfn: func() []byte {
				return PackTupleFieldString(nil, "testest testest")
			},
			outhex: "0f746573746573742074657374657374",
		},
		{
			name: "PackTupleFieldString(empty)",
			packfn: func() []byte {
				return PackTupleFieldString(nil, "")
			},
			outhex: "00",
		},
		{
			name: "PackTupleFieldBytes",
			packfn: func() []byte {
				return PackTupleFieldBytes(nil, []byte("testest testest"))
			},
			outhex: "0f746573746573742074657374657374",
		},
		{
			name: "PackTupleFieldBytes(nil)",
			packfn: func() []byte {
				return PackTupleFieldBytes(nil, nil)
			},
			outhex: "00",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := tc.packfn()
			gothex := hex.EncodeToString(b)
			if gothex != tc.outhex {
				t.Fatalf("got %q, want %q", gothex, tc.outhex)
			}
		})
	}
}

func BenchmarkPack(b *testing.B) {
	for _, tc := range testCases {
		name := fmt.Sprintf("%T", tc.in)
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				var err error
				_, err = Pack(nil, tc.in)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
