package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mozillazg/go-pinyin"
	"regexp"
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

// ExtractComment 提取TAPD评论相关信息
func ExtractComment(s string) string {
	s = strings.Trim(s, "\n  ")
	s = strings.ReplaceAll(s, "\n  ", ",")
	return s
}

// ExtractCommentCuePeople 提取评论里面@的人员
func ExtractCommentCuePeople(s string) []string {
	res := make([]string, 0)
	commentRegex := regexp.MustCompile("\\s@\\S+\\s")
	slice := commentRegex.FindAllString(s, -1)
	peopleRegex := regexp.MustCompile("@[^\\(\\)\\s]+")
	for _, sub := range slice {
		people := peopleRegex.FindAllString(sub, 1)
		if len(people) > 0 {
			res = append(res, strings.ReplaceAll(people[0], "@", ""))
		}
	}
	return res
}
