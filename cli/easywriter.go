package cli

import (
	"fmt"
	"io"
)

type easyWriter struct {
	w   io.Writer
	err error
}

func (ew *easyWriter) Write(data []byte) {
	if ew.err != nil {
		return
	}

	_, err := ew.w.Write(data)
	if err != nil {
		ew.err = err
	}
}

func (ew *easyWriter) WriteString(s string) {
	ew.Write([]byte(s))
}

func (ew *easyWriter) Writef(format string, a ...interface{}) {
	if ew.err != nil {
		return
	}

	var err error
	if len(a) == 0 {
		_, err = ew.w.Write([]byte(format))
	} else {
		_, err = fmt.Fprintf(ew.w, format, a...)
	}

	if err != nil {
		ew.err = err
	}
}

func (ew *easyWriter) Err() error {
	return ew.err
}
