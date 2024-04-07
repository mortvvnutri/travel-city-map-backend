package utils

// Makes a set of unique values out of the array
func ToSetInt32(array []int32) []int32 {
	m := make(map[int32]bool)
	ret := []int32{}
	for _, v := range array {
		if !m[v] {
			// not in the array
			// add to map
			m[v] = true
			// add to output
			ret = append(ret, v)
		}
	}
	return ret
}

// Efficiently removes the last elemet of the array. Unordered.
func RemoveSingleInt32(array []int32, value int32) []int32 {
	for i, v := range array {
		if v == value {
			// replace it with the last value of the array
			array[i] = array[len(array)-1]
			// return trimmed array without the last value to avoid duplicating it
			return array[:len(array)-1]
		}
	}

	// nothing was found
	return array
}
