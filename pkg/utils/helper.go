package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

func MapGet[T any](m map[string]any, key string) T {
	if val, ok := m[key]; ok {
		if v, ok := val.(T); ok {
			return v
		}
	}

	return *new(T)
}

func StrToTime(str string) (time.Time, error) {
	var res time.Time
	if num, err := strconv.Atoi(str); err == nil {
		res = time.Unix(int64(num), 0)
	} else {
		res, err = time.Parse("2006-01-02 15:04:05", str)

		if err != nil {
			return time.Time{}, err
		}
	}

	return res, nil
}

func Getenv(key string, defaultVal string) string {
	if val := os.Getenv(key); val == "" {
		return defaultVal
	} else {
		return val
	}
}

func Md5hash(buf []byte) (string, error) {
	hash := md5.New()

	length, err := hash.Write(buf)

	if err != nil {
		return "", err
	}

	if length != len(buf) {
		return "", errors.New("generate md5 hash failed")
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func MaskEmail(email string) string {
	if email == "" {
		return ""
	}

	parts := strings.SplitN(email, "@", 2)
	if len(parts) != 2 {
		return email
	}

	length := len(parts[0])
	if length == 1 {
		return "*" + "@" + parts[1]
	} else if length == 2 {
		return parts[0][0:1] + "*@" + parts[1]
	}

	mask := make([]byte, length-2)
	for i := range mask {
		mask[i] = '*'
	}

	return string(parts[0][:1]) + string(mask) + parts[0][length-1:] + "@" + parts[1]
}

func RandStr(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func Replace(source string, old, new []string) string {
	for i := range old {
		source = strings.ReplaceAll(source, old[i], new[i])
	}

	return source
}

func GetIP(r *http.Request) string {
	for _, header := range []string{"X-Real-IP", "X-Forwarded-For"} {
		if ip := r.Header.Get(header); ip != "" {
			return ip
		}
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

func Unique[T ~int | ~string](list []T) []T {
	hash := map[T]struct{}{}
	for _, val := range list {
		hash[val] = struct{}{}
	}

	res := make([]T, len(list))
	for key := range hash {
		res = append(res, key)
	}

	return res
}

func JsonEncode(data any) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	if err := jsoniter.NewEncoder(buf).Encode(&data); err != nil {
		return nil, err
	}

	return buf, nil
}

func JsonDecode[T any](reader io.Reader) (T, error) {
	var res T
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}
	if err := jsoniter.NewDecoder(reader).Decode(&res); err != nil {
		return res, err
	}

	return res, nil
}

func ToPtr[T any](v T) *T {
	return &v
}
