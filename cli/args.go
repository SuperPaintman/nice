package cli

type Arg struct {
	Value     Value
	Name      string
	Usage     Usager
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

func (a *Arg) String() string {
	return "Arg(" + a.Type() + "," + a.Name + ")"
}

//go:generate python ./generate_args.py
