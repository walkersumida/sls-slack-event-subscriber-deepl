package ja_test

import (
	"testing"

	"github.com/walkersumida/sls-slack-event-subscriber-template/event/action/strings/ja"
)

func TestIsJA(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "is Hiragana",
			input: "こんにちは",
			want:  true,
		},
		{
			name:  "are English and Hiragana",
			input: "Walker こんにちは",
			want:  true,
		},
		{
			name:  "is Kanji",
			input: "漢字",
			want:  true,
		},
		{
			name:  "is Katakana",
			input: "カタカナ",
			want:  true,
		},
		{
			name:  "is not ja",
			input: "hello",
			want:  false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ja.IsJA(tc.input); got != tc.want {
				t.Errorf("isJA(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}
