package utils

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"
)

var (
	StructTimeType     = reflect.TypeOf(time.Now())
	StructBsonObjectId = reflect.TypeOf(bson.NewObjectId())
)

const (
	STRUCT_DATE_TIME_FORMAT_LAYOUT = "2006-01-02 15:04:05"
	STRUCT_DATE_FORMAT_LAYOUT      = "2006-01-02"
)

type StructConvFunc func(interface{}, reflect.Type) interface{}

func CopyStructDefault(src, dest interface{}) error {
	return CopyStruct(src, dest, BaseConvert)
}

func CopyStructsDefault(src, dest interface{}) error {
	return CopyStructs(src, dest, BaseConvert)
}

//dest 必须是指针
func CopyStruct(src, dest interface{}, f StructConvFunc) error {
	ptrDestType := reflect.TypeOf(dest)
	if ptrDestType.Kind() != reflect.Ptr {
		return errors.New("dest type must be ptr")
	}

	destType := ptrDestType.Elem()

	srcValue := reflect.Indirect(reflect.ValueOf(src))
	destValue := reflect.Indirect(reflect.ValueOf(dest))

	for i := 0; i < destType.NumField(); i++ {
		dstField := destType.Field(i)

		srcField := srcValue.FieldByName(dstField.Name)
		if !srcField.IsValid() {
			continue
		}

		dstFieldValue := destValue.FieldByIndex(dstField.Index)
		if !dstFieldValue.CanSet() {
			continue
		}

		if srcField.Type() == dstField.Type {
			dstFieldValue.Set(srcField)
		} else if f != nil {
			convSrcElemField := f(srcField.Interface(), dstField.Type)
			if convSrcElemField != nil {
				dstFieldValue.Set(reflect.ValueOf(convSrcElemField))
			}
		}
	}

	return nil
}

//dest 必须是数组的指针
func CopyStructs(src, dest interface{}, f StructConvFunc) error {
	srcType := reflect.TypeOf(src)
	if srcType.Kind() != reflect.Ptr && srcType.Kind() != reflect.Slice {
		return errors.New("src type must be slice or a slice address")
	}

	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr || destType.Elem().Kind() != reflect.Slice {
		return errors.New("dest type must be a slice address")
	}

	ptrDestValue := reflect.ValueOf(dest)
	destValue := reflect.Indirect(ptrDestValue)     //.Elem()
	destValue = destValue.Slice(0, destValue.Cap()) //destValue的slice可能是nil

	srcValue := reflect.Indirect(reflect.ValueOf(src))

	destElemType := destType.Elem().Elem() //destType is slice address

	var isSliceElemPtr = false
	if destElemType.Kind() == reflect.Ptr {
		destElemType = destElemType.Elem()
		isSliceElemPtr = true
	}

	for i := 0; i < srcValue.Len(); i++ {
		srcElemValue := reflect.Indirect(srcValue.Index(i))

		destElemValuePtr := reflect.New(destElemType)
		destElemValue := reflect.Indirect(destElemValuePtr)

		for j := 0; j < destElemType.NumField(); j++ {
			dstElemTypeField := destElemType.Field(j)

			srcElemField := srcElemValue.FieldByName(dstElemTypeField.Name)
			if !srcElemField.IsValid() {
				continue
			}

			dstFieldValue := destElemValue.FieldByIndex(dstElemTypeField.Index)
			if !dstFieldValue.CanSet() {
				continue
			}

			if srcElemField.Type() == dstElemTypeField.Type {
				dstFieldValue.Set(srcElemField)
			} else if f != nil {
				convSrcElemField := f(srcElemField.Interface(), dstElemTypeField.Type)
				if convSrcElemField != nil {
					dstFieldValue.Set(reflect.ValueOf(convSrcElemField))
				}
			}
		}

		if isSliceElemPtr {
			destValue = reflect.Append(destValue, destElemValuePtr)
		} else {
			destValue = reflect.Append(destValue, destElemValue)
		}
	}

	ptrDestValue.Elem().Set(destValue.Slice(0, destValue.Len()))
	return nil
}

