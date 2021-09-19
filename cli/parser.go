package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ErrDuplicate = errors.New("duplicate")

	ErrMissingName = errors.New("missing name")

	ErrInvalidName = errors.New("invalid name")

	ErrNotProvided = errors.New("not provided")

	ErrRequiredAfterOptional = errors.New("required after optional")

	ErrArgAfterRest = errors.New("arg after rest")

	ErrUnknownArg = errors.New("unknown arg")
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

	return fmt.Sprintf("parse arg error: '%s': %s", e.Arg, msg)
}

func (e *ParseArgError) Unwrap() error { return e.Err }

func (e *ParseArgError) Is(err error) bool {
	pe, ok := err.(*ParseArgError)
	return ok && pe.Arg == e.Arg && errors.Is(pe.Err, e.Err)
}

type ParseFlagError struct {
	Name string
	Err  error
}

func (e *ParseFlagError) Error() string {
	msg := "unknown error"
	if e.Err != nil {
		msg = e.Err.Error()
	}

	return fmt.Sprintf("parse flag error: '%s': %s", e.Name, msg)
}

func (e *ParseFlagError) Unwrap() error { return e.Err }

func (e *ParseFlagError) Is(err error) bool {
	pe, ok := err.(*ParseFlagError)
	return ok && pe.Name == e.Name && errors.Is(pe.Err, e.Err)
}

type FlagError struct {
	Short string
	Long  string
	Err   error
}

func (e *FlagError) Error() string {
	name := e.Short
	if e.Long != "" {
		if name != "" {
			name += "' '"
		}
		name += e.Long
	}

	msg := "unknown error"
	if e.Err != nil {
		msg = e.Err.Error()
	}

	if name == "" {
		return fmt.Sprintf("flag error: %s", msg)
	} else {
		return fmt.Sprintf("flag error: '%s': %s", name, msg)
	}
}

func (e *FlagError) Is(err error) bool {
	pe, ok := err.(*FlagError)
	return ok && pe.Short == e.Short && pe.Long == e.Long && errors.Is(pe.Err, e.Err)
}

type ArgError struct {
	Name string
	Err  error
}

func (e *ArgError) Error() string {
	msg := "unknown error"
	if e.Err != nil {
		msg = e.Err.Error()
	}

	if e.Name == "" {
		return fmt.Sprintf("arg error: %s", msg)
	} else {
		return fmt.Sprintf("arg error: '%s': %s", e.Name, msg)
	}
}

func (e *ArgError) Is(err error) bool {
	pe, ok := err.(*ArgError)
	return ok && pe.Name == e.Name && errors.Is(pe.Err, e.Err)
}

type RestArgsError struct {
	Name string
	Err  error
}

func (e *RestArgsError) Error() string {
	msg := "unknown error"
	if e.Err != nil {
		msg = e.Err.Error()
	}

	if e.Name == "" {
		return fmt.Sprintf("rest args error: %s", msg)
	} else {
		return fmt.Sprintf("rest args error: '%s': %s", e.Name, msg)
	}
}

func (e *RestArgsError) Is(err error) bool {
	pe, ok := err.(*RestArgsError)
	return ok && pe.Name == e.Name && errors.Is(pe.Err, e.Err)
}

type Register interface {
	RegisterFlag(flag Flag) error
	RegisterArg(arg Arg) error
	RegisterRestArgs(rest RestArgs) error
}

type Commander interface {
	IsCommand(name string) bool
	SetCommand(name string) error
}

type Parser interface {
	Register
	Parse(commander Commander, arguments []string) error
	Args() []Arg
	Rest() *RestArgs
	Flags() []Flag
	FormatLongFlag(name string) string
	FormatShortFlag(name string) string
}

type flags struct {
	data  []Flag
	set   []bool         // Markers if flags were set.
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
			f.set[idx] = false
		}
	}

	if !found {
		if flag.Short != "" {
			idx, found = f.short[flag.Short]
			if found {
				f.data[idx] = flag
				f.set[idx] = false
			}
		}
	}

	if found {
		return
	}

	// Append a new flag.
	f.data = append(f.data, flag)
	f.set = append(f.set, false)
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
	f.set = f.set[:0]
	f.long = nil
	f.short = nil
}

type args struct {
	data  []Arg
	set   []bool         // Markers if args were set.
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
	// Find already added arg.
	idx, found := a.index[arg.Name]
	if found {
		a.data[idx] = arg
		a.set[idx] = false
		return
	}

	// Append a new arg.
	a.data = append(a.data, arg)
	a.set = append(a.set, false)
	idx = len(a.data) - 1

	if a.index == nil {
		a.index = make(map[string]int)
	}

	a.index[arg.Name] = idx
}

