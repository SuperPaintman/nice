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

func (a *Arg) Type() string {
	if t, ok := a.Value.(Typer); ok {
		return t.Type()
	}

	return ""
}

func (a *Arg) Required() bool {
	// By default args are required.
	return a.Necessary != Optional
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

func StringArgVar(parser Register, p *string, name string, options ...ArgOptionApplyer) error {
	var opts ArgOptions
	opts.applyName(name)
	opts.applyArgOptions(options)

	return parser.RegisterArg(newArg(newStringValue(p), opts))
}

func StringArg(register Register, name string, options ...ArgOptionApplyer) *string {
	p := new(string)
	_ = StringArgVar(register, p, name, options...)
	return p
}
