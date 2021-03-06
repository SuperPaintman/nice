// Code generated by generate_values.py; DO NOT EDIT.

package cli

import (
	"strings"
	"time"
)

// []bool

var (
	_ Value   = (*boolValues)(nil)
	_ Getter  = (*boolValues)(nil)
	_ Emptier = (*boolValues)(nil)
	_ Typer   = (*boolValues)(nil)
)

type boolValues []bool

func newBoolValues(p *[]bool) *boolValues {
	return (*boolValues)(p)
}

func (vs *boolValues) Set(val string) error {
	rest := val
	for rest != "" {
		idx := strings.IndexByte(rest, ',')
		if idx != -1 {
			val = rest[:idx]
			rest = rest[idx+1:]
		} else {
			val = rest
			rest = ""
		}

		var def bool
		*vs = append(*vs, def)
		if err := (*boolValue)(&(*vs)[len(*vs)-1]).Set(val); err != nil {
			return err
		}
	}

	return nil
}

func (vs *boolValues) String() string {
	if len(*vs) == 0 {
		return ""
	}

	var buf strings.Builder
	_, _ = buf.WriteString((*boolValue)(&(*vs)[0]).String())

	for i := 1; i < len(*vs); i++ {
		_ = buf.WriteByte(',')
		_, _ = buf.WriteString((*boolValue)(&(*vs)[i]).String())
	}

	return buf.String()
}

func (v *boolValues) Empty() bool { return len(*v) == 0 }

func (v *boolValues) Get() interface{} { return []bool(*v) }

func (*boolValues) Type() string { return "[]bool" }

// []uint8

var (
	_ Value   = (*uint8Values)(nil)
	_ Getter  = (*uint8Values)(nil)
	_ Emptier = (*uint8Values)(nil)
	_ Typer   = (*uint8Values)(nil)
)

type uint8Values []uint8

func newUint8Values(p *[]uint8) *uint8Values {
	return (*uint8Values)(p)
}

func (vs *uint8Values) Set(val string) error {
	rest := val
	for rest != "" {
		idx := strings.IndexByte(rest, ',')
		if idx != -1 {
			val = rest[:idx]
			rest = rest[idx+1:]
		} else {
			val = rest
			rest = ""
		}

		var def uint8
		*vs = append(*vs, def)
		if err := (*uint8Value)(&(*vs)[len(*vs)-1]).Set(val); err != nil {
			return err
		}
	}

	return nil
}

func (vs *uint8Values) String() string {
	if len(*vs) == 0 {
		return ""
	}

	var buf strings.Builder
	_, _ = buf.WriteString((*uint8Value)(&(*vs)[0]).String())

	for i := 1; i < len(*vs); i++ {
		_ = buf.WriteByte(',')
		_, _ = buf.WriteString((*uint8Value)(&(*vs)[i]).String())
	}

	return buf.String()
}

func (v *uint8Values) Empty() bool { return len(*v) == 0 }

func (v *uint8Values) Get() interface{} { return []uint8(*v) }

func (*uint8Values) Type() string { return "[]uint8" }

// []uint16

var (
	_ Value   = (*uint16Values)(nil)
	_ Getter  = (*uint16Values)(nil)
	_ Emptier = (*uint16Values)(nil)
	_ Typer   = (*uint16Values)(nil)
)

type uint16Values []uint16

func newUint16Values(p *[]uint16) *uint16Values {
	return (*uint16Values)(p)
}

func (vs *uint16Values) Set(val string) error {
	rest := val
	for rest != "" {
		idx := strings.IndexByte(rest, ',')
		if idx != -1 {
			val = rest[:idx]
			rest = rest[idx+1:]
		} else {
			val = rest
			rest = ""
		}

		var def uint16
		*vs = append(*vs, def)
		if err := (*uint16Value)(&(*vs)[len(*vs)-1]).Set(val); err != nil {
			return err
		}
	}

	return nil
}

func (vs *uint16Values) String() string {
	if len(*vs) == 0 {
		return ""
	}

	var buf strings.Builder
	_, _ = buf.WriteString((*uint16Value)(&(*vs)[0]).String())

	for i := 1; i < len(*vs); i++ {
		_ = buf.WriteByte(',')
		_, _ = buf.WriteString((*uint16Value)(&(*vs)[i]).String())
	}

	return buf.String()
}

