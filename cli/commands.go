package cli

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
