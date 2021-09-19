#!/usr/bin/env python

types = [
    ("bool", "Bool", "strconv.FormatBool(bool(*%s))"),
    ("uint8", "Uint8", "strconv.FormatUint(uint64(*%s), 10)"),
    ("uint16", "Uint16", "strconv.FormatUint(uint64(*%s), 10)"),
    ("uint32", "Uint32", "strconv.FormatUint(uint64(*%s), 10)"),
    ("uint64", "Uint64", "strconv.FormatUint(uint64(*%s), 10)"),
    ("int8", "Int8", "strconv.FormatInt(int64(*%s), 10)"),
    ("int16", "Int16", "strconv.FormatInt(int64(*%s), 10)"),
    ("int32", "Int32", "strconv.FormatInt(int64(*%s), 10)"),
    ("int64", "Int64", "strconv.FormatInt(int64(*%s), 10)"),
    ("float32", "Float32", "strconv.FormatFloat(float64(*%s), 'g', -1, 32)"),
    ("float64", "Float64", "strconv.FormatFloat(float64(*%s), 'g', -1, 64)"),
    ("string", "String", "string(*%s)"),
    ("int", "Int", "strconv.Itoa(int(*%s))"),
    ("uint", "Uint", "strconv.FormatUint(uint64(*%s), 10)"),
    # TODO: Duration
    # TODO: Func
]

imports_stringer = [
    "strconv"
]
