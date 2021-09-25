package cli

import (
	"io/ioutil"
	"os/exec"
	"strings"
	"testing"
)

const (
	simpleHelp = `Usage: simple [options...] <name> [real-name] [role]
       simple [options...] [command]

Simple app for simple tasks

Commands:
  list-buckets    Print buckets list for user
  help            Show information about a command

Arguments:
  <name> string         Name of the user
  [real-name] string    Real name of the user
  [role] string         Role of the user (default: user)

Options:
      --hide
  -u, --update     Update user if exists
      --uid int    (required, default: 1000)
      --gid int    User's group ID (default: 1000)
  -h, --help       Show information about a command
  -v, --version    Print version information and quit
`

	listBucketsHelp = `Usage: simple list-buckets [options...] <name> [buckets...]

Print buckets list for user

Arguments:
  <name> string         Name of the user
  [buckets...] []int    Bucket IDs

Options:
      --show-hidden
  -h, --help       Show information about a command
  -v, --version    Print version information and quit
`
)

func TestDefaultHelper_Help(t *testing.T) {
	app := App{
		Name:  "simple",
		Usage: Usage("Simple app for simple tasks"),
		Action: ActionFunc(func(cmd *Command) ActionRunner {
			_ = StringArg(cmd, "name",
				Usage("Name of the user"),
			)

			_ = StringArg(cmd, "real-name",
				Usage("Real name of the user"),
				Optional,
			)

			role := StringArg(cmd, "role",
				Usage("Role of the user"),
				Optional,
			)
			*role = "user"

			_ = Bool(cmd, "hide")

			_ = Bool(cmd, "update",
				WithShort("u"),
				Usage("Update user if exists"),
			)

			// _ = String(cmd, "comment") // ShowEmptyDefault,

			uid := Int(cmd, "uid",
				Required,
			)
			*uid = 1000

			gid := Int(cmd, "gid",
				Usage("User's group ID"),
			)
			*gid = 1000

			return func(cmd *Command) error { panic("not implemented") }
		}),
		Commands: []Command{
			{
				Name:  "list-buckets",
				Usage: Usage("Print buckets list for user"),
				Action: ActionFunc(func(cmd *Command) ActionRunner {
					_ = StringArg(cmd, "name",
						Usage("Name of the user"),
					)

					_ = RestInts(cmd, "buckets",
						Usage("Bucket IDs"),
					)

					_ = Bool(cmd, "show-hidden")

					return func(cmd *Command) error { panic("not implemented") }
				}),
			},
			HelpCommand(),
		},
		CommandFlags: []CommandFlag{
			HelpCommandFlag(),
			VersionCommandFlag("0.0.0-test"),
		},
	}

	tt := []struct {
		name string
		path []string
		want string
	}{
		{
			name: "simple",
			path: []string{"simple"},
			want: simpleHelp,
		},
		{
			name: "simple list-buckets",
			path: []string{"simple", "list-buckets"},
			want: listBucketsHelp,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cmd, err := app.Command(tc.path...)
			if err != nil {
				t.Fatalf("Command(%v): failed to get command: %s", tc.path, err)
			}

			var (
				helper DefaultHelper
				buf    strings.Builder
			)
			if err := helper.Help(cmd, &buf); err != nil {
				t.Fatalf("Help(): failed to write help: %s", err)
			}

			assertStringsDiff(t, buf.String(), tc.want)
		})
	}
}

func assertStringsDiff(t *testing.T, got, want string) {
	t.Helper()

	if err := exec.Command("which", "diff").Run(); err == nil {
		wantTMP, err := ioutil.TempFile(t.TempDir(), "want-*")
		if err != nil {
			t.Errorf("Failed to create a temp file for want: %s", err)
			return
		}
		defer wantTMP.Close()

		if _, err := wantTMP.WriteString(want); err != nil {
			t.Errorf("Failed to write content into the want temp file: %s", err)
			return
		}

		gotTMP, err := ioutil.TempFile(t.TempDir(), "got-*")
		if err != nil {
			t.Errorf("Failed to create a temp file for got: %s", err)
			return
		}
		defer gotTMP.Close()

		if _, err := gotTMP.WriteString(got); err != nil {
			t.Errorf("Failed to write content into the got temp file: %s", err)
			return
		}

		cmd := exec.Command("diff", "-u", wantTMP.Name(), gotTMP.Name())
		cmd.Stdout = &logPipe{t}
		cmd.Stderr = &logPipe{t}

		if err := cmd.Run(); err != nil {
			t.Errorf("Diff returned error: %s", err)
			return
		}
	}

	if got != want {
		t.Errorf("got:\n```\n%s```\n\nwant:\n```\n%s```\n", got, want)
	}
}

type logPipe struct {
	t testing.TB
}

func (lp *logPipe) Write(p []byte) (n int, err error) {
	lp.t.Logf("\n%s", p)
	return 0, nil
}
