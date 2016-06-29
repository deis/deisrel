package git

import (
	"errors"
)

var (
	// ErrSHANotLongEnough is the error returned from various functions in this package when a
	// given SHA is not long enough for use
	ErrSHANotLongEnough = errors.New("SHA not long enough")
)

// NoTransform returns a git sha without any transformation on it. This function exists as a
// pair with ShortShaTransform. It's purposefully the identity function
func NoTransform(s string) string { return s }

// ShortSHATransform returns the shortened version of the given SHA given in s. If the given
// string is not long enough, returns the empty string and ErrSHANotLongEnough
func ShortSHATransform(s string) (string, error) {
	if len(s) < 7 {
		return "", ErrSHANotLongEnough
	}
	return s[:7], nil
}

// ShortSHATransformNoErr calls ShortSHATransform. If that call returns an error, return s.
// Otherwise return its string return value
func ShortSHATransformNoErr(s string) string {
	ret, err := ShortSHATransform(s)
	if err != nil {
		return s
	}
	return ret
}
