package helper

import (
	"fmt"
)

// GetPath - получить путь к сохранению json http запроса и ответа
func GetPath(folder, filename, postfix string) string {
	return fmt.Sprintf("./json/%s/%s_%s.json", folder, filename, postfix)
}

// GetUint32FromInterfaceFloat64Nullable - получить Uint32 из interface{}.float64
func GetUint32FromInterfaceFloat64Nullable(value interface{}) uint32 {
	if v, ok := value.(float64); ok {
		return uint32(v)
	} else {
		return 0
	}
}

// GetStringFromInterfaceStringNullable - получить String из interface{}.string
func GetStringFromInterfaceStringNullable(value interface{}) string {
	if v, ok := value.(string); ok {
		return v
	} else {
		return ""
	}
}

// GetUint64FromInterfaceFloat64Nullable - получить Uint64 из interface{}.float64
func GetUint64FromInterfaceFloat64Nullable(value interface{}) uint64 {
	if v, ok := value.(float64); ok {
		return uint64(v)
	} else {
		return 0
	}
}
