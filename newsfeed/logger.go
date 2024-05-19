package main

import (
	"reflect"
	"runtime"
	"strings"
)

func getFunctionName(i interface{}) string {

    // Get caller function details
    pc, _, _, ok := runtime.Caller(1)
    if !ok {
        return "unknown"
    }

    // Get function name

    fullName := runtime.FuncForPC(pc).Name()
    parts := strings.Split(fullName, ".")
    length := len(parts)
    
    if length < 2 {
        return fullName
    }
    
    // Assume second last part is the name
    
    structType := reflect.TypeOf(i)
    if structType.Kind() == reflect.Ptr {
        structType = structType.Elem()
    }
    
    return structType.Name() + "." + parts[length-1]
}