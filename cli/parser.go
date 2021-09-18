package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type ParseArgError struct {
	Arg string
	Err error
}

func (e *ParseArgError) Error() string {
	msg := "unknown error"
	if e.Err != nil {
		msg = e.Err.Error()
	}

	return fmt.Sprintf("parse arg error: %s: %s", e.Arg, msg)
}

func (e *ParseArgError) Unwrap() error { return e.Err }

func (e *ParseArgError) Is(err error) bool {
	pe, ok := err.(*ParseArgError)
	return ok && pe.Arg == e.Arg && errors.Is(pe.Err, e.Err)
}

type DuplicatedFlagError struct {
	Flag *Flag
}

func (e *DuplicatedFlagError) Is(err error) bool {
	pe, ok := err.(*DuplicatedFlagError)
	return ok && pe.Flag == e.Flag
}

func (e *DuplicatedFlagError) Error() string {
	return fmt.Sprintf("duplicated flag: %s", e.Flag.String())
}

type DuplicatedArgError struct {
	Arg *Arg
}

func (e *DuplicatedArgError) Is(err error) bool {
	pe, ok := err.(*DuplicatedArgError)
	return ok && pe.Arg == e.Arg
}

func (e *DuplicatedArgError) Error() string {
	return fmt.Sprintf("duplicated arg: %s", e.Arg.String())
}

type Register interface {
	RegisterFlag(flag Flag) error
	RegisterArg(arg Arg) error
}

type Commander interface {
	IsCommand(name string) bool
	SetCommand(name string) error
}

type Parser interface {
	Register
	Parse(commander Commander, arguments []string) error
	Args() []Arg
	Flags() []Flag
	FormatLongFlag(name string) string
	FormatShortFlag(name string) string
}

type flags struct {
	data  []Flag
	long  map[string]int // Indexes of flags in the data. It contains all long names.
	short map[string]int // Indexes of flags in the data. It contains all short names.
}

func (f *flags) GetLong(name string) (idx int, flag *Flag, ok bool) {
	idx, ok = f.long[name]
	if ok {
		flag = &f.data[idx]
	}

	return
}

func (f *flags) GetShort(name string) (idx int, flag *Flag, ok bool) {
	idx, ok = f.short[name]
	if ok {
		flag = &f.data[idx]
	}

	return
}

func (f *flags) Find(long, short string) (idx int, flag *Flag, ok bool) {
	if long != "" {
		if idx, ok = f.long[long]; ok {
			flag = &f.data[idx]
			return
		}
	}

	if short != "" {
		if idx, ok = f.short[short]; ok {
			flag = &f.data[idx]
			return
		}
	}

	return
}

func (f *flags) Add(flag Flag) {
	// Find already added flag.
	var (
		idx   int
		found bool
	)
	if flag.Long != "" {
		idx, found = f.long[flag.Long]
		if found {
			f.data[idx] = flag
		}
	}

	if !found {
		if flag.Short != "" {
			idx, found = f.short[flag.Short]
			if found {
				f.data[idx] = flag
			}
		}
	}

	if found {
		return
	}

	// Append a new flag.
	f.data = append(f.data, flag)
	idx = len(f.data) - 1

	if flag.Long != "" {
		if f.long == nil {
			f.long = make(map[string]int)
		}

		f.long[flag.Long] = idx
	}

	if flag.Short != "" {
		if f.short == nil {
			f.short = make(map[string]int)
		}

		f.short[flag.Short] = idx
	}
}

func (f *flags) Reset() {
	f.data = f.data[:0]
	f.long = nil
	f.short = nil
}

type args struct {
	data  []Arg
	index map[string]int // Indexes of named arg in the data.
}

func (a *args) Get(name string) (idx int, arg *Arg, ok bool) {
	idx, ok = a.index[name]
	if ok {
		arg = &a.data[idx]
	}
	return
}

func (a *args) Nth(i int) (arg *Arg, ok bool) {
	if i >= len(a.data) {
		return
	}

	return &a.data[i], true
}

func (a *args) Add(arg Arg) {
	if arg.Name == "" {
		return
	}

	// Find already added arg.
	idx, found := a.index[arg.Name]
	if found {
		a.data[idx] = arg
		return
	}

	// Append a new arg.
	a.data = append(a.data, arg)
	idx = len(a.data) - 1

	if a.index == nil {
		a.index = make(map[string]int)
	}

	a.index[arg.Name] = idx
}

func (a *args) Reset() {
	a.data = a.data[:0]
	a.index = nil
}

var _ Parser = (*DefaultParser)(nil)

type DefaultParser struct {
	Universal     bool
	OverrideFlags bool
	OverrideArgs  bool
	// TODO(SuperPaintman): disable POSIX-style short flag combining (-a -b -> -ab).
	// TODO(SuperPaintman): disable Short-flag+parameter combining (-a parm -> -aparm).

	flags           flags
	args            args
	unknown         []string             // Unknown flags (without named flags).
	rest            []string             // Other arguments (without named args).
	registerFlagErr *DuplicatedFlagError // RegisterFlag first error.
	registerArgErr  *DuplicatedArgError  // RegisterArg first error.
}

