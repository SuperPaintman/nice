package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
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

	if e.Name == "" {
		return fmt.Sprintf("cli: broken command: %s", msg)
	}

	return fmt.Sprintf("cli: broken command: '%s': %s", e.Name, msg)
}

func (e *InvalidCommandError) Is(err error) bool {
	pe, ok := err.(*InvalidCommandError)
	return ok && pe.Name == e.Name && errors.Is(pe.Err, e.Err)
}

type ExitCode int

func (e ExitCode) Error() string {
	return fmt.Sprintf("cli: exit code: %d", e)
}

type CommandError struct {
	Command *Command
	Err     error

	exitCode ExitCode
}

func (e *CommandError) Error() string {
	msg := "unknown error"
	if e.Err != nil {
		msg = e.Err.Error()
	}

	if e.Command == nil {
		return fmt.Sprintf("cli: command error: %s", msg)
	}

	if path := e.Command.Path(); len(path) != 0 {
		return fmt.Sprintf("cli: command error: '%s': %s", strings.Join(path, " "), msg)
	}

	if e.Command.Name == "" {
		return fmt.Sprintf("cli: command error: %s", msg)
	}

	return fmt.Sprintf("cli: command error: '%s': %s", e.Command.Name, msg)
}

func (e *CommandError) Unwrap() error { return e.Err }

func (e *CommandError) Is(err error) bool {
	ce, ok := err.(*CommandError)
	return ok && ce.Command == e.Command && ce.exitCode == e.exitCode && errors.Is(ce.Err, e.Err)
}

func (e *CommandError) ExitCode() ExitCode {
	if e.exitCode == 0 {
		return 1
	}

	return e.exitCode
}

var _ Commander = (*commander)(nil)

type commander struct {
	app *App
	use func(*Command) (Register, error)

	command *Command
}

func (c *commander) IsCommand(name string) bool {
	var commands []Command
	if c.command != nil {
		commands = c.command.Commands
	} else if c.app != nil {
		commands = c.app.Commands
	}

	for i := range commands {
		cmd := &commands[i]

		if cmd.Name == name {
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

	// Find a command.
	var commands []Command
	if c.command != nil {
		commands = c.command.Commands
	} else if c.app != nil {
		commands = c.app.Commands
	}

	var found *Command
	for i := range commands {
		cmd := &commands[i]

		if cmd.Name == name {
			found = cmd
			break
		}
	}

	if found == nil {
		return nil, &InvalidCommandError{
			Name: name,
			Err:  ErrUnknown,
		}
	}

	c.command = found

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
				if err := f.Action.Setup(cmd); err != nil {
					return err
				}

				if err := f.Action.Run(cmd); err != nil {
					return err
				}
			}

			return nil
		}

		cmdWithCFS = cmdWithCFS.parent
	}

	// Run action.
	if cmd.Action != nil {
		if err := cmd.Action.Run(cmd); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) Run() error {
	return app.RunContext(context.Background())
}

func (app *App) RootCommand(path ...string) (*Command, error) {
	cmd, err := app.command()
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func (app *App) Command(path ...string) (*Command, error) {
	if len(path) == 0 || path[0] != app.Name {
		return nil, nil
	}

	cmd, err := app.command()
	if err != nil {
		return nil, err
	}

	for j, name := range path[1:] {
		// Find a subcommand.
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
			return nil, nil
		}
	}

	return cmd, nil
}

func (app *App) Help(cmd *Command, w io.Writer) error {
	if app.Helper != nil {
		return app.Helper.Help(cmd, w)
	}

	return (DefaultHelper{}).Help(cmd, w)
}

func (app *App) HandleError(err error) {
	code := app.handleError(err, app.stderr())
	if code != 0 {
		os.Exit(int(code))
	}
}

