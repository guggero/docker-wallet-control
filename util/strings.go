package util

import "reflect"

func ArrayContains(arr interface{}, item interface{}) (bool) {
    arrayValue := reflect.ValueOf(arr)

    if arrayValue.Kind() == reflect.Slice {
        for i := 0; i < arrayValue.Len(); i++ {
            if arrayValue.Index(i).Interface() == item {
                return true
            }
        }
    }
    return false
}