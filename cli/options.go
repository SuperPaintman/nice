package cli

import "unicode/utf8"

// Option interfaces.

type FlagOptionApplyer interface {
	FlagOptionApply(*FlagOptions)
}

var _ FlagOptionApplyer = (FlagOptionFunc)(nil)

type FlagOptionFunc func(*FlagOptions)

func (fn FlagOptionFunc) FlagOptionApply(o *FlagOptions) {
	fn(o)
}

type ArgOptionApplyer interface {
	ArgOptionApply(*ArgOptions)
}

var _ ArgOptionApplyer = (ArgOptionFunc)(nil)

type ArgOptionFunc func(*ArgOptions)

func (fn ArgOptionFunc) ArgOptionApply(o *ArgOptions) {
	fn(o)
}

// Common options.

var (
	_ FlagOptionApplyer = NoopOption{}
	_ ArgOptionApplyer  = NoopOption{}
)

type NoopOption struct{}

func (opt NoopOption) FlagOptionApply(o *FlagOptions) {}
func (opt NoopOption) ArgOptionApply(o *ArgOptions)   {}

func WithNoop() NoopOption {
	return NoopOption{}
}

var (
	_ FlagOptionApplyer = UsageOption("")
	_ ArgOptionApplyer  = UsageOption("")
)

type UsageOption string

func (opt UsageOption) FlagOptionApply(o *FlagOptions) {
	if opt != "" {
		o.Usage = string(opt)
	}
}

func (opt UsageOption) ArgOptionApply(o *ArgOptions) {
	if opt != "" {
		o.Usage = string(opt)
	}
}

func WithUsage(usage string) UsageOption {
	return UsageOption(usage)
}

// Flag options.

var _ FlagOptionApplyer = FlagOptions{}

type FlagOptions struct {
	Value   Value
	Short   string
	Long    string
	Aliases []Alias
	Usage   string
}

func (o FlagOptions) FlagOptionApply(opts *FlagOptions) {
	if o.Short != "" {
		opts.Short = o.Short
	}

	if o.Long != "" {
		opts.Long = o.Long
	}

	if o.Usage != "" {
		opts.Usage = o.Usage
	}

	for _, alias := range o.Aliases {
		opts.Aliases = append(opts.Aliases, alias)
	}
}

func (o *FlagOptions) applyName(name string) {
	nameCount := utf8.RuneCountInString(name)
	if nameCount > 1 {
		o.Long = name
	} else if nameCount == 1 {
		o.Short = name
	}
}

func (o *FlagOptions) applyFlagOptions(options []FlagOptionApplyer) {
	for _, opt := range options {
		if opt != nil {
			opt.FlagOptionApply(o)
		}
	}
}

// Arg options.

var _ ArgOptionApplyer = ArgOptions{}

type ArgOptions struct {
	Value Value
	Name  string
	Usage string
}

func (o ArgOptions) ArgOptionApply(opts *ArgOptions) {
	if o.Name != "" {
		opts.Name = o.Name
	}

	if o.Usage != "" {
		opts.Usage = o.Usage
	}
}

func (o *ArgOptions) applyName(name string) {
	o.Name = name
}

func (o *ArgOptions) applyArgOptions(options []ArgOptionApplyer) {
	for _, opt := range options {
		if opt != nil {
			opt.ArgOptionApply(o)
		}
	}
}

func WithShort(name string) FlagOptionFunc {
	return func(o *FlagOptions) {
		o.Short = name
	}
}

func WithLong(name string) FlagOptionFunc {
	return func(o *FlagOptions) {
		o.Long = name
	}
}

func WithAliases(aliases ...Alias) FlagOptionFunc {
	return func(o *FlagOptions) {
		o.Aliases = aliases
	}
}
