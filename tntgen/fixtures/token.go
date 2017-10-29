package fixtures

import (
	"fmt"
	"time"

	"github.com/lomik/go-tnt"
	"github.com/narqo/tarantool15"
)

type Service struct {
	ID   uint32
	Name string
}

//go:generate tntgen -type Token -out token_tnt15.go

type Token struct {
	Email        string
	Type         uint32
	AppID        int64
	CreatedAt    time.Time
	ExpiresAt    time.Time
	Token        string
	RefreshToken string
	Data         []byte
	RawStatus    []byte
	ValidTill    time.Time
}

func (t Token) String() string {
	return fmt.Sprintf("%s/%d/%d/%d/%d/%s/%s/%x/%x/%d",
		t.Email,
		t.Type,
		t.AppID,
		t.CreatedAt.Unix(),
		t.ExpiresAt.Unix(),
		t.Token,
		t.RefreshToken,
		t.Data,
		t.RawStatus,
		t.ValidTill.Unix(),
	)
}

func (v *Token) PackSilly() (ret tnt.Tuple) {
	return tnt.Tuple{
		tnt.Bytes(v.Email),
		tnt.PackInt(v.Type),
		tnt.PackLong(uint64(v.AppID)),
		tnt.PackInt(uint32(v.CreatedAt.Unix())),
		tnt.PackInt(uint32(v.ExpiresAt.Unix())),
		tnt.Bytes(v.Token),
		tnt.Bytes(v.RefreshToken),
		v.Data,
		v.RawStatus,
		tnt.PackInt(uint32(v.ValidTill.Unix())),
	}
}

func (v *Token) PackIdx() (ret tnt.Tuple) {
	var idx [10]int

	data := make([]byte, 0, 100)

	data = tarantool15.PackString(data, v.Email)
	idx[0] = len(v.Email)
	data = tarantool15.PackUint32(data, v.Type)
	idx[1] = 4
	data = tarantool15.PackUint64(data, uint64(v.AppID))
	idx[2] = 8
	if !v.CreatedAt.IsZero() {
		data = tarantool15.PackUint32(data, uint32(v.CreatedAt.Unix()))
		idx[3] = 4
	}
	if !v.ExpiresAt.IsZero() {
		data = tarantool15.PackUint32(data, uint32(v.ExpiresAt.Unix()))
		idx[4] = 4
	}
	data = tarantool15.PackString(data, v.Token)
	idx[5] = len(v.Token)
	data = tarantool15.PackString(data, v.RefreshToken)
	idx[6] = len(v.RefreshToken)
	data = tarantool15.PackRawBytes(data, v.Data)
	idx[7] = len(v.Data)
	data = tarantool15.PackRawBytes(data, v.RawStatus)
	idx[8] = len(v.RawStatus)
	if !v.ValidTill.IsZero() {
		data = tarantool15.PackUint32(data, uint32(v.ValidTill.Unix()))
		idx[9] = 4
	}

	ret = make(tnt.Tuple, 10)
	var left int
	for i, n := range idx {
		if n == 0 {
			continue
		}
		ret[i] = data[left : left+n]
		left += n
	}
	return
}

var TestToken = Token{
	Email:        "user@example.com",
	Type:         8,
	AppID:        32,
	CreatedAt:    time.Now(),
	ExpiresAt:    time.Now().AddDate(0, 0, 1),
	Token:        "testest tokentoken",
	RefreshToken: "test refresh",
	RawStatus:    []byte("testest statusstatus"),
}
