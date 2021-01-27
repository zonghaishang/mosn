package wasmer

import (
	"reflect"

	wasmerGo "github.com/wasmerio/wasmer-go/wasmer"
	"mosn.io/mosn/pkg/log"
)

func convertFromGoType(t reflect.Type) *wasmerGo.ValueType {
	switch t.Kind() {
	case reflect.Int32:
		return wasmerGo.NewValueType(wasmerGo.I32)
	case reflect.Int64:
		return wasmerGo.NewValueType(wasmerGo.I64)
	case reflect.Float32:
		return wasmerGo.NewValueType(wasmerGo.F32)
	case reflect.Float64:
		return wasmerGo.NewValueType(wasmerGo.F64)
	default:
		log.DefaultLogger.Errorf("[wasmer][type] convertFromGoType unsupported type: %v", t.Kind().String())
	}
	return nil
}

func convertToGoTypes(in []wasmerGo.Value) []reflect.Value {
	res := make([]reflect.Value, len(in))
	for i, v := range in {
		switch v.Kind() {
		case wasmerGo.I32:
			res[i] = reflect.ValueOf(v.I32())
		case wasmerGo.I64:
			res[i] = reflect.ValueOf(v.I64())
		case wasmerGo.F32:
			res[i] = reflect.ValueOf(v.F32())
		case wasmerGo.F64:
			res[i] = reflect.ValueOf(v.F64())
		}
	}
	return res
}

func convertFromGoValue(val reflect.Value) wasmerGo.Value {
	switch val.Kind() {
	case reflect.Int32:
		return wasmerGo.NewI32(int32(val.Int()))
	case reflect.Int64:
		return wasmerGo.NewI64(int64(val.Int()))
	case reflect.Float32:
		return wasmerGo.NewF32(float32(val.Float()))
	case reflect.Float64:
		return wasmerGo.NewF64(float64(val.Float()))
	default:
		log.DefaultLogger.Errorf("[wasmer][type] convertFromGoValue unsupported val type: %v", val.Kind().String())
	}
	return wasmerGo.Value{}
}