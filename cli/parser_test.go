package cli

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"testing"
	"time"
	"unsafe"
)

const (
	maxUint = ^uint(0)
	minUint = 0
	maxInt  = int(maxUint >> 1)
	minInt  = -maxInt - 1
)

func deref(t *testing.T, ptr interface{}) interface{} {
	t.Helper()

	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr {
		t.Fatalf("flag is not a pointer: %T", v.String())
	}

	return v.Elem().Interface()
}

type commonValue struct {
	value string
	want  interface{}
}

var (
	commonBoolValues = []commonValue{
		{"true", true},
		{"Y", true},
		{"1", true},
		{"false", false},
		{"F", false},
		{"0", false},
	}

	commonFloat64Values = []commonValue{
		{"0", 0.0},
		{"-0", 0.0},
		{"1337", 1337.0},
		{"-7331", -7331.0},
		{strconv.FormatFloat(math.MaxFloat64, 'g', -1, 64), math.MaxFloat64},
		{strconv.FormatFloat(math.SmallestNonzeroFloat64, 'g', -1, 64), math.SmallestNonzeroFloat64},
	}

	commonIntValues = []commonValue{
		{"0", 0},
		{"-0", 0},
		{"1337", 1337},
		{"-7331", -7331},
		{"0xABC", 0xABC},
		{"-0xCBA", -0xCBA},
		{"0b10111011", 0b10111011},
		{"-0b11011101", -0b11011101},
		{strconv.FormatInt(int64(maxInt), 10), maxInt},
		{strconv.FormatInt(int64(minInt), 10), minInt},
	}

	commonUintValues = []commonValue{
		{"0", uint(0)},
		{"1337", uint(1337)},
		{"0xABC", uint(0xABC)},
		{"0b10111011", uint(0b10111011)},
		{strconv.FormatUint(uint64(maxUint), 10), uint(maxUint)},
		{strconv.FormatUint(uint64(minUint), 10), uint(minUint)},
	}

	commonDurationValues = []commonValue{
		{"0ns", 0 * time.Nanosecond},
		{"123ns", 123 * time.Nanosecond},
		{"0s", 0 * time.Second},
		{"-5s", -5 * time.Second},
		{"2h45m", 2*time.Hour + 45*time.Minute},
		{"+2h45m", 2*time.Hour + 45*time.Minute},
		{"4.1s", 4*time.Second + 100*time.Millisecond},
	}
)

type commonBroken struct {
	name  string
	value string
	want  error
}

const (
	float32MaxOverflowValue = "3.40282e+39"          // 3.40282e+38
	float64MaxOverflowValue = "1.79769e+309"         // 1.79769e+308
	int32MaxOverflowValue   = "2147483648"           // 2147483647
	int64MaxOverflowValue   = "9223372036854775808"  // 9223372036854775807
	int32MinOverflowValue   = "-2147483649"          // -2147483648
	int64MinOverflowValue   = "9223372036854775809"  // -9223372036854775808
	uint32MaxOverflowValue  = "4294967296"           // 4294967295
	uint64MaxOverflowValue  = "18446744073709551616" // 18446744073709551615
)

func intMaxOverflowValue() string {
	if unsafe.Sizeof(int(0)) == unsafe.Sizeof(int32(0)) {
		return int32MaxOverflowValue
	}

	return int64MaxOverflowValue
}

func intMinOverflowValue() string {
	if unsafe.Sizeof(int(0)) == unsafe.Sizeof(int32(0)) {
		return int32MinOverflowValue
	}

	return int64MinOverflowValue
}

func uintMaxOverflowValue() string {
	if unsafe.Sizeof(uint(0)) == unsafe.Sizeof(uint32(0)) {
		return uint32MaxOverflowValue
	}

	return uint64MaxOverflowValue
}

var (
	commonFloat64Brokens = []commonBroken{
		{"empty", "", &ParseValueError{Type: "float64", Err: ErrSyntax}},
		{"not float64-like", "abcd", &ParseValueError{Type: "float64", Err: ErrSyntax}},
		{"broken float64", "12.43a", &ParseValueError{Type: "float64", Err: ErrSyntax}},
		{"true", "true", &ParseValueError{Type: "float64", Err: ErrSyntax}},
		{"false", "false", &ParseValueError{Type: "float64", Err: ErrSyntax}},
		{"float64 max overflow", float64MaxOverflowValue, &ParseValueError{Type: "float64", Err: ErrRange}},
	}

	commonIntBrokens = []commonBroken{
		{"empty", "", &ParseValueError{Type: "int", Err: ErrSyntax}},
		{"not int-like", "abcd", &ParseValueError{Type: "int", Err: ErrSyntax}},
		{"broken int", "1337a", &ParseValueError{Type: "int", Err: ErrSyntax}},
		{"true", "true", &ParseValueError{Type: "int", Err: ErrSyntax}},
		{"false", "false", &ParseValueError{Type: "int", Err: ErrSyntax}},
		{"float", "12.34", &ParseValueError{Type: "int", Err: ErrSyntax}},
		{"negative float", "-43.21", &ParseValueError{Type: "int", Err: ErrSyntax}},
		{"int max overflow", intMaxOverflowValue(), &ParseValueError{Type: "int", Err: ErrRange}},
		{"int min overflow", intMinOverflowValue(), &ParseValueError{Type: "int", Err: ErrRange}},
	}

	commonUintBrokens = []commonBroken{
		{"empty", "", &ParseValueError{Type: "uint", Err: ErrSyntax}},
		{"not uint-like", "abcd", &ParseValueError{Type: "uint", Err: ErrSyntax}},
		{"broken uint", "1337a", &ParseValueError{Type: "uint", Err: ErrSyntax}},
		{"true", "true", &ParseValueError{Type: "uint", Err: ErrSyntax}},
		{"false", "false", &ParseValueError{Type: "uint", Err: ErrSyntax}},
		{"negative int", "-7331", &ParseValueError{Type: "uint", Err: ErrSyntax}},
		{"float", "12.34", &ParseValueError{Type: "uint", Err: ErrSyntax}},
		{"negative float", "-43.21", &ParseValueError{Type: "uint", Err: ErrSyntax}},
		{"uint max overflow", uintMaxOverflowValue(), &ParseValueError{Type: "uint", Err: ErrRange}},
		{"uint min overflow", "-0", &ParseValueError{Type: "uint", Err: ErrSyntax}},
	}

	commonDurationBrokens = []commonBroken{
		{"empty", "", &ParseValueError{Type: "time.Duration", Err: ErrSyntax}},
		{"not duration-like", "100", &ParseValueError{Type: "time.Duration", Err: ErrSyntax}},
		{"broken duration", "100sm", &ParseValueError{Type: "time.Duration", Err: ErrSyntax}},
		{"true", "true", &ParseValueError{Type: "time.Duration", Err: ErrSyntax}},
		{"false", "false", &ParseValueError{Type: "time.Duration", Err: ErrSyntax}},
	}
)

