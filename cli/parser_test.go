package cli

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func deref(t *testing.T, ptr interface{}) interface{} {
	t.Helper()

	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr {
		t.Fatalf("flag is not a pointer: %T", v.String())
	}

	return v.Elem().Interface()
}

func TestParseFlags(t *testing.T) {
	type testValue struct {
		name string
		args []string
		want interface{}
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
					name: "true value",
					args: []string{"-t=true"},
					want: true,
				},
				{
					name: "false value",
					args: []string{"-t=false"},
					want: false,
				},
				{
					name: "false next arg",
					args: []string{"-t", "false"},
					want: false,
				},
				{
					name: "true next arg",
					args: []string{"-t", "true"},
					want: true,
				},
				{
					name: "skip not bool-like next arg",
					args: []string{"-t", "abcd"},
					want: true,
				},
			},
		},
		{
			name:  "Int",
			setup: func(r Register) interface{} { return Int(r, "t") },
			tests: []testValue{
				{
					name: "0 value",
					args: []string{"-t=0"},
					want: 0,
				},
				{
					name: "-0 value",
					args: []string{"-t=-0"},
					want: 0,
				},
				{
					name: "1337 value",
					args: []string{"-t=1337"},
					want: 1337,
				},
				{
					name: "-7331 value",
					args: []string{"-t=-7331"},
					want: -7331,
				},
				{
					name: "0 next arg",
					args: []string{"-t", "0"},
					want: 0,
				},
				{
					name: "-0 next arg",
					args: []string{"-t", "-0"},
					want: 0,
				},
				{
					name: "1337 next arg",
					args: []string{"-t", "1337"},
					want: 1337,
				},
				{
					name: "-7331 next arg",
					args: []string{"-t", "-7331"},
					want: -7331,
				},
			},
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
					name: "test next arg",
					args: []string{"-t", "test"},
					want: "test",
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
					name: "next flag",
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
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			for _, tvc := range tc.tests {
				t.Run(tvc.name, func(t *testing.T) {
					var parser DefaultParser

					f := tc.setup(&parser)

					if err := parser.Parse(nil, tvc.args); err != nil {
						t.Fatalf("Parse(%v): failed to parse flags: %s", tvc.args, err)
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

func TestParseFlags_broken_value(t *testing.T) {
	type testValue struct {
		name string
		args []string
		want error
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
					want: &ParseError{
						Type: "bool",
						Err:  ErrSyntax,
					},
				},
			},
		},
		{
			name:  "Int",
			setup: func(r Register) interface{} { return Int(r, "t") },
			tests: []testValue{
				{
					name: "not int-like next arg",
					args: []string{"-t", "abcd"},
					want: &ParseError{
						Type: "int",
						Err:  ErrSyntax,
					},
				},
				{
					name: "int overflow",
					args: []string{"-t", "99999999999999999999999999"},
					want: &ParseError{
						Type: "int",
						Err:  ErrRange,
					},
				},
				{
					name: "int negative overflow",
					args: []string{"-t", "-99999999999999999999999999"},
					want: &ParseError{
						Type: "int",
						Err:  ErrRange,
					},
				},
				{
					name: "bool next arg",
					args: []string{"-t", "true"},
					want: &ParseError{
						Type: "int",
						Err:  ErrSyntax,
					},
				},
				{
					name: "float next arg",
					args: []string{"-t", "12.34"},
					want: &ParseError{
						Type: "int",
						Err:  ErrSyntax,
					},
				},
				{
					name: "without value",
					args: []string{"-t"},
					want: &ParseError{
						Type: "int",
						Err:  ErrSyntax,
					},
				},
				{
					name: "empty value",
					args: []string{"-t", ""},
					want: &ParseError{
						Type: "int",
						Err:  ErrSyntax,
					},
				},
				{
					name: "next flag",
					args: []string{"-t", "-b"},
					want: &ParseError{
						Type: "int",
						Err:  ErrSyntax,
					},
				},
			},
		},
		{
			name:  "String",
			setup: func(r Register) interface{} { return String(r, "t") },
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			for _, tvc := range tc.tests {
				t.Run(tvc.name, func(t *testing.T) {
					var parser DefaultParser

					_ = tc.setup(&parser)

					err := parser.Parse(nil, tvc.args)
					if !errors.Is(err, tvc.want) {
						t.Fatalf("Parse(%v): got error = %q, want error = %q", tvc.args, err, tvc.want)
					}
				})
			}
		})
	}
}

func TestParseArgs(t *testing.T) {
	type testValue struct {
		name string
		args []string
		want interface{}
	}

	tt := []struct {
		name  string
		setup func(Register) interface{}
		tests []testValue
	}{
		{
			name:  "IntArg",
			setup: func(r Register) interface{} { return IntArg(r, "t") },
			tests: []testValue{
				{
					name: "0",
					args: []string{"0"},
					want: 0,
				},
				{
					name: "-0",
					args: []string{"-0"},
					want: 0,
				},
				{
					name: "1337",
					args: []string{"1337"},
					want: 1337,
				},
				{
					name: "-7331",
					args: []string{"-7331"},
					want: -7331,
				},
			},
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
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			for _, tvc := range tc.tests {
				t.Run(tvc.name, func(t *testing.T) {
					var parser DefaultParser

					f := tc.setup(&parser)

					if err := parser.Parse(nil, tvc.args); err != nil {
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

	tt := []struct {
		name  string
		setup func(Register) interface{}
		tests []testValue
	}{
		{
			name:  "IntArg",
			setup: func(r Register) interface{} { return IntArg(r, "t") },
			tests: []testValue{
				{
					name: "not int-like",
					args: []string{"abcd"},
					want: &ParseError{
						Type: "int",
						Err:  ErrSyntax,
					},
				},
				{
					name: "int overflow",
					args: []string{"99999999999999999999999999"},
					want: &ParseError{
						Type: "int",
						Err:  ErrRange,
					},
				},
				{
					name: "int negative overflow",
					args: []string{"-99999999999999999999999999"},
					want: &ParseError{
						Type: "int",
						Err:  ErrRange,
					},
				},
				{
					name: "bool",
					args: []string{"true"},
					want: &ParseError{
						Type: "int",
						Err:  ErrSyntax,
					},
				},
				{
					name: "float",
					args: []string{"12.34"},
					want: &ParseError{
						Type: "int",
						Err:  ErrSyntax,
					},
				},
				{
					name: "empty",
					args: []string{""},
					want: &ParseError{
						Type: "int",
						Err:  ErrSyntax,
					},
				},
			},
		},
		{
			name:  "StringArg",
			setup: func(r Register) interface{} { return StringArg(r, "t") },
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			for _, tvc := range tc.tests {
				t.Run(tvc.name, func(t *testing.T) {
					var parser DefaultParser

					_ = tc.setup(&parser)

					err := parser.Parse(nil, tvc.args)
					if !errors.Is(err, tvc.want) {
						t.Fatalf("Parse(%v): got error = %q, want error = %q", tvc.args, err, tvc.want)
					}
				})
			}
		})
	}
}

func TestRegisterDuplicatedFlag(t *testing.T) {
	var parser DefaultParser

	_ = Bool(&parser, "a")
	flag := &parser.Flags()[0]

	_ = Int(&parser, "a")

	args := []string{"-a", "100"}

	got := parser.Parse(nil, args)
	want := &DuplicatedFlagError{
		Flag: flag,
	}
	if !errors.Is(got, want) {
		t.Fatalf("Parse(): got error = %q, want error = %q", got, want)
	}
}

func TestRegisterOverrideFlag(t *testing.T) {
	parser := DefaultParser{OverrideFlags: true}

	oldA := Bool(&parser, "a")
	a := Int(&parser, "a")

	args := []string{"-a", "100"}
	if err := parser.Parse(nil, args); err != nil {
		t.Fatalf("Parse(): failed to parse args: %s", err)
	}

	const wantOldA = false
	if *oldA != wantOldA {
		t.Errorf("Parse(): oldA: got = %v, want = %v", *oldA, wantOldA)
	}

	const wantA = 100
	if *a != wantA {
		t.Errorf("Parse(): a: got = %v, want = %v", *a, wantA)
	}
}

func TestRegisterDuplicatedArg(t *testing.T) {
	var parser DefaultParser

	_ = StringArg(&parser, "a")
	arg := &parser.Args()[0]

	_ = IntArg(&parser, "a")

	args := []string{"100"}

	got := parser.Parse(nil, args)
	want := &DuplicatedArgError{
		Arg: arg,
	}
	if !errors.Is(got, want) {
		t.Fatalf("Parse(): got error = %q, want error = %q", got, want)
	}
}

func TestRegisterOverrideArg(t *testing.T) {
	parser := DefaultParser{OverrideArgs: true}

	oldA := StringArg(&parser, "a")
	a := IntArg(&parser, "a")

	args := []string{"100"}
	if err := parser.Parse(nil, args); err != nil {
		t.Fatalf("Parse(): failed to parse args: %s", err)
	}

	const wantOldA = ""
	if *oldA != wantOldA {
		t.Errorf("Parse(): oldA: got = %v, want = %v", *oldA, wantOldA)
	}

	const wantA = 100
	if *a != wantA {
		t.Errorf("Parse(): a: got = %v, want = %v", *a, wantA)
	}
}

func TestParse_Parse_posix_style_short_flags(t *testing.T) {
	var parser DefaultParser

	a := Bool(&parser, "a")
	b := Bool(&parser, "b")
	c := Bool(&parser, "c")
	d := Bool(&parser, "d")
	e := Bool(&parser, "e")
	f := Bool(&parser, "f")
	g := Bool(&parser, "g")

	args := []string{"-ab", "-def", "-g"}

	if err := parser.Parse(nil, args); err != nil {
		t.Fatalf("Parse(): failed to parse args: %s", err)
	}

	// Check flags.
	const (
		wantA = true
		wantB = true
		wantC = false
		wantD = true
		wantE = true
		wantF = true
		wantG = true
	)

	assertParseBoolFlags(t, "a", *a, wantA)
	assertParseBoolFlags(t, "b", *b, wantB)
	assertParseBoolFlags(t, "c", *c, wantC)
	assertParseBoolFlags(t, "d", *d, wantD)
	assertParseBoolFlags(t, "e", *e, wantE)
	assertParseBoolFlags(t, "f", *f, wantF)
	assertParseBoolFlags(t, "g", *g, wantG)
}

func TestParser_Parse(t *testing.T) {
	var parser DefaultParser

	show := Bool(&parser, "show",
		WithShort("s"),
		Usage("Show the resuld of the function"),
	)

	recreate := Bool(&parser, "recreate",
		Usage("Re-create the current user"),
	)

	update := Bool(&parser, "update",
		Usage("Update the DB"),
	)

	unused := Bool(&parser, "unused")

	count := Int(&parser, "count",
		WithShort("c"),
	)

	userID := IntArg(&parser, "user-id",
		Usage("Current User ID"),
	)

	args := []string{
		"--show", "--recreate=false", "-c", "100500", "1337", "--update", "true",
		"--first-unknown", "other", "vals", "--second-unknown", "in", "args",
	}

	if err := parser.Parse(nil, args); err != nil {
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
	wantUnknown := []string{"--first-unknown", "--second-unknown"}
	if !reflect.DeepEqual(parser.unknown, wantUnknown) {
		t.Errorf("Parse(): unknown: got = %#v, want = %#v", parser.unknown, wantUnknown)
	}

	// Check unknown.
	wantRest := []string{"other", "vals", "in", "args"}
	if !reflect.DeepEqual(parser.rest, wantRest) {
		t.Errorf("Parse(): rest: got = %#v, want = %#v", parser.rest, wantRest)
	}
}

type testCommander struct {
	commands []string
	next     func() error

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

func (c *testCommander) SetCommand(name string) error {
	if c.i >= len(c.commands) {
		return fmt.Errorf("command not found: %s", name)
	}

	cmd := c.commands[c.i]
	if cmd != name {
		return fmt.Errorf("command not found: %s", name)
	}

	c.i++
	c.path = append(c.path, cmd)

	if err := c.next(); err != nil {
		return err
	}

	return nil
}

func (c *testCommander) Path() []string { return c.path }

func TestParser_Parse_with_commands(t *testing.T) {
	var parser DefaultParser

	show := new(bool)
	recreate := new(bool)
	update := new(bool)
	unused := new(bool)
	count := new(int)
	userID := new(int)

	commander := testCommander{
		commands: []string{"first", "second", "third"},
		next: func() error {
			show = Bool(&parser, "show",
				WithShort("s"),
				Usage("Show the resuld of the function"),
			)

			recreate = Bool(&parser, "recreate",
				Usage("Re-create the current user"),
			)

			update = Bool(&parser, "update",
				Usage("Update the DB"),
			)

			unused = Bool(&parser, "unused")

			count = Int(&parser, "count",
				WithShort("c"),
			)

			userID = IntArg(&parser, "user-id",
				Usage("Current User ID"),
			)

			return nil
		},
	}

	args := []string{
		"first", "second",
		"1337", "--show", "--recreate=false", "-c", "100500", "--update", "true",
		"--first-unknown", "other", "vals", "--second-unknown", "in", "args",
	}

	if err := parser.Parse(&commander, args); err != nil {
		t.Fatalf("Parse(): failed to parse args: %s", err)
	}

	// Chack path.
	wantPath := []string{"first", "second"}
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
	wantUnknown := []string{"--first-unknown", "--second-unknown"}
	if !reflect.DeepEqual(parser.unknown, wantUnknown) {
		t.Errorf("Parse(): unknown: got = %#v, want = %#v", parser.unknown, wantUnknown)
	}

	// Check unknown.
	wantRest := []string{"other", "vals", "in", "args"}
	if !reflect.DeepEqual(parser.rest, wantRest) {
		t.Errorf("Parse(): rest: got = %#v, want = %#v", parser.rest, wantRest)
	}
}

func assertParseBoolFlags(t *testing.T, name string, got, want bool) {
	t.Helper()

	if got != want {
		t.Errorf("Parse(): %s: got = %v, want = %v", name, got, want)
	}
}
