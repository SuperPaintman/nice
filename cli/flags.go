package cli

import (
	"flag"
	"strconv"
	"unicode/utf8"
)

var _ flag.Value = (Value)(nil)

type Value interface {
	String() string
	Set(string) error
}

var _ flag.Value = (Getter)(nil)

type Getter interface {
	Value
	Get() interface{}
}

type boolValue bool

func newBoolValue(p *bool) *boolValue {
	return (*boolValue)(p)
}

func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		// err = errParse
	}
	*b = boolValue(v)
	return err
}

func (b *boolValue) Get() interface{} { return bool(*b) }

func (b *boolValue) String() string { return strconv.FormatBool(bool(*b)) }

func (b *boolValue) IsBoolFlag() bool { return true }

type intValue int

func newIntValue(p *int) *intValue {
	return (*intValue)(p)
}

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		// err = numError(err)
	}
	*i = intValue(v)
	return err
}

func (i *intValue) Get() interface{} { return int(*i) }

func (i *intValue) String() string { return strconv.Itoa(int(*i)) }

type Alias struct {
	Short string
	Long  string
}

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

type FlagOptionApplyer interface {
	FlagOptionApply(*FlagOptions)
}

type FlagOptionFunc func(*FlagOptions)

func (fn FlagOptionFunc) FlagOptionApply(o *FlagOptions) {
	fn(o)
}

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

type ArgOptionApplyer interface {
	ArgOptionApply(*ArgOptions)
}

func WithNoop() FlagOptionFunc {
	return nil
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

type UsageOption string

func (opt UsageOption) FlagOptionApply(o *FlagOptions) {
	if opt == "" {
		return
	}

	o.Usage = string(opt)
}

func (opt UsageOption) ArgOptionApply(o *ArgOptions) {
	if opt == "" {
		return
	}

	o.Usage = string(opt)
}

func WithUsage(usage string) UsageOption {
	return UsageOption(usage)
}

func BoolVar(parser Parser, p *bool, name string, options ...FlagOptionApplyer) error {
	var opts FlagOptions
	opts.applyName(name)
	opts.applyFlagOptions(options)

	return parser.RegisterFlag(newFlag(newBoolValue(p), opts))
}

func Bool(parser Parser, name string, options ...FlagOptionApplyer) *bool {
	p := new(bool)
	_ = BoolVar(parser, p, name, options...)
	return p
}

func IntVar(parser Parser, p *int, name string, options ...FlagOptionApplyer) error {
	var opts FlagOptions
	opts.applyName(name)
	opts.applyFlagOptions(options)

	return parser.RegisterFlag(newFlag(newIntValue(p), opts))
}

func Int(parser Parser, name string, options ...FlagOptionApplyer) *int {
	p := new(int)
	_ = IntVar(parser, p, name, options...)
	return p
}

func IntArgVar(parser Parser, p *int, name string, options ...ArgOptionApplyer) error {
	var opts ArgOptions
	opts.applyName(name)
	opts.applyArgOptions(options)

	return parser.RegisterArg(newArg(newIntValue(p), opts))
}

func IntArg(parser Parser, name string, options ...ArgOptionApplyer) *int {
	p := new(int)
	_ = IntArgVar(parser, p, name, options...)
	return p
}
