package controller

import (
	"reflect"
	"strconv"
	"strings"
	"time"
	"zhaoxin2025/model"
)

// 将netid转换为session需要的int类型
func NetID(netid string) int {
	var convErr error
	var netidInt int
	// 传入参数时做过校验，故不应出现错误
	netidInt, convErr = strconv.Atoi(netid)
	if convErr != nil {
		return 0
	}
	return netidInt
}

// 将管理员级别super normal转换成session的role数字
func Level(level model.AdminLevel) int {
	switch level {
	case model.Super:
		return 3
	case model.Normal:
		return 2
	default:
		return 0 // 默认值，处理未知的管理员级别
	}
}

// MapStruct 将结构体中的非零值字段转换为 map
// 支持处理嵌套结构体、指针、时间类型等
// 使用 json 标签作为 map 的键名
func Struct2Map(obj any) map[string]any {
	result := make(map[string]any)
	if obj == nil {
		return result
	}

	// 获取结构体的反射值
	val := reflect.ValueOf(obj)

	// 处理指针类型
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return result
		}
		val = val.Elem()
	}

	// 确保是结构体
	if val.Kind() != reflect.Struct {
		return result
	}

	typ := val.Type()

	// 遍历结构体字段
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// 跳过未导出字段
		if !fieldType.IsExported() {
			continue
		}

		// 获取 json 标签作为键名
		tag := fieldType.Tag.Get("json")
		if tag == "-" {
			continue
		}

		// 从 json 标签中提取字段名
		jsonName := fieldType.Name
		if tag != "" {
			parts := strings.Split(tag, ",")
			if parts[0] != "" {
				jsonName = parts[0]
			}
		}

		// 处理不同类型的字段
		switch field.Kind() {
		case reflect.Ptr, reflect.Interface:
			// 处理指针和接口
			if !field.IsNil() {
				// 如果是指针指向结构体，递归处理
				if field.Elem().Kind() == reflect.Struct &&
					field.Elem().Type() != reflect.TypeOf(time.Time{}) {
					nestedMap := Struct2Map(field.Interface())
					for k, v := range nestedMap {
						result[jsonName+"."+k] = v
					}
				} else {
					// 其他非空指针
					result[jsonName] = field.Interface()
				}
			}
		case reflect.Struct:
			// 处理结构体字段
			if field.Type() == reflect.TypeOf(time.Time{}) {
				// 特殊处理时间类型
				t := field.Interface().(time.Time)
				if !t.IsZero() {
					result[jsonName] = t
				}
			} else {
				// 处理其他结构体
				nestedMap := Struct2Map(field.Interface())
				for k, v := range nestedMap {
					result[jsonName+"."+k] = v // 嵌套字段使用点分隔
				}
			}
		default:
			// 处理基本类型
			if !isZeroValue(field) {
				result[jsonName] = field.Interface()
			}
		}
	}

	return result
}

// isZeroValue 判断字段是否为零值
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Array, reflect.Slice, reflect.Map:
		return v.Len() == 0
	case reflect.Struct:
		// 结构体类型的零值检测较复杂，这里简化处理
		// 需要注意，对于自定义结构体，可能需要更精确的判断
		return false
	default:
		return false
	}
}
