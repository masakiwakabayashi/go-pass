// go-pass はパスワード用のランダムな文字列を生成するCLIツール。
package main

import (
	"fmt"
	"os"

	"github.com/masakiwakabayashi/go-pass/internal/password"
)

func main() {
	pw, err := password.Generate()
	if err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}
	fmt.Println(pw)
}
