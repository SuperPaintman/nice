package cli

import (
	"io"

	"github.com/SuperPaintman/nice/colors"
)

type Helper interface {
	Help(ctx Context, w io.Writer) error
}

var _ Helper = (HelperFunc)(nil)

type HelperFunc func(ctx Context, w io.Writer) error

func (fn HelperFunc) Help(ctx Context, w io.Writer) error {
	return fn(ctx, w)
}

var _ Helper = noopHelper{}

type noopHelper struct{}

func (n noopHelper) Help(ctx Context, w io.Writer) error {
	return nil
}

func DisableHelp() Helper {
	return noopHelper{}
}

var _ Helper = DefaultHelper{}

type DefaultHelper struct{}

func (h DefaultHelper) Help(ctx Context, w io.Writer) error {
	const (
		colorName     = colors.Blue
		colorCommand  = colors.Magenta
		colorArgument = colors.Magenta
		colorOption   = colors.Yellow
		colorType     = colors.Green
	)

	ew := easyWriter{w: w}

	args := ctx.Args()
	flags := ctx.Flags()
	cmd := ctx.App().Command()
	path := ctx.Path()
	for _, name := range path[1:] {
		// Find a sub command.
		var found bool
		for i := range cmd.Commands {
			c := &cmd.Commands[i]

			if c.Name == name {
				found = true
				cmd = c
				break
			}
		}

		if !found {
			// return fmt.Errorf("cli: command not found: %s", name)
			break
		}
	}

	// Usage with a command.
	if len(cmd.Commands) > 0 {
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
		if len(cmd.Commands) == 0 {
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
			if arg.Required() {
				ew.Writef(" %s<%s>%s", colorArgument, arg.Name, colorArgument.Reset())
			} else {
				ew.Writef(" %s[%s]%s", colorArgument, arg.Name, colorArgument.Reset())
			}
		}

		ew.Writef("\n")

		if err := ew.Err(); err != nil {
			return err
		}
	} else if len(cmd.Commands) == 0 {
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

	// Description from Usage field.
	if cmd.Usage != "" {
		ew.Writef("\n")
		ew.Writef(cmd.Usage)

		if len(cmd.Usage) > 0 && cmd.Usage[len(cmd.Usage)-1] != '\n' {
			ew.Writef("\n")
		}
	}

	// Commands.
	if len(cmd.Commands) > 0 {
		ew.Writef("\n")
		ew.Writef("Commands:\n")

		for _, cmd := range cmd.Commands {
			// Name.
			ew.Writef("  %s%s%s", colorCommand, cmd.Name, colorCommand.Reset())

			// Usage.
			if cmd.Usage != "" {
				ew.Writef("\t\t%s", cmd.Usage)
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
			if arg.Required() {
				ew.Writef("  %s<%s>%s", colorArgument, arg.Name, colorArgument.Reset())
			} else {
				ew.Writef("  %s[%s]%s", colorArgument, arg.Name, colorArgument.Reset())
			}

			// Type.
			if t := arg.Type(); t != "bool" {
				if t == "" {
					t = "(unknown)"
				}

				ew.Writef(" %s%s%s", colorType, t, colorType.Reset())
			}

			// Usage.
			if arg.Usage != "" {
				ew.Writef("\t\t%s", arg.Usage)
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

			// Type.
			if t := flag.Type(); t != "bool" {
				if t == "" {
					t = "(unknown)"
				}

				ew.Writef(" %s%s%s", colorType, t, colorType.Reset())
			}

			// Usage.
			if flag.Usage != "" {
				ew.Writef("\t\t%s", flag.Usage)
			}

			ew.Writef("\n")
		}

		if err := ew.Err(); err != nil {
			return err
		}
	}

	return nil
}
