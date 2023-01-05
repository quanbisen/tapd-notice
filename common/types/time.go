package types

import (
	"tapd-notice/common"
	"time"
)

type LocalTime time.Time

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	// 空值不进行解析
	if len(data) == 2 {
		*t = LocalTime(time.Time{})
		return
	}

	// 指定解析的格式
	now, err := time.Parse(`"`+common.StrTimeFormat+`"`, string(data))
	*t = LocalTime(now)
	return
}

func (t *LocalTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(common.StrTimeFormat)+2)
	b = append(b, '"')
	b = time.Time(*t).AppendFormat(b, common.StrTimeFormat)
	b = append(b, '"')
	return b, nil
}
