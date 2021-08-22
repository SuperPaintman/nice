package cli

import (
	"reflect"
	"testing"
)

func TestParser(t *testing.T) {
	var parser DefaultParser

	show := Bool(&parser, "show",
		WithShort("s"),
		WithUsage("Shows the resuld of the function."),
	)

	recreate := Bool(&parser, "recreate",
		WithUsage("Re-create the current user."),
	)

	update := Bool(&parser, "update",
		WithUsage("Update the DB."),
	)

	unused := Bool(&parser, "unused")

	count := Int(&parser, "count",
		WithShort("c"),
	)

	userID := IntArg(&parser, "user-id",
		WithUsage("Current User ID."),
	)

	args := []string{
		"--show", "--recreate=false", "-c", "100500", "1337", "--update", "true",
		"--first-unknown", "other", "vals", "--second-unknown", "in", "args",
	}
	if err := parser.Parse(args); err != nil {
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