func (app *App) handleError(err error, w io.Writer) (exitCode ExitCode) {
	if err == nil {
		return
	}

	// TODO(SuperPaintman):
	//     I think in might be better to add `Friendly()` method to call it
	//     directly on the error.

	exitCode = 1

	ew := easyWriter{w: w}

	// NOTE(SuperPaintman): ParseValueError is not a top level error.

	cmdErr := &CommandError{}
	invalidCommandErr := &InvalidCommandError{}
	parseArgErr := &ParseArgError{}
	parseFlagErr := &ParseFlagError{}
	flagErr := &FlagError{}
	argErr := &ArgError{}
	restArgsErr := &RestArgsError{}
	switch {
	case errors.As(err, &exitCode):
		// Nothing to do.

	case errors.As(err, &cmdErr):
		exitCode = cmdErr.ExitCode()

	case errors.As(err, &invalidCommandErr):
		switch {
		case errors.Is(invalidCommandErr.Err, ErrMissingName):
			ew.Writef("Missing command name\n")

		case errors.Is(invalidCommandErr.Err, ErrInvalidName):
			ew.Writef("Invalie command name: %s\n", invalidCommandErr.Name)

		default:
			ew.WriteString(err.Error())
			ew.WriteString("\n")
		}

	case errors.As(err, &parseArgErr):
		switch {
		case errors.Is(parseArgErr.Err, ErrUnknown):
			ew.Writef("Unknown %s argument: %s\n", nthNumber(parseArgErr.Index), parseArgErr.Arg)

		default:
			ew.WriteString(err.Error())
			ew.WriteString("\n")
		}

	case errors.As(err, &parseFlagErr):
		switch {
		case errors.Is(parseFlagErr.Err, ErrSyntax):
			ew.Writef("Invalid flag syntax: %s\n", parseFlagErr.Name)

		case errors.Is(parseFlagErr.Err, ErrUnknown):
			ew.Writef("Unknown flag: %s\n", parseFlagErr.Name)

		default:
			ew.WriteString(err.Error())
			ew.WriteString("\n")
		}

	case errors.As(err, &flagErr):
		parser := app.parser()
		writeFlagName := func(short, long string) {
			if long != "" {
				if short != "" {
					ew.WriteString(parser.FormatShortFlag(short))
					ew.WriteString(" ")
				}

				ew.WriteString(parser.FormatLongFlag(long))
			} else if short != "" {
				ew.WriteString(parser.FormatShortFlag(short))
			}
		}

		parseValueError := &ParseValueError{}
		switch {
		case errors.Is(flagErr.Err, ErrMissingName):
			ew.Writef("Unable to register a flag without a short or long name\n")

		case errors.Is(flagErr.Err, ErrInvalidName):
			if flagErr.Long != "" {
				ew.Writef("Unable to register a flag with an invalid long name: %s\n", flagErr.Long)
			} else {
				ew.Writef("Unable to register a flag with an invalid short name: %s\n", flagErr.Short)
			}

		case errors.Is(flagErr.Err, ErrDuplicate):
			ew.WriteString("Unable to register a duplicate flag: ")
			writeFlagName(flagErr.Short, flagErr.Long)
			ew.WriteString("\n")

		case errors.Is(flagErr.Err, ErrNotProvided):
			ew.WriteString("Flag is required: ")
			writeFlagName(flagErr.Short, flagErr.Long)
			ew.WriteString("\n")

		case errors.As(flagErr.Err, &parseValueError):
			ew.WriteString("Invalid ")
			writeFlagName(flagErr.Short, flagErr.Long)
			ew.WriteString(" flag value: ")
			ew.WriteString(parseValueError.Error())
			ew.WriteString("\n")

		default:
			ew.WriteString(err.Error())
			ew.WriteString("\n")
		}

	case errors.As(err, &argErr):
		parseValueError := &ParseValueError{}
		switch {
		case errors.Is(argErr.Err, ErrRequiredAfterOptional):
			ew.Writef("Unable to register a required %s argument (%s) after another optional argument",
				nthNumber(argErr.Index), argErr.Name,
			)

		case errors.Is(argErr.Err, ErrArgAfterRest):
			ew.Writef("Unable to register the %s argument (%s) after the rest arguments\n",
				nthNumber(argErr.Index), argErr.Name,
			)

		case errors.Is(argErr.Err, ErrMissingName):
			ew.Writef("Unable to register the %s argument without a name\n", nthNumber(argErr.Index))

		case errors.Is(argErr.Err, ErrInvalidName):
			ew.Writef("Unable to register the %s argument with an invalid name: %s\n",
				nthNumber(argErr.Index), argErr.Name,
			)

		case errors.Is(argErr.Err, ErrDuplicate):
			ew.Writef("Unable to register a duplicate %s argument: %s\n",
				nthNumber(argErr.Index), argErr.Name,
			)

		case errors.Is(argErr.Err, ErrNotProvided):
			ew.Writef("%s argument (%s) is required\n",
				nthNumber(argErr.Index), argErr.Name,
			)

		case errors.As(argErr.Err, &parseValueError):
			ew.Writef("Invalid %s argument (%s) value: ",
				nthNumber(argErr.Index), argErr.Name,
			)
			ew.WriteString(parseValueError.Error())
			ew.WriteString("\n")

		default:
			ew.WriteString(err.Error())
			ew.WriteString("\n")
		}

	case errors.As(err, &restArgsErr):
		switch {
		case errors.Is(restArgsErr.Err, ErrInvalidName):
			ew.Writef("Unable to register the rest arguments with an invalid name: %s\n",
				restArgsErr.Name,
			)

		case errors.Is(restArgsErr.Err, ErrDuplicate):
			ew.Writef("Unable to register another rest arguments: %s\n", restArgsErr.Name)

		default:
			ew.WriteString(err.Error())
			ew.WriteString("\n")
		}

	default:
		ew.WriteString(err.Error())
		ew.WriteString("\n")
	}

	return
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

func (c *Command) App() *App { return c.app }

func (c *Command) Path() []string { return c.path }

func (c *Command) Context() context.Context {
	if c.ctx == nil {
		return context.Background()
	}

	return c.ctx
}

func (c *Command) Parser() Parser { return c.app.parser() }

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

func (c *Command) Stdout() io.Writer { return c.app.stdout() }

func (c *Command) Stderr() io.Writer { return c.app.stderr() }

func (c *Command) Stdin() io.Reader { return c.app.stdin() }

func (c *Command) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(c.app.stdout(), format, a...)
}

