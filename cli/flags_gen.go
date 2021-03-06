// Code generated by generate_flags.py; DO NOT EDIT.

package cli

import (
	"time"
)

// bool

// BoolVar defines a bool flag with specified name.
// The argument p points to a bool variable in which to store the value of the flag.
// The return value will be an error from the register.RegisterFlag if it
// failed to register the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.BoolVar(register, &p, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.BoolVar(register, &p, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.BoolVar(register, &p, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.BoolVar(register, &p, "name", cli.Required)
//
// All options can be used together.
func BoolVar(register Register, p *bool, name string, options ...FlagOptionApplyer) error {
	return Var(register, newBoolValue(p), name, options...)
}

// Bool defines a bool flag with specified name.
// The return value is the address of a bool variable that stores the value of the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Bool(register, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Bool(register, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Bool(register, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Bool(register, "name", cli.Required)
//
// All options can be used together.
func Bool(register Register, name string, options ...FlagOptionApplyer) *bool {
	p := new(bool)
	_ = BoolVar(register, p, name, options...)
	return p
}

// uint8

// Uint8Var defines a uint8 flag with specified name.
// The argument p points to a uint8 variable in which to store the value of the flag.
// The return value will be an error from the register.RegisterFlag if it
// failed to register the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Uint8Var(register, &p, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Uint8Var(register, &p, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Uint8Var(register, &p, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Uint8Var(register, &p, "name", cli.Required)
//
// All options can be used together.
func Uint8Var(register Register, p *uint8, name string, options ...FlagOptionApplyer) error {
	return Var(register, newUint8Value(p), name, options...)
}

// Uint8 defines a uint8 flag with specified name.
// The return value is the address of a uint8 variable that stores the value of the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Uint8(register, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Uint8(register, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Uint8(register, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Uint8(register, "name", cli.Required)
//
// All options can be used together.
func Uint8(register Register, name string, options ...FlagOptionApplyer) *uint8 {
	p := new(uint8)
	_ = Uint8Var(register, p, name, options...)
	return p
}

// uint16

// Uint16Var defines a uint16 flag with specified name.
// The argument p points to a uint16 variable in which to store the value of the flag.
// The return value will be an error from the register.RegisterFlag if it
// failed to register the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Uint16Var(register, &p, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Uint16Var(register, &p, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Uint16Var(register, &p, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Uint16Var(register, &p, "name", cli.Required)
//
// All options can be used together.
func Uint16Var(register Register, p *uint16, name string, options ...FlagOptionApplyer) error {
	return Var(register, newUint16Value(p), name, options...)
}

// Uint16 defines a uint16 flag with specified name.
// The return value is the address of a uint16 variable that stores the value of the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Uint16(register, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Uint16(register, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Uint16(register, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Uint16(register, "name", cli.Required)
//
// All options can be used together.
func Uint16(register Register, name string, options ...FlagOptionApplyer) *uint16 {
	p := new(uint16)
	_ = Uint16Var(register, p, name, options...)
	return p
}

// uint32

// Uint32Var defines a uint32 flag with specified name.
// The argument p points to a uint32 variable in which to store the value of the flag.
// The return value will be an error from the register.RegisterFlag if it
// failed to register the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Uint32Var(register, &p, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Uint32Var(register, &p, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Uint32Var(register, &p, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Uint32Var(register, &p, "name", cli.Required)
//
// All options can be used together.
func Uint32Var(register Register, p *uint32, name string, options ...FlagOptionApplyer) error {
	return Var(register, newUint32Value(p), name, options...)
}

// Uint32 defines a uint32 flag with specified name.
// The return value is the address of a uint32 variable that stores the value of the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Uint32(register, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Uint32(register, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Uint32(register, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Uint32(register, "name", cli.Required)
//
// All options can be used together.
func Uint32(register Register, name string, options ...FlagOptionApplyer) *uint32 {
	p := new(uint32)
	_ = Uint32Var(register, p, name, options...)
	return p
}

// uint64

// Uint64Var defines a uint64 flag with specified name.
// The argument p points to a uint64 variable in which to store the value of the flag.
// The return value will be an error from the register.RegisterFlag if it
// failed to register the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Uint64Var(register, &p, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Uint64Var(register, &p, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Uint64Var(register, &p, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Uint64Var(register, &p, "name", cli.Required)
//
// All options can be used together.
func Uint64Var(register Register, p *uint64, name string, options ...FlagOptionApplyer) error {
	return Var(register, newUint64Value(p), name, options...)
}

// Uint64 defines a uint64 flag with specified name.
// The return value is the address of a uint64 variable that stores the value of the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Uint64(register, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Uint64(register, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Uint64(register, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Uint64(register, "name", cli.Required)
//
// All options can be used together.
func Uint64(register Register, name string, options ...FlagOptionApplyer) *uint64 {
	p := new(uint64)
	_ = Uint64Var(register, p, name, options...)
	return p
}