func TestParseFlags(t *testing.T) {
	type testValue struct {
		name       string
		extraSetup func(Register) interface{}
		extraCheck func(t *testing.T, v interface{})
		args       []string
		want       interface{}
	}

	mergeTestValues := func(tvss ...[]testValue) []testValue {
		t.Helper()

		var all []testValue

		for _, tvs := range tvss {
			all = append(all, tvs...)
		}

		return all
	}

	commonValuesToTestValues := func(vals []commonValue) []testValue {
		t.Helper()

		tvs := make([]testValue, 0, len(vals)*2)

		// value.
		for _, v := range vals {
			tvs = append(tvs, testValue{
				name: v.value + " value",
				args: []string{"-t=" + v.value},
				want: v.want,
			})
		}

		// next arg
		for _, v := range vals {
			tvs = append(tvs, testValue{
				name: v.value + " next arg",
				args: []string{"-t", v.value},
				want: v.want,
			})
		}

		return tvs
	}

	tt := []struct {
		name  string
		setup func(Register) interface{}
		tests []testValue
	}{
		{
			name:  "Bool",
			setup: func(r Register) interface{} { return Bool(r, "t") },
			tests: mergeTestValues(
				[]testValue{
					{
						name: "without value",
						args: []string{"-t"},
						want: true,
					},
					{
						name: "empty value",
						args: []string{"-t="},
						want: false,
					},
					{
						name:       "skip not bool-like next arg",
						extraSetup: func(r Register) interface{} { return StringArg(r, "extra") },
						extraCheck: func(t *testing.T, got interface{}) {
							t.Helper()

							const want = "abcd"
							if !reflect.DeepEqual(got, want) {
								t.Errorf("Parse(): extra: got = %#v, want = %#v", got, want)
							}
						},
						args: []string{"-t", "abcd"},
						want: true,
					},
				},
				commonValuesToTestValues(commonBoolValues),
			),
		},
		{
			name:  "Float64",
			setup: func(r Register) interface{} { return Float64(r, "t") },
			tests: mergeTestValues(
				commonValuesToTestValues(commonFloat64Values),
			),
		},
		{
			name:  "Int",
			setup: func(r Register) interface{} { return Int(r, "t") },
			tests: mergeTestValues(
				commonValuesToTestValues(commonIntValues),
			),
		},
		{
			name:  "Uint",
			setup: func(r Register) interface{} { return Uint(r, "t") },
			tests: mergeTestValues(
				commonValuesToTestValues(commonUintValues),
			),
		},
		{
			name:  "String",
			setup: func(r Register) interface{} { return String(r, "t") },
			tests: []testValue{
				{
					name: "test value",
					args: []string{"-t=test"},
					want: "test",
				},
				{
					name: "empty value",
					args: []string{"-t="},
					want: "",
				},
				{
					name: "single dash value",
					args: []string{"-t=-"},
					want: "-",
				},
				{
					name: "test next arg",
					args: []string{"-t", "test"},
					want: "test",
				},
				{
					name: "single dash next arg",
					args: []string{"-t", "-"},
					want: "-",
				},
				{
					name: "without value",
					args: []string{"-t"},
					want: "",
				},
				{
					name: "empty next arg",
					args: []string{"-t", ""},
					want: "",
				},
				{
					name:       "next flag",
					extraSetup: func(r Register) interface{} { return Bool(r, "b") },
					extraCheck: func(t *testing.T, got interface{}) {
						t.Helper()

						const want = true
						if !reflect.DeepEqual(got, want) {
							t.Errorf("Parse(): b: got = %#v, want = %#v", got, want)
						}
					},
					args: []string{"-t", "-b"},
					want: "",
				},
				{
					name: "with dash value",
					args: []string{"-t=go-test"},
					want: "go-test",
				},
				{
					name: "with start dash value",
					args: []string{"-t=-go-test"},
					want: "-go-test",
				},
				{
					name: "with equals value",
					args: []string{"-t=go=test"},
					want: "go=test",
				},
				{
					name: "with start equals value",
					args: []string{"-t==go=test"},
					want: "=go=test",
				},
				{
					name: "with dash next arg",
					args: []string{"-t", "go-test"},
					want: "go-test",
				},
				{
					name: "with equals next arg",
					args: []string{"-t", "go=test"},
					want: "go=test",
				},
				{
					name: "with equals start next arg",
					args: []string{"-t", "=go=test"},
					want: "=go=test",
				},
			},
		},
		{
			name:  "Duration",
			setup: func(r Register) interface{} { return Duration(r, "t") },
			tests: mergeTestValues(
				commonValuesToTestValues(commonDurationValues),
			),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			for _, tvc := range tc.tests {
				t.Run(tvc.name, func(t *testing.T) {
					var (
						register DefaultRegister
						parser   DefaultParser
					)

					if len(tvc.args) > 1 && tvc.args[1] == "-5s" {
						t.Log("g")
					}

					f := tc.setup(&register)

					var extra interface{}
					if tvc.extraSetup != nil {
						extra = tvc.extraSetup(&register)
					}

					if err := parser.Parse(nil, &register, tvc.args); err != nil {
						t.Fatalf("Parse(%v): failed to parse flags: %s", tvc.args, err)
					}

					got := deref(t, f)
					if !reflect.DeepEqual(got, tvc.want) {
						t.Errorf("Parse(%v): got = %#v, want = %#v", tvc.args, got, tvc.want)
					}

					if tvc.extraCheck != nil {
						gotExtra := deref(t, extra)
						tvc.extraCheck(t, gotExtra)
					}
				})
			}
		})
	}
}

func TestParseFlags_broken_value(t *testing.T) {
	type testValue struct {
		name string
		args []string
		want error
	}

	mergeTestValues := func(tvss ...[]testValue) []testValue {
		t.Helper()

		var all []testValue

		for _, tvs := range tvss {
			all = append(all, tvs...)
		}

		return all
	}

	commonBrokensToTestValues := func(vals []commonBroken) []testValue {
		t.Helper()

		tvs := make([]testValue, 0, len(vals)*2)

		// value.
		for _, v := range vals {
			tvs = append(tvs, testValue{
				name: v.name + " value",
				args: []string{"-t=" + v.value},
				want: &FlagError{
					Short: "t",
					Err:   v.want,
				},
			})
		}

		// next arg
		for _, v := range vals {
			tvs = append(tvs, testValue{
				name: v.name + " next arg",
				args: []string{"-t", v.value},
				want: &FlagError{
					Short: "t",
					Err:   v.want,
				},
			})
		}

		return tvs
	}

	tt := []struct {
		name  string
		setup func(Register) interface{}
		tests []testValue
	}{
		{
			name:  "Bool",
			setup: func(r Register) interface{} { return Bool(r, "t") },
			tests: []testValue{
				{
					name: "not bool-like value",
					args: []string{"-t=abcd"},
					want: &FlagError{
						Short: "t",
						Err: &ParseValueError{
							Type: "bool",
							Err:  ErrSyntax,
						},
					},
				},
				{
					name: "not bool-like value 2",
					args: []string{"-t=2"},
					want: &FlagError{
						Short: "t",
						Err: &ParseValueError{
							Type: "bool",
							Err:  ErrSyntax,
						},
					},
				},
			},
		},
		{
			name:  "Float64",
			setup: func(r Register) interface{} { return Float64(r, "t") },
			tests: mergeTestValues(
				[]testValue{
					{
						name: "without value",
						args: []string{"-t"},
						want: &FlagError{
							Short: "t",
							Err: &ParseValueError{
								Type: "float64",
								Err:  ErrSyntax,
							},
						},
					},
				},
				commonBrokensToTestValues(commonFloat64Brokens),
			),
		},
		{
			name:  "Int",
			setup: func(r Register) interface{} { return Int(r, "t") },
			tests: mergeTestValues(
				[]testValue{
					{
						name: "without value",
						args: []string{"-t"},
						want: &FlagError{
							Short: "t",
							Err: &ParseValueError{
								Type: "int",
								Err:  ErrSyntax,
							},
						},
					},
				},
				commonBrokensToTestValues(commonIntBrokens),
			),
		},
		{
			name:  "Uint",
			setup: func(r Register) interface{} { return Uint(r, "t") },
			tests: mergeTestValues(
				[]testValue{
					{
						name: "without value",
						args: []string{"-t"},
						want: &FlagError{
							Short: "t",
							Err: &ParseValueError{
								Type: "uint",
								Err:  ErrSyntax,
							},
						},
					},
				},
				commonBrokensToTestValues(commonUintBrokens),
			),
		},
		{
			name:  "String",
			setup: func(r Register) interface{} { return String(r, "t") },
		},
		{
			name:  "Duration",
			setup: func(r Register) interface{} { return Duration(r, "t") },
			tests: mergeTestValues(
				[]testValue{
					{
						name: "without value",
						args: []string{"-t"},
						want: &FlagError{
							Short: "t",
							Err: &ParseValueError{
								Type: "time.Duration",
								Err:  ErrSyntax,
							},
						},
					},
				},
				commonBrokensToTestValues(commonDurationBrokens),
			),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			for _, tvc := range tc.tests {
				t.Run(tvc.name, func(t *testing.T) {
					var (
						register DefaultRegister
						parser   DefaultParser
					)

					_ = tc.setup(&register)

					err := parser.Parse(nil, &register, tvc.args)
					if !errors.Is(err, tvc.want) {
						t.Fatalf("Parse(%v): got error = %q, want error = %q", tvc.args, err, tvc.want)
					}
				})
			}
		})
	}
}

