package cli

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
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

func TestParser_with_commands(t *testing.T) {
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
