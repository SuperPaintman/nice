package cli

import (
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
					name: "false",
					args: []string{"-t", "false"},
					want: false,
				},
				{
					name: "true",
					args: []string{"-t", "true"},
					want: true,
				},
				{
					name: "skip unknown value",
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
					name: "0",
					args: []string{"-t", "0"},
					want: 0,
				},
				{
					name: "1337",
					args: []string{"-t", "1337"},
					want: 1337,
				},
				// TODO(SuperPaintman): add negative numbers.
				// {
				// 	name: "-7331",
				// 	args: []string{"-t", "-7331"},
				// 	want: -7331,
				// },
			},
		},
		{
			name:  "String",
			setup: func(r Register) interface{} { return String(r, "t") },
			tests: []testValue{
				{
					name: "test",
					args: []string{"-t", "test"},
					want: "test",
				},
				{
					name: "empty",
					args: []string{"-t", ""},
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
						t.Fatalf("Parse(): failed to parse flags: %s", err)
					}

					got := deref(t, f)
					if !reflect.DeepEqual(got, tvc.want) {
						t.Errorf("Parse(): got = %#v, want = %#v", got, tvc.want)
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
			name:  "Int",
			setup: func(r Register) interface{} { return IntArg(r, "t") },
			tests: []testValue{
				{
					name: "0",
					args: []string{"0"},
					want: 0,
				},
				{
					name: "1337",
					args: []string{"1337"},
					want: 1337,
				},
				// TODO(SuperPaintman): add negative numbers.
				// {
				// 	name: "-7331",
				// 	args: []string{"-t", "-7331"},
				// 	want: -7331,
				// },
			},
		},
		{
			name:  "String",
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
						t.Fatalf("Parse(): failed to parse args: %s", err)
					}

					got := deref(t, f)
					if !reflect.DeepEqual(got, tvc.want) {
						t.Errorf("Parse(): got = %#v, want = %#v", got, tvc.want)
					}
				})
			}
		})
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
