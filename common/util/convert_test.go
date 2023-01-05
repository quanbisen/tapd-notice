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
	fmt.Println(ConvertUnicodeToCharacter("\\u5168\\u78a7\\u68ee"))
}

func TestCharToPinyin(t *testing.T) {
	assert.Equal(t, "quanbisen", CharToPinyin("全碧森"))
	assert.Equal(t, "quanbisen123", CharToPinyin("全碧森123"))
}
