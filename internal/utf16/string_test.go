package utf16_test

import (
	"reflect"
	"testing"

	"github.com/NeowayLabs/abad/internal/utf16"
)

var S = utf16.S

func TestStrIndex(t *testing.T) {
	for _, tc := range []struct {
		str      utf16.Str
		sub      utf16.Str
		index    int
		contains bool
	}{
		{
			str:      S("hello world"),
			sub:      S("h"),
			index:    0,
			contains: true,
		},
		{
			str:      S("hello world"),
			sub:      S("e"),
			index:    1,
			contains: true,
		},
		{
			str:      S("hello world"),
			sub:      S("l"),
			index:    2,
			contains: true,
		},
		{
			str:      S("hello world"),
			sub:      S(""),
			index:    0,
			contains: true,
		},
		{
			str:      S("hello world"),
			sub:      S("world"),
			index:    6,
			contains: true,
		},
		{
			str:      S(""),
			sub:      S("world"),
			index:    -1,
			contains: false,
		},
		{
			str:      S("hello"),
			sub:      S("hello "),
			index:    -1,
			contains: false,
		},
		{
			str:      S("hello evil world"),
			sub:      S("evil"),
			index:    6,
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

func TestAppendStrings(t *testing.T) {
	type tcase struct {
		s1   utf16.Str
		s2   utf16.Str
		want utf16.Str
	}

	cases := []tcase{
		{s1: S("a"), s2: S("bad"), want: S("abad")},
		{s1: S(""), s2: S("abad"), want: S("abad")},
		{s1: S("abad"), s2: S(""), want: S("abad")},
		{s1: S(""), s2: S(""), want: S("")},
	}

	for _, c := range cases {
		got := c.s1.Append(c.s2)
		if !c.want.Equal(got) {
			t.Fatalf("got[%s] !=  want[%s]", got, c.want)
		}
	}
}

func TestPrependStrings(t *testing.T) {
	type tcase struct {
		s1   utf16.Str
		s2   utf16.Str
		want utf16.Str
	}

	cases := []tcase{
		{s1: S("bad"), s2: S("a"), want: S("abad")},
		{s1: S(""), s2: S("abad"), want: S("abad")},
		{s1: S("abad"), s2: S(""), want: S("abad")},
		{s1: S(""), s2: S(""), want: S("")},
	}

	for _, c := range cases {
		got := c.s1.Prepend(c.s2)
		if !c.want.Equal(got) {
			t.Fatalf("got[%s] !=  want[%s]", got, c.want)
		}
	}
}
