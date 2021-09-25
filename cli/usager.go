package cli

import "io"

type Usager interface {
	Usage(app *App, cmd *Command, w io.Writer) error
}

var _ Usager = (UsagerFunc)(nil)

type UsagerFunc func(app *App, cmd *Command, w io.Writer) error

func (fn UsagerFunc) Usage(app *App, cmd *Command, w io.Writer) error {
	return fn(app, cmd, w)
}

var _ Usager = Usage("")

type Usage string

func (s Usage) Usage(app *App, cmd *Command, w io.Writer) error {
	_, err := w.Write([]byte(s))

	return err
}
