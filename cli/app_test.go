package cli

import (
	"errors"
	"reflect"
	"testing"
)

func TestApp_Command(t *testing.T) {
	app := &App{
		Name: "test",
		Commands: []Command{
			{
				Name: "test-first",
			},
			{
				Name: "test-second",
				Commands: []Command{
					{
						Name: "test-second-first",
						Commands: []Command{
							{
								Name: "test-second-first-first",
							},
						},
					},
					{
						Name: "test-second-second",
					},
				},
			},
			{
				Name: "test-third",
			},
		},
	}

	tt := []struct {
		name string
		path []string
		want *Command
	}{
		{
			name: "test",
			path: []string{"test"},
			want: &Command{
				Name:     "test",
				Commands: app.Commands,
			},
		},
		{
			name: "test-first",
			path: []string{"test", "test-first"},
			want: &app.Commands[0],
		},
		{
			name: "test-second",
			path: []string{"test", "test-second"},
			want: &app.Commands[1],
		},
		{
			name: "test-second-first-first",
			path: []string{"test", "test-second", "test-second-first", "test-second-first-first"},
			want: &app.Commands[1].Commands[0].Commands[0],
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := app.Command(tc.path...)
			if err != nil {
				t.Fatalf("Command(%v): failed to get command: %s", tc.path, err)
			}

			if got.Name != tc.want.Name {
				t.Errorf("Command(%v): Name: got = %q, want = %q", tc.path, got.Name, tc.want.Name)
			}

			if !reflect.DeepEqual(got.Commands, tc.want.Commands) {
				t.Errorf("Command(%v): Commands: got = %#v, want = %#v", tc.path, got.Commands, tc.want.Commands)
			}
		})
	}
}

func TestApp_Command_not_found(t *testing.T) {
	app := &App{
		Name: "test",
		Commands: []Command{
			{
				Name: "test-first",
			},
			{
				Name: "test-second",
				Commands: []Command{
					{
						Name: "test-second-first",
					},
				},
			},
			{
				Name: "test-third",
			},
		},
	}

	tt := []struct {
		name string
		path []string
	}{
		{
			name: "empty path",
			path: nil,
		},
		{
			name: "empty string",
			path: []string{""},
		},
		{
			name: "non-existent root",
			path: []string{"unknown"},
		},
		{
			name: "non-existent child",
			path: []string{"test", "unknown"},
		},
		{
			name: "empty root with child",
			path: []string{"", "test-second", "test-second-first"},
		},
		{
			name: "non-existent root with child",
			path: []string{"unknown", "test-second", "test-second-first"},
		},
		{
			name: "without root with child",
			path: []string{"test-second"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, got := app.Command(tc.path...)
			want := ErrCommandNotFound
			if !errors.Is(got, want) {
				t.Fatalf("Command(%v): got error = %q, want error = %q", tc.path, got, want)
			}
		})
	}
}

var _ Action = (*actionSpy)(nil)

type actionSpy struct {
	setuped int
	runned  int
}

func (as *actionSpy) Setup(cmd *Command) error {
	as.setuped++
	return nil
}

func (as *actionSpy) Run(cmd *Command) error {
	as.runned++
	return nil
}

func (as *actionSpy) AssertSetuped(t *testing.T, name string, want int) {
	t.Helper()

	if as.setuped != want {
		t.Errorf("Command(): setuped %s: got = %v, want = %v", name, as.setuped, want)
	}
}

