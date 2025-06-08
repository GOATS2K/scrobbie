package util

func GetItemOrNil[T any](c []T, predicate func(*T) bool) *T {
	for _, item := range c {
		if predicate(&item) {
			return &item
		}
	}

	return nil
}

func Filter[T any](c []T, predicate func(*T) bool) []*T {
	var resultList []*T
	for _, item := range c {
		if predicate(&item) {
			resultList = append(resultList, &item)
		}
	}

	return resultList
}
