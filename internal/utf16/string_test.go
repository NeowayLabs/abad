package utf16_test

import (
	"testing"
	"reflect"

	"github.com/NeowayLabs/abad/internal/utf16"
)

var S = utf16.S

func TestStrIndex(t *testing.T) {
	for _, tc := range []struct {
		str      utf16.Str
		sub      utf16.Str
		index int
		contains bool
	}{
		{
			str:      S("hello world"),
			sub:      S("h"),
			index: 0,
			contains: true,
		},
		{
			str:      S("hello world"),
			sub:      S("e"),
			index: 1,
			contains: true,
		},
		{
			str:      S("hello world"),
			sub:      S("l"),
			index: 2,
			contains: true,
		},
		{
			str:      S("hello world"),
			sub:      S(""),
			index: 0,
			contains: true,
		},
		{
			str:      S("hello world"),
			sub:      S("world"),
			index: 6,
			contains: true,
		},
		{
			str:      S(""),
			sub:      S("world"),
			index: -1,
			contains: false,
		},
		{
			str:      S("hello"),
			sub:      S("hello "),
			index: -1,
			contains: false,
		},
		{
			str:      S("hello evil world"),
			sub:      S("evil"),
			index: 6,
			contains: true,
		},
	} {
		got := tc.str.Index(tc.sub)
		if got != tc.index {
			t.Fatalf("index differs: %d != %d",
				got, tc.index)
		}

		gotcontains := tc.str.Contains(tc.sub)
		if gotcontains != tc.contains {
			t.Fatalf("contains differs: %v != %v",
				gotcontains, tc.contains)
		}
	}
}


func TestCreateStringFromRunes(t *testing.T) {
	wantstr := "hello world"
	want := []rune(wantstr)
	
	str := utf16.NewFromRunes(want)
	got := str.Runes()
	
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want[%v] != got[%v]", want, got)
	}
	
	if str.String() != wantstr {
		t.Fatalf("want[%s] != got[%s]", wantstr, str.String())
	}
}
