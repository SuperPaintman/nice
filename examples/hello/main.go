package main

import (
	"log"

	"github.com/SuperPaintman/nice/cli"
)

func main() {
	app := cli.App{
		Name:  "hello",
		Usage: cli.Usage("Print a friendly greeting"),
		Action: cli.ActionFunc(func(cmd *cli.Command) cli.ActionRunner {
			name := cli.StringArg(cmd, "name",
				cli.Usage("Who we say hello to"),
				cli.Optional,
			)
			*name = "Nice" // Default value.

			return func(cmd *cli.Command) error {
				cmd.Printf("Hello, %s!\n", *name)

				return nil
			}
		}),
		CommandFlags: []cli.CommandFlag{
			cli.HelpCommandFlag(),
			cli.VersionCommandFlag("0.0.0"),
		},
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
