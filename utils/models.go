package utils

import (
	"fmt"
	"reflect"
)

type modelType struct {
	name      string
	model     interface{}
	modelType reflect.Type // *struct
	sliceType reflect.Type // *[]struct 而不是*[]*struct
}

var gRegModels = make([]interface{}, 0)
var gModelTypes = make([]*modelType, 0)
var gModelTypeMap = make(map[string]*modelType)

//--------------------------------------------------------------------------------------------------

func GetRegModel(name string) interface{} {
	mtype, ok := gModelTypeMap[name]
	if ok {
		return mtype.model
	}

	return nil
}

func GetRegModels() []interface{} {
	return gRegModels
}

func GetRegModelType(modelName string) reflect.Type {
	mtype, ok := gModelTypeMap[modelName]
	if ok {
		return mtype.modelType
	}
	return nil
}

func GetModelNames() []string {
	var keys []string
	for k, _ := range gModelTypeMap {
		keys = append(keys, k)
	}

	return keys
}

//--------------------------------------------------------------------------------------------------

func CreateModel(name string) interface{} {
	mtype, ok := gModelTypeMap[name]
	if ok {
		return reflect.New(mtype.modelType).Interface()
	}
	return nil
}

//返回的是 结构体的数组指针 即 []*struct
func CreateModels(name string) interface{} {
	mtype, ok := gModelTypeMap[name]
	if ok {
		return reflect.New(mtype.sliceType).Interface()
	}
	return nil
}

func IsModelHasField(modelName string, fieldName string) bool {
	var result bool = false

	mtype, ok := gModelTypeMap[modelName]
	if ok {
		_, result = mtype.modelType.FieldByName(fieldName)
	}
	return result
}

//--------------------------------------------------------------------------------------------------

func registerModel(model interface{}) {
	var val = reflect.ValueOf(model)
	var typ = reflect.Indirect(val).Type()

	var typName string = typ.Name()
	_, ok := gModelTypeMap[typName]
	if ok {
		return //====>>>>
	}

	if val.Kind() != reflect.Ptr {
		panic(fmt.Errorf("<models.RegisterModels> cannot use non-ptr model struct `%s`", typ.Name()))
	}

	var sliceType = reflect.MakeSlice(reflect.SliceOf(typ), 0, 0).Type()
	var mtype = &modelType{name: typName, model: model, modelType: typ, sliceType: sliceType}

	gModelTypeMap[typName] = mtype
	gRegModels = append(gRegModels, model)
	gModelTypes = append(gModelTypes, mtype)
}

func RegisterModels(models ...interface{}) {
	for _, model := range models {
		registerModel(model)
	}
}
