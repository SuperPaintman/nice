package cli

type Flag struct {
	Value     Value
	Short     string
	Long      string
	Usage     Usager
	Necessary Necessary

	set          bool
	defaultSaved bool
	defaultValue string
	defaultEmpty bool
	commandFlag  bool

	// NOTE(SuperPaintman):
	//     The first version had "Aliases" for flags. It's quite handy to have
	//     (e.g. --dry and --dry-run) but at the same time makes API a bit
	//     confusing because of duplication logic.
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

		commandFlag: opts.commandFlag,
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

func (f *Flag) Set() bool {
	return f.set
}

func (f *Flag) MarkSet() {
	f.set = true
}

func (f *Flag) Default() (v string, empty bool) {
	if !f.defaultSaved {
		return "", true
	}

	return f.defaultValue, f.defaultEmpty
}

func (f *Flag) SaveDefault() {
	f.defaultValue = f.Value.String()
	if ev, ok := f.Value.(Emptier); ok {
		f.defaultEmpty = ev.Empty()
	} else {
		f.defaultEmpty = f.defaultValue == ""
	}

	f.defaultSaved = true
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

func Var(register Register, value Value, name string, options ...FlagOptionApplyer) error {
	var opts FlagOptions
	opts.applyName(name)
	opts.applyFlagOptions(options)

	return register.RegisterFlag(newFlag(value, opts))
}

//go:generate python ./generate_flags.py

//go:generate python ./generate_multi_flags.py
