package chains

func ChainForEach[A any](collection []A, chainfunc func(int, A)) {
	for i, v := range collection {
		chainfunc(i, v)
	}
}

func ChainMap[A any](collection []A, chainfunc func(int, A) A) []A {
	result := make([]A, len(collection))
	for i, v := range collection {
		result[i] = chainfunc(i, v)
	}
	return result
}

func ChainTransform[A any, B any](collection []A, chainfunc func(int, A) B) []B {
	result := make([]B, len(collection))
	for i, v := range collection {
		result[i] = chainfunc(i, v)
	}
	return result
}

func ChainReduce[A, B any](collection []A, accumulator func(B, A) B, initialValue B) B {
	var result = initialValue
	for _, x := range collection {
		result = accumulator(result, x)
	}
	return result
}

func ChainExists[A comparable](collection []A, value A) bool {
	for _, v := range collection {
		if v == value {
			return true
		}
	}
	return false
}
