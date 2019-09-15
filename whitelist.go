package fh4server

// Whitelist is used to represent a set of strings that are allowed.
type Whitelist = func(string) bool

func constant(val bool) Whitelist {
	return func(string) bool {
		return val
	}
}

// AllowAll returns a Whitelist function which accepts any string.
func AllowAll() Whitelist {
	return constant(true)
}

// BanAll returns a Whitelist function which rejects all strings.
func BanAll() Whitelist {
	return constant(false)
}

// AllowList returns a Whitelist function which accepts any of the strings in
// `labels`
func AllowList(labels []string) Whitelist {
	set := make(map[string]struct{})

	for _, label := range labels {
		set[label] = struct{}{}
	}

	return func(label string) bool {
		_, ok := set[label]
		return ok
	}
}