func (v *uint16Values) Empty() bool { return len(*v) == 0 }

func (v *uint16Values) Get() interface{} { return []uint16(*v) }

func (*uint16Values) Type() string { return "[]uint16" }

// []uint32

var (
	_ Value   = (*uint32Values)(nil)
	_ Getter  = (*uint32Values)(nil)
	_ Emptier = (*uint32Values)(nil)
	_ Typer   = (*uint32Values)(nil)
)

type uint32Values []uint32

func newUint32Values(p *[]uint32) *uint32Values {
	return (*uint32Values)(p)
}

func (vs *uint32Values) Set(val string) error {
	rest := val
	for rest != "" {
		idx := strings.IndexByte(rest, ',')
		if idx != -1 {
			val = rest[:idx]
			rest = rest[idx+1:]
		} else {
			val = rest
			rest = ""
		}

		var def uint32
		*vs = append(*vs, def)
		if err := (*uint32Value)(&(*vs)[len(*vs)-1]).Set(val); err != nil {
			return err
		}
	}

	return nil
}

func (vs *uint32Values) String() string {
	if len(*vs) == 0 {
		return ""
	}

	var buf strings.Builder
	_, _ = buf.WriteString((*uint32Value)(&(*vs)[0]).String())

	for i := 1; i < len(*vs); i++ {
		_ = buf.WriteByte(',')
		_, _ = buf.WriteString((*uint32Value)(&(*vs)[i]).String())
	}

	return buf.String()
}

func (v *uint32Values) Empty() bool { return len(*v) == 0 }

func (v *uint32Values) Get() interface{} { return []uint32(*v) }

func (*uint32Values) Type() string { return "[]uint32" }

// []uint64

var (
	_ Value   = (*uint64Values)(nil)
	_ Getter  = (*uint64Values)(nil)
	_ Emptier = (*uint64Values)(nil)
	_ Typer   = (*uint64Values)(nil)
)

type uint64Values []uint64

func newUint64Values(p *[]uint64) *uint64Values {
	return (*uint64Values)(p)
}

func (vs *uint64Values) Set(val string) error {
	rest := val
	for rest != "" {
		idx := strings.IndexByte(rest, ',')
		if idx != -1 {
			val = rest[:idx]
			rest = rest[idx+1:]
		} else {
			val = rest
			rest = ""
		}

		var def uint64
		*vs = append(*vs, def)
		if err := (*uint64Value)(&(*vs)[len(*vs)-1]).Set(val); err != nil {
			return err
		}
	}

	return nil
}

func (vs *uint64Values) String() string {
	if len(*vs) == 0 {
		return ""
	}

	var buf strings.Builder
	_, _ = buf.WriteString((*uint64Value)(&(*vs)[0]).String())

	for i := 1; i < len(*vs); i++ {
		_ = buf.WriteByte(',')
		_, _ = buf.WriteString((*uint64Value)(&(*vs)[i]).String())
	}

	return buf.String()
}

func (v *uint64Values) Empty() bool { return len(*v) == 0 }

func (v *uint64Values) Get() interface{} { return []uint64(*v) }

func (*uint64Values) Type() string { return "[]uint64" }

// []int8

var (
	_ Value   = (*int8Values)(nil)
	_ Getter  = (*int8Values)(nil)
	_ Emptier = (*int8Values)(nil)
	_ Typer   = (*int8Values)(nil)
)

type int8Values []int8

func newInt8Values(p *[]int8) *int8Values {
	return (*int8Values)(p)
}

func (vs *int8Values) Set(val string) error {
	rest := val
	for rest != "" {
		idx := strings.IndexByte(rest, ',')
		if idx != -1 {
			val = rest[:idx]
			rest = rest[idx+1:]
		} else {
			val = rest
			rest = ""
		}

		var def int8
		*vs = append(*vs, def)
		if err := (*int8Value)(&(*vs)[len(*vs)-1]).Set(val); err != nil {
			return err
		}
	}

	return nil
}

func (vs *int8Values) String() string {
	if len(*vs) == 0 {
		return ""
	}

	var buf strings.Builder
	_, _ = buf.WriteString((*int8Value)(&(*vs)[0]).String())

	for i := 1; i < len(*vs); i++ {
		_ = buf.WriteByte(',')
		_, _ = buf.WriteString((*int8Value)(&(*vs)[i]).String())
	}

	return buf.String()
}

