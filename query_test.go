package tarantool15

import (
	"bytes"
	"testing"

	"github.com/lomik/go-tnt"
)

func TestSelect_Pack(t *testing.T) {
	cases := []struct {
		query  Select
		golden *tnt.Select
	}{
		{
			query: Select{
				Space:  9,
				Index:  2,
				Limit:  5,
				Offset: 7,
				Tuples: []Tuple{
					{"test", 4},
				},
			},
			golden: &tnt.Select{
				Space:  9,
				Index:  2,
				Limit:  5,
				Offset: 7,
				Tuples: []tnt.Tuple{
					{tnt.Bytes("test"), tnt.PackInt(4)},
				},
			},
		},
		{
			query: Select{
				Space:  9,
				Index:  2,
				Limit:  5,
				Offset: 7,
				Tuples: []Tuple{
					{},
				},
			},
			golden: &tnt.Select{
				Space:  9,
				Index:  2,
				Limit:  5,
				Offset: 7,
				Tuples: []tnt.Tuple{
					{},
				},
			},
		},
		{
			query: Select{
				Limit:  5,
				Offset: 7,
			},
			golden: &tnt.Select{
				Limit:  5,
				Offset: 7,
			},
		},
	}

	for _, tc := range cases {
		testQuery(t, tc.query, tc.golden)
	}
}

func TestInsert_Pack(t *testing.T) {
	cases := []struct {
		query  Insert
		golden *tnt.Insert
	}{
		{
			query: Insert{
				Space: 9,
				Tuple: Tuple{
					"me@example.com",
					"user1@example.com",
					2,
					uint64(5),
				},
			},
			golden: &tnt.Insert{
				Space: 9,
				Tuple: tnt.Tuple{
					tnt.Bytes("me@example.com"),
					tnt.Bytes("user1@example.com"),
					tnt.PackInt(2),
					tnt.PackLong(uint64(5)),
				},
			},
		},
	}

	for _, tc := range cases {
		testQuery(t, tc.query, tc.golden)
	}
}

func TestUpdate_Pack(t *testing.T) {
	cases := []struct {
		query  Update
		golden *tnt.Update
	}{
		{
			query: Update{
				Space: 9,
				Tuple: Tuple{
					"me@example.com",
					"user1@example.com",
					2,
					uint64(5),
				},
				Ops: []Operation{
					OpSet(1, "OK"),
				},
			},
			golden: &tnt.Update{
				Space: 9,
				Tuple: tnt.Tuple{
					tnt.Bytes("me@example.com"),
					tnt.Bytes("user1@example.com"),
					tnt.PackInt(2),
					tnt.PackLong(uint64(5)),
				},
				Ops: []tnt.Operator{
					tnt.OpSet(1, tnt.Bytes("OK")),
				},
			},
		},
	}

	for _, tc := range cases {
		testQuery(t, tc.query, tc.golden)
	}
}

func TestDelete_Pack(t *testing.T) {
	cases := []struct {
		query  Delete
		golden *tnt.Delete
	}{
		{
			query: Delete{
				Space: 9,
				Tuple: Tuple{
					"me@example.com",
					"user1@example.com",
					2,
					uint64(5),
				},
			},
			golden: &tnt.Delete{
				Space: 9,
				Tuple: tnt.Tuple{
					tnt.Bytes("me@example.com"),
					tnt.Bytes("user1@example.com"),
					tnt.PackInt(2),
					tnt.PackLong(uint64(5)),
				},
			},
		},
	}

	for _, tc := range cases {
		testQuery(t, tc.query, tc.golden)
	}
}

func TestCall_Pack(t *testing.T) {
	cmdDoString := []byte("box.dostring")
	rawlua := "local x,t = box.space[1].index[3]:next_equal('G','OK')\n\tif x == nil then return end\n\tif (os.time() - box.unpack('i',t[13])) >= 600 then return t end\n\treturn"

	cases := []struct {
		query  Call
		golden *tnt.Call
	}{
		{
			query: Call{
				Name: cmdDoString,
				Tuple: Tuple{
					rawlua,
				},
			},
			golden: &tnt.Call{
				Name: cmdDoString,
				Tuple: tnt.Tuple{
					tnt.Bytes(rawlua),
				},
			},
		},
	}

	for _, tc := range cases {
		testQuery(t, tc.query, tc.golden)
	}
}