func TestApp_Command_setup_once(t *testing.T) {
	var (
		testAction            actionSpy
		testFirstAction       actionSpy
		testSecondAction      actionSpy
		testSecondFirstAction actionSpy
		testThirdAction       actionSpy
	)

	app := &App{
		Name:   "test",
		Action: &testAction,
		Commands: []Command{
			{
				Name:   "test-first",
				Action: &testFirstAction,
			},
			{
				Name:   "test-second",
				Action: &testSecondAction,
				Commands: []Command{
					{
						Name:   "test-second-first",
						Action: &testSecondFirstAction,
					},
				},
			},
			{
				Name:   "test-third",
				Action: &testThirdAction,
			},
		},
	}

	// Root.
	if _, err := app.Command("test"); err != nil {
		t.Fatalf("Command(): failed to get command: %s", err)
	}

	testAction.AssertSetuped(t, "test", 1)
	testFirstAction.AssertSetuped(t, "test-first", 0)
	testSecondAction.AssertSetuped(t, "test-second", 0)
	testSecondFirstAction.AssertSetuped(t, "test-second-first", 0)
	testThirdAction.AssertSetuped(t, "test-third", 0)

	// Deep child.
	if _, err := app.Command("test", "test-second", "test-second-first"); err != nil {
		t.Fatalf("Command(): failed to get command: %s", err)
	}

	testAction.AssertSetuped(t, "test", 1)
	testFirstAction.AssertSetuped(t, "test-first", 0)
	testSecondAction.AssertSetuped(t, "test-second", 1)
	testSecondFirstAction.AssertSetuped(t, "test-second-first", 1)
	testThirdAction.AssertSetuped(t, "test-third", 0)

	// Second hild.
	if _, err := app.Command("test", "test-second"); err != nil {
		t.Fatalf("Command(): failed to get command: %s", err)
	}

	testAction.AssertSetuped(t, "test", 1)
	testFirstAction.AssertSetuped(t, "test-first", 0)
	testSecondAction.AssertSetuped(t, "test-second", 1)
	testSecondFirstAction.AssertSetuped(t, "test-second-first", 1)
	testThirdAction.AssertSetuped(t, "test-third", 0)
}

func TestCommanderInvalidNameApp(t *testing.T) {
	tt := []struct {
		name    string
		appName string
		want    error
	}{
		{
			name: "empty app name",
			want: &InvalidCommandError{Err: ErrMissingName},
		},
		{
			name:    "start dash in app name",
			appName: "-help",
			want:    &InvalidCommandError{Name: "-help", Err: ErrInvalidName},
		},
		{
			name:    "ignore non-start dash in app name",
			appName: "go-help",
			want:    nil,
		},
		{
			name:    "ignore end dash in app name",
			appName: "help-",
			want:    nil,
		},
		{
			name:    "equal in app name",
			appName: "help=test",
			want:    &InvalidCommandError{Name: "help=test", Err: ErrInvalidName},
		},
		{
			name:    "space in app name",
			appName: "help test",
			want:    &InvalidCommandError{Name: "help test", Err: ErrInvalidName},
		},
		{
			name:    "comma in app name",
			appName: "help,test",
			want:    &InvalidCommandError{Name: "help,test", Err: ErrInvalidName},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			app := App{
				Name: tc.appName,
			}

			_, got := app.command()
			if !errors.Is(got, tc.want) {
				t.Fatalf("command(): got error = %q, want error = %q", got, tc.want)
			}
		})
	}
}

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
		{
			name:    "comma in command name",
			command: "help,test",
			want:    &InvalidCommandError{Name: "help,test", Err: ErrInvalidName},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cmdr := commander{
				command: &Command{
					Name: tc.command,
					Commands: []Command{
						{Name: "go-help"},
						{Name: "help-"},
					},
				},
				use: func(*Command) (Register, error) { return &DefaultRegister{}, nil },
			}

			_, got := cmdr.SetCommand(tc.command)
			if !errors.Is(got, tc.want) {
				t.Fatalf("SetCommand(): got error = %q, want error = %q", got, tc.want)
			}
		})
	}
}

func TestCommanderUnknownCommand(t *testing.T) {
	cmdr := commander{
		command: &Command{
			Name: "test",
		},
		use: func(*Command) (Register, error) { return &DefaultRegister{}, nil },
	}

	_, got := cmdr.SetCommand("help")
	want := &InvalidCommandError{Name: "help", Err: ErrUnknown}
	if !errors.Is(got, want) {
		t.Fatalf("SetCommand(): got error = %q, want error = %q", got, want)
	}
}