func (v *int8Values) Empty() bool { return len(*v) == 0 }

func (v *int8Values) Get() interface{} { return []int8(*v) }

func (*int8Values) Type() string { return "[]int8" }

// []int16

var (
	_ Value   = (*int16Values)(nil)
	_ Getter  = (*int16Values)(nil)
	_ Emptier = (*int16Values)(nil)
	_ Typer   = (*int16Values)(nil)
)

type int16Values []int16

func newInt16Values(p *[]int16) *int16Values {
	return (*int16Values)(p)
}

func (vs *int16Values) Set(val string) error {
	rest := val
	for rest != "" {
		idx := strings.IndexByte(rest, ',')
		if idx != -1 {
			val = rest[:idx]
			rest = rest[idx+1:]
		} else {
			val = rest
			rest = ""
		}

		var def int16
		*vs = append(*vs, def)
		if err := (*int16Value)(&(*vs)[len(*vs)-1]).Set(val); err != nil {
			return err
		}
	}

	return nil
}

func (vs *int16Values) String() string {
	if len(*vs) == 0 {
		return ""
	}

	var buf strings.Builder
	_, _ = buf.WriteString((*int16Value)(&(*vs)[0]).String())

	for i := 1; i < len(*vs); i++ {
		_ = buf.WriteByte(',')
		_, _ = buf.WriteString((*int16Value)(&(*vs)[i]).String())
	}

	return buf.String()
}

func (v *int16Values) Empty() bool { return len(*v) == 0 }

func (v *int16Values) Get() interface{} { return []int16(*v) }

func (*int16Values) Type() string { return "[]int16" }

// []int32

var (
	_ Value   = (*int32Values)(nil)
	_ Getter  = (*int32Values)(nil)
	_ Emptier = (*int32Values)(nil)
	_ Typer   = (*int32Values)(nil)
)

type int32Values []int32

func newInt32Values(p *[]int32) *int32Values {
	return (*int32Values)(p)
}

func (vs *int32Values) Set(val string) error {
	rest := val
	for rest != "" {
		idx := strings.IndexByte(rest, ',')
		if idx != -1 {
			val = rest[:idx]
			rest = rest[idx+1:]
		} else {
			val = rest
			rest = ""
		}

		var def int32
		*vs = append(*vs, def)
		if err := (*int32Value)(&(*vs)[len(*vs)-1]).Set(val); err != nil {
			return err
		}
	}

	return nil
}

func (vs *int32Values) String() string {
	if len(*vs) == 0 {
		return ""
	}

	var buf strings.Builder
	_, _ = buf.WriteString((*int32Value)(&(*vs)[0]).String())

	for i := 1; i < len(*vs); i++ {
		_ = buf.WriteByte(',')
		_, _ = buf.WriteString((*int32Value)(&(*vs)[i]).String())
	}

	return buf.String()
}

func (v *int32Values) Empty() bool { return len(*v) == 0 }

func (v *int32Values) Get() interface{} { return []int32(*v) }

func (*int32Values) Type() string { return "[]int32" }

// []int64

var (
	_ Value   = (*int64Values)(nil)
	_ Getter  = (*int64Values)(nil)
	_ Emptier = (*int64Values)(nil)
	_ Typer   = (*int64Values)(nil)
)

type int64Values []int64

func newInt64Values(p *[]int64) *int64Values {
	return (*int64Values)(p)
}

func (vs *int64Values) Set(val string) error {
	rest := val
	for rest != "" {
		idx := strings.IndexByte(rest, ',')
		if idx != -1 {
			val = rest[:idx]
			rest = rest[idx+1:]
		} else {
			val = rest
			rest = ""
		}

		var def int64
		*vs = append(*vs, def)
		if err := (*int64Value)(&(*vs)[len(*vs)-1]).Set(val); err != nil {
			return err
		}
	}

	return nil
}

func (vs *int64Values) String() string {
	if len(*vs) == 0 {
		return ""
	}

	var buf strings.Builder
	_, _ = buf.WriteString((*int64Value)(&(*vs)[0]).String())

	for i := 1; i < len(*vs); i++ {
		_ = buf.WriteByte(',')
		_, _ = buf.WriteString((*int64Value)(&(*vs)[i]).String())
	}

	return buf.String()
}

func (v *int64Values) Empty() bool { return len(*v) == 0 }

