package main

import (
	"net/mail"
)

// Requires input to be sorted
func IntersectionOf(a, b []string) []string {

	out := []string{}
	i, j := 0, 0

	for i < len(a) && j < len(b) {
		if a[i] == b[j] {
			out = append(out, a[i])
			i++
			j++
		} else if a[i] < b[j] {
			i++
		} else {
			j++
		}
	}

	return out
}

// Returns the Address portion of the input addresses
func AddressesOnly(full []*mail.Address) []string {
	addrs := []string{}
	for _, addr := range full {
		addrs = append(addrs, addr.Address)
	}
	return addrs
}