func MergeStructs(src, dest interface{}, f StructConvFunc, keyField string, fieldMapping ...string) error {
	srcType := reflect.TypeOf(src)
	if srcType.Kind() != reflect.Ptr && srcType.Kind() != reflect.Slice {
		return errors.New("src type must be slice or a slice address") //[] *[]
	}

	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr || destType.Elem().Kind() != reflect.Slice ||
		destType.Elem().Elem().Kind() != reflect.Ptr {
		return errors.New("dest type must be a slice of ptr address") // *[]*a
	}

	srcValue := reflect.Indirect(reflect.ValueOf(src))
	destValue := reflect.Indirect(reflect.ValueOf(dest))

	//有一个没数据不需要取合并
	if srcValue.Len() == 0 || destValue.Len() == 0 {
		return nil
	}

	var srcFName, dstFName string
	keyFields := strings.Split(keyField, ":")
	if len(keyFields) < 2 {
		return errors.New("keyField format is srcFieldName:dstFieldName")
	}
	srcFName = keyFields[0]
	dstFName = keyFields[1]

	srcElemType := reflect.Indirect(srcValue.Index(0)).Type()
	srcKeyFieldType, ok := srcElemType.FieldByName(srcFName)
	if !ok {
		return errors.New("src can not found field: " + srcFName)
	}

	destElemType := reflect.Indirect(destValue.Index(0)).Type()
	destKeyFieldType, ok := destElemType.FieldByName(dstFName)
	if !ok {
		return errors.New("dest can not found field: " + dstFName)
	}

	tupleInts := fieldMappingToIndex(srcElemType, destElemType, fieldMapping...)

	srcMap := make(map[interface{}]interface{})
	for i := 0; i < srcValue.Len(); i++ {
		srcElemValueRaw := srcValue.Index(i)
		srcElemField := reflect.Indirect(srcElemValueRaw).FieldByIndex(srcKeyFieldType.Index)
		if !srcElemField.IsValid() {
			continue
		}

		srcMap[srcElemField.Interface()] = srcElemValueRaw.Interface()
	}

	for i := 0; i < destValue.Len(); i++ {
		dstElemValuePtr := destValue.Index(i)
		dstElemField := reflect.Indirect(dstElemValuePtr).FieldByIndex(destKeyFieldType.Index)
		if !dstElemField.IsValid() {
			continue
		}

		if destKeyFieldType.Type == srcKeyFieldType.Type {
			if srcElemValue, ok := srcMap[dstElemField.Interface()]; ok {
				copyStructFields(srcElemValue, dstElemValuePtr.Interface(), f, tupleInts)
			}
		} else if f != nil {
			dstElemField2 := f(dstElemField.Interface(), srcKeyFieldType.Type)
			if dstElemField2 != nil {
				if srcElemValue, ok := srcMap[dstElemField2]; ok {
					copyStructFields(srcElemValue, dstElemValuePtr.Interface(), f, tupleInts)
				}
			}
		}
	}

	return nil
}

func fieldMappingToIndex(srcElemType, destElemType reflect.Type, fieldMapping ...string) []TupleInts {
	tupleInts := make([]TupleInts, 0)
	for _, item := range fieldMapping {
		keyFields := strings.Split(item, ":")
		if len(keyFields) < 2 {
			continue
		}

		srcFName := keyFields[0]
		dstFName := keyFields[1]

		srcField, ok := srcElemType.FieldByName(srcFName)
		if !ok {
			continue
		}

		dstField, ok := destElemType.FieldByName(dstFName)
		if !ok {
			continue
		}

		tupleInts = append(tupleInts, TupleInts{Ints1: srcField.Index, Ints2: dstField.Index})
	}
	return tupleInts
}

func copyStructFields(src, dest interface{}, f StructConvFunc, tupleInts []TupleInts) {
	srcValue := reflect.Indirect(reflect.ValueOf(src))
	destValue := reflect.Indirect(reflect.ValueOf(dest))

	for _, tupleInt := range tupleInts {
		srcField := srcValue.FieldByIndex(tupleInt.Ints1)
		if !srcField.IsValid() {
			continue
		}

		dstField := destValue.FieldByIndex(tupleInt.Ints2)
		if !dstField.IsValid() {
			continue
		}

		if !dstField.CanSet() {
			continue
		}

		if srcField.Type() == dstField.Type() {
			dstField.Set(srcField)
		} else if f != nil {
			convSrcElemField := f(srcField.Interface(), dstField.Type())
			if convSrcElemField != nil {
				dstField.Set(reflect.ValueOf(convSrcElemField))
			}
		}
	}
}

func BaseConvert(src interface{}, dstType reflect.Type) interface{} {
	if bid, ok := src.(bson.ObjectId); ok && dstType.Kind() == reflect.String {
		return bid.Hex()
	} else if bid, ok := src.(string); ok && dstType == StructBsonObjectId {
		if bid != "" {
			return bson.ObjectIdHex(bid)
		}
	} else if srcData, ok := src.(time.Time); ok && dstType.Kind() == reflect.String {
		if !srcData.IsZero() {
			return srcData.Format(STRUCT_DATE_TIME_FORMAT_LAYOUT)
		} else {
			return ""
		}
	} else if srcData, ok := src.(string); ok && dstType == StructTimeType {
		if srcData != "" {
			if len(srcData) > 10 {
				t, _ := time.ParseInLocation(STRUCT_DATE_TIME_FORMAT_LAYOUT, srcData, time.Local)
				return t
			} else {
				t, _ := time.ParseInLocation(STRUCT_DATE_FORMAT_LAYOUT, srcData, time.Local)
				return t
			}
		} else {
			return time.Time{}
		}
	}
	return nil
}
