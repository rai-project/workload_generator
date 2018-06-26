package workload

import (
	"strings"
)

func IsValidDistribution(dist string) bool {
	dist = strings.ToLower(dist)
	for _, d := range ValidDistributions {
		if d == dist {
			return true
		}
	}
	return false
}
