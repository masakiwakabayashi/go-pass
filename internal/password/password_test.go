package password

import (
	"strings"
	"testing"
)

const iterations = 1000

func TestGenerate(t *testing.T) {
	allowed := upperChars + lowerChars + digitChars + symbolChars

	for range iterations {
		pw, err := Generate()
		if err != nil {
			t.Fatalf("Generate() がエラーを返しました: %v", err)
		}

		if len(pw) != defaultLength {
			t.Fatalf("桁数 = %d, want %d（%q）", len(pw), defaultLength, pw)
		}

		for _, cs := range []string{upperChars, lowerChars, digitChars, symbolChars} {
			if !strings.ContainsAny(pw, cs) {
				t.Fatalf("%q に文字種 %q の文字が含まれていません", pw, cs)
			}
		}

		for _, r := range pw {
			if !strings.ContainsRune(allowed, r) {
				t.Fatalf("%q に想定外の文字 %q が含まれています", pw, r)
			}
		}
	}
}

func TestGenerateLength(t *testing.T) {
	charsets := []string{upperChars, lowerChars, digitChars, symbolChars}

	for _, length := range []int{4, 8, 12, 32, 64} {
		for range iterations {
			pw, err := generate(length, charsets)
			if err != nil {
				t.Fatalf("generate(%d) がエラーを返しました: %v", length, err)
			}
			if len(pw) != length {
				t.Fatalf("桁数 = %d, want %d（%q）", len(pw), length, pw)
			}
		}
	}
}

func TestGenerateExcludesDisabledCharsets(t *testing.T) {
	enabled := []string{lowerChars, digitChars}
	disabled := upperChars + symbolChars

	for range iterations {
		pw, err := generate(defaultLength, enabled)
		if err != nil {
			t.Fatalf("generate がエラーを返しました: %v", err)
		}
		if strings.ContainsAny(pw, disabled) {
			t.Fatalf("%q に無効にした文字種の文字が含まれています", pw)
		}
		for _, cs := range enabled {
			if !strings.ContainsAny(pw, cs) {
				t.Fatalf("%q に文字種 %q の文字が含まれていません", pw, cs)
			}
		}
	}
}

func TestGenerateInvalid(t *testing.T) {
	all := []string{upperChars, lowerChars, digitChars, symbolChars}

	tests := []struct {
		name     string
		length   int
		charsets []string
	}{
		{"文字種なし", 12, nil},
		{"空の文字種", 12, []string{lowerChars, ""}},
		{"桁数が文字種の数より少ない", 3, all},
		{"桁数0", 0, all},
		{"桁数が負", -1, all},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pw, err := generate(tt.length, tt.charsets)
			if err == nil {
				t.Fatalf("エラーが返されませんでした（結果 %q）", pw)
			}
			if pw != "" {
				t.Fatalf("エラー時に空でない文字列 %q が返されました", pw)
			}
		})
	}
}
