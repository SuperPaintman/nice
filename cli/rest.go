package cli

type RestArgs struct {
	Values *[]string
	Name   string
	Usage  Usager
}

func newRest(values *[]string, opts RestOptions) RestArgs {
	return RestArgs{
		Values: values,
		Name:   opts.Name,
		Usage:  opts.Usage,
	}
}

func (ra *RestArgs) IsZero() bool {
	return ra.Values == nil
}

func (ra *RestArgs) Add(arg string) {
	if ra.Values != nil {
		*ra.Values = append(*ra.Values, arg)
	}
}

func RestVar(register Register, p *[]string, name string, options ...RestOptionApplyer) error {
	var opts RestOptions
	opts.applyName(name)
	opts.applyRestOptions(options)

	return register.RegisterRestArgs(newRest(p, opts))
}

func Rest(register Register, name string, options ...RestOptionApplyer) *[]string {
	p := new([]string)
	_ = RestVar(register, p, name, options...)
	return p
}