func (a *args) Reset() {
	a.data = a.data[:0]
	a.set = a.set[:0]
	a.index = nil
}

var _ Parser = (*DefaultParser)(nil)

type DefaultParser struct {
	Universal     bool
	OverrideFlags bool
	OverrideArgs  bool
	// TODO(SuperPaintman): disable POSIX-style short flag combining (-a -b -> -ab).
	// TODO(SuperPaintman): disable Short-flag+parameter combining (-a parm -> -aparm).

	flags               flags
	args                args
	rest                RestArgs // Other arguments (without named args).
	unknown             []string // Unknown flags (without named flags).
	lastArgOptional     bool     // Is last arg optional.
	registerFlagErr     error    // RegisterFlag first error.
	registerArgErr      error    // RegisterArg first error.
	registerRestArgsErr error    // RegisterRestArgs first error.
}

func (p *DefaultParser) RegisterFlag(flag Flag) (err error) {
	defer func() {
		if err != nil && p.registerFlagErr == nil {
			p.registerFlagErr = err
		}
	}()

	// Check if short of long net is set.
	if flag.Short == "" && flag.Long == "" {
		return &FlagError{Err: ErrMissingName}
	}

	// Validate short flag.
	if flag.Short != "" {
		if !validShortFlag(flag.Short) {
			return &FlagError{
				Short: flag.Short,
				Long:  flag.Long,
				Err:   ErrInvalidName,
			}
		}
	}

	// Validate long flag.
	if flag.Long != "" {
		if !validLongFlag(flag.Long) {
			return &FlagError{
				Short: flag.Short,
				Long:  flag.Long,
				Err:   ErrInvalidName,
			}
		}
	}

	if !p.OverrideFlags {
		if _, f, ok := p.flags.Find(flag.Long, flag.Short); ok {
			return &FlagError{
				Long:  f.Long,
				Short: f.Short,
				Err:   ErrDuplicate,
			}
		}
	}

	p.flags.Add(flag)

	return nil
}

func validShortFlag(name string) bool {
	if len(name) != 1 {
		return false
	}

	if name[0] == '-' || name[0] == '=' || name[0] == ' ' || name[0] == ',' {
		return false
	}

	return true
}

func validLongFlag(name string) bool {
	if len(name) < 1 {
		return false
	}

	if name[0] == '-' || name[0] == '=' || name[0] == ' ' || name[0] == ',' {
		return false
	}

	// We need to iterate by bytes not by runes.
	var foundValid bool
	for i := 1; i < len(name); i++ {
		c := name[i]

		if !foundValid && c == '-' {
			return false
		}

		if c == '=' || c == ' ' || c == ',' {
			return false
		}

		foundValid = true
	}

	return true
}

func (p *DefaultParser) RegisterArg(arg Arg) (err error) {
	defer func() {
		if err != nil && p.registerArgErr == nil {
			p.registerArgErr = err
		}
	}()

	if arg.Required() {
		if p.lastArgOptional {
			return &ArgError{
				Name: arg.Name,
				Err:  ErrRequiredAfterOptional,
			}
		}
	} else {
		p.lastArgOptional = true
	}

	if !p.rest.IsZero() {
		return &ArgError{
			Name: arg.Name,
			Err:  ErrArgAfterRest,
		}
	}

	if arg.Name == "" {
		return &ArgError{Err: ErrMissingName}
	}

	if !validArg(arg.Name) {
		return &ArgError{
			Name: arg.Name,
			Err:  ErrInvalidName,
		}
	}

	if !p.OverrideArgs {
		if _, _, ok := p.args.Get(arg.Name); ok {
			return &ArgError{
				Name: arg.Name,
				Err:  ErrDuplicate,
			}
		}
	}

	p.args.Add(arg)

	return nil
}

func validArg(name string) bool {
	if len(name) < 1 {
		return false
	}

	if name[0] == '-' || name[0] == '=' || name[0] == ' ' || name[0] == ',' {
		return false
	}

	// We need to iterate by bytes not by runes.
	var foundValid bool
	for i := 1; i < len(name); i++ {
		c := name[i]

		if !foundValid && c == '-' {
			return false
		}

		if c == '=' || c == ' ' || c == ',' {
			return false
		}

		foundValid = true
	}

	return true
}

