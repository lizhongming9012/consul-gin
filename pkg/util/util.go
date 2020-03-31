package util

import (
	"NULL/consul-gin/pkg/setting"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"github.com/mozillazg/go-pinyin"
	r "math/rand"
	"strings"
	"time"
	"unicode"
)

// Setup Initialize the util
func Setup() {
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
}
func TransToCharacter(s string) string {
	a := pinyin.NewArgs()
	tmp := make([]string, 0)
	s = strings.TrimSpace(s)
	beg := strings.Index(s, "（")
	end := strings.Index(s, "）")
	if beg > 0 {
		if end == len(s)-3 {
			s = s[:beg]
		} else {
			if beg == 0 {
				s = s[end+3:]
			} else {
				s1 := s[:beg+3]
				s2 := s[end+3:]
				s = s1 + s2
			}
		}
	}
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")
	s = strings.ReplaceAll(s, "*", "")
	s = strings.ReplaceAll(s, "+", "")
	length := len(s)
	for i := 0; i < length; {
		c := s[i : i+1]
		rs := []rune(c)
		if unicode.IsDigit(rs[0]) || unicode.IsLetter(rs[0]) {
			tmp = append(tmp, string(c))
			i++
			continue
		} else {
			if i+3 > length {
				return s
			}
			c = s[i : i+3]
			quan := pinyin.Pinyin(string(c), a)
			for _, z := range quan {
				tmp = append(tmp, z[0][:1])
			}
			i = i + 3
		}
	}
	res := strings.Join(tmp, "")
	return res
}

// 随机字符串
func RandomString(n int, alphabets ...byte) string {
	const alphanum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	var randby bool
	if num, err := rand.Read(bytes); num != n || err != nil {
		r.Seed(time.Now().UnixNano())
		randby = true
	}
	for i, b := range bytes {
		if len(alphabets) == 0 {
			if randby {
				bytes[i] = alphanum[r.Intn(len(alphanum))]
			} else {
				bytes[i] = alphanum[b%byte(len(alphanum))]
			}
		} else {
			if randby {
				bytes[i] = alphabets[r.Intn(len(alphabets))]
			} else {
				bytes[i] = alphabets[b%byte(len(alphabets))]
			}
		}
	}
	return string(bytes)
}
func Sha1Sign(s string) string {
	// The pattern for generating a hash is `sha1.New()`,
	// `sha1.Write(bytes)`, then `sha1.Sum([]byte{})`.
	// Here we start with a new hash.
	h := sha1.New()

	// `Write` expects bytes. If you have a string `s`,
	// use `[]byte(s)` to coerce it to bytes.
	h.Write([]byte(s))

	// This gets the finalized hash result as a byte
	// slice. The argument to `Sum` can be used to append
	// to an existing byte slice: it usually isn't needed.
	bs := h.Sum(nil)

	// SHA1 values are often printed in hex, for example
	// in git commits. Use the `%x` format verb to convert
	// a hash results to a hex string.
	return fmt.Sprintf("%x", bs)
}