func TestParseMultiFlags(t *testing.T) {
	var (
		register DefaultRegister
		parser   DefaultParser
	)

	a := Bool(&register, "a")
	b := Ints(&register, "b")
	c := String(&register, "c")
	d := Int(&register, "d")

	args := []string{
		"-a", "-b", "1", "-d=99", "-b=2,-3,4", "-c", "test", "-b", "5,6", "-b", "7",
		"-d", "8", "-b=9",
	}

	if err := parser.Parse(nil, &register, args); err != nil {
		t.Fatalf("Parse(): failed to parse args: %s", err)
	}

	// Check flags.
	const (
		wantA = true
		wantC = "test"
		wantD = 8
	)
	wantB := []int{1, 2, -3, 4, 5, 6, 7, 9}

	if *a != wantA {
		t.Errorf("Parse(): a: got = %v, want = %v", *a, wantA)
	}

	if !reflect.DeepEqual(*b, wantB) {
		t.Errorf("Parse(): b: got = %#v, want = %#v", *b, wantB)
	}

	if *c != wantC {
		t.Errorf("Parse(): c: got = %v, want = %v", *c, wantC)
	}

	if *d != wantD {
		t.Errorf("Parse(): a: got = %v, want = %v", *d, wantD)
	}
}

func TestParseArgs(t *testing.T) {
	type testValue struct {
		name string
		args []string
		want interface{}
	}

	mergeTestValues := func(tvss ...[]testValue) []testValue {
		t.Helper()

		var all []testValue

		for _, tvs := range tvss {
			all = append(all, tvs...)
		}

		return all
	}

	commonValuesToTestValues := func(vals []commonValue) []testValue {
		t.Helper()

		tvs := make([]testValue, 0, len(vals))

		for _, v := range vals {
			tvs = append(tvs, testValue{
				name: v.value,
				args: []string{v.value},
				want: v.want,
			})
		}

		return tvs
	}

	tt := []struct {
		name  string
		setup func(Register) interface{}
		tests []testValue
	}{
		{
			name:  "BoolArg",
			setup: func(r Register) interface{} { return BoolArg(r, "t") },
			tests: mergeTestValues(
				commonValuesToTestValues(commonBoolValues),
			),
		},
		{
			name:  "Float64Arg",
			setup: func(r Register) interface{} { return Float64Arg(r, "t") },
			tests: mergeTestValues(
				commonValuesToTestValues(commonFloat64Values),
			),
		},
		{
			name:  "IntArg",
			setup: func(r Register) interface{} { return IntArg(r, "t") },
			tests: mergeTestValues(
				commonValuesToTestValues(commonIntValues),
			),
		},
		{
			name:  "UintArg",
			setup: func(r Register) interface{} { return UintArg(r, "t") },
			tests: mergeTestValues(
				commonValuesToTestValues(commonUintValues),
			),
		},
		{
			name:  "StringArg",
			setup: func(r Register) interface{} { return StringArg(r, "t") },
			tests: []testValue{
				{
					name: "test",
					args: []string{"test"},
					want: "test",
				},
				{
					name: "empty",
					args: []string{""},
					want: "",
				},
				{
					name: "single dash",
					args: []string{"-"},
					want: "-",
				},
			},
		},
		{
			name:  "DurationArg",
			setup: func(r Register) interface{} { return DurationArg(r, "t") },
			tests: mergeTestValues(
				commonValuesToTestValues(commonDurationValues),
			),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			for _, tvc := range tc.tests {
				t.Run(tvc.name, func(t *testing.T) {
					var (
						register DefaultRegister
						parser   DefaultParser
					)

					f := tc.setup(&register)

					if err := parser.Parse(nil, &register, tvc.args); err != nil {
						t.Fatalf("Parse(%v): failed to parse args: %s", tvc.args, err)
					}

					got := deref(t, f)
					if !reflect.DeepEqual(got, tvc.want) {
						t.Errorf("Parse(%v): got = %#v, want = %#v", tvc.args, got, tvc.want)
					}
				})
			}
		})
	}
}

func TestParseArgs_broken_value(t *testing.T) {
	type testValue struct {
		name string
		args []string
		want error
	}

	mergeTestValues := func(tvss ...[]testValue) []testValue {
		t.Helper()

		var all []testValue

		for _, tvs := range tvss {
			all = append(all, tvs...)
		}

		return all
	}

	commonBrokensToTestValues := func(vals []commonBroken) []testValue {
		t.Helper()

		tvs := make([]testValue, 0, len(vals))

		for _, v := range vals {
			tvs = append(tvs, testValue{
				name: v.name,
				args: []string{v.value},
				want: &ArgError{
					Name:  "t",
					Index: 0,
					Err:   v.want,
				},
			})
		}

		return tvs
	}

	tt := []struct {
		name  string
		setup func(Register) interface{}
		tests []testValue
	}{
		{
			name:  "BoolArg",
			setup: func(r Register) interface{} { return BoolArg(r, "t") },
			tests: []testValue{
				{
					name: "empty",
					args: []string{""},
					want: &ArgError{
						Name:  "t",
						Index: 0,
						Err: &ParseValueError{
							Type: "bool",
							Err:  ErrSyntax,
						},
					},
				},
				{
					name: "not bool-like",
					args: []string{"abcd"},
					want: &ArgError{
						Name:  "t",
						Index: 0,
						Err: &ParseValueError{
							Type: "bool",
							Err:  ErrSyntax,
						},
					},
				},
				{
					name: "not bool-like 2",
					args: []string{"2"},
					want: &ArgError{
						Name:  "t",
						Index: 0,
						Err: &ParseValueError{
							Type: "bool",
							Err:  ErrSyntax,
						},
					},
				},
			},
		},
		{
			name:  "Float64Arg",
			setup: func(r Register) interface{} { return Float64Arg(r, "t") },
			tests: mergeTestValues(
				commonBrokensToTestValues(commonFloat64Brokens),
			),
		},
		{
			name:  "IntArg",
			setup: func(r Register) interface{} { return IntArg(r, "t") },
			tests: mergeTestValues(
				commonBrokensToTestValues(commonIntBrokens),
			),
		},
		{
			name:  "UintArg",
			setup: func(r Register) interface{} { return UintArg(r, "t") },
			tests: mergeTestValues(
				commonBrokensToTestValues(commonUintBrokens),
			),
		},
		{
			name:  "StringArg",
			setup: func(r Register) interface{} { return StringArg(r, "t") },
		},
		{
			name:  "DurationArg",
			setup: func(r Register) interface{} { return DurationArg(r, "t") },
			tests: mergeTestValues(
				commonBrokensToTestValues(commonDurationBrokens),
			),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			for _, tvc := range tc.tests {
				t.Run(tvc.name, func(t *testing.T) {
					var (
						register DefaultRegister
						parser   DefaultParser
					)

					_ = tc.setup(&register)

					err := parser.Parse(nil, &register, tvc.args)
					if !errors.Is(err, tvc.want) {
						t.Fatalf("Parse(%v): got error = %q, want error = %q", tvc.args, err, tvc.want)
					}
				})
			}
		})
	}
}

