package git

// NoTransform returns a git sha without any transformation on it. This function exists as a pair with ShortShaTransform. It is simply the identity function
func NoTransform(s string) string { return s }

// ShortSHATransform returns the shortened version of the given SHA given in s
func ShortSHATransform(s string) string { return s[:7] }
