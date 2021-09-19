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
		ra.Values.Set(val)
	}

	return nil
}

//go:generate python ./generate_rest.py
