package gmime

/*
#cgo pkg-config: gmime-3.0
#include <stdlib.h>
#include <gmime/gmime.h>
*/
import "C"

type GMimeParamsCallback func(name string, value string)

type Parametrized interface {
	SetParameter(name, value string)
	Parameter(name string) string
	ForEachParam(callback GMimeParamsCallback)
}

func forEachParam(params *C.GMimeParamList, callback GMimeParamsCallback) {
	i := 0
	for {
		param := C.g_mime_param_list_get_parameter_at(params, C.int(i))
		if param == nil {
			return
		}
		cName := C.g_mime_param_get_name(param)
		name := C.GoString(cName)
		cValue := C.g_mime_param_get_value(param)
		value := C.GoString(cValue)
		callback(name, value)
		i ++
	}
}
