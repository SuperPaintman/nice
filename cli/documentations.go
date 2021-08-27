package cli

import "io"

type DocumantationGenerator interface {
	DocumantationGenerate(ctx Context, w io.Writer) error
}

type DocumantationGeneratorFunc func(ctx Context, w io.Writer) error

func (fn DocumantationGeneratorFunc) DocumantationGenerate(ctx Context, w io.Writer) error {
	return fn(ctx, w)
}

func ManDocumantation() DocumantationGenerator {
	return DocumantationGeneratorFunc(func(ctx Context, w io.Writer) error {
		ew := easyWriter{w: w}

		// app := ctx.App()
		// name := app.Name

		ew.Writef(`.TH "" "" "" "" ""` + "\n")
		ew.Writef(`.SH "NAME"` + "\n")

		if err := ew.Err(); err != nil {
			return err
		}

		return nil
	})
}

func MarkdownDocumantation() DocumantationGenerator {
	return DocumantationGeneratorFunc(func(ctx Context, w io.Writer) error {
		ew := easyWriter{w: w}

		app := ctx.App()
		name := app.Name
		args := ctx.Args()
		flags := ctx.Flags()

		ew.Writef("# %s\n", name)

		if err := ew.Err(); err != nil {
			return err
		}

		// Description.
		if app.Usage != "" {
			usage := normalizeUsage(app.Usage)

			if len(usage) > 0 {
				ew.Writef("\n")
				ew.Writef("%s.\n", usage)
			}

			if err := ew.Err(); err != nil {
				return err
			}
		}

		// Usage with argumens.
		ew.Writef("\n")
		ew.Writef("## Usage\n")
		ew.Writef("\n")
		ew.Writef("```\n")
		ew.Writef("%s", name)

		if len(flags) > 0 {
			ew.Writef(" [options]")
		}

		ew.Writef("\n")
		ew.Writef("```\n")

		if err := ew.Err(); err != nil {
			return err
		}

		// Arguments.
		if len(args) > 0 {
			ew.Writef("\n")
			ew.Writef("### Arguments\n")

			ew.Writef("\n")
			ew.Writef("| Name | Type | Required | Description |\n")
			ew.Writef("|------|------|----------|-------------|\n")

			for _, arg := range args {
				ew.Writef("|")

				// Name.
				ew.Writef(" **%s** |", arg.Name)

				// Type.
				if t := arg.Type(); t != "" {
					ew.Writef(" `%s` |", t)
				} else {
					ew.Writef(" _(unknown)_ |")
				}

				// Required.
				if arg.Required() {
					ew.Writef(" **yes** |")
				} else {
					ew.Writef(" |")
				}

				// Description.
				if usage := normalizeUsage(arg.Usage); usage != "" {
					ew.Writef(" %s |", usage)
				} else {
					ew.Writef(" |")
				}

				ew.Writef("\n")
			}

			ew.Writef("\n")

			if err := ew.Err(); err != nil {
				return err
			}
		}

		// Options.
		if len(flags) > 0 {
			ew.Writef("\n")
			ew.Writef("### Options\n")

			ew.Writef("\n")
			ew.Writef("| Short | Long | Type | Required | Description |\n")
			ew.Writef("|-------|------|------|----------|-------------|\n")

			for _, flag := range flags {
				ew.Writef("|")

				// Short.
				if flag.Short != "" {
					ew.Writef(" **%s** |", ctx.Parser().FormatShortFlag(flag.Short))
				} else {
					ew.Writef(" |")
				}

				// Long.
				if flag.Long != "" {
					ew.Writef(" **%s** |", ctx.Parser().FormatLongFlag(flag.Long))
				} else {
					ew.Writef(" |")
				}

				// Type.
				if t := flag.Type(); t != "" {
					ew.Writef(" `%s` |", t)
				} else {
					ew.Writef(" _(unknown)_ |")
				}

				// Required.
				if flag.Required() {
					ew.Writef(" **yes** |")
				} else {
					ew.Writef(" |")
				}

				// Description.
				if usage := normalizeUsage(flag.Usage); usage != "" {
					ew.Writef(" %s |", usage)
				} else {
					ew.Writef(" |")
				}

				ew.Writef("\n")
			}

			ew.Writef("\n")

			if err := ew.Err(); err != nil {
				return err
			}
		}

		return nil
	})
}

func normalizeUsage(usage string) string {
	// Leading newlines.
	for {
		if len(usage) > 0 && usage[0] == '\n' {
			usage = usage[1:]
		} else {
			break
		}
	}

	// Ending newlines.
	for {
		if len(usage) > 0 && usage[len(usage)-1] == '\n' {
			usage = usage[len(usage)-1:]
		} else {
			break
		}
	}

	if len(usage) > 0 && usage[len(usage)-1] == '.' {
		usage = usage[len(usage)-1:]
	}

	return usage
}
