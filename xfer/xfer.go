package xfer

import (
	"fmt"
	"reflect"

	"github.com/opennox/opennox-lib/binenc"
)

// Type is a decoding type for XFER data.
type Type string

// Xfer is type that can be encoded or decoded to/from Nox XFER data.
type Xfer interface {
	XferType() Type
	DecodeXfer(reg ObjectRegistry, r *binenc.Reader) error
}

// ObjectRegistry is an interface for looking up XFER Type based on object type name or ID.
type ObjectRegistry interface {
	// XferByObjectType finds XFER Type for a given object type name.
	XferByObjectType(typ string) Type
	// XferByObjectTypeID finds XFER Type for a given object type ID.
	XferByObjectTypeID(id int) Type
}

var byType = make(map[Type]reflect.Type)

// Register a new XFER data type.
func Register(x Xfer) {
	typ := x.XferType()
	if _, ok := byType[typ]; ok {
		panic("already registered")
	}
	byType[typ] = reflect.TypeOf(x).Elem()
}

// Decode XFER data with a specific Type name.
func Decode(reg ObjectRegistry, xfer Type, r *binenc.Reader) (Xfer, error) {
	if reg == nil {
		reg = DefaultRegistry
	}
	rt, ok := byType[xfer]
	if !ok {
		return nil, fmt.Errorf("unsupported xfer: %q", xfer)
	}
	x := reflect.New(rt).Interface().(Xfer)
	err := x.DecodeXfer(reg, r)
	return x, err
}

// DecodeByObjectType is similar to Decode, but accepts object type name, instead of XFER Type.
func DecodeByObjectType(reg ObjectRegistry, typ string, r *binenc.Reader) (Xfer, error) {
	if reg == nil {
		reg = DefaultRegistry
	}
	xfer := reg.XferByObjectType(typ)
	if xfer == "" {
		return nil, fmt.Errorf("cannot find xfer for object type: %q", typ)
	}
	return Decode(reg, xfer, r)
}

// DecodeByObjectTypeID is similar to Decode, but accepts object type ID, instead of XFER Type.
func DecodeByObjectTypeID(reg ObjectRegistry, id int, r *binenc.Reader) (Xfer, error) {
	if reg == nil {
		reg = DefaultRegistry
	}
	xfer := reg.XferByObjectTypeID(id)
	if xfer == "" {
		return nil, fmt.Errorf("cannot find xfer for object type id: %d", id)
	}
	return Decode(reg, xfer, r)
}
