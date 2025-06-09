package util

import (
	"strconv"
	"time"
)

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

func FromUnixTimestamp(s string) (time.Time, error) {
	seconds, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	unixTime := time.Unix(seconds, 0)
	return unixTime, nil
}
