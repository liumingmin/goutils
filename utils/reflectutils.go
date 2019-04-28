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
