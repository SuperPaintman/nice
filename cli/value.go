package cli

import (
	"errors"
	"fmt"
	"strconv"
)

var ErrSyntax = errors.New("invalid syntax")

var ErrRange = errors.New("value out of range")

type ParseError struct {
	Type string
	Err  error
}

func (e *ParseError) Error() string {
	msg := "unknown error"
	if e.Err != nil {
		msg = e.Err.Error()
	}

	return fmt.Sprintf("parse %s error: %s", e.Type, msg)
}

func (e *ParseError) Unwrap() error { return e.Err }

func (e *ParseError) Is(err error) bool {
	pe, ok := err.(*ParseError)
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

	return &ParseError{
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
	Value
	Type() string
}

// bool

func (b *boolValue) Set(s string) error {
	v, err := parseBool(s)
	if err != nil {
		err = &ParseError{
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
// TODO

// uint16
// TODO

// uint32
// TODO

// uint64
// TODO

// int8
// TODO

// int16
// TODO

// int32
// TODO

// int64
// TODO

// float32
// TODO

// float64
// TODO

// string

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
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
