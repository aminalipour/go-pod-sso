package pkg

import "reflect"

func isZeroValue(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.String() == ""
	case reflect.Int, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0.0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Slice, reflect.Array:
		return value.Len() == 0
	default:
		return value.IsZero()
	}
}
