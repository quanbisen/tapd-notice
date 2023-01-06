package util

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestConvertEncode1(t *testing.T) {
	s := "\\u6d4b\\u8bd5"
	assert.Equal(t, "测试", ConvertUnicodeToCharacter(s))
}

func TestConvertEncode2(t *testing.T) {
	s := "test\\u6d4b\\u8bd5"
	assert.Equal(t, "test测试", ConvertUnicodeToCharacter(s))
}

func TestConvertEncode3(t *testing.T) {
	fmt.Println(ConvertUnicodeToCharacter("\\u003cp\\u003e测试评论\\u003c/p\\u003e"))
}

func TestCharToPinyin(t *testing.T) {
	assert.Equal(t, "quanbisen", CharToPinyin("全碧森"))
	assert.Equal(t, "quanbisen123", CharToPinyin("全碧森123"))
}

func TestExtractComment(t *testing.T) {
	s := "\n  123\n  456\n  "
	fmt.Println(ExtractComment(s))
}

func TestExtractCommentCuePeople(t *testing.T) {
	s := "\n   @全碧森(全碧森)  @邹思瑶(邹思瑶) 麻烦帮忙解决\n  \n  "
	assert.Equal(t, ExtractCommentCuePeople(s), []string{"全碧森", "邹思瑶"})
	s = "\n   @全碧森(全碧森) @邹\n  \n  "
	assert.Equal(t, ExtractCommentCuePeople(s), []string{"全碧森"})
}
