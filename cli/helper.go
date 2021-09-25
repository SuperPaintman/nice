package cli

import (
	"bytes"
	"io"

	"github.com/SuperPaintman/nice/colors"
)

type Helper interface {
	Help(cmd *Command, w io.Writer) error
}

var _ Helper = (HelperFunc)(nil)

type HelperFunc func(cmd *Command, w io.Writer) error

func (fn HelperFunc) Help(cmd *Command, w io.Writer) error {
	return fn(cmd, w)
}

var _ Helper = noopHelper{}

type noopHelper struct{}

func (n noopHelper) Help(cmd *Command, w io.Writer) error {
	return nil
}

func DisableHelp() Helper {
	return noopHelper{}
}

var _ Helper = DefaultHelper{}

type DefaultHelper struct{}

func (h DefaultHelper) Help(cmd *Command, w io.Writer) error {
	const (
		colorName     = colors.Blue
		colorCommand  = colors.Magenta
		colorArgument = colors.Magenta
		colorOption   = colors.Yellow
		colorType     = colors.Green
		colorDefault  = colors.Blue
	)

	ew := easyWriter{w: w}

	path := cmd.Path()
	args := cmd.Args()
	rest := cmd.Rest()
	flags := cmd.Flags()

	// Usage with argumens.
	ew.Writef("Usage:")

	for _, name := range path {
		ew.Writef(" %s%s%s", colorName, name, colorName.Reset())
	}

	if len(flags) > 0 {
		ew.Writef(" %s[options...]%s", colorOption, colorOption.Reset())
	}

	for _, arg := range args {
		if arg.Required() {
			ew.Writef(" %s<%s>%s", colorArgument, arg.Name, colorArgument.Reset())
		} else {
			ew.Writef(" %s[%s]%s", colorArgument, arg.Name, colorArgument.Reset())
		}
	}

	if rest != nil {
		ew.Writef(" %s[%s...]%s", colorArgument, rest.Name, colorArgument.Reset())
	}

	ew.Writef("\n")

	if err := ew.Err(); err != nil {
		return err
	}

	// Usage with a command.
	if len(cmd.Commands) > 0 {
		ew.Writef("      ")

		for _, name := range path {
			ew.Writef(" %s%s%s", colorName, name, colorName.Reset())
		}

		if len(flags) > 0 {
			ew.Writef(" %s[options...]%s", colorOption, colorOption.Reset())
		}

		ew.Writef(" %s[command]%s", colorCommand, colorCommand.Reset())

		ew.Writef("\n")
	}

	// Description from Usage field.
	if cmd.Usage != nil {
		// TODO(SuperPaintman): optimize it.
		var buf bytes.Buffer
		if err := cmd.Usage.Usage(cmd, &buf); err != nil {
			return err
		}
		usage := buf.String()

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

		for i := range cmd.Commands {
			cmd := &cmd.Commands[i]

			// Name.
			ew.Writef("  %s%s%s", colorCommand, cmd.Name, colorCommand.Reset())

			// Usage.
			if cmd.Usage != nil {
				// TODO(SuperPaintman): optimize it.
				var buf bytes.Buffer
				if err := cmd.Usage.Usage(cmd, &buf); err != nil {
					return err
				}
				usage := buf.String()

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
	var argMaxLen int
	if len(args) > 0 {
		ew.Writef("\n")
		ew.Writef("Arguments:\n")

		for _, arg := range args {
			if l := len(arg.Name) + 2 + len(arg.Type()) + 1; l > argMaxLen {
				argMaxLen = l
			}
		}

		if rest != nil {
			if l := len(rest.Name) + 2 + len(rest.Type()) + 1 + 3; l > argMaxLen {
				argMaxLen = l
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
			var hasUsage bool
			if arg.Usage != nil {
				// TODO(SuperPaintman): optimize it.
				var buf bytes.Buffer
				if err := arg.Usage.Usage(cmd, &buf); err != nil {
					return err
				}
				usage := buf.String()

				if usage != "" {
					indent := 4 + argMaxLen - (len(arg.Name) + 2 + len(arg.Type()) + 1)
					for i := 0; i < indent; i++ {
						ew.WriteString(" ")
					}

					ew.Writef("%s", usage)

					hasUsage = true
				}
			}

			// Default.
			if value, empty := arg.Default(); !empty {
				if !hasUsage {
					indent := 4 + argMaxLen - (len(arg.Name) + 2 + len(arg.Type()) + 1)
					for i := 0; i < indent; i++ {
						ew.WriteString(" ")
					}
				} else {
					ew.WriteString(" ")
				}

				ew.Writef("(default: %s%v%s)", colorDefault, value, colorDefault.Reset())
			}

			ew.Writef("\n")
		}

		if err := ew.Err(); err != nil {
			return err
		}
	}

	// Rest.
	if rest != nil {
		if len(args) == 0 {
			ew.Writef("\n")
			ew.Writef("Arguments:\n")
		}

		if l := len(rest.Name) + 2 + len(rest.Type()) + 1 + 3; l > argMaxLen {
			argMaxLen = l
		}

		// Name.
		ew.Writef("  %s[%s...]%s", colorArgument, rest.Name, colorArgument.Reset())

		// Type.
		if t := rest.Type(); t != "bool" {
			if t == "" {
				t = "(unknown)"
			}

			ew.Writef(" %s%s%s", colorType, t, colorType.Reset())
		}

		// Usage.
		var hasUsage bool
		if rest.Usage != nil {
			// TODO(SuperPaintman): optimize it.
			var buf bytes.Buffer
			if err := rest.Usage.Usage(cmd, &buf); err != nil {
				return err
			}
			usage := buf.String()

			if usage != "" {
				indent := 4 + argMaxLen - ((len(rest.Name) + 2 + len(rest.Type()) + 1) + 3)
				for i := 0; i < indent; i++ {
					ew.WriteString(" ")
				}

				ew.Writef("%s", usage)
			}
		}

		// Default.
		if value, empty := rest.Default(); !empty {
			if !hasUsage {
				indent := 4 + argMaxLen - ((len(rest.Name) + 2 + len(rest.Type())) + 1 + 3)
				for i := 0; i < indent; i++ {
					ew.WriteString(" ")
				}
			} else {
				ew.WriteString(" ")
			}

			ew.Writef("(default: %s%v%s)", colorDefault, value, colorDefault.Reset())
		}

		ew.Writef("\n")

		if err := ew.Err(); err != nil {
			return err
		}
	}

	// Options.
	if len(flags) > 0 {
		ew.Writef("\n")
		ew.Writef("Options:\n")

		var (
			maxLenShort int
			maxLen      int
		)
		for _, flag := range flags {
			var l int
			if flag.Short != "" {
				shortLen := len(cmd.Parser().FormatShortFlag(flag.Short))

				if shortLen > maxLenShort {
					maxLenShort = shortLen
				}

				l += shortLen
			}

			if flag.Long != "" {
				if l != 0 {
					l += 2
				}

				l += len(cmd.Parser().FormatLongFlag(flag.Long))
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
					cmd.Parser().FormatShortFlag(flag.Short),
					colorOption.Reset(),
				)
				ew.Writef(", ")
			} else {
				if maxLenShort > 0 {
					for i := 0; i < maxLenShort; i++ {
						ew.WriteString(" ")
					}

					ew.Writef("  ")
				}
			}

			// Long.
			if flag.Long != "" {
				ew.Writef("%s%s%s",
					colorOption,
					cmd.Parser().FormatLongFlag(flag.Long),
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
			var hasUsage bool
			if flag.Usage != nil {
				// TODO(SuperPaintman): optimize it.
				var buf bytes.Buffer
				if err := flag.Usage.Usage(cmd, &buf); err != nil {
					return err
				}
				usage := buf.String()

				if usage != "" {
					l := len(cmd.Parser().FormatShortFlag(flag.Short))
					if l == 0 {
						l += maxLenShort
					}

					if flag.Long != "" {
						if l != 0 {
							l += 2
						}

						l += len(cmd.Parser().FormatLongFlag(flag.Long))
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

					hasUsage = true
				}
			}

			// Default.
			if value, empty := flag.Default(); !empty {
				if !hasUsage {
					l := len(cmd.Parser().FormatShortFlag(flag.Short))
					if l == 0 {
						l += maxLenShort
					}

					if flag.Long != "" {
						if l != 0 {
							l += 2
						}

						l += len(cmd.Parser().FormatLongFlag(flag.Long))
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
				} else {
					ew.WriteString(" ")
				}

				if flag.Required() {
					ew.Writef("(required, default: %s%v%s)", colorDefault, value, colorDefault.Reset())
				} else {
					ew.Writef("(default: %s%v%s)", colorDefault, value, colorDefault.Reset())
				}
			} else if flag.Required() {
				if !hasUsage {
					l := len(cmd.Parser().FormatShortFlag(flag.Short))
					if l == 0 {
						l += maxLenShort
					}

					if flag.Long != "" {
						if l != 0 {
							l += 2
						}

						l += len(cmd.Parser().FormatLongFlag(flag.Long))
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
				} else {
					ew.WriteString(" ")
				}

				ew.Writef("(required)")
			}

			ew.Writef("\n")
		}

		if err := ew.Err(); err != nil {
			return err
		}
	}

	return nil
}
