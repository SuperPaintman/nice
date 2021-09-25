package cli

import "io"

type Usager interface {
	Usage(cmd *Command, w io.Writer) error
}

var _ Usager = (UsagerFunc)(nil)

type UsagerFunc func(cmd *Command, w io.Writer) error

func (fn UsagerFunc) Usage(cmd *Command, w io.Writer) error {
	return fn(cmd, w)
}

var _ Usager = Usage("")

type Usage string

func (s Usage) Usage(cmd *Command, w io.Writer) error {
	_, err := w.Write([]byte(s))

	return err
}
