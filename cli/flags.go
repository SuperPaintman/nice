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

func BoolVar(parser Parser, p *bool, name string, options ...FlagOptionApplyer) error {
	var opts FlagOptions
	opts.applyName(name)
	opts.applyFlagOptions(options)

	return parser.RegisterFlag(newFlag(newBoolValue(p), opts))
}

func Bool(parser Parser, name string, options ...FlagOptionApplyer) *bool {
	p := new(bool)
	_ = BoolVar(parser, p, name, options...)
	return p
}

func IntVar(parser Parser, p *int, name string, options ...FlagOptionApplyer) error {
	var opts FlagOptions
	opts.applyName(name)
	opts.applyFlagOptions(options)

	return parser.RegisterFlag(newFlag(newIntValue(p), opts))
}

func Int(parser Parser, name string, options ...FlagOptionApplyer) *int {
	p := new(int)
	_ = IntVar(parser, p, name, options...)
	return p
}
