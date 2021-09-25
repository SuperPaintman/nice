package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrCommandNotFound = fmt.Errorf("cli: command not found")
)

type InvalidCommandError struct {
	Name string
	Err  error
}

func (e *InvalidCommandError) Error() string {
	msg := "unknown error"
	if e.Err != nil {
		msg = e.Err.Error()
	}

	if e.Name != "" {
		return fmt.Sprintf("broken command: %s", msg)
	} else {
		return fmt.Sprintf("broken command: '%s': %s", e.Name, msg)
	}
}

func (e *InvalidCommandError) Is(err error) bool {
	pe, ok := err.(*InvalidCommandError)
	return ok && pe.Name == e.Name && errors.Is(pe.Err, e.Err)
}

var _ Commander = (*commander)(nil)

type commander struct {
	app *App
	use func(*Command) (Register, error)

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

func (c *commander) SetCommand(name string) (Register, error) {
	if name == "" {
		return nil, &InvalidCommandError{Err: ErrMissingName}
	}

	if !validCommandName(name) {
		return nil, &InvalidCommandError{
			Name: name,
			Err:  ErrInvalidName,
		}
	}

	if c.found.Name != name {
		// Internal error. Something went wrong in IsCommand.
		return nil, ErrCommandNotFound
	}

	c.command = c.found
	c.found = nil

	register, err := c.use(c.command)
	if err != nil {
		return nil, err
	}

	return register, nil
}

func validCommandName(name string) bool {
	return validArg(name)
}

type App struct {
	Name         string
	Usage        Usager
	Action       Action
	CommandFlags []CommandFlag
	Commands     []Command
	Args         []string
	Stdout       io.Writer
	Stderr       io.Writer
	Stdin        io.Reader
	Parser       Parser
	NewRegister  func() Register
	Helper       Helper

	ctx           context.Context
	rootCmd       *Command
	defaultParser *DefaultParser
}

func (app *App) RunContext(ctx context.Context) error {
	// Inject context into the app.
	app.ctx = ctx

	// Build the root command.
	cmd, err := app.command()
	if err != nil {
		return err
	}

	cmder := commander{
		app: app,
		use: func(c *Command) (Register, error) {
			// Update path and init new command.
			// Do not mutate the previous path.
			path := cmd.Path()

			newPath := make([]string, len(path)+1)
			copy(newPath, path)
			newPath[len(newPath)-1] = c.Name
			path = newPath

			c.init(ctx, app, cmd, app.newRegister(), path)

			// Setup a child command.
			if err := c.setup(); err != nil {
				return nil, err
			}

			// Set new cmd.
			cmd = c

			return cmd.register, nil
		},
	}

	if err := app.parser().Parse(&cmder, cmd.register, app.args()); err != nil {
		return err
	}

	// Find and run command flag.
	cmdWithCFS := cmd
	for cmdWithCFS != nil {
		for i := range cmdWithCFS.CommandFlags {
			f := &cmdWithCFS.CommandFlags[i]
			if !f.value {
				continue
			}

			if f.Action != nil {
				if err := f.Action.Setup(app, cmd); err != nil {
					return err
				}

				if err := f.Action.Run(app, cmd); err != nil {
					return err
				}
			}

			return nil
		}

		cmdWithCFS = cmdWithCFS.parent
	}

	// Run action.
	if cmd.Action != nil {
		if err := cmd.Action.Run(app, cmd); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) Run() error {
	return app.RunContext(context.Background())
}

func (app *App) Command(path ...string) (*Command, error) {
	if len(path) == 0 || path[0] != app.Name {
		return nil, ErrCommandNotFound
	}

	cmd, err := app.command()
	if err != nil {
		return nil, err
	}

	for j, name := range path[1:] {
		// Find a sub command.
		var found bool
		for i := range cmd.Commands {
			c := &cmd.Commands[i]

			if c.Name == name {
				found = true
				c.init(app.ctx, app, cmd, app.newRegister(), path[:j+2])

				// Setup a command.
				if err := c.setup(); err != nil {
					return nil, err
				}

				cmd = c
				break
			}
		}

		if !found {
			return nil, ErrCommandNotFound
		}
	}

	return cmd, nil
}

func (app *App) Help(cmd *Command, w io.Writer) error {
	if app.Helper != nil {
		return app.Helper.Help(app, cmd, w)
	}

	return (DefaultHelper{}).Help(app, cmd, w)
}

func (app *App) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(app.stdout(), format, a...)
}

func (app *App) Print(a ...interface{}) (n int, err error) {
	return fmt.Fprint(app.stdout(), a...)
}

func (app *App) Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(app.stdout(), a...)
}

func (app *App) Warnf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(app.stderr(), format, a...)
}