func TestRegisterInvalidNameFlag(t *testing.T) {
	tt := []struct {
		name  string
		short string
		long  string
		want  error
	}{
		{
			name: "empty short and long names",
			want: &FlagError{Err: ErrMissingName},
		},
		{
			name:  "too long short name",
			short: "he",
			long:  "help",
			want:  &FlagError{Short: "he", Long: "help", Err: ErrInvalidName},
		},
		{
			name:  "dash in short name",
			short: "-",
			want:  &FlagError{Short: "-", Err: ErrInvalidName},
		},
		{
			name:  "equal in short name",
			short: "=",
			want:  &FlagError{Short: "=", Err: ErrInvalidName},
		},
		{
			name:  "space in short name",
			short: " ",
			want:  &FlagError{Short: " ", Err: ErrInvalidName},
		},
		{
			name:  "comma in short name",
			short: ",",
			want:  &FlagError{Short: ",", Err: ErrInvalidName},
		},
		{
			name: "start dash in long name",
			long: "-help",
			want: &FlagError{Long: "-help", Err: ErrInvalidName},
		},
		{
			name: "ignore non-start dash in long name",
			long: "go-help",
			want: nil,
		},
		{
			name: "ignore end dash in long name",
			long: "help-",
			want: nil,
		},
		{
			name: "equal in long name",
			long: "help=test",
			want: &FlagError{Long: "help=test", Err: ErrInvalidName},
		},
		{
			name: "space in long name",
			long: "help test",
			want: &FlagError{Long: "help test", Err: ErrInvalidName},
		},
		{
			name: "comma in long name",
			long: "help,test",
			want: &FlagError{Long: "help,test", Err: ErrInvalidName},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var register DefaultRegister

			_ = Bool(&register, "", FlagOptions{
				Short: tc.short,
				Long:  tc.long,
			})

			got := register.Err()
			if !errors.Is(got, tc.want) {
				t.Fatalf("Parse(): got error = %q, want error = %q", got, tc.want)
			}
		})
	}
}

func TestRegisterDuplicatedFlag(t *testing.T) {
	var register DefaultRegister

	_ = Bool(&register, "a")
	_ = Int(&register, "a")

	got := register.Err()
	want := &FlagError{
		Short: "a",
		Err:   ErrDuplicate,
	}
	if !errors.Is(got, want) {
		t.Fatalf("Parse(): got error = %q, want error = %q", got, want)
	}
}

func TestRegisterInvalidNameArg(t *testing.T) {
	tt := []struct {
		name string
		arg  string
		want error
	}{
		{
			name: "empty arg name",
			want: &ArgError{Err: ErrMissingName},
		},
		{
			name: "start dash in arg name",
			arg:  "-help",
			want: &ArgError{Name: "-help", Err: ErrInvalidName},
		},
		{
			name: "ignore non-start dash in arg name",
			arg:  "go-help",
			want: nil,
		},
		{
			name: "ignore end dash in arg name",
			arg:  "help-",
			want: nil,
		},
		{
			name: "equal in arg name",
			arg:  "help=test",
			want: &ArgError{Name: "help=test", Err: ErrInvalidName},
		},
		{
			name: "space in arg name",
			arg:  "help test",
			want: &ArgError{Name: "help test", Err: ErrInvalidName},
		},
		{
			name: "comma in arg name",
			arg:  "help,test",
			want: &ArgError{Name: "help,test", Err: ErrInvalidName},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var register DefaultRegister

			_ = BoolArg(&register, tc.arg)

			got := register.Err()
			if !errors.Is(got, tc.want) {
				t.Fatalf("Parse(): got error = %q, want error = %q", got, tc.want)
			}
		})
	}
}

func TestRegisterDuplicatedArg(t *testing.T) {
	var register DefaultRegister

	_ = StringArg(&register, "a")
	_ = IntArg(&register, "a")

	got := register.Err()
	want := &ArgError{
		Name: "a",
		Err:  ErrDuplicate,
	}
	if !errors.Is(got, want) {
		t.Fatalf("Parse(): got error = %q, want error = %q", got, want)
	}
}

func TestParse_Parse_invalid_flags_syntax(t *testing.T) {
	tt := []struct {
		name string
		arg  string
		want error
	}{
		{
			name: "extra dash",
			arg:  "---test",
			want: &ParseFlagError{Name: "-test", Err: ErrSyntax},
		},
		{
			name: "equals after dash",
			arg:  "--=val",
			want: &ParseFlagError{Name: "=val", Err: ErrSyntax},
		},
		{
			name: "space after dash",
			arg:  "-- val",
			want: &ParseFlagError{Name: " val", Err: ErrSyntax},
		},
		{
			name: "comma after dash",
			arg:  "--,val",
			want: &ParseFlagError{Name: ",val", Err: ErrSyntax},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var (
				register DefaultRegister
				parser   DefaultParser
			)

			args := []string{tc.arg}

			err := parser.Parse(nil, &register, args)
			if !errors.Is(err, tc.want) {
				t.Fatalf("Parse(%v): got error = %q, want error = %q", args, err, tc.want)
			}
		})
	}
}

func TestParse_Parse_universal(t *testing.T) {
	var register DefaultRegister
	parser := DefaultParser{
		Universal: true,
	}

	show := Bool(&register, "show")
	help := Bool(&register, "help", WithShort("h"))
	c := Int(&register, "c")
	unused := Bool(&register, "unused", WithShort("u"))
	userID := Int(&register, "user-id")
	delay := Duration(&register, "delay")

	args := []string{"-show", "-h", "--c", "100", "-user-id", "200", "--delay=2s"}

	if err := parser.Parse(nil, &register, args); err != nil {
		t.Fatalf("Parse(): failed to parse args: %s", err)
	}

	// Check flags.
	const (
		wantShow   = true
		wantHelp   = true
		wantC      = 100
		wantUnused = false
		wantUserID = 200
		wantDelay  = 2 * time.Second
	)

	assertParseBoolFlags(t, "show", *show, wantShow)
	assertParseBoolFlags(t, "help", *help, wantHelp)

	if *c != wantC {
		t.Errorf("Parse(): c: got = %v, want = %v", *c, wantC)
	}

	assertParseBoolFlags(t, "unused", *unused, wantUnused)

	if *userID != wantUserID {
		t.Errorf("Parse(): userID: got = %v, want = %v", *userID, wantUserID)
	}

	if *delay != wantDelay {
		t.Errorf("Parse(): delay: got = %v, want = %v", *delay, wantDelay)
	}
}

