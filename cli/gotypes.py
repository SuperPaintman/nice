#!/usr/bin/env python

types = [
    ("bool", "Bool", "strconv.FormatBool(bool(*%s))", "*%s == false"),
    ("uint8", "Uint8", "strconv.FormatUint(uint64(*%s), 10)", "*%s == 0"),
    ("uint16", "Uint16", "strconv.FormatUint(uint64(*%s), 10)", "*%s == 0"),
    ("uint32", "Uint32", "strconv.FormatUint(uint64(*%s), 10)", "*%s == 0"),
    ("uint64", "Uint64", "strconv.FormatUint(uint64(*%s), 10)", "*%s == 0"),
    ("int8", "Int8", "strconv.FormatInt(int64(*%s), 10)", "*%s == 0"),
    ("int16", "Int16", "strconv.FormatInt(int64(*%s), 10)", "*%s == 0"),
    ("int32", "Int32", "strconv.FormatInt(int64(*%s), 10)", "*%s == 0"),
    ("int64", "Int64", "strconv.FormatInt(int64(*%s), 10)", "*%s == 0"),
    ("float32", "Float32", "strconv.FormatFloat(float64(*%s), 'g', -1, 32)", "*%s == 0.0"),
    ("float64", "Float64", "strconv.FormatFloat(float64(*%s), 'g', -1, 64)", "*%s == 0.0"),
    ("string", "String", "string(*%s)", "*%s == \"\""),
    ("int", "Int", "strconv.Itoa(int(*%s))", "*%s == 0"),
    ("uint", "Uint", "strconv.FormatUint(uint64(*%s), 10)", "*%s == 0"),
    ("time.Duration", "Duration", "(*time.Duration)(%s).String()", "*%s == 0"),
    # TODO: Func
]

imports = [
    "time"
]

imports_stringer = [
    "strconv"
]
