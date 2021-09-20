package cli

func HelpCommand() Command {
	return Command{
		Name:  "help",
		Usage: Usage("Show information about a command"),
		Action: ActionFunc(func(ctx Context) ActionRunner {
			path := ctx.Path()
			path = path[:len(path)-1]
			_ = RestStringsVar(ctx, &path, "command")

			return func(ctx Context) error {
				return ctx.Help(ctx, path, ctx.Stdout())
			}
		}),
	}
}

func HelpCommandFlag() CommandFlag {
	return CommandFlag{
		Long:  "help",
		Short: "h",
		Usage: Usage("Show information about a command"),
		Action: ActionRunner(func(ctx Context) error {
			return ctx.Help(ctx, ctx.Path(), ctx.Stdout())
		}),
	}
}

func VersionFlag(version string) CommandFlag {
	return CommandFlag{
		Long:  "version",
		Short: "v",
		Usage: Usage("Print version information and quit"),
		Action: ActionRunner(func(ctx Context) error {
			if _, err := ctx.Printf(version); err != nil {
				return err
			}

			if _, err := ctx.Printf("\n"); err != nil {
				return err
			}

			return nil
		}),
	}
}