func (p *DefaultParser) RegisterFlag(flag Flag) error {
	if !p.OverrideFlags {
		if _, f, ok := p.flags.Find(flag.Long, flag.Short); ok {
			err := &DuplicatedFlagError{
				Flag: f,
			}

			if p.registerFlagErr == nil {
				p.registerFlagErr = err
			}

			return err
		}
	}

	p.flags.Add(flag)

	return nil
}

func (p *DefaultParser) RegisterArg(arg Arg) error {
	if !p.OverrideArgs {
		if _, a, ok := p.args.Get(arg.Name); ok {
			err := &DuplicatedArgError{
				Arg: a,
			}

			if p.registerFlagErr == nil {
				p.registerArgErr = err
			}

			return err
		}
	}

	p.args.Add(arg)

	return nil
}

func (p *DefaultParser) Parse(commander Commander, arguments []string) error {
	if p.registerFlagErr != nil {
		return p.registerFlagErr
	}
	if p.registerArgErr != nil {
		return p.registerArgErr
	}

	var (
		argMode bool
		argIdx  int
	)
	for {
		if len(arguments) == 0 {
			break
		}

		arg := arguments[0]
		arguments = arguments[1:]

		// Commands or Args.
		if len(arg) == 0 || arg[0] != '-' || isNumber(arg) {
			// Check if the arg is a command.
			if !argMode && commander != nil && commander.IsCommand(arg) {
				// Reset previous flags and args.
				p.flags.Reset()
				p.args.Reset()

				if err := commander.SetCommand(arg); err != nil {
					return err
				}

				continue
			}

			// Parse rest as args.
			argMode = true

			a, ok := p.args.Nth(argIdx)
			if ok {
				if err := a.Value.Set(arg); err != nil {
					return err
				}
			} else {
				p.rest = append(p.rest, arg)
			}

			argIdx++

			continue
		}

		// TODO(SuperPaintman): add required flags.
		// TODO(SuperPaintman): add optional args.

		// Flags.
		numMinuses := 1
		if arg[1] == '-' {
			numMinuses++

			// TODO(SuperPaintman): add the "--" bypass.
		}

		shortFlag := numMinuses == 1 && !p.Universal

		name := arg[numMinuses:]
		if len(name) == 0 || name[0] == '-' || name[0] == '=' || name[0] == ' ' {
			return &ParseArgError{
				Arg: arg,
				Err: ErrSyntax,
			}
		}

		// Find a value.
		var (
			value    string
			hasValue bool
		)
		// Equals cannot be first.
		for i := 1; i < len(name); i++ {
			if name[i] == '=' {
				value = name[i+1:]
				hasValue = true
				name = name[0:i]
				break
			}
		}

		// Find a known flag.
		restName := name
		lastHasValue := hasValue
		lastValue := value
		for len(restName) > 0 {
			// Parse POSIX-style short flag combining (-a -b -> -ab).
			var (
				flag          *Flag
				knownflag     bool
				lastShortFlag bool
			)
			if shortFlag {
				name = restName[:1]
				restName = restName[1:]

				if len(restName) == 0 {
					hasValue = lastHasValue
					value = lastValue
					lastShortFlag = true
				} else {
					hasValue = false
					value = ""
				}

				_, flag, knownflag = p.flags.GetShort(name)
				if knownflag {
					if fv, ok := flag.Value.(boolFlag); ok && fv.IsBoolFlag() {
					}
				}

				// TODO(SuperPaintman): add Short-flag+parameter combining (-a parm -> -aparm).
			} else {
				restName = ""

				_, flag, knownflag = p.flags.GetLong(name)
			}

			if (!shortFlag || lastShortFlag) && !hasValue && len(arguments) > 0 {
				next := arguments[0]
				if len(next) > 0 && (next[0] != '-' || isNumber(next)) {
					setValue := knownflag
					if knownflag {
						// Special case for bool flags. Allow only bool-like values.
						if fv, ok := flag.Value.(boolFlag); ok && fv.IsBoolFlag() {
							setValue = isBoolValue(next)
						}
					}

					if setValue {
						value = next
						hasValue = true
						arguments = arguments[1:]
					}
				}
			}

			if !knownflag {
				prefix := strings.Repeat("-", numMinuses)
				p.unknown = append(p.unknown, prefix+name)
				continue
			}

			// Set Value.
			// Special case for bool flags which doesn't need a value.
			if fv, ok := flag.Value.(boolFlag); ok && fv.IsBoolFlag() {
				if !hasValue {
					value = "true"
				} else if value == "" {
					value = "false"
				}
			}

			if err := flag.Value.Set(value); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *DefaultParser) Args() []Arg {
	return p.args.data
}

func (p *DefaultParser) Flags() []Flag {
	return p.flags.data
}

func (p *DefaultParser) FormatLongFlag(name string) string {
	if name == "" {
		return ""
	}

	if p.Universal {
		return "-" + name
	}

	return "--" + name
}

func (p *DefaultParser) FormatShortFlag(name string) string {
	if name == "" {
		return ""
	}

	return "-" + name
}

func isNumber(s string) bool {
	// TODO(SuperPaintman): optimize it.

	if _, err := strconv.ParseInt(s, 0, 0); err == nil {
		return true
	}

	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return true
	}

	if _, err := strconv.ParseUint(s, 0, 0); err == nil {
		return true
	}

	return false
}
