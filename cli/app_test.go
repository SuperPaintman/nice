package cli

import (
	"errors"
	"testing"
)

func TestCommanderInvalidNameCommand(t *testing.T) {
	tt := []struct {
		name    string
		command string
		want    error
	}{
		{
			name: "empty command name",
			want: &InvalidCommandError{Err: ErrMissingName},
		},
		{
			name:    "start dash in command name",
			command: "-help",
			want:    &InvalidCommandError{Name: "-help", Err: ErrInvalidName},
		},
		{
			name:    "ignore non-start dash in command name",
			command: "go-help",
			want:    nil,
		},
		{
			name:    "ignore end dash in command name",
			command: "help-",
			want:    nil,
		},
		{
			name:    "equal in command name",
			command: "help=test",
			want:    &InvalidCommandError{Name: "help=test", Err: ErrInvalidName},
		},
		{
			name:    "space in command name",
			command: "help test",
			want:    &InvalidCommandError{Name: "help test", Err: ErrInvalidName},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cmdr := commander{
				found: &Command{Name: tc.command},
				use:   func(*Command) error { return nil },
			}

			got := cmdr.SetCommand(tc.command)
			if !errors.Is(got, tc.want) {
				t.Fatalf("SetCommand(): got error = %q, want error = %q", got, tc.want)
			}
		})
	}
}
