package cli

type Alias struct {
	Short string
	Long  string
}

type Flag struct {
	Value   Value
	Short   string
	Long    string
	Aliases []Alias
	Usage   string
}

func newFlag(value Value, opts FlagOptions) Flag {
	return Flag{
		Value:   value,
		Short:   opts.Short,
		Long:    opts.Long,
		Aliases: opts.Aliases,
		Usage:   opts.Usage,
	}
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