func testQuery(t *testing.T, q1 Query, q2 tnt.Query) {
	const reqid uint32 = 10
	const defsp uint32 = 1

	got, err := q1.Pack(reqid, defsp)
	if err != nil {
		t.Fatal(err)
	}
	want, _ := q2.Pack(reqid, defsp)
	if !bytes.Equal(got, want) {
		t.Fatalf("Pack: \ngot\n=====\n%x\nwant\n======\n%x\n", got, want)
	}
}

func BenchmarkQuery_Pack(b *testing.B) {
	cases := []struct {
		name  string
		query func() Query
	}{
		{
			"Select",
			func() Query {
				return Select{
					Space:  9,
					Index:  2,
					Limit:  5,
					Offset: 7,
					Tuples: []Tuple{
						{"me@example.com", "user1@example.com", 2, uint64(5)},
					},
				}
			},
		},
		{
			"tnt.Select",
			func() Query {
				return &tnt.Select{
					Space:  9,
					Index:  2,
					Limit:  5,
					Offset: 7,
					Tuples: []tnt.Tuple{
						{tnt.Bytes("me@example.com"), tnt.Bytes("user1@example.com"), tnt.PackInt(2), tnt.PackLong(uint64(5))},
					},
				}
			},
		},
		{
			"Insert",
			func() Query {
				return Insert{
					Space: 9,
					Tuple: Tuple{
						"me@example.com",
						"user1@example.com",
						uint32(2),
						uint64(5),
					},
				}
			},
		},
		{
			"tnt.Insert",
			func() Query {
				return &tnt.Insert{
					Space: 9,
					Tuple: tnt.Tuple{
						tnt.Bytes("me@example.com"),
						tnt.Bytes("user1@example.com"),
						tnt.PackInt(2),
						tnt.PackLong(uint64(5)),
					},
				}
			},
		},
		{
			"Update",
			func() Query {
				return Update{
					Space: 9,
					Tuple: Tuple{
						"me@example.com",
						"user1@example.com",
						uint32(2),
						uint64(5),
					},
					Ops: []Operation{
						OpSet(1, "OK"),
						OpSet(9, 17),
					},
				}
			},
		},
		{
			"tnt.Update",
			func() Query {
				return &tnt.Update{
					Space: 9,
					Tuple: tnt.Tuple{
						tnt.Bytes("me@example.com"),
						tnt.Bytes("user1@example.com"),
						tnt.PackInt(2),
						tnt.PackLong(uint64(5)),
					},
					Ops: []tnt.Operator{
						tnt.OpSet(1, tnt.Bytes("OK")),
						tnt.OpSet(9, tnt.PackInt(17)),
					},
				}
			},
		},
		{
			"Delete",
			func() Query {
				return Delete{
					Space: 9,
					Tuple: Tuple{
						"me@example.com",
						"user1@example.com",
						uint32(2),
						uint64(5),
					},
				}
			},
		},
		{
			"tnt.Delete",
			func() Query {
				return &tnt.Delete{
					Space: 9,
					Tuple: tnt.Tuple{
						tnt.Bytes("me@example.com"),
						tnt.Bytes("user1@example.com"),
						tnt.PackInt(2),
						tnt.PackLong(uint64(5)),
					},
				}
			},
		},
		{
			"Call",
			func() Query {
				return Call{
					Name: []byte("box.dostring"),
					Tuple: Tuple{
						"local x,t = box.space[2].index[1]:next_equal('G','OK') if x == nil then return",
					},
				}
			},
		},
		{
			"tnt.Call",
			func() Query {
				return &tnt.Call{
					Name: []byte("box.dostring"),
					Tuple: tnt.Tuple{
						tnt.Bytes("local x,t = box.space[2].index[1]:next_equal('G','OK') if x == nil then return"),
					},
				}
			},
		},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			query := tc.query()
			for i := 0; i < b.N; i++ {
				p, err := query.Pack(10, 1)
				if err != nil {
					b.Fatal(err)
				}
				if len(p) == 0 {
					b.Fatal("unexpected Pack response")
				}
			}
		})
	}
}
