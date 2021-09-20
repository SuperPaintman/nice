package cli

import (
	"context"
	"fmt"
	"io"
	"time"
)

type Context interface {
	context.Context
	Register

	App() *App
	Command() *Command
	Path() []string
	Help(ctx Context, path []string, w io.Writer) error
	Parser() Parser
	Args() []Arg
	Rest() *RestArgs
	Flags() []Flag
	Stdout() io.Writer
	Stderr() io.Writer
	Stdin() io.Reader
	Printf(format string, a ...interface{}) (n int, err error)
	Print(a ...interface{}) (n int, err error)
	Println(a ...interface{}) (n int, err error)
	Warnf(format string, a ...interface{}) (n int, err error)
	Warn(a ...interface{}) (n int, err error)
	Warnln(a ...interface{}) (n int, err error)
}

var _ Context = (*commandContext)(nil)

type commandContext struct {
	parent context.Context

	app     *App
	command *Command
	path    []string
}

func newCommandContext(ctx context.Context, app *App, command *Command, path []string) *commandContext {
	return &commandContext{
		parent:  ctx,
		app:     app,
		command: command,
		path:    path,
	}
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

func (c *commandContext) RegisterRestArgs(rest RestArgs) error {
	return c.app.parser().RegisterRestArgs(rest)
}

func (c *commandContext) App() *App { return c.app }

func (c *commandContext) Command() *Command { return c.command }

func (c *commandContext) Path() []string { return c.path }

func (c *commandContext) Help(ctx Context, path []string, w io.Writer) error {
	return c.app.help(ctx, path, w)
}

func (c *commandContext) Parser() Parser { return c.app.parser() }

func (c *commandContext) Args() []Arg { return c.app.parser().Args() }

func (c *commandContext) Rest() *RestArgs { return c.app.parser().Rest() }

func (c *commandContext) Flags() []Flag { return c.app.parser().Flags() }

func (c *commandContext) Stdout() io.Writer { return c.app.stdout() }

func (c *commandContext) Stderr() io.Writer { return c.app.stderr() }

func (c *commandContext) Stdin() io.Reader { return c.app.stdin() }

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
