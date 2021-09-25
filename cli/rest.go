package cli

type RestArgs struct {
	Values Value
	Name   string
	Usage  Usager

	defaultSaved bool
	defaultValue string
	defaultEmpty bool
}

func newRest(values Value, opts RestOptions) RestArgs {
	return RestArgs{
		Values: values,
		Name:   opts.Name,
		Usage:  opts.Usage,
	}
}

func (ra *RestArgs) Type() string {
	if ra.Values != nil {
		if t, ok := ra.Values.(Typer); ok {
			return t.Type()
		}
	}

	return ""
}

func (ra *RestArgs) IsZero() bool {
	return ra.Values == nil
}

func (ra *RestArgs) Add(val string) error {
	if ra.Values != nil {
		if err := ra.Values.Set(val); err != nil {
			return err
		}
	}

	return nil
}

func (ra *RestArgs) Default() (v string, empty bool) {
	if !ra.defaultSaved {
		return "", true
	}

	return ra.defaultValue, ra.defaultEmpty
}

func (ra *RestArgs) SaveDefault() {
	if ra.Values != nil {
		ra.defaultValue = ra.Values.String()
		if ev, ok := ra.Values.(Emptier); ok {
			ra.defaultEmpty = ev.Empty()
		} else {
			ra.defaultEmpty = ra.defaultValue == ""
		}
	}

	ra.defaultSaved = true
}

func RestVar(register Register, value Value, name string, options ...RestOptionApplyer) error {
	var opts RestOptions
	opts.applyName(name)
	opts.applyRestOptions(options)

	return register.RegisterRestArgs(newRest(value, opts))
}

//go:generate python ./generate_rest.py
