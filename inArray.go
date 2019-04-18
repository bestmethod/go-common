package gocommon

import "reflect"

// provided an array and an element, checks if element exists in array by performing a reflect.DeepEqual comparison
// returns first matching element index, or -1 if not found
func inArray(array interface{}, element interface{}) (index int) {
	index = -1
	s := reflect.ValueOf(array)
	for i := 0; i < s.Len(); i++ {
		if reflect.DeepEqual(element, s.Index(i).Interface()) == true {
			index = i
			return
		}
	}
	return
}
