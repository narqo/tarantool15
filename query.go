package tarantool15

import (
	"encoding/binary"
	"math"

	"github.com/lomik/go-tnt"
)

type Query interface {
	Pack(reqID uint32, defaultSpace uint32) ([]byte, error)
}

const (
	requestTypeInsert = 13
	requestTypeSelect = 17
	requestTypeUpdate = 19
	requestTypeDelete = 21
	requestTypeCall   = 22
)

var _ Query = (*Select)(nil)

type Select struct {
	Space interface{}
	Index uint32

	Offset uint32
	Limit  uint32

	Tuples []Tuple
}

func (q Select) Pack(reqID uint32, defaultSpace uint32) ([]byte, error) {
	b := make([]byte, 0, 256)

	b = PackUint32(b, requestTypeSelect)
	// request size, we will put the actual size value later
	b = PackUint32(b, 0)
	b = PackUint32(b, reqID)
	headSize := len(b)

	if q.Space != nil {
		var err error
		b, err = Pack(b, q.Space)
		if err != nil {
			return nil, err
		}
	} else {
		b = PackUint32(b, defaultSpace)
	}

	b = PackUint32(b, q.Index)
	b = PackUint32(b, q.Offset)

	limit := q.Limit
	if limit == 0 {
		limit = 0xffffffff
	}
	b = PackUint32(b, limit)

	if q.Tuples != nil {
		b = PackUint32(b, uint32(len(q.Tuples)))
		for i := 0; i < len(q.Tuples); i++ {
			b = q.Tuples[i].PackTo(b)
		}
	} else {
		b = PackUint32(b, 0) // count
	}

	dataSize := len(b) - headSize
	// write the real size of the request to the second (uint32, 4 bytes) position of the header
	binary.LittleEndian.PutUint32(b[4:], uint32(dataSize))

	return b, nil
}

var _ Query = (*Insert)(nil)

type Insert struct {
	Space       interface{}
	ReturnTuple bool
	Tuple       PackerTo
}

func (q Insert) Pack(reqID uint32, defaultSpace uint32) ([]byte, error) {
	return packQuery(requestTypeInsert, reqID, q.Space, defaultSpace, q.ReturnTuple, nil, q.Tuple)
}

var _ Query = (*Delete)(nil)

type Delete struct {
	ReturnTuple bool
	Space       interface{}
	Tuple       PackerTo
}

func (q Delete) Pack(reqID uint32, defaultSpace uint32) ([]byte, error) {
	return packQuery(requestTypeDelete, reqID, q.Space, defaultSpace, q.ReturnTuple, nil, q.Tuple)
}

type Call struct {
	ReturnTuple bool
	Name        []byte
	Tuple       PackerTo
}

func (q Call) Pack(reqID uint32, defaultSpace uint32) ([]byte, error) {
	return packQuery(requestTypeCall, reqID, nil, noDefaultSpace, q.ReturnTuple, q.Name, q.Tuple)
}

type OpCode uint8

const (
	opSet OpCode = iota
	opAdd
	opAnd
	opXor
	opOr
	opSplice
	opDelete
	opInsert
)

func OpSet(field uint32, value interface{}) Operation {
	return Operation{field, opSet, value}
}

func OpDelete(field uint32, value interface{}) Operation {
	return Operation{field, opDelete, value}
}

func OpInsert(field uint32, value interface{}) Operation {
	return Operation{field, opInsert, value}
}

type Operation struct {
	Field  uint32
	OpCode OpCode
	Value  interface{}
}

func (op Operation) PackTo(data []byte) []byte {
	data = PackUint32(data, op.Field)
	data = PackUint8(data, uint8(op.OpCode))
	data = packTupleField(data, op.Value)
	return data
}

type Update struct {
	Space       interface{}
	ReturnTuple bool
	Tuple       PackerTo
	Ops         []Operation
}

func (q Update) Pack(reqID uint32, defaultSpace uint32) ([]byte, error) {
	b := make([]byte, 0, 256)

	b = PackUint32(b, requestTypeUpdate)
	// request size, we will put the actual size value later
	b = PackUint32(b, 0)
	b = PackUint32(b, reqID)
	headSize := len(b)

	if q.Space != nil {
		var err error
		b, err = Pack(b, q.Space)
		if err != nil {
			return nil, err
		}
	} else {
		b = PackUint32(b, defaultSpace)
	}

	if q.ReturnTuple {
		b = append(b, packedInt1[:]...)
	} else {
		b = append(b, packedInt0[:]...)
	}

	b = q.Tuple.PackTo(b)

	if len(q.Ops) != 0 {
		b = PackUint32(b, uint32(len(q.Ops)))
		for i := 0; i < len(q.Ops); i++ {
			b = q.Ops[i].PackTo(b)
		}
	}

	dataSize := len(b) - headSize
	// write the real size of the request to the second (uint32, 4 bytes) position of the header
	binary.LittleEndian.PutUint32(b[4:], uint32(dataSize))

	return b, nil
}

type PackerTo interface {
	PackTo(data []byte) []byte
}

type Tuple []interface{}

var _ PackerTo = (Tuple)(nil)

func (t Tuple) PackTo(data []byte) []byte {
	data = PackUint32(data, uint32(len(t)))
	for j := 0; j < len(t); j++ {
		f := t[j]
		if f == nil {
			continue
		}
		data = packTupleField(data, f)
	}
	return data
}

func packTupleField(data []byte, f interface{}) []byte {
	switch fld := f.(type) {
	case int:
		data = PackTupleFieldUint32(data, uint32(fld))
	case int32:
		data = PackTupleFieldUint32(data, uint32(fld))
	case int64:
		data = PackTupleFieldUint64(data, uint64(fld))
	case uint:
		data = PackTupleFieldUint32(data, uint32(fld))
	case uint32:
		data = PackTupleFieldUint32(data, uint32(fld))
	case uint64:
		data = PackTupleFieldUint64(data, uint64(fld))
	case string:
		data = PackTupleFieldString(data, fld)
	case []byte:
		data = PackTupleFieldBytes(data, fld)
	case tnt.Bytes:
		data = PackTupleFieldBytes(data, fld)
	}
	return data
}

var (
	packedInt0 = [...]byte{0, 0, 0, 0}
	packedInt1 = [...]byte{1, 0, 0, 0}
)

const noDefaultSpace = math.MaxUint32

func packQuery(reqType, reqID uint32, space interface{}, defaultSpace uint32, retTupleFlag bool, procName []byte, tuple PackerTo) (b []byte, err error) {
	b = make([]byte, 0, 256)

	b = PackUint32(b, reqType)
	// request size, we will put the actual size value later
	b = PackUint32(b, 0)
	b = PackUint32(b, reqID)
	headSize := len(b)

	if space != nil {
		b, err = Pack(b, space)
		if err != nil {
			return nil, err
		}
	} else if defaultSpace != noDefaultSpace {
		b = PackUint32(b, defaultSpace)
	}

	if retTupleFlag {
		b = append(b, packedInt1[:]...)
	} else {
		b = append(b, packedInt0[:]...)
	}

	if procName != nil {
		b = PackTupleFieldBytes(b, procName)
	}
	b = tuple.PackTo(b)

	dataSize := len(b) - headSize
	// write the real size of the request to the second (uint32, 4 bytes) position of the header
	binary.LittleEndian.PutUint32(b[4:], uint32(dataSize))

	return b, nil
}
