package lib

func AbsInt(num int) int {
	if num >= 0 {
		return num
	}

	return -num
}

func NumberOfItemsLookupMap[T comparable](arr []T) map[T]int {
	lookupMap := make(map[T]int)
	for _, v := range arr {
		existingValue := lookupMap[v]
		lookupMap[v] = existingValue + 1
	}
	return lookupMap
}
