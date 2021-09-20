package cli

type RestArgs struct {
	Values Value
	Name   string
	Usage  Usager
}

func newRest(values Value, opts RestOptions) RestArgs {
	return RestArgs{
		Values: values,
		Name:   opts.Name,
		Usage:  opts.Usage,
	}
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

func RestVar(register Register, value Value, name string, options ...RestOptionApplyer) error {
	var opts RestOptions
	opts.applyName(name)
	opts.applyRestOptions(options)

	return register.RegisterRestArgs(newRest(value, opts))
}

//go:generate python ./generate_rest.py