func TestParse_Parse_posix_style_short_flags(t *testing.T) {
	var (
		register DefaultRegister
		parser   DefaultParser
	)

	a := Bool(&register, "a")
	b := Bool(&register, "b")
	c := Bool(&register, "c")
	d := Bool(&register, "d")
	e := Bool(&register, "e")
	f := Int(&register, "f")
	g := Bool(&register, "g")

	args := []string{"-ab", "-def", "100", "-g"}

	if err := parser.Parse(nil, &register, args); err != nil {
		t.Fatalf("Parse(): failed to parse args: %s", err)
	}

	// Check flags.
	const (
		wantA = true
		wantB = true
		wantC = false
		wantD = true
		wantE = true
		wantF = 100
		wantG = true
	)

	assertParseBoolFlags(t, "a", *a, wantA)
	assertParseBoolFlags(t, "b", *b, wantB)
	assertParseBoolFlags(t, "c", *c, wantC)
	assertParseBoolFlags(t, "d", *d, wantD)
	assertParseBoolFlags(t, "e", *e, wantE)

	if *f != wantF {
		t.Errorf("Parse(): f: got = %v, want = %v", *f, wantF)
	}

	assertParseBoolFlags(t, "g", *g, wantG)
}

func TestParse_Parse_posix_style_short_flags_unknown(t *testing.T) {
	var (
		register DefaultRegister
		parser   DefaultParser
	)

	_ = Bool(&register, "a")
	_ = Bool(&register, "b")
	_ = Bool(&register, "c")
	_ = Bool(&register, "d")

	args := []string{"-a", "-deb", "-g"}

	got := parser.Parse(nil, &register, args)
	want := &ParseFlagError{Name: "-e", Err: ErrUnknown}
	if !errors.Is(got, want) {
		t.Fatalf("Parse(): got error = %q, want error = %q", got, want)
	}
}

func TestParse_Parse_posix_style_short_flags_ignore_unknown(t *testing.T) {
	var register DefaultRegister
	parser := DefaultParser{
		IgnoreUnknownFlags: true,
	}

	a := Bool(&register, "a")
	b := Bool(&register, "b")
	c := Bool(&register, "c")
	d := Bool(&register, "d")

	args := []string{"-a", "-deb", "-g"}

	if err := parser.Parse(nil, &register, args); err != nil {
		t.Fatalf("Parse(): failed to parse args: %s", err)
	}

	// Check flags.
	const (
		wantA = true
		wantB = true
		wantC = false
		wantD = true
	)

	assertParseBoolFlags(t, "a", *a, wantA)
	assertParseBoolFlags(t, "b", *b, wantB)
	assertParseBoolFlags(t, "c", *c, wantC)
	assertParseBoolFlags(t, "d", *d, wantD)
}

func TestParse_Parse_disable_posix_style_short_flags(t *testing.T) {
	var register DefaultRegister
	parser := DefaultParser{
		DisablePosixStyle: true,
	}

	_ = Bool(&register, "a")
	_ = Bool(&register, "b")
	_ = Bool(&register, "c")
	_ = Bool(&register, "d")
	_ = Bool(&register, "e")
	_ = Bool(&register, "f")
	_ = Bool(&register, "g")

	args := []string{"-ab", "-def", "-g"}

	got := parser.Parse(nil, &register, args)
	want := &ParseFlagError{Name: "-ab", Err: ErrUnknown}
	if !errors.Is(got, want) {
		t.Fatalf("Parse(): got error = %q, want error = %q", got, want)
	}
}

func TestParse_Parse_universal_and_posix_style_short_flags(t *testing.T) {
	var register DefaultRegister
	parser := DefaultParser{
		Universal: true,
	}

	_ = Bool(&register, "a")
	_ = Bool(&register, "b")
	_ = Bool(&register, "c")
	_ = Bool(&register, "d")
	_ = Bool(&register, "e")
	_ = Bool(&register, "f")
	_ = Bool(&register, "g")

	args := []string{"-ab", "-def", "-g"}

	got := parser.Parse(nil, &register, args)
	want := &ParseFlagError{Name: "-ab", Err: ErrUnknown}
	if !errors.Is(got, want) {
		t.Fatalf("Parse(): got error = %q, want error = %q", got, want)
	}
}

func TestParse_Parse_short_flag_inline_value(t *testing.T) {
	var (
		register DefaultRegister
		parser   DefaultParser
	)

	a := String(&register, "a")
	b := Bool(&register, "b")
	c := Bool(&register, "c")
	d := Bool(&register, "d")
	e := Bool(&register, "e")

	args := []string{"-c", "-aboot", "-d"}

	if err := parser.Parse(nil, &register, args); err != nil {
		t.Fatalf("Parse(): failed to parse args: %s", err)
	}

	// Check flags.
	const (
		wantA = "boot"
		wantB = false
		wantC = true
		wantD = true
		wantE = false
	)

	if *a != wantA {
		t.Errorf("Parse(): a: got = %v, want = %v", *a, wantA)
	}

	assertParseBoolFlags(t, "b", *b, wantB)
	assertParseBoolFlags(t, "c", *c, wantC)
	assertParseBoolFlags(t, "d", *d, wantD)
	assertParseBoolFlags(t, "e", *e, wantE)
}

func TestParse_Parse_short_flag_inline_value_with_equals(t *testing.T) {
	var (
		register DefaultRegister
		parser   DefaultParser
	)

	a := String(&register, "a")
	b := Bool(&register, "b")
	c := Bool(&register, "c")
	d := Bool(&register, "d")
	e := Bool(&register, "e")

	args := []string{"-c", "-aboot=now", "-d"}

	if err := parser.Parse(nil, &register, args); err != nil {
		t.Fatalf("Parse(): failed to parse args: %s", err)
	}

	// Check flags.
	const (
		wantA = "boot=now"
		wantB = false
		wantC = true
		wantD = true
		wantE = false
	)

	if *a != wantA {
		t.Errorf("Parse(): a: got = %v, want = %v", *a, wantA)
	}

	assertParseBoolFlags(t, "b", *b, wantB)
	assertParseBoolFlags(t, "c", *c, wantC)
	assertParseBoolFlags(t, "d", *d, wantD)
	assertParseBoolFlags(t, "e", *e, wantE)
}

func TestParse_Parse_short_flag_inline_value_with_invalid_value(t *testing.T) {
	var register DefaultRegister
	parser := DefaultParser{
		DisableInlineValue: true,
	}

	_ = Int(&register, "a")
	_ = Bool(&register, "b")
	_ = Bool(&register, "c")
	_ = Bool(&register, "d")
	_ = Bool(&register, "e")

	args := []string{"-c", "-aboot", "-d"}

	got := parser.Parse(nil, &register, args)
	want := &FlagError{
		Short: "a",
		Err: &ParseValueError{
			Type: "int",
			Err:  ErrSyntax,
		},
	}
	if !errors.Is(got, want) {
		t.Fatalf("Parse(): got error = %q, want error = %q", got, want)
	}
}