// int8

// Int8Var defines a int8 flag with specified name.
// The argument p points to a int8 variable in which to store the value of the flag.
// The return value will be an error from the register.RegisterFlag if it
// failed to register the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Int8Var(register, &p, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Int8Var(register, &p, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Int8Var(register, &p, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Int8Var(register, &p, "name", cli.Required)
//
// All options can be used together.
func Int8Var(register Register, p *int8, name string, options ...FlagOptionApplyer) error {
	return Var(register, newInt8Value(p), name, options...)
}

// Int8 defines a int8 flag with specified name.
// The return value is the address of a int8 variable that stores the value of the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Int8(register, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Int8(register, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Int8(register, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Int8(register, "name", cli.Required)
//
// All options can be used together.
func Int8(register Register, name string, options ...FlagOptionApplyer) *int8 {
	p := new(int8)
	_ = Int8Var(register, p, name, options...)
	return p
}

// int16

// Int16Var defines a int16 flag with specified name.
// The argument p points to a int16 variable in which to store the value of the flag.
// The return value will be an error from the register.RegisterFlag if it
// failed to register the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Int16Var(register, &p, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Int16Var(register, &p, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Int16Var(register, &p, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Int16Var(register, &p, "name", cli.Required)
//
// All options can be used together.
func Int16Var(register Register, p *int16, name string, options ...FlagOptionApplyer) error {
	return Var(register, newInt16Value(p), name, options...)
}

// Int16 defines a int16 flag with specified name.
// The return value is the address of a int16 variable that stores the value of the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Int16(register, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Int16(register, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Int16(register, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Int16(register, "name", cli.Required)
//
// All options can be used together.
func Int16(register Register, name string, options ...FlagOptionApplyer) *int16 {
	p := new(int16)
	_ = Int16Var(register, p, name, options...)
	return p
}

// int32

// Int32Var defines a int32 flag with specified name.
// The argument p points to a int32 variable in which to store the value of the flag.
// The return value will be an error from the register.RegisterFlag if it
// failed to register the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Int32Var(register, &p, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Int32Var(register, &p, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Int32Var(register, &p, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Int32Var(register, &p, "name", cli.Required)
//
// All options can be used together.
func Int32Var(register Register, p *int32, name string, options ...FlagOptionApplyer) error {
	return Var(register, newInt32Value(p), name, options...)
}

// Int32 defines a int32 flag with specified name.
// The return value is the address of a int32 variable that stores the value of the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Int32(register, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Int32(register, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Int32(register, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Int32(register, "name", cli.Required)
//
// All options can be used together.
func Int32(register Register, name string, options ...FlagOptionApplyer) *int32 {
	p := new(int32)
	_ = Int32Var(register, p, name, options...)
	return p
}

// int64

// Int64Var defines a int64 flag with specified name.
// The argument p points to a int64 variable in which to store the value of the flag.
// The return value will be an error from the register.RegisterFlag if it
// failed to register the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Int64Var(register, &p, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Int64Var(register, &p, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Int64Var(register, &p, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Int64Var(register, &p, "name", cli.Required)
//
// All options can be used together.
func Int64Var(register Register, p *int64, name string, options ...FlagOptionApplyer) error {
	return Var(register, newInt64Value(p), name, options...)
}

// Int64 defines a int64 flag with specified name.
// The return value is the address of a int64 variable that stores the value of the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Int64(register, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Int64(register, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Int64(register, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Int64(register, "name", cli.Required)
//
// All options can be used together.
func Int64(register Register, name string, options ...FlagOptionApplyer) *int64 {
	p := new(int64)
	_ = Int64Var(register, p, name, options...)
	return p
}

// float32

// Float32Var defines a float32 flag with specified name.
// The argument p points to a float32 variable in which to store the value of the flag.
// The return value will be an error from the register.RegisterFlag if it
// failed to register the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Float32Var(register, &p, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Float32Var(register, &p, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Float32Var(register, &p, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Float32Var(register, &p, "name", cli.Required)
//
// All options can be used together.
func Float32Var(register Register, p *float32, name string, options ...FlagOptionApplyer) error {
	return Var(register, newFloat32Value(p), name, options...)
}

// Float32 defines a float32 flag with specified name.
// The return value is the address of a float32 variable that stores the value of the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Float32(register, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Float32(register, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Float32(register, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Float32(register, "name", cli.Required)
//
// All options can be used together.
func Float32(register Register, name string, options ...FlagOptionApplyer) *float32 {
	p := new(float32)
	_ = Float32Var(register, p, name, options...)
	return p
}

// float64

