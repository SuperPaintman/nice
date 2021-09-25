package cli

import (
	"os"
)

func HelpCommand() Command {
	return Command{
		Name:  "help",
		Usage: Usage("Show information about a command"),
		Action: ActionFunc(func(app *App, cmd *Command) ActionRunner {
			// Do not mutate the previous path.
			path := cmd.Path()
			newPath := make([]string, len(path)-1)
			copy(newPath, path[:len(path)-1])
			path = newPath

			_ = RestStringsVar(cmd, &path, "command")

			return func(app *App, cmd *Command) error {
				cmd, err := app.Command(path...)
				if err != nil {
					return err
				}

				if app.Stdout == nil {
					return app.Help(cmd, os.Stdout)
				} else {
					return app.Help(cmd, app.Stdout)
				}
			}
		}),
	}
}

func HelpCommandFlag() CommandFlag {
	return CommandFlag{
		Long:  "help",
		Short: "h",
		Usage: Usage("Show information about a command"),
		Action: ActionRunner(func(app *App, cmd *Command) error {
			if app.Stdout == nil {
				return app.Help(cmd, os.Stdout)
			} else {
				return app.Help(cmd, app.Stdout)
			}
		}),
	}
}

func VersionCommandFlag(version string) CommandFlag {
	return CommandFlag{
		Long:  "version",
		Short: "v",
		Usage: Usage("Print version information and quit"),
		Action: ActionRunner(func(app *App, cmd *Command) error {
			if _, err := app.Printf(version); err != nil {
				return err
			}

			if _, err := app.Printf("\n"); err != nil {
				return err
			}

			return nil
		}),
	}
}
