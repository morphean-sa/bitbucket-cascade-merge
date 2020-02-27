package main

import (
	"math"
	"sort"
	"strconv"
	"strings"
)

type Cascade struct {
	Branches []string
	Current  int
}

type CascadeOptions struct {
	DevelopmentName string
	ReleasePrefix   string
}

// It returns the next branch in the cascade or an empty string if it reached the end.
func (c *Cascade) Next() string {
	if len(c.Branches) > c.Current+1 {
		c.Current += 1
		return c.Branches[c.Current]
	}
	return ""
}

// Add a branch to the cascade and sort branches. If the cascade already contains a branch named identically,
// the cascade will remain unmodified.
func (c *Cascade) Append(branchName string) {
	for _, b := range c.Branches {
		if b == branchName {
			return
		}
	}
	c.Branches = append(c.Branches, branchName)
	sort.Sort(ByVersion(c.Branches))
}

// Extract an int representation of the version found in the given branch. Branch must be named accordingly to the
// following format :
//     <kind>/<version>
// The part following the slash must be an int.
// It returns the version or MaxInt32 if it comply to the format.
func extractVersion(branch string) int {
	parts := strings.Split(branch, "/")
	if len(parts) > 0 {
		version, err := strconv.Atoi(parts[len(parts)-1])
		if err == nil {
			return version
		}
	}
	return math.MaxInt32
}

type ByVersion []string

func (b ByVersion) Len() int {
	return len(b)
}

func (b ByVersion) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByVersion) Less(i, j int) bool {
	return extractVersion(b[i]) < extractVersion(b[j])
}
