package cli

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrSyntax = errors.New("invalid syntax")

	ErrRange = errors.New("value out of range")
)

type ParseValueError struct {
	Type string
	Err  error
}

func (e *ParseValueError) Error() string {
	msg := "unknown error"
	if e.Err != nil {
		msg = e.Err.Error()
	}

	return fmt.Sprintf("parse %s error: %s", e.Type, msg)
}

func (e *ParseValueError) Unwrap() error { return e.Err }

func (e *ParseValueError) Is(err error) bool {
	pe, ok := err.(*ParseValueError)
	return ok && pe.Type == e.Type && errors.Is(pe.Err, e.Err)
}

func numError(typ string, err error) error {
	ne, ok := err.(*strconv.NumError)
	if ok {
		if ne.Err == strconv.ErrSyntax {
			err = ErrSyntax
		} else if ne.Err == strconv.ErrRange {
			err = ErrRange
		}
	}

	return &ParseValueError{
		Type: typ,
		Err:  err,
	}
}

// var _ flag.Value = (Value)(nil)

type Value interface {
	String() string
	Set(string) error
}

// var _ flag.Value = (Getter)(nil)

type Getter interface {
	Value
	Get() interface{}
}

type Typer interface {
	Type() string
}

// bool

func (b *boolValue) Set(s string) error {
	v, err := parseBool(s)
	if err != nil {
		err = &ParseValueError{
			Type: "bool",
			Err:  err,
		}
	}
	*b = boolValue(v)
	return err
}

func (b *boolValue) IsBoolFlag() bool { return true }

type boolFlag interface {
	Value
	IsBoolFlag() bool
}

const maxBoolStringLen = len("false") // "1 byte", no, yes, true, false

func boolToLower(src []byte) {
	for i, b := range src {
		if b >= 'A' && b <= 'Z' {
			src[i] = b - 'A' + 'a'
		}
	}
}

func isBoolValue(str string) bool {
	if len(str) <= maxBoolStringLen {
		// Inline optimized and alloc-free "to lower" converter.
		var buf [maxBoolStringLen]byte
		s := buf[:len(str)]
		copy(s, str)
		boolToLower(s)

		// A little bit optimized value checking switch.
		switch len(s) {
		case 1: // 1, t, y, 0, f, n
			return s[0] == '1' || s[0] == 't' || s[0] == 'y' ||
				s[0] == '0' || s[0] == 'f' || s[0] == 'n'
		case 2: // no
			return string(s) == "no"
		case 3: // yes
			return string(s) == "yes"
		case 4: // true
			return string(s) == "true"
		case 5: // false
			return string(s) == "false"
		}
	}

	return false
}

func parseBool(str string) (bool, error) {
	if len(str) <= maxBoolStringLen {
		// Inline optimized and alloc-free "to lower" converter.
		var buf [maxBoolStringLen]byte
		s := buf[:len(str)]
		copy(s, str)
		boolToLower(s)

		// A little bit optimized value checking switch.
		switch len(s) {
		case 1: // 1, t, y, 0, f, n
			if s[0] == '1' || s[0] == 't' || s[0] == 'y' {
				return true, nil
			} else if s[0] == '0' || s[0] == 'f' || s[0] == 'n' {
				return false, nil
			}
		case 2: // no
			if string(s) == "no" {
				return false, nil
			}
		case 3: // yes
			if string(s) == "yes" {
				return true, nil
			}
		case 4: // true
			if string(s) == "true" {
				return true, nil
			}
		case 5: // false
			if string(s) == "false" {
				return false, nil
			}
		}
	}

	return false, ErrSyntax
}

// uint8

func (u *uint8Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 8)
	if err != nil {
		err = numError("uint8", err)
	}
	*u = uint8Value(v)
	return err
}

// uint16

func (u *uint16Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 16)
	if err != nil {
		err = numError("uint16", err)
	}
	*u = uint16Value(v)
	return err
}

// uint32

func (u *uint32Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 32)
	if err != nil {
		err = numError("uint32", err)
	}
	*u = uint32Value(v)
	return err
}

// uint64

func (u *uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		err = numError("uint64", err)
	}
	*u = uint64Value(v)
	return err
}

// int8

func (i *int8Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 8)
	if err != nil {
		err = numError("int8", err)
	}
	*i = int8Value(v)
	return err
}

// int16

func (i *int16Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 16)
	if err != nil {
		err = numError("int16", err)
	}
	*i = int16Value(v)
	return err
}

// int32

func (i *int32Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 32)
	if err != nil {
		err = numError("int32", err)
	}
	*i = int32Value(v)
	return err
}

// int64

func (i *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		err = numError("int64", err)
	}
	*i = int64Value(v)
	return err
}

// float32

func (i *float32Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		err = numError("float32", err)
	}
	*i = float32Value(v)
	return err
}

// float64

func (i *float64Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		err = numError("float64", err)
	}
	*i = float64Value(v)
	return err
}

// string

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) IsStringFlag() bool { return true }

type stringFlag interface {
	Value
	IsStringFlag() bool
}

// int

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		err = numError("int", err)
	}
	*i = intValue(v)
	return err
}

// uint

func (u *uintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, strconv.IntSize)
	if err != nil {
		err = numError("uint", err)
	}
	*u = uintValue(v)
	return err
}

//go:generate python ./generate_value.py

//go:generate python ./generate_values.py
