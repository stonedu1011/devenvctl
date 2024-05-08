package utils

import (
	"reflect"
)

func ConvertSlice[SRC any, DST any](src []SRC, converter func (SRC) DST) []DST {
	dst := make([]DST, 0, len(src))
	for i := range src {
		v := converter(src[i])
		if !reflect.ValueOf(v).IsZero() {
			dst = append(dst, v)
		}
	}
	return dst
}