func (app *App) Warn(a ...interface{}) (n int, err error) {
	return fmt.Fprint(app.stderr(), a...)
}

func (app *App) Warnln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(app.stderr(), a...)
}

func (app *App) command() (*Command, error) {
	if app.rootCmd == nil {
		if app.Name == "" {
			return nil, &InvalidCommandError{Err: ErrMissingName}
		}

		if !validCommandName(app.Name) {
			return nil, &InvalidCommandError{
				Name: app.Name,
				Err:  ErrInvalidName,
			}
		}

		cmd := &Command{
			Name:         app.Name,
			Usage:        app.Usage,
			Action:       app.Action,
			CommandFlags: app.CommandFlags,
			Commands:     app.Commands,
		}

		path := []string{cmd.Name}
		cmd.init(app.ctx, app, nil, app.newRegister(), path)

		if err := cmd.setup(); err != nil {
			return nil, err
		}

		app.rootCmd = cmd
	}

	return app.rootCmd, nil
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

func (app *App) newRegister() Register {
	if app.NewRegister != nil {
		if register := app.NewRegister(); register != nil {
			return register
		}
	}

	return &DefaultRegister{}
}

var _ Register = (*Command)(nil)

type Command struct {
	Name         string
	Usage        Usager
	Action       Action
	CommandFlags []CommandFlag
	Commands     []Command

	ctx        context.Context
	app        *App
	parent     *Command
	register   Register
	path       []string
	initilized bool
	setuped    bool
}

func (c *Command) Path() []string { return c.path }

func (c *Command) Context() context.Context {
	if c.ctx == nil {
		return context.Background()
	}

	return c.ctx
}

func (c *Command) RegisterFlag(flag Flag) error {
	return c.register.RegisterFlag(flag)
}

func (c *Command) RegisterArg(arg Arg) error {
	return c.register.RegisterArg(arg)
}

func (c *Command) RegisterRestArgs(rest RestArgs) error {
	return c.register.RegisterRestArgs(rest)
}

func (c *Command) Arg(i int) (*Arg, bool) { return c.register.Arg(i) }

func (c *Command) ShortFlag(name string) (*Flag, bool) { return c.register.ShortFlag(name) }

func (c *Command) LongFlag(name string) (*Flag, bool) { return c.register.LongFlag(name) }

func (c *Command) Args() []Arg { return c.register.Args() }

func (c *Command) Rest() *RestArgs { return c.register.Rest() }

func (c *Command) Flags() []Flag { return c.register.Flags() }

func (c *Command) Err() error { return c.register.Err() }

func (c *Command) init(ctx context.Context, app *App, parent *Command, register Register, path []string) {
	if c.initilized {
		return
	}

	c.ctx = ctx
	c.app = app
	c.parent = parent
	c.register = register
	c.path = path

	c.initilized = true
}

func (c *Command) setup() error {
	if c.setuped {
		return nil
	}

	if c.Action != nil {
		if err := c.Action.Setup(c.app, c); err != nil {
			return err
		}
	}

	// Add command flags.
	var parents []*Command

	cmdParent := c.parent
	for cmdParent != nil {
		parents = append(parents, cmdParent)
		cmdParent = cmdParent.parent
	}

	addCommandFlags := func(cfs []CommandFlag) error {
		for i := range cfs {
			f := &cfs[i]

			err := BoolVar(c, &f.value, f.Long,
				WithShort(f.Short),
				WithUsage(f.Usage),
				commandFlag(true), // Mark this flag as "magic" command flag.
			)
			if err != nil {
				return err
			}
		}

		return nil
	}

	for i := len(parents) - 1; i >= 0; i-- {
		if err := addCommandFlags(parents[i].CommandFlags); err != nil {
			return err
		}
	}

	if err := addCommandFlags(c.CommandFlags); err != nil {
		return err
	}

	c.setuped = true

	return nil
}

type Action interface {
	Setup(app *App, cmd *Command) error
	Run(app *App, cmd *Command) error
}

// TODO(SuperPaintman): add Before.
// TODO(SuperPaintman): add After.

type ActionRunner func(app *App, cmd *Command) error

type ActionBuilder func(app *App, cmd *Command) ActionRunner

var _ Action = (ActionRunner)(nil)

func (fn ActionRunner) Setup(app *App, cmd *Command) error { return nil }

func (fn ActionRunner) Run(app *App, cmd *Command) error { return fn(app, cmd) }

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

func (a *actionFunc) Setup(app *App, cmd *Command) error {
	a.runner = a.builder(app, cmd)

	return nil
}

func (a *actionFunc) Run(app *App, cmd *Command) error {
	if a.runner == nil {
		return nil
	}

	return a.runner(app, cmd)
}

type CommandFlag struct {
	Short  string
	Long   string
	Usage  Usager
	Action Action

	value bool
}
