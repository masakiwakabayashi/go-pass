// Package password はパスワードの生成ロジックを提供する。
package password

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

// 使用する文字種。
const (
	upperChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerChars  = "abcdefghijklmnopqrstuvwxyz"
	digitChars  = "0123456789"
	symbolChars = "!@#$%^&*()-_=+[]{}<>?"
)

// 桁数の既定値と上限。
const (
	DefaultLength = 12
	MaxLength     = 25
)

// Options はパスワード生成の設定。
type Options struct {
	// Length は生成するパスワードの桁数。
	Length int
	// NoSymbols が true のとき、記号を使わない。
	NoSymbols bool
	// NoUpper が true のとき、大文字を使わない。
	NoUpper bool
}

// DefaultOptions は既定の設定を返す。
func DefaultOptions() Options {
	return Options{Length: DefaultLength}
}

// Generate は opts に従ってパスワードを生成する。
// 有効な文字種（既定では大文字・小文字・数字・記号）は、それぞれ最低1文字ずつ含まれる。
func Generate(opts Options) (string, error) {
	if opts.Length > MaxLength {
		return "", fmt.Errorf("桁数（%d）は%d以下にしてください", opts.Length, MaxLength)
	}

	charsets := []string{lowerChars, digitChars}
	if !opts.NoUpper {
		charsets = append(charsets, upperChars)
	}
	if !opts.NoSymbols {
		charsets = append(charsets, symbolChars)
	}
	return generate(opts.Length, charsets)
}

// generate は charsets の各文字種を最低1文字ずつ含む、length 桁のパスワードを生成する。
func generate(length int, charsets []string) (string, error) {
	if len(charsets) == 0 {
		return "", errors.New("有効な文字種がありません")
	}
	for _, cs := range charsets {
		if cs == "" {
			return "", errors.New("空の文字種が指定されています")
		}
	}
	if length < len(charsets) {
		return "", fmt.Errorf("桁数（%d）は有効な文字種の数（%d）以上にしてください", length, len(charsets))
	}

	buf := make([]byte, 0, length)

	// 有効な文字種から1文字ずつ必ず入れる。
	for _, cs := range charsets {
		c, err := randomChar(cs)
		if err != nil {
			return "", err
		}
		buf = append(buf, c)
	}

	// 残りはすべての文字種から均等に選ぶ。
	all := strings.Join(charsets, "")
	for len(buf) < length {
		c, err := randomChar(all)
		if err != nil {
			return "", err
		}
		buf = append(buf, c)
	}

	// 必須文字の位置に偏りが出ないようにシャッフルする（Fisher-Yates）。
	for i := len(buf) - 1; i > 0; i-- {
		j, err := randomInt(i + 1)
		if err != nil {
			return "", err
		}
		buf[i], buf[j] = buf[j], buf[i]
	}

	return string(buf), nil
}

// randomChar は chars から1文字を均等に選ぶ。
func randomChar(chars string) (byte, error) {
	i, err := randomInt(len(chars))
	if err != nil {
		return 0, err
	}
	return chars[i], nil
}

// randomInt は crypto/rand を使って [0, n) の整数を均等に返す。
func randomInt(n int) (int, error) {
	v, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		return 0, fmt.Errorf("乱数の生成に失敗しました: %w", err)
	}
	return int(v.Int64()), nil
}