func TestParse_Parse_disable_short_flag_inline_value(t *testing.T) {
	var register DefaultRegister
	parser := DefaultParser{
		DisableInlineValue: true,
	}

	_ = String(&register, "a")
	_ = Bool(&register, "b")
	_ = Bool(&register, "c")
	_ = Bool(&register, "d")
	_ = Bool(&register, "e")

	args := []string{"-c", "-aboot", "-d"}

	got := parser.Parse(nil, &register, args)
	want := &ParseFlagError{Name: "-o", Err: ErrUnknown}
	if !errors.Is(got, want) {
		t.Fatalf("Parse(): got error = %q, want error = %q", got, want)
	}
}

func TestParse_Parse_universal_and_short_flag_inline_value(t *testing.T) {
	var register DefaultRegister
	parser := DefaultParser{
		Universal: true,
	}

	_ = String(&register, "a")
	_ = Bool(&register, "b")
	_ = Bool(&register, "c")
	_ = Bool(&register, "d")
	_ = Bool(&register, "e")

	args := []string{"-c", "-aboot", "-d"}

	got := parser.Parse(nil, &register, args)
	want := &ParseFlagError{Name: "-aboot", Err: ErrUnknown}
	if !errors.Is(got, want) {
		t.Fatalf("Parse(): got error = %q, want error = %q", got, want)
	}
}

func TestParse_Parse_posix_style_and_short_flag_inline_value(t *testing.T) {
	var (
		register DefaultRegister
		parser   DefaultParser
	)

	a := String(&register, "a")
	b := Bool(&register, "b")
	c := Bool(&register, "c")
	d := Bool(&register, "d")
	e := String(&register, "e")
	f := Bool(&register, "f")
	g := Bool(&register, "g")

	args := []string{"-c", "-aboot", "-degif"}

	if err := parser.Parse(nil, &register, args); err != nil {
		t.Fatalf("Parse(): failed to parse args: %s", err)
	}

	// Check flags.
	const (
		wantA = "boot"
		wantB = false
		wantC = true
		wantD = true
		wantE = "gif"
		wantF = false
		wantG = false
	)

	if *a != wantA {
		t.Errorf("Parse(): a: got = %v, want = %v", *a, wantA)
	}

	assertParseBoolFlags(t, "b", *b, wantB)
	assertParseBoolFlags(t, "c", *c, wantC)
	assertParseBoolFlags(t, "d", *d, wantD)

	if *e != wantE {
		t.Errorf("Parse(): e: got = %v, want = %v", *e, wantE)
	}

	assertParseBoolFlags(t, "f", *f, wantF)
	assertParseBoolFlags(t, "g", *g, wantG)
}

func TestParse_Parse_disable_posix_style_and_short_flag_inline_value(t *testing.T) {
	var register DefaultRegister
	parser := DefaultParser{
		DisablePosixStyle: true,
	}

	a := String(&register, "a")
	b := Bool(&register, "b")
	c := Bool(&register, "c")
	d := Bool(&register, "d")
	e := Bool(&register, "e")
	f := Bool(&register, "f")

	args := []string{"-c", "-aboot", "-def"}

	got := parser.Parse(nil, &register, args)
	want := &ParseFlagError{Name: "-def", Err: ErrUnknown}
	if !errors.Is(got, want) {
		t.Fatalf("Parse(): got error = %q, want error = %q", got, want)
	}

	// Check flags.
	const (
		wantA = "boot"
		wantB = false
		wantC = true
		wantD = false
		wantE = false
		wantF = false
	)

	if *a != wantA {
		t.Errorf("Parse(): a: got = %v, want = %v", *a, wantA)
	}

	assertParseBoolFlags(t, "b", *b, wantB)
	assertParseBoolFlags(t, "c", *c, wantC)
	assertParseBoolFlags(t, "d", *d, wantD)
	assertParseBoolFlags(t, "e", *e, wantE)
	assertParseBoolFlags(t, "f", *f, wantF)
}

func TestParser_Parse(t *testing.T) {
	var (
		register DefaultRegister
		parser   DefaultParser
	)

	show := Bool(&register, "show",
		WithShort("s"),
		Usage("Show the resuld of the function"),
	)

	recreate := Bool(&register, "recreate",
		Usage("Re-create the current user"),
	)

	update := Bool(&register, "update",
		Usage("Update the DB"),
	)

	unused := Bool(&register, "unused")

	count := Int(&register, "count",
		WithShort("c"),
	)

	userID := IntArg(&register, "user-id",
		Usage("Current User ID"),
	)

	rest := RestStrings(&register, "rest")

	args := []string{
		"--show", "--recreate=false", "-c", "100500", "1337",
		"other", "vals", "--update", "true", "in", "args",
	}

	if err := parser.Parse(nil, &register, args); err != nil {
		t.Fatalf("Parse(): failed to parse args: %s", err)
	}

	// Check flags.
	const (
		wantShow     = true
		wantRecreate = false
		wantUpdate   = true
		wantUnused   = false
		wantCount    = 100500
	)

	if *show != wantShow {
		t.Errorf("Parse(): show: got = %v, want = %v", *show, wantShow)
	}

	if *recreate != wantRecreate {
		t.Errorf("Parse(): recreate: got = %v, want = %v", *recreate, wantRecreate)
	}

	if *update != wantUpdate {
		t.Errorf("Parse(): update: got = %v, want = %v", *update, wantUpdate)
	}

	if *unused != wantUnused {
		t.Errorf("Parse(): unused: got = %v, want = %v", *unused, wantUnused)
	}

	if *count != wantCount {
		t.Errorf("Parse(): count: got = %v, want = %v", *count, wantCount)
	}

	// Check args.
	const (
		wantUserID = 1337
	)

	if *userID != wantUserID {
		t.Errorf("Parse(): userID: got = %v, want = %v", *userID, wantUserID)
	}

	// Check unknown.
	wantRest := []string{"other", "vals", "in", "args"}
	if !reflect.DeepEqual(*rest, wantRest) {
		t.Errorf("Parse(): rest: got = %#v, want = %#v", *rest, wantRest)
	}
}

func TestParser_Parse_unknown_flags(t *testing.T) {
	var (
		register DefaultRegister
		parser   DefaultParser
	)

	a := Bool(&register, "a")
	b := Bool(&register, "b")

	args := []string{"-a", "-c", "false", "-d", "-b"}

	got := parser.Parse(nil, &register, args)
	want := &ParseFlagError{Name: "-c", Err: ErrUnknown}
	if !errors.Is(got, want) {
		t.Fatalf("Parse(): got error = %q, want error = %q", got, want)
	}

	const (
		wantA = true
		wantB = false
	)

	assertParseBoolFlags(t, "a", *a, wantA)
	assertParseBoolFlags(t, "b", *b, wantB)
}

func TestParser_Parse_unknown_flags_with_value(t *testing.T) {
	var (
		register DefaultRegister
		parser   DefaultParser
	)

	_ = Bool(&register, "a")
	_ = Bool(&register, "b")

	args := []string{"-a", "-c=100", "-d", "-b"}

	got := parser.Parse(nil, &register, args)
	want := &ParseFlagError{Name: "-c", Err: ErrUnknown}
	if !errors.Is(got, want) {
		t.Fatalf("Parse(): got error = %q, want error = %q", got, want)
	}
}

func TestParser_Parse_ignore_unknown_flags(t *testing.T) {
	var register DefaultRegister
	parser := DefaultParser{
		IgnoreUnknownFlags: true,
	}

	a := Bool(&register, "a")
	b := Bool(&register, "b")
	v := BoolArg(&register, "v")

	args := []string{"-a", "-c=200", "false", "-d", "-e", "-b"}

	if err := parser.Parse(nil, &register, args); err != nil {
		t.Fatalf("Parse(): failed to parse args: %s", err)
	}

	const (
		wantA = true
		wantB = true
		wantV = false
	)

	assertParseBoolFlags(t, "a", *a, wantA)
	assertParseBoolFlags(t, "b", *b, wantB)

	if *v != wantV {
		t.Errorf("Parse(): v: got = %v, want = %v", *v, wantV)
	}
}