// Float64Var defines a float64 flag with specified name.
// The argument p points to a float64 variable in which to store the value of the flag.
// The return value will be an error from the register.RegisterFlag if it
// failed to register the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Float64Var(register, &p, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Float64Var(register, &p, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Float64Var(register, &p, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Float64Var(register, &p, "name", cli.Required)
//
// All options can be used together.
func Float64Var(register Register, p *float64, name string, options ...FlagOptionApplyer) error {
	return Var(register, newFloat64Value(p), name, options...)
}

// Float64 defines a float64 flag with specified name.
// The return value is the address of a float64 variable that stores the value of the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Float64(register, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Float64(register, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Float64(register, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Float64(register, "name", cli.Required)
//
// All options can be used together.
func Float64(register Register, name string, options ...FlagOptionApplyer) *float64 {
	p := new(float64)
	_ = Float64Var(register, p, name, options...)
	return p
}

// string

// StringVar defines a string flag with specified name.
// The argument p points to a string variable in which to store the value of the flag.
// The return value will be an error from the register.RegisterFlag if it
// failed to register the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.StringVar(register, &p, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.StringVar(register, &p, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.StringVar(register, &p, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.StringVar(register, &p, "name", cli.Required)
//
// All options can be used together.
func StringVar(register Register, p *string, name string, options ...FlagOptionApplyer) error {
	return Var(register, newStringValue(p), name, options...)
}

// String defines a string flag with specified name.
// The return value is the address of a string variable that stores the value of the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.String(register, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.String(register, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.String(register, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.String(register, "name", cli.Required)
//
// All options can be used together.
func String(register Register, name string, options ...FlagOptionApplyer) *string {
	p := new(string)
	_ = StringVar(register, p, name, options...)
	return p
}

// int

// IntVar defines a int flag with specified name.
// The argument p points to a int variable in which to store the value of the flag.
// The return value will be an error from the register.RegisterFlag if it
// failed to register the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.IntVar(register, &p, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.IntVar(register, &p, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.IntVar(register, &p, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.IntVar(register, &p, "name", cli.Required)
//
// All options can be used together.
func IntVar(register Register, p *int, name string, options ...FlagOptionApplyer) error {
	return Var(register, newIntValue(p), name, options...)
}

// Int defines a int flag with specified name.
// The return value is the address of a int variable that stores the value of the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Int(register, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Int(register, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Int(register, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Int(register, "name", cli.Required)
//
// All options can be used together.
func Int(register Register, name string, options ...FlagOptionApplyer) *int {
	p := new(int)
	_ = IntVar(register, p, name, options...)
	return p
}

// uint

// UintVar defines a uint flag with specified name.
// The argument p points to a uint variable in which to store the value of the flag.
// The return value will be an error from the register.RegisterFlag if it
// failed to register the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.UintVar(register, &p, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.UintVar(register, &p, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.UintVar(register, &p, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.UintVar(register, &p, "name", cli.Required)
//
// All options can be used together.
func UintVar(register Register, p *uint, name string, options ...FlagOptionApplyer) error {
	return Var(register, newUintValue(p), name, options...)
}

// Uint defines a uint flag with specified name.
// The return value is the address of a uint variable that stores the value of the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Uint(register, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Uint(register, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Uint(register, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Uint(register, "name", cli.Required)
//
// All options can be used together.
func Uint(register Register, name string, options ...FlagOptionApplyer) *uint {
	p := new(uint)
	_ = UintVar(register, p, name, options...)
	return p
}

// time.Duration

// DurationVar defines a time.Duration flag with specified name.
// The argument p points to a time.Duration variable in which to store the value of the flag.
// The return value will be an error from the register.RegisterFlag if it
// failed to register the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.DurationVar(register, &p, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.DurationVar(register, &p, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.DurationVar(register, &p, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.DurationVar(register, &p, "name", cli.Required)
//
// All options can be used together.
func DurationVar(register Register, p *time.Duration, name string, options ...FlagOptionApplyer) error {
	return Var(register, newDurationValue(p), name, options...)
}

// Duration defines a time.Duration flag with specified name.
// The return value is the address of a time.Duration variable that stores the value of the flag.
//
// If a name contains only one rune, it will be a short name, otherwise a long name.
// To set a short name, you pass a cli.WithShort.
//
//   _ = cli.Duration(register, "name", cli.WithShort("n"))
//
// To set a long name, you pass a cli.WithLong.
//
//   _ = cli.Duration(register, "n", cli.WithLong("name"))
//
// A usage may be set by passing a cli.Usage.
//
//   _ = cli.Duration(register, "name", cli.Usage("The name of user"))
//
// The flag is optional by default.
// This may be changed by passing the cli.Required.
//
//   _ = cli.Duration(register, "name", cli.Required)
//
// All options can be used together.
func Duration(register Register, name string, options ...FlagOptionApplyer) *time.Duration {
	p := new(time.Duration)
	_ = DurationVar(register, p, name, options...)
	return p
}
