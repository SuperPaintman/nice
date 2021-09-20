package cli

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/SuperPaintman/nice/colors"
)

type Helper interface {
	Help(ctx Context, path []string, w io.Writer) error
}

var _ Helper = (HelperFunc)(nil)

type HelperFunc func(ctx Context, path []string, w io.Writer) error

func (fn HelperFunc) Help(ctx Context, path []string, w io.Writer) error {
	return fn(ctx, path, w)
}

var _ Helper = noopHelper{}

type noopHelper struct{}

func (n noopHelper) Help(ctx Context, path []string, w io.Writer) error {
	return nil
}

func DisableHelp() Helper {
	return noopHelper{}
}

var _ Helper = DefaultHelper{}

type DefaultHelper struct{}

func (h DefaultHelper) Help(ctx Context, path []string, w io.Writer) error {
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
			return fmt.Errorf("cli: command not found: %s", strings.Join(path, " "))
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
	if cmd.Usage != nil {
		// TODO(SuperPaintman): optimize it.
		var buf bytes.Buffer
		if err := cmd.Usage.Usage(ctx, &buf); err != nil {
			return err
		}
		usage := string(buf.Bytes())

		ew.Writef("\n")
		ew.Writef(usage)

		if len(usage) > 0 && usage[len(usage)-1] != '\n' {
			ew.Writef("\n")
		}
	}

	// Commands.
	if len(cmd.Commands) > 0 {
		ew.Writef("\n")
		ew.Writef("Commands:\n")

		var maxLen int
		for _, cmd := range cmd.Commands {
			if len(cmd.Name) > maxLen {
				maxLen = len(cmd.Name)
			}
		}

		for _, cmd := range cmd.Commands {
			// Name.
			ew.Writef("  %s%s%s", colorCommand, cmd.Name, colorCommand.Reset())

			// Usage.
			if cmd.Usage != nil {
				// TODO(SuperPaintman): optimize it.
				var buf bytes.Buffer
				if err := cmd.Usage.Usage(ctx, &buf); err != nil {
					return err
				}
				usage := string(buf.Bytes())

				if usage != "" {
					indent := 4 + maxLen - len(cmd.Name)
					for i := 0; i < indent; i++ {
						ew.WriteString(" ")
					}

					ew.Writef("%s", usage)
				}
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

		var maxLen int
		for _, arg := range args {
			if l := len(arg.Name) + 2 + len(arg.Type()); l > maxLen {
				maxLen = l
			}
		}

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
			if arg.Usage != nil {
				// TODO(SuperPaintman): optimize it.
				var buf bytes.Buffer
				if err := arg.Usage.Usage(ctx, &buf); err != nil {
					return err
				}
				usage := string(buf.Bytes())

				if usage != "" {
					indent := 4 + maxLen - (len(arg.Name) + 2 + len(arg.Type()))
					for i := 0; i < indent; i++ {
						ew.WriteString(" ")
					}

					ew.Writef("%s", usage)
				}
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

		var maxLen int
		for _, flag := range flags {
			l := len(ctx.Parser().FormatShortFlag(flag.Short))

			if flag.Long != "" {
				if l != 0 {
					l += 2
				}

				l += len(ctx.Parser().FormatLongFlag(flag.Long))
			}

			if t := flag.Type(); t != "bool" {
				if t == "" {
					l += len("(unknown)") + 1
				} else {
					l += len(t) + 1
				}
			}

			if l > maxLen {
				maxLen = l
			}
		}

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
			if flag.Usage != nil {
				// TODO(SuperPaintman): optimize it.
				var buf bytes.Buffer
				if err := flag.Usage.Usage(ctx, &buf); err != nil {
					return err
				}
				usage := string(buf.Bytes())

				if usage != "" {
					l := len(ctx.Parser().FormatShortFlag(flag.Short))

					if flag.Long != "" {
						if l != 0 {
							l += 2
						}

						l += len(ctx.Parser().FormatLongFlag(flag.Long))
					}

					if t := flag.Type(); t != "bool" {
						if t == "" {
							l += len("(unknown)") + 1
						} else {
							l += len(t) + 1
						}
					}

					indent := 4 + maxLen - l
					for i := 0; i < indent; i++ {
						ew.WriteString(" ")
					}

					ew.Writef("%s", usage)
				}
			}

			ew.Writef("\n")
		}

		if err := ew.Err(); err != nil {
			return err
		}
	}

	return nil
}