func TestParser_Parse_unknown_rest(t *testing.T) {
	var (
		register DefaultRegister
		parser   DefaultParser
	)

	_ = BoolArg(&register, "a")
	_ = BoolArg(&register, "b", Optional)

	args := []string{"true", "false", "c", "d", "1337", "false", "e"}

	got := parser.Parse(nil, &register, args)
	want := &ParseArgError{Arg: "c", Index: 2, Err: ErrUnknown}
	if !errors.Is(got, want) {
		t.Fatalf("Parse(): got error = %q, want error = %q", got, want)
	}
}

func TestParser_Parse_ignore_unknown_rest(t *testing.T) {
	var register DefaultRegister
	parser := DefaultParser{
		IgnoreUnknownArgs: true,
	}

	a := BoolArg(&register, "a")
	b := BoolArg(&register, "b")

	args := []string{"true", "false", "c", "d", "1337", "false", "e"}

	if err := parser.Parse(nil, &register, args); err != nil {
		t.Fatalf("Parse(): failed to parse args: %s", err)
	}

	const (
		wantA = true
		wantB = false
	)

	if *a != wantA {
		t.Errorf("Parse(): a: got = %v, want = %v", *a, wantA)
	}

	if *b != wantB {
		t.Errorf("Parse(): b: got = %v, want = %v", *b, wantB)
	}
}

func TestParser_Parse_rest_with_ignore_unknown_rest(t *testing.T) {
	var register DefaultRegister
	parser := DefaultParser{
		IgnoreUnknownArgs: true,
	}

	a := BoolArg(&register, "a")
	b := BoolArg(&register, "b")

	var rest []string
	_ = RestStringsVar(&register, &rest, "rest")

	args := []string{"true", "false", "c", "d", "1337", "false", "e"}

	if err := parser.Parse(nil, &register, args); err != nil {
		t.Fatalf("Parse(): failed to parse args: %s", err)
	}

	const (
		wantA = true
		wantB = false
	)

	wantRest := []string{"c", "d", "1337", "false", "e"}

	if *a != wantA {
		t.Errorf("Parse(): a: got = %v, want = %v", *a, wantA)
	}

	if *b != wantB {
		t.Errorf("Parse(): b: got = %v, want = %v", *b, wantB)
	}

	if !reflect.DeepEqual(rest, wantRest) {
		t.Errorf("Parse(): rest: got = %#v, want = %#v", rest, wantRest)
	}
}

var _ Commander = (*testCommander)(nil)

type testCommander struct {
	commands []string
	use      func() (Register, error)

	path []string
	i    int
}

func (c *testCommander) IsCommand(name string) bool {
	if c.i >= len(c.commands) {
		return false
	}

	cmd := c.commands[c.i]

	return cmd == name
}

func (c *testCommander) SetCommand(name string) (Register, error) {
	if c.i >= len(c.commands) {
		return nil, fmt.Errorf("command not found: %s", name)
	}

	cmd := c.commands[c.i]
	if cmd != name {
		return nil, fmt.Errorf("command not found: %s", name)
	}

	c.i++
	c.path = append(c.path, cmd)

	register, err := c.use()
	if err != nil {
		return nil, err
	}

	return register, nil
}

func (c *testCommander) Path() []string { return c.path }

func TestParser_Parse_with_commands(t *testing.T) {
	var (
		register DefaultRegister
		parser   DefaultParser
	)

	show := new(bool)
	recreate := new(bool)
	update := new(bool)
	unused := new(bool)
	count := new(int)
	userID := new(int)
	rest := new([]string)

	commander := testCommander{
		commands: []string{"first", "second", "third", "fourth"},
		use: func() (Register, error) {
			var register DefaultRegister

			show = Bool(&register, "show",
				WithShort("s"),
				Usage("Show the resuld of the function"),
			)

			recreate = Bool(&register, "recreate",
				Usage("Re-create the current user"),
				Required,
			)

			update = Bool(&register, "update",
				Usage("Update the DB"),
			)

			unused = Bool(&register, "unused")

			count = Int(&register, "count",
				WithShort("c"),
			)

			userID = IntArg(&register, "user-id",
				Usage("Current User ID"),
			)

			rest = RestStrings(&register, "rest")

			return &register, nil
		},
	}

	args := []string{
		"first", "second", "third",
		"1337", "--show", "--recreate=false", "-c", "100500",
		"other", "vals", "--update", "true", "in", "args",
	}

	if err := parser.Parse(&commander, &register, args); err != nil {
		t.Fatalf("Parse(): failed to parse args: %s", err)
	}

	// Chack path.
	wantPath := []string{"first", "second", "third"}
	if !reflect.DeepEqual(commander.path, wantPath) {
		t.Fatalf("Parse(): path: got = %v, want = %v", commander.path, wantPath)
	}

	// Check flags.
	const (
		wantShow     = true
		wantRecreate = false
		wantUpdate   = true
		wantUnused   = false
		wantCount    = 100500
	)

	if *show != wantShow {
		t.Errorf("Parse(): show: got = %v, want = %v", *show, wantShow)
	}

	if *recreate != wantRecreate {
		t.Errorf("Parse(): recreate: got = %v, want = %v", *recreate, wantRecreate)
	}

	if *update != wantUpdate {
		t.Errorf("Parse(): update: got = %v, want = %v", *update, wantUpdate)
	}

	if *unused != wantUnused {
		t.Errorf("Parse(): unused: got = %v, want = %v", *unused, wantUnused)
	}

	if *count != wantCount {
		t.Errorf("Parse(): count: got = %v, want = %v", *count, wantCount)
	}

	// Check args.
	const (
		wantUserID = 1337
	)

	if *userID != wantUserID {
		t.Errorf("Parse(): userID: got = %v, want = %v", *userID, wantUserID)
	}

	// Check unknown.
	wantRest := []string{"other", "vals", "in", "args"}
	if !reflect.DeepEqual(*rest, wantRest) {
		t.Errorf("Parse(): rest: got = %#v, want = %#v", *rest, wantRest)
	}
}

func TestParser_Parse_required_flag(t *testing.T) {
	var (
		register DefaultRegister
		parser   DefaultParser
	)

	_ = Bool(&register, "a")
	_ = Bool(&register, "b", Required)
	_ = Bool(&register, "c")

	args := []string{"-a"}

	got := parser.Parse(nil, &register, args)
	want := &FlagError{Short: "b", Err: ErrNotProvided}
	if !errors.Is(got, want) {
		t.Fatalf("Parse(): got error = %q, want error = %q", got, want)
	}
}

func TestParser_Parse_required_multi_flag(t *testing.T) {
	var (
		register DefaultRegister
		parser   DefaultParser
	)

	_ = Bool(&register, "a")
	_ = Bools(&register, "b")
	_ = Bools(&register, "c", Required)
	_ = Bools(&register, "d", Required)
	_ = Bools(&register, "e")

	args := []string{"-a", "-c"}

	got := parser.Parse(nil, &register, args)
	want := &FlagError{Short: "d", Err: ErrNotProvided}
	if !errors.Is(got, want) {
		t.Fatalf("Parse(): got error = %q, want error = %q", got, want)
	}
}

