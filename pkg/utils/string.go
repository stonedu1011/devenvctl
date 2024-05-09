package utils

import (
	"github.com/cisco-open/go-lanai/pkg/utils"
	"strings"
)

func SnakeCase(s string) string {
	ret := utils.CamelToSnakeCase(s)
	ret = strings.ReplaceAll(ret, " ", "-")
	return ret
}
