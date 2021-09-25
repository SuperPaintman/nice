package cli

func HelpCommand() Command {
	return Command{
		Name:  "help",
		Usage: Usage("Show information about a command"),
		Action: ActionFunc(func(cmd *Command) ActionRunner {
			// Do not mutate the previous path.
			path := cmd.Path()
			newPath := make([]string, len(path)-1)
			copy(newPath, path[:len(path)-1])
			path = newPath

			_ = RestStringsVar(cmd, &path, "command")

			return func(cmd *Command) error {
				cmd, err := cmd.App().Command(path...)
				if err != nil {
					return err
				}

				return cmd.App().Help(cmd, cmd.Stdout())
			}
		}),
	}
}

func HelpCommandFlag() CommandFlag {
	return CommandFlag{
		Long:  "help",
		Short: "h",
		Usage: Usage("Show information about a command"),
		Action: ActionRunner(func(cmd *Command) error {
			return cmd.App().Help(cmd, cmd.Stdout())
		}),
	}
}

func VersionCommandFlag(version string) CommandFlag {
	return CommandFlag{
		Long:  "version",
		Short: "v",
		Usage: Usage("Print version information and quit"),
		Action: ActionRunner(func(cmd *Command) error {
			if _, err := cmd.Printf(version); err != nil {
				return err
			}

			if _, err := cmd.Printf("\n"); err != nil {
				return err
			}

			return nil
		}),
	}
}
