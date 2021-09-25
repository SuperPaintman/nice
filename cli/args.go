package cli

type Arg struct {
	Value     Value
	Name      string
	Usage     Usager
	Necessary Necessary

	set          bool
	defaultSaved bool
	defaultValue string
	defaultEmpty bool
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

func (a *Arg) Set() bool {
	return a.set
}

func (a *Arg) MarkSet() {
	a.set = true
}

func (a *Arg) Default() (v string, empty bool) {
	if !a.defaultSaved {
		return "", true
	}

	return a.defaultValue, a.defaultEmpty
}

func (a *Arg) SaveDefault() {
	a.defaultValue = a.Value.String()
	if ev, ok := a.Value.(Emptier); ok {
		a.defaultEmpty = ev.Empty()
	} else {
		a.defaultEmpty = a.defaultValue == ""
	}

	a.defaultSaved = true
}

func (a *Arg) String() string {
	return "Arg(" + a.Type() + "," + a.Name + ")"
}

func ArgVar(register Register, value Value, name string, options ...ArgOptionApplyer) error {
	var opts ArgOptions
	opts.applyName(name)
	opts.applyArgOptions(options)

	return register.RegisterArg(newArg(value, opts))
}

//go:generate python ./generate_args.py