func (v *int64Values) Get() interface{} { return []int64(*v) }

func (*int64Values) Type() string { return "[]int64" }

// []float32

var (
	_ Value   = (*float32Values)(nil)
	_ Getter  = (*float32Values)(nil)
	_ Emptier = (*float32Values)(nil)
	_ Typer   = (*float32Values)(nil)
)

type float32Values []float32

func newFloat32Values(p *[]float32) *float32Values {
	return (*float32Values)(p)
}

func (vs *float32Values) Set(val string) error {
	rest := val
	for rest != "" {
		idx := strings.IndexByte(rest, ',')
		if idx != -1 {
			val = rest[:idx]
			rest = rest[idx+1:]
		} else {
			val = rest
			rest = ""
		}

		var def float32
		*vs = append(*vs, def)
		if err := (*float32Value)(&(*vs)[len(*vs)-1]).Set(val); err != nil {
			return err
		}
	}

	return nil
}

func (vs *float32Values) String() string {
	if len(*vs) == 0 {
		return ""
	}

	var buf strings.Builder
	_, _ = buf.WriteString((*float32Value)(&(*vs)[0]).String())

	for i := 1; i < len(*vs); i++ {
		_ = buf.WriteByte(',')
		_, _ = buf.WriteString((*float32Value)(&(*vs)[i]).String())
	}

	return buf.String()
}

func (v *float32Values) Empty() bool { return len(*v) == 0 }

func (v *float32Values) Get() interface{} { return []float32(*v) }

func (*float32Values) Type() string { return "[]float32" }

// []float64

var (
	_ Value   = (*float64Values)(nil)
	_ Getter  = (*float64Values)(nil)
	_ Emptier = (*float64Values)(nil)
	_ Typer   = (*float64Values)(nil)
)

type float64Values []float64

func newFloat64Values(p *[]float64) *float64Values {
	return (*float64Values)(p)
}

func (vs *float64Values) Set(val string) error {
	rest := val
	for rest != "" {
		idx := strings.IndexByte(rest, ',')
		if idx != -1 {
			val = rest[:idx]
			rest = rest[idx+1:]
		} else {
			val = rest
			rest = ""
		}

		var def float64
		*vs = append(*vs, def)
		if err := (*float64Value)(&(*vs)[len(*vs)-1]).Set(val); err != nil {
			return err
		}
	}

	return nil
}

func (vs *float64Values) String() string {
	if len(*vs) == 0 {
		return ""
	}

	var buf strings.Builder
	_, _ = buf.WriteString((*float64Value)(&(*vs)[0]).String())

	for i := 1; i < len(*vs); i++ {
		_ = buf.WriteByte(',')
		_, _ = buf.WriteString((*float64Value)(&(*vs)[i]).String())
	}

	return buf.String()
}

func (v *float64Values) Empty() bool { return len(*v) == 0 }

func (v *float64Values) Get() interface{} { return []float64(*v) }

func (*float64Values) Type() string { return "[]float64" }

// []string

var (
	_ Value   = (*stringValues)(nil)
	_ Getter  = (*stringValues)(nil)
	_ Emptier = (*stringValues)(nil)
	_ Typer   = (*stringValues)(nil)
)

type stringValues []string

func newStringValues(p *[]string) *stringValues {
	return (*stringValues)(p)
}

func (vs *stringValues) Set(val string) error {
	rest := val
	for rest != "" {
		idx := strings.IndexByte(rest, ',')
		if idx != -1 {
			val = rest[:idx]
			rest = rest[idx+1:]
		} else {
			val = rest
			rest = ""
		}

		var def string
		*vs = append(*vs, def)
		if err := (*stringValue)(&(*vs)[len(*vs)-1]).Set(val); err != nil {
			return err
		}
	}

	return nil
}

func (vs *stringValues) String() string {
	if len(*vs) == 0 {
		return ""
	}

	var buf strings.Builder
	_, _ = buf.WriteString((*stringValue)(&(*vs)[0]).String())

	for i := 1; i < len(*vs); i++ {
		_ = buf.WriteByte(',')
		_, _ = buf.WriteString((*stringValue)(&(*vs)[i]).String())
	}

	return buf.String()
}

func (v *stringValues) Empty() bool { return len(*v) == 0 }

func (v *stringValues) Get() interface{} { return []string(*v) }

func (*stringValues) Type() string { return "[]string" }

// []int

