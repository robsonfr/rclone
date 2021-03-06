package accounting

import (
	"sort"
	"strings"
	"sync"
)

// stringSet holds a set of strings
type stringSet struct {
	mu    sync.RWMutex
	items map[string]struct{}
}

// newStringSet creates a new empty string set of capacity size
func newStringSet(size int) *stringSet {
	return &stringSet{
		items: make(map[string]struct{}, size),
	}
}

// add adds remote to the set
func (ss *stringSet) add(remote string) {
	ss.mu.Lock()
	ss.items[remote] = struct{}{}
	ss.mu.Unlock()
}

// del removes remote from the set
func (ss *stringSet) del(remote string) {
	ss.mu.Lock()
	delete(ss.items, remote)
	ss.mu.Unlock()
}

// empty returns whether the set has any items
func (ss *stringSet) empty() bool {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	return len(ss.items) == 0
}

// Strings returns all the strings in the stringSet
func (ss *stringSet) Strings() []string {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	strings := make([]string, 0, len(ss.items))
	for name := range ss.items {
		var out string
		if acc := Stats.inProgress.get(name); acc != nil {
			out = acc.String()
		} else {
			out = name
		}
		strings = append(strings, " * "+out)
	}
	sorted := sort.StringSlice(strings)
	sorted.Sort()
	return sorted
}

// String returns all the file names in the stringSet joined by newline
func (ss *stringSet) String() string {
	return strings.Join(ss.Strings(), "\n")
}
