package cli

type Arg struct {
	Value     Value
	Name      string
	Usage     string
	Necessary Necessary
}

func newArg(value Value, opts ArgOptions) Arg {
	return Arg{
		Value:     value,
		Name:      opts.Name,
		Usage:     opts.Usage,
		Necessary: opts.Necessary,
	}
}

func (a *Arg) required() bool {
	return a.Necessary == Required || a.Necessary == necessaryUnset
}

func IntArgVar(register Register, p *int, name string, options ...ArgOptionApplyer) error {
	var opts ArgOptions
	opts.applyName(name)
	opts.applyArgOptions(options)

	return register.RegisterArg(newArg(newIntValue(p), opts))
}

func IntArg(register Register, name string, options ...ArgOptionApplyer) *int {
	p := new(int)
	_ = IntArgVar(register, p, name, options...)
	return p
}