var (
	_ Value   = (*intValues)(nil)
	_ Getter  = (*intValues)(nil)
	_ Emptier = (*intValues)(nil)
	_ Typer   = (*intValues)(nil)
)

type intValues []int

func newIntValues(p *[]int) *intValues {
	return (*intValues)(p)
}

func (vs *intValues) Set(val string) error {
	rest := val
	for rest != "" {
		idx := strings.IndexByte(rest, ',')
		if idx != -1 {
			val = rest[:idx]
			rest = rest[idx+1:]
		} else {
			val = rest
			rest = ""
		}

		var def int
		*vs = append(*vs, def)
		if err := (*intValue)(&(*vs)[len(*vs)-1]).Set(val); err != nil {
			return err
		}
	}

	return nil
}

func (vs *intValues) String() string {
	if len(*vs) == 0 {
		return ""
	}

	var buf strings.Builder
	_, _ = buf.WriteString((*intValue)(&(*vs)[0]).String())

	for i := 1; i < len(*vs); i++ {
		_ = buf.WriteByte(',')
		_, _ = buf.WriteString((*intValue)(&(*vs)[i]).String())
	}

	return buf.String()
}

func (v *intValues) Empty() bool { return len(*v) == 0 }

func (v *intValues) Get() interface{} { return []int(*v) }

func (*intValues) Type() string { return "[]int" }

// []uint

var (
	_ Value   = (*uintValues)(nil)
	_ Getter  = (*uintValues)(nil)
	_ Emptier = (*uintValues)(nil)
	_ Typer   = (*uintValues)(nil)
)

type uintValues []uint

func newUintValues(p *[]uint) *uintValues {
	return (*uintValues)(p)
}

func (vs *uintValues) Set(val string) error {
	rest := val
	for rest != "" {
		idx := strings.IndexByte(rest, ',')
		if idx != -1 {
			val = rest[:idx]
			rest = rest[idx+1:]
		} else {
			val = rest
			rest = ""
		}

		var def uint
		*vs = append(*vs, def)
		if err := (*uintValue)(&(*vs)[len(*vs)-1]).Set(val); err != nil {
			return err
		}
	}

	return nil
}

func (vs *uintValues) String() string {
	if len(*vs) == 0 {
		return ""
	}

	var buf strings.Builder
	_, _ = buf.WriteString((*uintValue)(&(*vs)[0]).String())

	for i := 1; i < len(*vs); i++ {
		_ = buf.WriteByte(',')
		_, _ = buf.WriteString((*uintValue)(&(*vs)[i]).String())
	}

	return buf.String()
}

func (v *uintValues) Empty() bool { return len(*v) == 0 }

func (v *uintValues) Get() interface{} { return []uint(*v) }

func (*uintValues) Type() string { return "[]uint" }

// []time.Duration

var (
	_ Value   = (*timeDurationValues)(nil)
	_ Getter  = (*timeDurationValues)(nil)
	_ Emptier = (*timeDurationValues)(nil)
	_ Typer   = (*timeDurationValues)(nil)
)

type timeDurationValues []time.Duration

func newDurationValues(p *[]time.Duration) *timeDurationValues {
	return (*timeDurationValues)(p)
}

func (vs *timeDurationValues) Set(val string) error {
	rest := val
	for rest != "" {
		idx := strings.IndexByte(rest, ',')
		if idx != -1 {
			val = rest[:idx]
			rest = rest[idx+1:]
		} else {
			val = rest
			rest = ""
		}

		var def time.Duration
		*vs = append(*vs, def)
		if err := (*timeDurationValue)(&(*vs)[len(*vs)-1]).Set(val); err != nil {
			return err
		}
	}

	return nil
}

func (vs *timeDurationValues) String() string {
	if len(*vs) == 0 {
		return ""
	}

	var buf strings.Builder
	_, _ = buf.WriteString((*timeDurationValue)(&(*vs)[0]).String())

	for i := 1; i < len(*vs); i++ {
		_ = buf.WriteByte(',')
		_, _ = buf.WriteString((*timeDurationValue)(&(*vs)[i]).String())
	}

	return buf.String()
}

func (v *timeDurationValues) Empty() bool { return len(*v) == 0 }

func (v *timeDurationValues) Get() interface{} { return []time.Duration(*v) }

func (*timeDurationValues) Type() string { return "[]time.Duration" }
