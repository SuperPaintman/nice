package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

type Context interface {
	context.Context
	Register

	App() *App
	Command() *Command
	Parser() Parser
	Args() []Arg
	Flags() []Flag
	Printf(format string, a ...interface{}) (n int, err error)
	Print(a ...interface{}) (n int, err error)
	Println(a ...interface{}) (n int, err error)
	Warnf(format string, a ...interface{}) (n int, err error)
	Warn(a ...interface{}) (n int, err error)
	Warnln(a ...interface{}) (n int, err error)
}

var _ (Context) = (*commandContext)(nil)

type commandContext struct {
	parent context.Context

	app     *App
	command *Command
	path    []string
}

func (c *commandContext) Deadline() (deadline time.Time, ok bool) {
	if c.parent != nil {
		return c.parent.Deadline()
	}

	return
}

func (c *commandContext) Done() <-chan struct{} {
	if c.parent != nil {
		return c.parent.Done()
	}

	return nil
}

func (c *commandContext) Err() error {
	if c.parent != nil {
		return c.parent.Err()
	}

	return nil
}

func (c *commandContext) Value(key interface{}) interface{} {
	if c.parent != nil {
		return c.parent.Value(key)
	}

	return nil
}

func (c *commandContext) RegisterFlag(flag Flag) error {
	return c.app.parser().RegisterFlag(flag)
}

func (c *commandContext) RegisterArg(arg Arg) error {
	return c.app.parser().RegisterArg(arg)
}

func (c *commandContext) App() *App { return c.app }

func (c *commandContext) Command() *Command { return c.command }

func (c *commandContext) Parser() Parser { return c.app.parser() }

func (c *commandContext) Path() []string { return c.path }

func (c *commandContext) Args() []Arg { return c.app.parser().Args() }

func (c *commandContext) Flags() []Flag { return c.app.parser().Flags() }

func (c *commandContext) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(c.app.stdout(), format, a...)
}

func (c *commandContext) Print(a ...interface{}) (n int, err error) {
	return fmt.Fprint(c.app.stdout(), a...)
}

func (c *commandContext) Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(c.app.stdout(), a...)
}

func (c *commandContext) Warnf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(c.app.stderr(), format, a...)
}

func (c *commandContext) Warn(a ...interface{}) (n int, err error) {
	return fmt.Fprint(c.app.stderr(), a...)
}

func (c *commandContext) Warnln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(c.app.stderr(), a...)
}

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

func (app *App) RunContext(ctx context.Context) error {
	// TODO(SuperPaintman): add fast parser.
	var cmdParser DefaultParser

	if err := cmdParser.Parse(app.args()); err != nil {
		return nil
	}

	// lastName := app.Name
	path := []string{app.Name}

	// rootCmd := true
	cmd := app.root()
	for _, name := range cmdParser.rest {
		// lastName = name

		// Find a sub command.
		var found bool
		for i := range cmd.Commands {
			// NOTE(SuperPaintman): inherit parent's flags?

			c := &cmd.Commands[i]

			if c.Name == name {
				found = true

				path = append(path, c.Name)
				// rootCmd = false
				cmd = c
				break
			}
		}

		if !found {
			// return fmt.Errorf("cli: command not found: %s", name)
			break
		}
	}

	// Run command.
	cmdCtx := &commandContext{
		parent:  ctx,
		app:     app,
		command: cmd,
	}

	// Help flag.
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

	if err := app.parser().Parse(app.args()); err != nil {
		return err
	}

	if *showHelp {
		return app.help(cmdCtx, app.stdout(), path)
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

func (app *App) root() *Command {
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

func (app *App) help(ctx Context, w io.Writer, path []string) error {
	if app.Helper != nil {
		return app.Helper.Help(ctx, w, path)
	}

	return (DefaultHelper{}).Help(ctx, w, path)
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
