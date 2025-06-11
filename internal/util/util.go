package util

import (
	"strconv"
	"strings"
	"time"
)

type UnixTimestamp time.Time

func (u *UnixTimestamp) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	unixTime, err := FromUnixTimestamp(s)
	if err != nil {
		return err
	}
	*u = UnixTimestamp(unixTime)
	return nil
}

func (u *UnixTimestamp) Time() time.Time {
	return time.Time(*u)
}

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
