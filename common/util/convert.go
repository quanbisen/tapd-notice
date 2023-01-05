package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mozillazg/go-pinyin"
	"strconv"
	"strings"
	"unicode"
)

func ConvertUnicodeToCharacter(unicodeStr string) string {
	unicodeArray := strings.Split(unicodeStr, "\\u")
	var result string
	for _, v := range unicodeArray {
		if len(v) < 1 {
			continue
		}
		temp, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			result += fmt.Sprintf("%s", v)
		} else {
			result += fmt.Sprintf("%c", temp)
		}
	}
	return result
}

func Int64Addr(a int64) *int64 {
	return &a
}

func GetIntByInterface(a interface{}) (int, error) {
	id := 0
	var err error
	switch val := a.(type) {
	case int:
		id = val
	case int32:
		id = int(val)
	case int64:
		id = int(val)
	case json.Number:
		var tmpID int64
		tmpID, err = val.Int64()
		id = int(tmpID)
	case float64:
		id = int(val)
	case float32:
		id = int(val)
	case string:
		var tmpID int64
		tmpID, err = strconv.ParseInt(a.(string), 10, 64)
		if err != nil {
			return 0, err
		}
		id = int(tmpID)
	default:
		err = errors.New("not numeric")
	}
	return id, err
}

func CharToPinyin(s string) string {
	res := ""
	args := pinyin.NewArgs()
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			res += pinyin.LazyPinyin(fmt.Sprintf("%v", string(r)), args)[0]
		} else {
			res += fmt.Sprintf("%v", string(r))
		}
	}
	return res
}