func TestParser_Parse_required_arg(t *testing.T) {
	var (
		register DefaultRegister
		parser   DefaultParser
	)

	// Args are implicitly required.
	_ = BoolArg(&register, "a")
	_ = BoolArg(&register, "b")
	_ = BoolArg(&register, "c")
	_ = BoolArg(&register, "d")

	args := []string{"true", "false"}

	got := parser.Parse(nil, &register, args)
	want := &ArgError{Name: "c", Err: ErrNotProvided}
	if !errors.Is(got, want) {
		t.Fatalf("Parse(): got error = %q, want error = %q", got, want)
	}
}

func TestParser_Parse_optional_arg(t *testing.T) {
	var (
		register DefaultRegister
		parser   DefaultParser
	)

	// Args are implicitly required.
	_ = BoolArg(&register, "a")
	_ = BoolArg(&register, "b")
	_ = BoolArg(&register, "c", Optional)
	_ = BoolArg(&register, "d", Optional)

	args := []string{"true", "false"}

	if err := parser.Parse(nil, &register, args); err != nil {
		t.Fatalf("Parse(): failed to parse args: %s", err)
	}
}

func TestParser_Parse_optional_arg_after_required(t *testing.T) {
	var (
		register DefaultRegister
		parser   DefaultParser
	)

	_ = BoolArg(&register, "a")
	_ = BoolArg(&register, "b", Optional)
	_ = BoolArg(&register, "c")
	_ = BoolArg(&register, "d", Optional)

	args := []string{"true", "false", "true", "true"}

	got := parser.Parse(nil, &register, args)
	want := &ArgError{Name: "c", Err: ErrRequiredAfterOptional}
	if !errors.Is(got, want) {
		t.Fatalf("Parse(): got error = %q, want error = %q", got, want)
	}
}

func TestParser_Parse_rest(t *testing.T) {
	var (
		register DefaultRegister
		parser   DefaultParser
	)

	_ = BoolArg(&register, "a")
	_ = BoolArg(&register, "b")

	var rest []string
	_ = RestStringsVar(&register, &rest, "rest")

	args := []string{"true", "false", "c", "d", "1337", "false", "e,d,f", "g"}

	if err := parser.Parse(nil, &register, args); err != nil {
		t.Fatalf("Parse(): failed to parse args: %s", err)
	}

	want := []string{"c", "d", "1337", "false", "e", "d", "f", "g"}
	if !reflect.DeepEqual(rest, want) {
		t.Errorf("Parse(): rest: got = %#v, want = %#v", rest, want)
	}
}

func TestRegisterInvalidNameRestArgs(t *testing.T) {
	tt := []struct {
		name     string
		restArgs string
		want     error
	}{
		{
			name: "empty rest args name",
			want: &RestArgsError{Err: ErrMissingName},
		},
		{
			name:     "start dash in rest args name",
			restArgs: "-help",
			want:     &RestArgsError{Name: "-help", Err: ErrInvalidName},
		},
		{
			name:     "ignore non-start dash in rest args name",
			restArgs: "go-help",
			want:     nil,
		},
		{
			name:     "ignore end dash in rest args name",
			restArgs: "help-",
			want:     nil,
		},
		{
			name:     "equal in rest args name",
			restArgs: "help=test",
			want:     &RestArgsError{Name: "help=test", Err: ErrInvalidName},
		},
		{
			name:     "space in rest args name",
			restArgs: "help test",
			want:     &RestArgsError{Name: "help test", Err: ErrInvalidName},
		},
		{
			name:     "comma in rest args name",
			restArgs: "help,test",
			want:     &RestArgsError{Name: "help,test", Err: ErrInvalidName},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var register DefaultRegister

			_ = RestStrings(&register, tc.restArgs)

			got := register.Err()
			if !errors.Is(got, tc.want) {
				t.Fatalf("Parse(): got error = %q, want error = %q", got, tc.want)
			}
		})
	}
}

func TestParser_Parse_rest_after_optional_arg(t *testing.T) {
	var (
		register DefaultRegister
		parser   DefaultParser
	)

	_ = BoolArg(&register, "a")
	_ = BoolArg(&register, "b", Optional)

	var rest []string
	_ = RestStringsVar(&register, &rest, "rest")

	args := []string{"true", "false", "c", "d", "1337", "false", "e"}

	if err := parser.Parse(nil, &register, args); err != nil {
		t.Fatalf("Parse(): failed to parse args: %s", err)
	}

	want := []string{"c", "d", "1337", "false", "e"}
	if !reflect.DeepEqual(rest, want) {
		t.Errorf("Parse(): rest: got = %#v, want = %#v", rest, want)
	}
}

func TestParser_Parse_multi_rest(t *testing.T) {
	var register DefaultRegister

	_ = RestStrings(&register, "rest")
	_ = RestStrings(&register, "other")

	got := register.Err()
	want := &RestArgsError{Name: "other", Err: ErrDuplicate}
	if !errors.Is(got, want) {
		t.Fatalf("Parse(): got error = %q, want error = %q", got, want)
	}
}

func TestParser_Parse_arg_after_rest(t *testing.T) {
	var register DefaultRegister

	_ = BoolArg(&register, "a")
	_ = RestStrings(&register, "rest")
	_ = BoolArg(&register, "b")

	got := register.Err()
	want := &ArgError{Name: "b", Err: ErrArgAfterRest}
	if !errors.Is(got, want) {
		t.Fatalf("Parse(): got error = %q, want error = %q", got, want)
	}
}

func TestParser_Parse_optional_arg_after_rest(t *testing.T) {
	var register DefaultRegister

	_ = BoolArg(&register, "a")
	_ = RestStrings(&register, "rest")
	_ = BoolArg(&register, "b", Optional)

	got := register.Err()
	want := &ArgError{Name: "b", Err: ErrArgAfterRest}
	if !errors.Is(got, want) {
		t.Fatalf("Parse(): got error = %q, want error = %q", got, want)
	}
}

func TestParser_Parse_flags_terminator(t *testing.T) {
	var (
		register DefaultRegister
		parser   DefaultParser
	)

	a := Bool(&register, "a")
	b := Bool(&register, "b")
	c := StringArg(&register, "c")
	d := StringArg(&register, "d")

	var rest []string
	_ = RestStringsVar(&register, &rest, "rest")

	args := []string{
		"-a", "true", "testC", "--", "-b", "-a=false", "-c", "--", "d",
	}

	if err := parser.Parse(nil, &register, args); err != nil {
		t.Fatalf("Parse(): failed to parse args: %s", err)
	}

	const (
		wantA = true
		wantB = false
		wantC = "testC"
		wantD = "-b"
	)

	assertParseBoolFlags(t, "a", *a, wantA)
	assertParseBoolFlags(t, "b", *b, wantB)

	if *c != wantC {
		t.Errorf("Parse(): c: got = %q, want = %q", *c, wantC)
	}

	if *d != wantD {
		t.Errorf("Parse(): d: got = %q, want = %q", *d, wantD)
	}

	want := []string{"-a=false", "-c", "--", "d"}
	if !reflect.DeepEqual(rest, want) {
		t.Errorf("Parse(): rest: got = %#v, want = %#v", rest, want)
	}
}

func assertParseBoolFlags(t *testing.T, name string, got, want bool) {
	t.Helper()

	if got != want {
		t.Errorf("Parse(): %s: got = %v, want = %v", name, got, want)
	}
}
