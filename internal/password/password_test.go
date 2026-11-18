package password

import (
	"strings"
	"testing"
)

const iterations = 1000

func TestGenerate(t *testing.T) {
	allowed := upperChars + lowerChars + digitChars + symbolChars

	for range iterations {
		pw, err := Generate(DefaultOptions())
		if err != nil {
			t.Fatalf("Generate() がエラーを返しました: %v", err)
		}

		if len(pw) != DefaultLength {
			t.Fatalf("桁数 = %d, want %d（%q）", len(pw), DefaultLength, pw)
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
	for _, length := range []int{4, 5, 8, 16, 24, MaxLength} {
		for range iterations {
			pw, err := Generate(Options{Length: length})
			if err != nil {
				t.Fatalf("Generate(Length=%d) がエラーを返しました: %v", length, err)
			}
			if len(pw) != length {
				t.Fatalf("桁数 = %d, want %d（%q）", len(pw), length, pw)
			}
			for _, cs := range []string{upperChars, lowerChars, digitChars, symbolChars} {
				if !strings.ContainsAny(pw, cs) {
					t.Fatalf("%q に文字種 %q の文字が含まれていません", pw, cs)
				}
			}
		}
	}
}

func TestGenerateInvalidLength(t *testing.T) {
	for _, length := range []int{3, 1, 0, -1, MaxLength + 1, 1000000000} {
		pw, err := Generate(Options{Length: length})
		if err == nil {
			t.Fatalf("Generate(Length=%d) でエラーが返されませんでした（結果 %q）", length, pw)
		}
		if pw != "" {
			t.Fatalf("エラー時に空でない文字列 %q が返されました", pw)
		}
	}
}

func TestGenerateNoSymbols(t *testing.T) {
	opts := Options{Length: DefaultLength, NoSymbols: true}

	for range iterations {
		pw, err := Generate(opts)
		if err != nil {
			t.Fatalf("Generate がエラーを返しました: %v", err)
		}
		if len(pw) != DefaultLength {
			t.Fatalf("桁数 = %d, want %d（%q）", len(pw), DefaultLength, pw)
		}
		if strings.ContainsAny(pw, symbolChars) {
			t.Fatalf("%q に記号が含まれています", pw)
		}
		for _, cs := range []string{upperChars, lowerChars, digitChars} {
			if !strings.ContainsAny(pw, cs) {
				t.Fatalf("%q に文字種 %q の文字が含まれていません", pw, cs)
			}
		}
	}
}

func TestGenerateNoSymbolsMinLength(t *testing.T) {
	// 記号なしでは有効な文字種が3つなので、3桁は生成でき、2桁はエラーになる。
	if _, err := Generate(Options{Length: 3, NoSymbols: true}); err != nil {
		t.Fatalf("3桁でエラーが返されました: %v", err)
	}
	if pw, err := Generate(Options{Length: 2, NoSymbols: true}); err == nil {
		t.Fatalf("2桁でエラーが返されませんでした（結果 %q）", pw)
	}
}

func TestGenerateNoUpper(t *testing.T) {
	opts := Options{Length: DefaultLength, NoUpper: true}

	for range iterations {
		pw, err := Generate(opts)
		if err != nil {
			t.Fatalf("Generate がエラーを返しました: %v", err)
		}
		if len(pw) != DefaultLength {
			t.Fatalf("桁数 = %d, want %d（%q）", len(pw), DefaultLength, pw)
		}
		if strings.ContainsAny(pw, upperChars) {
			t.Fatalf("%q に大文字が含まれています", pw)
		}
		for _, cs := range []string{lowerChars, digitChars, symbolChars} {
			if !strings.ContainsAny(pw, cs) {
				t.Fatalf("%q に文字種 %q の文字が含まれていません", pw, cs)
			}
		}
	}
}

func TestGenerateNoSymbolsNoUpper(t *testing.T) {
	opts := Options{Length: DefaultLength, NoSymbols: true, NoUpper: true}

	for range iterations {
		pw, err := Generate(opts)
		if err != nil {
			t.Fatalf("Generate がエラーを返しました: %v", err)
		}
		if len(pw) != DefaultLength {
			t.Fatalf("桁数 = %d, want %d（%q）", len(pw), DefaultLength, pw)
		}
		if strings.ContainsAny(pw, upperChars+symbolChars) {
			t.Fatalf("%q に大文字か記号が含まれています", pw)
		}
		for _, cs := range []string{lowerChars, digitChars} {
			if !strings.ContainsAny(pw, cs) {
				t.Fatalf("%q に文字種 %q の文字が含まれていません", pw, cs)
			}
		}
	}

	// 有効な文字種は2つなので、2桁は生成でき、1桁はエラーになる。
	if _, err := Generate(Options{Length: 2, NoSymbols: true, NoUpper: true}); err != nil {
		t.Fatalf("2桁でエラーが返されました: %v", err)
	}
	if pw, err := Generate(Options{Length: 1, NoSymbols: true, NoUpper: true}); err == nil {
		t.Fatalf("1桁でエラーが返されませんでした（結果 %q）", pw)
	}
}

func TestGenerateExcludesDisabledCharsets(t *testing.T) {
	enabled := []string{lowerChars, digitChars}
	disabled := upperChars + symbolChars

	for range iterations {
		pw, err := generate(DefaultLength, enabled)
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
