package chains

// Applies function for the collection.
// If chainfunc returns false, chain will be stopped
func ChainForEach[A any](collection []A, chainfunc func(int, A) bool) {
	for i, v := range collection {
		if chainfunc(i, v) {
			return
		}
	}
}

// Applies function for the collection. Returns collection of the same type
func ChainMap[A any](collection []A, chainfunc func(int, A) A) []A {
	result := make([]A, len(collection))
	for i, v := range collection {
		result[i] = chainfunc(i, v)
	}
	return result
}

// Applies function for the collection. Returns collection of another type
func ChainMapTransform[A, B any](collection []A, chainfunc func(int, A) B) []B {
	result := make([]B, len(collection))
	for i, v := range collection {
		result[i] = chainfunc(i, v)
	}
	return result
}

// Applies accumulative function to the collection
func ChainReduce[A, B any](collection []A, accumulator func(B, A) B, initialValue B) B {
	var result = initialValue
	for _, x := range collection {
		result = accumulator(result, x)
	}
	return result
}

// Searching value in the collection
func ChainExists[A comparable](collection []A, value A) bool {
	for _, v := range collection {
		if v == value {
			return true
		}
	}
	return false
}

// Searching value in the collection and returns index of element or -1 if not found
func ChainIndexOf[A comparable](collection []A, value A) int {
	for index, v := range collection {
		if v == value {
			return index
		}
	}
	return -1
}
