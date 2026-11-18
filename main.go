// go-pass はパスワード用のランダムな文字列を生成するCLIツール。
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/masakiwakabayashi/go-pass/internal/password"
)

func main() {
	opts := password.DefaultOptions()

	// 解析エラー時も終了コード1にするため、ContinueOnError で自前で処理する。
	fs := flag.NewFlagSet("go-pass", flag.ContinueOnError)
	fs.IntVar(&opts.Length, "l", opts.Length, fmt.Sprintf("パスワードの桁数（最大%d）", password.MaxLength))
	fs.BoolVar(&opts.NoSymbols, "ns", opts.NoSymbols, "記号を使わない")
	fs.BoolVar(&opts.NoUpper, "nu", opts.NoUpper, "大文字を使わない")
	if err := fs.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		os.Exit(1) // エラー内容と使い方は fs.Parse が標準エラー出力に表示済み
	}

	if fs.NArg() > 0 {
		fail(fmt.Errorf("不明な引数があります: %v", fs.Args()))
	}

	pw, err := password.Generate(opts)
	if err != nil {
		fail(err)
	}
	fmt.Println(pw)
}

// fail はエラーを標準エラー出力に出し、終了コード1で終了する。
func fail(err error) {
	fmt.Fprintln(os.Stderr, "エラー:", err)
	os.Exit(1)
}
