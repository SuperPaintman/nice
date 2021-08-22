package cli

type Arg struct {
	Value Value
	Name  string
	Usage string
}

func newArg(value Value, opts ArgOptions) Arg {
	return Arg{
		Value: value,
		Name:  opts.Name,
		Usage: opts.Usage,
	}
}

func IntArgVar(parser Parser, p *int, name string, options ...ArgOptionApplyer) error {
	var opts ArgOptions
	opts.applyName(name)
	opts.applyArgOptions(options)

	return parser.RegisterArg(newArg(newIntValue(p), opts))
}

func IntArg(parser Parser, name string, options ...ArgOptionApplyer) *int {
	p := new(int)
	_ = IntArgVar(parser, p, name, options...)
	return p
}