func (c *Command) Print(a ...interface{}) (n int, err error) {
	return fmt.Fprint(c.app.stdout(), a...)
}

func (c *Command) Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(c.app.stdout(), a...)
}

func (c *Command) Warnf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(c.app.stderr(), format, a...)
}

func (c *Command) Warn(a ...interface{}) (n int, err error) {
	return fmt.Fprint(c.app.stderr(), a...)
}

func (c *Command) Warnln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(c.app.stderr(), a...)
}

func (c *Command) WrapError(err error) error {
	if err == nil {
		return nil
	}

	return &CommandError{
		Command: c,
		Err:     err,

		exitCode: 1, // TODO
	}
}

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
		if err := c.Action.Setup(c); err != nil {
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

	// Save default values.
	flags := c.Flags()
	for i := range flags {
		flags[i].SaveDefault()
	}

	args := c.Args()
	for i := range args {
		args[i].SaveDefault()
	}

	if rest := c.Rest(); rest != nil {
		rest.SaveDefault()
	}

	c.setuped = true

	return nil
}

type Action interface {
	Setup(cmd *Command) error
	Run(cmd *Command) error
}

// TODO(SuperPaintman): add Before.
// TODO(SuperPaintman): add After.

type ActionRunner func(cmd *Command) error

type ActionBuilder func(cmd *Command) ActionRunner

var _ Action = (ActionRunner)(nil)

func (fn ActionRunner) Setup(cmd *Command) error { return nil }

func (fn ActionRunner) Run(cmd *Command) error { return fn(cmd) }

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

func (a *actionFunc) Setup(cmd *Command) error {
	a.runner = a.builder(cmd)

	return nil
}

func (a *actionFunc) Run(cmd *Command) error {
	if a.runner == nil {
		return nil
	}

	return a.runner(cmd)
}

type CommandFlag struct {
	Short  string
	Long   string
	Usage  Usager
	Action Action

	value bool
}
