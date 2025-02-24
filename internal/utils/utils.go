package utils

import (
	"golang.org/x/text/encoding/charmap"
	"os"
)

func CutStringByBytes(text string, limit int) string {
	runes := []rune(text)
	if len(runes) >= limit {
		return string(runes[:limit])
	}
	return text
}

func DecodeWindows1251(enc []byte) string {
	dec := charmap.Windows1251.NewDecoder()
	out, _ := dec.Bytes(enc)
	return string(out)
}

func EncodeWindows1251(inp string) string {
	enc := charmap.Windows1251.NewEncoder()
	out, _ := enc.String(inp)
	return out
}

// Exists - проверка, что файл существует
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
