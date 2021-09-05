package cli

import "testing"

func TestIsBoolValue(t *testing.T) {
	tt := []struct {
		value string
		want  bool
	}{
		// 1
		{
			value: "1",
			want:  true,
		},
		// t
		{
			value: "t",
			want:  true,
		},
		{
			value: "T",
			want:  true,
		},
		// true
		{
			value: "true",
			want:  true,
		},
		{
			value: "True",
			want:  true,
		},
		{
			value: "TRUE",
			want:  true,
		},
		{
			value: "TrUe",
			want:  true,
		},
		// y
		{
			value: "y",
			want:  true,
		},
		{
			value: "Y",
			want:  true,
		},
		// yes
		{
			value: "yes",
			want:  true,
		},
		{
			value: "Yes",
			want:  true,
		},
		{
			value: "YEs",
			want:  true,
		},
		{
			value: "yEs",
			want:  true,
		},
		// 0
		{
			value: "0",
			want:  true,
		},
		// f
		{
			value: "f",
			want:  true,
		},
		{
			value: "F",
			want:  true,
		},
		// false
		{
			value: "false",
			want:  true,
		},
		{
			value: "False",
			want:  true,
		},
		{
			value: "FALSE",
			want:  true,
		},
		{
			value: "fAlSe",
			want:  true,
		},
		// n
		{
			value: "n",
			want:  true,
		},
		{
			value: "N",
			want:  true,
		},
		// no
		{
			value: "no",
			want:  true,
		},
		{
			value: "No",
			want:  true,
		},
		{
			value: "NO",
			want:  true,
		},
		{
			value: "nO",
			want:  true,
		},
		// Wrong.
		{
			value: "",
		},
		{
			value: "000",
		},
		{
			value: "100",
		},
		{
			value: "da",
		},
		{
			value: "-1",
		},
		{
			value: "falseValue",
		},
	}

	for _, tc := range tt {
		t.Run(tc.value, func(t *testing.T) {
			got := isBoolValue(tc.value)
			if got != tc.want {
				t.Errorf("isBoolValue(): got = %v, want = %v", got, tc.want)
			}
		})
	}
}

var isBoolValueRes bool

func BenchmarkIsBoolValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, v := range [...]string{"t", "true", "TRUE", "Y", "Yes", "f", "false", "FALSE", "N", "No"} {
			isBoolValueRes = isBoolValue(v)
		}
	}
}

func TestParseBool(t *testing.T) {
	tt := []struct {
		value string
		want  bool
	}{
		// 1
		{
			value: "1",
			want:  true,
		},
		// t
		{
			value: "t",
			want:  true,
		},
		{
			value: "T",
			want:  true,
		},
		// true
		{
			value: "true",
			want:  true,
		},
		{
			value: "True",
			want:  true,
		},
		{
			value: "TRUE",
			want:  true,
		},
		{
			value: "TrUe",
			want:  true,
		},
		// y
		{
			value: "y",
			want:  true,
		},
		{
			value: "Y",
			want:  true,
		},
		// yes
		{
			value: "yes",
			want:  true,
		},
		{
			value: "Yes",
			want:  true,
		},
		{
			value: "YEs",
			want:  true,
		},
		{
			value: "yEs",
			want:  true,
		},
		// 0
		{
			value: "0",
			want:  false,
		},
		// f
		{
			value: "f",
			want:  false,
		},
		{
			value: "F",
			want:  false,
		},
		// false
		{
			value: "false",
			want:  false,
		},
		{
			value: "False",
			want:  false,
		},
		{
			value: "FALSE",
			want:  false,
		},
		{
			value: "fAlSe",
			want:  false,
		},
		// n
		{
			value: "n",
			want:  false,
		},
		{
			value: "N",
			want:  false,
		},
		// no
		{
			value: "no",
			want:  false,
		},
		{
			value: "No",
			want:  false,
		},
		{
			value: "NO",
			want:  false,
		},
		{
			value: "nO",
			want:  false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.value, func(t *testing.T) {
			got, err := parseBool(tc.value)
			if err != nil {
				t.Fatalf("parseBool(): err = %s", err)
			}

			if got != tc.want {
				t.Errorf("parseBool(): got = %v, want = %v", got, tc.want)
			}
		})
	}
}

func TestParseBool_wrong_value(t *testing.T) {
	tt := []struct {
		value string
	}{
		{
			value: "",
		},
		{
			value: "000",
		},
		{
			value: "100",
		},
		{
			value: "da",
		},
		{
			value: "-1",
		},
		{
			value: "falseValue",
		},
	}

	for _, tc := range tt {
		t.Run(tc.value, func(t *testing.T) {
			_, err := parseBool(tc.value)
			if err == nil {
				t.Fatalf("parseBool(): got nil, want err")
			}
		})
	}
}

var parseBoolRes bool

func BenchmarkParseBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, v := range [...]string{"t", "true", "TRUE", "Y", "Yes", "f", "false", "FALSE", "N", "No"} {
			parseBoolRes, _ = parseBool(v)
		}
	}
}
