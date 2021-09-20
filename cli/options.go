package cli

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

type RestOptionApplyer interface {
	RestOptionApply(*RestOptions)
}

var _ RestOptionApplyer = (RestOptionFunc)(nil)

type RestOptionFunc func(*RestOptions)

func (fn RestOptionFunc) RestOptionApply(o *RestOptions) {
	fn(o)
}

// Common options.

var (
	_ FlagOptionApplyer = NoopOption{}
	_ ArgOptionApplyer  = NoopOption{}
	_ RestOptionApplyer = NoopOption{}
)

type NoopOption struct{}

func (opt NoopOption) FlagOptionApply(o *FlagOptions) {}
func (opt NoopOption) ArgOptionApply(o *ArgOptions)   {}
func (opt NoopOption) RestOptionApply(o *RestOptions) {}

func WithNoop() NoopOption {
	return NoopOption{}
}

var (
	_ FlagOptionApplyer = Necessary(Optional)
	_ ArgOptionApplyer  = Necessary(Optional)
)

type Necessary uint8

const (
	necessaryUnset Necessary = iota
	Optional
	Required
)

func (opt Necessary) FlagOptionApply(o *FlagOptions) {
	o.Necessary = opt
}

func (opt Necessary) ArgOptionApply(o *ArgOptions) {
	o.Necessary = opt
}

// Usage option.

var (
	_ FlagOptionApplyer = (UsagerFunc)(nil)
	_ ArgOptionApplyer  = (UsagerFunc)(nil)
	_ RestOptionApplyer = (UsagerFunc)(nil)
)

func (fn UsagerFunc) FlagOptionApply(o *FlagOptions) {
	if fn != nil {
		o.Usage = fn
	}
}

func (fn UsagerFunc) ArgOptionApply(o *ArgOptions) {
	if fn != nil {
		o.Usage = fn
	}
}

func (fn UsagerFunc) RestOptionApply(o *RestOptions) {
	if fn != nil {
		o.Usage = fn
	}
}

var (
	_ FlagOptionApplyer = Usage("")
	_ ArgOptionApplyer  = Usage("")
	_ RestOptionApplyer = Usage("")
)

func (s Usage) FlagOptionApply(o *FlagOptions) {
	if s != "" {
		o.Usage = s
	}
}

func (s Usage) ArgOptionApply(o *ArgOptions) {
	if s != "" {
		o.Usage = s
	}
}

func (s Usage) RestOptionApply(o *RestOptions) {
	if s != "" {
		o.Usage = s
	}
}

var (
	_ FlagOptionApplyer = usager{}
	_ ArgOptionApplyer  = usager{}
	_ RestOptionApplyer = usager{}
)

type usager struct{ usager Usager }

func (u usager) FlagOptionApply(o *FlagOptions) {
	if u.usager != nil {
		o.Usage = u.usager
	}
}

func (u usager) ArgOptionApply(o *ArgOptions) {
	if u.usager != nil {
		o.Usage = u.usager
	}
}

func (u usager) RestOptionApply(o *RestOptions) {
	if u.usager != nil {
		o.Usage = u.usager
	}
}

type UsageOption interface {
	FlagOptionApplyer
	ArgOptionApplyer
	RestOptionApplyer
}

func WithUsage(u Usager) UsageOption {
	return usager{u}
}

// Flag options.

var _ FlagOptionApplyer = FlagOptions{}

type FlagOptions struct {
	Value     Value
	Short     string
	Long      string
	Usage     Usager
	Necessary Necessary // Optional if unset

	commandFlag bool

	// Global bool // TODO
}

func (o FlagOptions) FlagOptionApply(opts *FlagOptions) {
	if o.Short != "" {
		opts.Short = o.Short
	}

	if o.Long != "" {
		opts.Long = o.Long
	}

	if o.Usage != nil {
		opts.Usage = o.Usage
	}

	if o.Usage != nil {
		opts.Usage = o.Usage
	}

	opts.Necessary = o.Necessary

	opts.commandFlag = o.commandFlag
}

func (o *FlagOptions) applyName(name string) {
	nameCount := len(name)
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

// var _ FlagOptionApplyer = Global(false)
//
// type Global bool
//
// func (g Global) FlagOptionApply(o *FlagOptions) {
// 	o.Global = bool(g)
// }

type commandFlag bool

func (cf commandFlag) FlagOptionApply(o *FlagOptions) {
	o.commandFlag = bool(cf)
}

// Arg options.

var _ ArgOptionApplyer = ArgOptions{}

type ArgOptions struct {
	Value     Value
	Name      string
	Usage     Usager
	Necessary Necessary // Required if unset
	// NOTE(SuperPaintman):
	//     Usually when we use args in our CLIs they are required by default.
	//     So yes, it's a little bit counfusing (why it isn't Optional?) but
	//     it makes writing CLIs simpler with default options.
}

func (o ArgOptions) ArgOptionApply(opts *ArgOptions) {
	if o.Name != "" {
		opts.Name = o.Name
	}

	if o.Usage != nil {
		opts.Usage = o.Usage
	}

	opts.Necessary = o.Necessary
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

// Rest options.

var _ RestOptionApplyer = RestOptions{}

type RestOptions struct {
	Name  string
	Usage Usager
}

func (o RestOptions) RestOptionApply(opts *RestOptions) {
	if o.Name != "" {
		opts.Name = o.Name
	}

	if o.Usage != nil {
		opts.Usage = o.Usage
	}
}

func (o *RestOptions) applyName(name string) {
	o.Name = name
}

func (o *RestOptions) applyRestOptions(options []RestOptionApplyer) {
	for _, opt := range options {
		if opt != nil {
			opt.RestOptionApply(o)
		}
	}
}