func (p *DefaultParser) RegisterRestArgs(rest RestArgs) (err error) {
	defer func() {
		if err != nil && p.registerRestArgsErr == nil {
			p.registerRestArgsErr = err
		}
	}()

	if rest.Name == "" {
		return &RestArgsError{Err: ErrMissingName}
	}

	if !validRestArgs(rest.Name) {
		return &RestArgsError{
			Name: rest.Name,
			Err:  ErrInvalidName,
		}
	}

	if !p.rest.IsZero() {
		return &RestArgsError{
			Name: rest.Name,
			Err:  ErrDuplicate,
		}
	}

	p.rest = rest

	return nil
}

func validRestArgs(name string) bool {
	return validArg(name)
}

func (p *DefaultParser) Parse(commander Commander, arguments []string) error {
	if p.registerFlagErr != nil {
		return p.registerFlagErr
	}
	if p.registerArgErr != nil {
		return p.registerArgErr
	}
	if p.registerRestArgsErr != nil {
		return p.registerRestArgsErr
	}

	var (
		argMode         bool
		argIdx          int
		flagsTerminated bool
	)
	for {
		if len(arguments) == 0 {
			break
		}

		arg := arguments[0]
		arguments = arguments[1:]

		// Commands or Args.
		if len(arg) == 0 || flagsTerminated || arg[0] != '-' || isNumber(arg) || isDuration(arg) {
			// TODO: "-" for strings

			// Check if the arg is a command.
			if !argMode && commander != nil && commander.IsCommand(arg) {
				// Reset previous flags and args.
				p.flags.Reset()
				p.args.Reset()
				p.rest = RestArgs{}

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

				// Mark the arg as set.
				p.args.set[argIdx] = true
			} else {
				if p.rest.IsZero() {
					return &ParseArgError{
						Arg: arg,
						Err: ErrUnknownArg,
					}
				}

				p.rest.Add(arg)
			}

			argIdx++

			continue
		}

		// Flags.
		numMinuses := 1
		if arg[1] == '-' {
			numMinuses++

			// "--" terminates the flags.
			if len(arg) == 2 {
				flagsTerminated = true
				continue
			}
		}

		shortFlag := numMinuses == 1 && !p.Universal

		name := arg[numMinuses:]
		if len(name) == 0 || name[0] == '-' || name[0] == '=' || name[0] == ' ' || name[0] == ',' {
			return &ParseFlagError{
				Name: name,
				Err:  ErrSyntax,
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
				idx           int
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

				idx, flag, knownflag = p.flags.GetShort(name)
				if knownflag {
					if fv, ok := flag.Value.(boolFlag); ok && fv.IsBoolFlag() {
					}
				}

				// TODO(SuperPaintman): add Short-flag+parameter combining (-a parm -> -aparm).
			} else {
				restName = ""

				idx, flag, knownflag = p.flags.GetLong(name)
			}

			if (!shortFlag || lastShortFlag) && !hasValue && len(arguments) > 0 {
				next := arguments[0]

				var setValue bool
				if len(next) == 0 {
					if knownflag {
						// Special case for empty string flags.
						if fv, ok := flag.Value.(stringFlag); ok && fv.IsStringFlag() {
							setValue = true
						}
					}
				} else if len(next) > 0 && (next[0] != '-' || isNumber(next) || isDuration(next)) {
					// Set value if this is a known flag (if it is a bool we also check the value).
					setValue = knownflag

					if knownflag {
						// Special case for bool flags. Allow only bool-like values.
						if fv, ok := flag.Value.(boolFlag); ok && fv.IsBoolFlag() {
							setValue = isBoolValue(next)
						}
					}
				}

				if setValue {
					value = next
					hasValue = true
					arguments = arguments[1:]
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

			// Mark the flag as set.
			p.flags.set[idx] = true
		}
	}

	// Check required flags.
	for idx, v := range p.flags.set {
		if !v {
			flag := &p.flags.data[idx]

			if flag.Required() {
				return &FlagError{
					Short: flag.Short,
					Long:  flag.Long,
					Err:   ErrNotProvided,
				}
			}
		}
	}

	// Check required args.
	for idx, v := range p.args.set {
		if !v {
			arg := &p.args.data[idx]

			if arg.Required() {
				return &ArgError{
					Name: arg.Name,
					Err:  ErrNotProvided,
				}
			}
		}
	}

	return nil
}

func (p *DefaultParser) Args() []Arg {
	return p.args.data
}

func (p *DefaultParser) Rest() *RestArgs {
	return &p.rest
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

func isDuration(s string) bool {
	// TODO(SuperPaintman): optimize it.

	if _, err := time.ParseDuration(s); err == nil {
		return true
	}

	return false
}
