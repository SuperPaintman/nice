package cli

import (
	"fmt"
	"io"

	"github.com/SuperPaintman/nice/colors"
)

type Helper interface {
	Help(ctx Context, w io.Writer, path []string) error
}

var _ Helper = (HelperFunc)(nil)

type HelperFunc func(ctx Context, w io.Writer, path []string) error

func (fn HelperFunc) Help(ctx Context, w io.Writer, path []string) error {
	return fn(ctx, w, path)
}

var _ Helper = noopHelper{}

type noopHelper struct{}

func (n noopHelper) Help(ctx Context, w io.Writer, path []string) error {
	return nil
}

func DisableHelp() Helper {
	return noopHelper{}
}

var _ Helper = DefaultHelper{}

type DefaultHelper struct{}

func (h DefaultHelper) Help(ctx Context, w io.Writer, path []string) error {
	const (
		colorName     = colors.Blue
		colorCommand  = colors.Magenta
		colorArgument = colors.Magenta
		colorOption   = colors.Yellow
	)

	ew := easyWriter{w: w}

	args := ctx.Args()
	flags := ctx.Flags()
	commands := ctx.App().Commands
	for _, name := range path[1:] {
		// Find a sub command.
		var found bool
		for i := range commands {
			c := &commands[i]

			if c.Name == name {
				found = true
				commands = c.Commands
				break
			}
		}

		if !found {
			// return fmt.Errorf("cli: command not found: %s", name)
			break
		}
	}

	// Usage with a command.
	if len(commands) > 0 {
		ew.Writef("Usage:")

		for _, name := range path {
			ew.Writef(" %s%s%s", colorName, name, colorName.Reset())
		}

		if len(flags) > 0 {
			ew.Writef(" %s[options]%s", colorOption, colorOption.Reset())
		}

		ew.Writef(" %s[command]%s", colorCommand, colorCommand.Reset())

		ew.Writef("\n")
	}

	// Usage with argumens.
	if len(args) > 0 {
		if len(commands) == 0 {
			ew.Writef("Usage:")
		} else {
			ew.Writef("      ")
		}
		for _, name := range path {
			ew.Writef(" %s%s%s", colorName, name, colorName.Reset())
		}

		if len(flags) > 0 {
			ew.Writef(" %s[options]%s", colorOption, colorOption.Reset())
		}

		for _, arg := range args {
			if arg.required() {
				ew.Writef(" %s<%s>%s", colorArgument, arg.Name, colorArgument.Reset())
			} else {
				ew.Writef(" %s[%s]%s", colorArgument, arg.Name, colorArgument.Reset())
			}
		}

		ew.Writef("\n")

		if err := ew.Err(); err != nil {
			return err
		}
	} else if len(commands) == 0 {
		ew.Writef("Usage:")

		for _, name := range path {
			ew.Writef(" %s%s%s", colorName, name, colorName.Reset())
		}

		if len(flags) > 0 {
			ew.Writef(" %s[options]%s", colorOption, colorOption.Reset())
		}

		ew.Writef("\n")

		if err := ew.Err(); err != nil {
			return err
		}
	}

	// Commands.
	if len(commands) > 0 {
		ew.Writef("\n")
		ew.Writef("Commands:\n")

		for _, cmd := range commands {
			// Name.
			ew.Writef("  %s%s%s", colorCommand, cmd.Name, colorCommand.Reset())

			// Usage.
			if cmd.Usage != "" {
				ew.Writef("   %s", cmd.Usage)
			}

			ew.Writef("\n")
		}

		if err := ew.Err(); err != nil {
			return err
		}
	}

	// Arguments.
	if len(args) > 0 {
		ew.Writef("\n")
		ew.Writef("Arguments:\n")

		for _, arg := range args {
			// Name.
			ew.Writef("  %s%s%s", colorArgument, arg.Name, colorArgument.Reset())

			// Usage.
			if arg.Usage != "" {
				ew.Writef("   %s", arg.Usage)
			}

			ew.Writef("\n")
		}

		if err := ew.Err(); err != nil {
			return err
		}
	}

	// Options.
	if len(flags) > 0 {
		ew.Writef("\n")
		ew.Writef("Options:\n")

		for _, flag := range flags {
			ew.Writef("  ")

			// Short.
			if flag.Short != "" {
				ew.Writef("%s%s%s",
					colorOption,
					ctx.Parser().FormatShortFlag(flag.Short),
					colorOption.Reset(),
				)
				ew.Writef(", ")
			} else {
				ew.Writef("    ")
			}

			// Long.
			if flag.Long != "" {
				ew.Writef("%s%s%s",
					colorOption,
					ctx.Parser().FormatLongFlag(flag.Long),
					colorOption.Reset(),
				)
			} else {
				// TODO
			}

			// Usage.
			if flag.Usage != "" {
				ew.Writef("   %s", flag.Usage)
			}

			ew.Writef("\n")
		}

		if err := ew.Err(); err != nil {
			return err
		}
	}

	return nil
}

type easyWriter struct {
	w   io.Writer
	err error
}

func (ew *easyWriter) Write(data []byte) {
	if ew.err != nil {
		return
	}

	_, err := ew.w.Write(data)
	if err != nil {
		ew.err = err
	}
}

func (ew *easyWriter) WriteString(s string) {
	ew.Write([]byte(s))
}

func (ew *easyWriter) Writef(format string, a ...interface{}) {
	if ew.err != nil {
		return
	}

	var err error
	if len(a) == 0 {
		_, err = ew.w.Write([]byte(format))
	} else {
		_, err = fmt.Fprintf(ew.w, format, a...)
	}

	if err != nil {
		ew.err = err
	}
}

func (ew *easyWriter) Err() error {
	return ew.err
}
