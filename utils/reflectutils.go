package utils

import "reflect"

//func main() {
//	var m = map[string]interface{}{"one": 2222}
//	mv := reflect.ValueOf(m)
//	value := mv.MapIndex(reflect.ValueOf("one"))

//	fmt.Println(value.Kind())							//interface
//	fmt.Println(reflect.ValueOf(m["one"]).Kind())		//int
//	fmt.Println(AnyIndirect(value).Kind())				//int
//}

func AnyIndirect(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Ptr && v.Kind() != reflect.Interface {
		return v
	}
	return v.Elem()
}

func SafeIsNil(value *reflect.Value) bool {
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return value.IsNil()
	default:
		return false
	}
}
