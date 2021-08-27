package cli

import (
	"context"
	"fmt"
	"io"
	"os"
)

type App struct {
	Name     string
	Usage    string
	Action   Action
	Commands []Command
	Args     []string
	Stdout   io.Writer
	Stderr   io.Writer
	Stdin    io.Reader
	Parser   Parser
	Helper   Helper

	defaultParser *DefaultParser
}

type commander struct {
	app  *App
	next func(*Command) error

	command *Command
	found   *Command
}

func (c *commander) IsCommand(name string) bool {
	commands := c.app.Commands
	if c.command != nil {
		commands = c.command.Commands
	}

	for i := range commands {
		cmd := &commands[i]

		if cmd.Name == name {
			c.found = cmd
			return true
		}
	}

	return false
}

func (c *commander) SetCommand(name string) error {
	if c.found.Name != name {
		return fmt.Errorf("cli: command not found: %s", name)
	}

	c.command = c.found
	c.found = nil

	if err := c.next(c.command); err != nil {
		return err
	}

	return nil
}

func (app *App) RunContext(ctx context.Context) error {
	cmd := app.Command()
	path := []string{cmd.Name}

	// Run command.
	cmdCtx := newCommandContext(ctx, app, cmd, path)

	// Help flag.
	// TODO(SuperPaintman): move into global flags.
	showHelp := new(bool)
	if app.helpEnabled() {
		showHelp = Bool(cmdCtx, "help",
			WithShort("h"),
			WithUsage("Show information about a command"),
		)
	}

	if cmd.Action != nil {
		if err := cmd.Action.Setup(cmdCtx); err != nil {
			return err
		}
	}

	cmder := commander{
		app: app,
		next: func(c *Command) error {
			cmd = c

			// Update path and context.
			// Do not mutate previous contextes.
			newPath := make([]string, len(path)+1)
			copy(newPath, path)
			newPath[len(newPath)-1] = cmd.Name
			path = newPath

			cmdCtx = newCommandContext(cmdCtx, app, cmd, path)

			// Help flag.
			// TODO(SuperPaintman): move into global flags.
			if app.helpEnabled() {
				showHelp = Bool(cmdCtx, "help",
					WithShort("h"),
					WithUsage("Show information about a command"),
				)
			}

			// Setup a child command.
			if cmd.Action != nil {
				if err := cmd.Action.Setup(cmdCtx); err != nil {
					return err
				}
			}

			return nil
		},
	}

	if err := app.parser().Parse(&cmder, app.args()); err != nil {
		return err
	}

	if *showHelp {
		return app.help(cmdCtx, app.stdout())
	}

	if cmd.Action != nil {
		if err := cmd.Action.Run(cmdCtx); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) Run() error {
	return app.RunContext(context.Background())
}

// TODO(SuperPaintman):
// func (app *App) ShowHelpContext(ctx context.Context) error {
// 	cmdCtx := &commandContext{
// 		parent:  ctx,
// 		app:     app,
// 		command: app.root(),
// 	}
//
// 	path := []string{app.Name}
//
// 	return app.help(cmdCtx, path, app.stdout())
// }

// func (app *App) ShowHelp() error {
// 	return app.ShowHelpContext(context.Background())
// }

func (app *App) Command() *Command {
	return &Command{
		Name:     app.Name,
		Usage:    app.Usage,
		Action:   app.Action,
		Commands: app.Commands,
	}
}

func (app *App) args() []string {
	if app.Args != nil {
		return app.Args
	}

	return os.Args[1:]
}

func (app *App) stdout() io.Writer {
	if app.Stdout != nil {
		return app.Stdout
	}

	return os.Stdout
}

func (app *App) stderr() io.Writer {
	if app.Stderr != nil {
		return app.Stderr
	}

	return os.Stderr
}

func (app *App) stdin() io.Reader {
	if app.Stdin != nil {
		return app.Stdin
	}

	return os.Stdin
}

func (app *App) parser() Parser {
	if app.Parser != nil {
		return app.Parser
	}

	if app.defaultParser == nil {
		app.defaultParser = &DefaultParser{}
	}

	return app.defaultParser
}

func (app *App) helpEnabled() bool {
	var disabled bool
	if app.Helper != nil {
		_, disabled = app.Helper.(noopHelper)
	}

	return !disabled
}

func (app *App) help(ctx Context, w io.Writer) error {
	if app.Helper != nil {
		return app.Helper.Help(ctx, w)
	}

	return (DefaultHelper{}).Help(ctx, w)
}

type Command struct {
	Name     string
	Usage    string
	Action   Action
	Commands []Command
}

// TODO(SuperPaintman): fix this command.
// func HelpCommand() Command {
// 	return Command{
// 		Name:  "help",
// 		Usage: "Show information about a command",
// 		Action: SimpleActionFunc(func(ctx Context) error {
// 			return ctx.App().ShowHelpContext(ctx)
// 		}),
// 	}
// }

type Action interface {
	Setup(ctx Context) error
	Run(ctx Context) error
}

type ActionRunner func(ctx Context) error

type ActionBuilder func(ctx Context) ActionRunner

var _ Action = (SimpleActionFunc)(nil)

type SimpleActionFunc ActionRunner

func (fn SimpleActionFunc) Setup(ctx Context) error { return nil }

func (fn SimpleActionFunc) Run(ctx Context) error { return fn(ctx) }

var _ Action = (*actionFunc)(nil)

type actionFunc struct {
	builder ActionBuilder
	runner  ActionRunner
}

func ActionFunc(fn ActionBuilder) Action {
	return &actionFunc{
		builder: fn,
	}
}

func (a *actionFunc) Setup(ctx Context) error {
	a.runner = a.builder(ctx)

	return nil
}

func (a *actionFunc) Run(ctx Context) error {
	if a.runner == nil {
		return nil
	}

	return a.runner(ctx)
}
