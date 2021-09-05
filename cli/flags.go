package cli

type Flag struct {
	Value     Value
	Short     string
	Long      string
	Usage     Usager
	Necessary Necessary

	// NOTE(SuperPaintman):
	//     The first version had "Aliases" for flags. It's quite handy to have
	//     (e.g. --dry and --dry-run) but at the same time makes API a bit
	//     confusing because of duplication logic.
	//
	//     We can override flags on aliases collision or remove the alias from the
	//     original list but it makes API quite unpredictable for developers.
	//
	//     I decided to remove aliases. It's not so commonly used feature and
	//     developers can easely make a workaround if they need it.
}

func newFlag(value Value, opts FlagOptions) Flag {
	return Flag{
		Value:     value,
		Short:     opts.Short,
		Long:      opts.Long,
		Usage:     opts.Usage,
		Necessary: opts.Necessary,
	}
}

func (f *Flag) Type() string {
	if t, ok := f.Value.(Typer); ok {
		return t.Type()
	}

	return ""
}

func (f *Flag) Required() bool {
	return f.Necessary == Required
}

func (f *Flag) String() string {
	v := "Flag("
	v += f.Type()

	if f.Short != "" {
		v += ","
		v += "-" + f.Short
	}

	if f.Long != "" {
		if v == "" {
			v += ","
		} else {
			v += "/"
		}

		v += "--" + f.Long
	}
	v += ")"

	return v
}

func BoolVar(register Register, p *bool, name string, options ...FlagOptionApplyer) error {
	var opts FlagOptions
	opts.applyName(name)
	opts.applyFlagOptions(options)

	return register.RegisterFlag(newFlag(newBoolValue(p), opts))
}

func Bool(register Register, name string, options ...FlagOptionApplyer) *bool {
	p := new(bool)
	_ = BoolVar(register, p, name, options...)
	return p
}

func IntVar(parser Register, p *int, name string, options ...FlagOptionApplyer) error {
	var opts FlagOptions
	opts.applyName(name)
	opts.applyFlagOptions(options)

	return parser.RegisterFlag(newFlag(newIntValue(p), opts))
}

func Int(register Register, name string, options ...FlagOptionApplyer) *int {
	p := new(int)
	_ = IntVar(register, p, name, options...)
	return p
}

func StringVar(parser Register, p *string, name string, options ...FlagOptionApplyer) error {
	var opts FlagOptions
	opts.applyName(name)
	opts.applyFlagOptions(options)

	return parser.RegisterFlag(newFlag(newStringValue(p), opts))
}

func String(register Register, name string, options ...FlagOptionApplyer) *string {
	p := new(string)
	_ = StringVar(register, p, name, options...)
	return p
}
