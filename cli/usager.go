package cli

import "io"

type Usager interface {
	Usage(ctx Context, w io.Writer) error
}

var _ Usager = (UsagerFunc)(nil)

type UsagerFunc func(ctx Context, w io.Writer) error

func (fn UsagerFunc) Usage(ctx Context, w io.Writer) error {
	return fn(ctx, w)
}

var _ Usager = Usage("")

type Usage string

func (s Usage) Usage(ctx Context, w io.Writer) error {
	_, err := w.Write([]byte(s))

	return err
}
