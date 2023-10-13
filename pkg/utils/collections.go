package utils

import (
	"sort"
	"strings"
)

// ToLowercaseUniqueSorted transforms a list of strings in a list of sorted, lower case and unique strings.
func ToLowercaseUniqueSorted(list []string) []string {
	mapLowerUnique := make(map[string]struct{})
	lowerUniqueList := []string{}
	for _, s := range list {
		lower := strings.ToLower(s)
		_, ok := mapLowerUnique[lower]
		if !ok {
			mapLowerUnique[lower] = struct{}{}
			lowerUniqueList = append(lowerUniqueList, lower)
		}
	}
	sort.Strings(lowerUniqueList)
	return lowerUniqueList
}
