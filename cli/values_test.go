package cli

import (
	"reflect"
	"testing"
)

func TestValues_Set_comma_separated(t *testing.T) {
	tt := []struct {
		name  string
		setup func() Getter
		value string
		want  interface{}
	}{
		{
			name: "bools",
			setup: func() Getter {
				var v boolValues
				return &v
			},
			value: "true,y,F,no,T",
			want:  []bool{true, true, false, false, true},
		},
		{
			name: "strings",
			setup: func() Getter {
				var v stringValues
				return &v
			},
			value: "a,b,1337,true,e",
			want:  []string{"a", "b", "1337", "true", "e"},
		},
		{
			name: "ints",
			setup: func() Getter {
				var v intValues
				return &v
			},
			value: "0,-7331,1337,0xABC,0b10101110",
			want:  []int{0, -7331, 1337, 0xABC, 0b10101110},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			v := tc.setup()

			if err := v.Set(tc.value); err != nil {
				t.Fatalf("Set(%q): failed to set the value", tc.value)
			}

			got := v.Get()
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Set(%q): got = %#v, want = %#v", tc.value, got, tc.want)
			}
		})
	}
}

func TestValues_String(t *testing.T) {
	tt := []struct {
		name  string
		value Value
		want  string
	}{
		{
			name:  "bools",
			value: &boolValues{true, true, false, false, true},
			want:  "true,true,false,false,true",
		},
		{
			name:  "strings",
			value: &stringValues{"a", "b", "1337", "true", "e"},
			want:  "a,b,1337,true,e",
		},
		{
			name:  "ints",
			value: &intValues{0, -7331, 1337, 0xABC, 0b10101110},
			want:  "0,-7331,1337,2748,174",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.value.String()
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("String(): got = %q, want = %q", got, tc.want)
			}
		})
	}
}
