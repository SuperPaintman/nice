package cli

import "strconv"

// var _ flag.Value = (Value)(nil)

type Value interface {
	String() string
	Set(string) error
}

// var _ flag.Value = (Getter)(nil)

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

type boolFlag interface {
	Value
	IsBoolFlag() bool
}

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
