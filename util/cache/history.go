package cache

import "sort"

// Record is a single history record.
type Record struct {
	// Rank how high should this be suggested.
	Rank int

	// Query is the search query.
	Query string
}

// Records is a slice of records with convenience methods.
type Records []*Record

// Sort the records, in descending order based on their Rank.
func (r *Records) Sort() {
	if r == nil {
		return
	}
	sort.SliceStable(*r, func(i, j int) bool {
		return (*r)[i].Rank > (*r)[j].Rank
	})
}

// Get the records as a slice of strings.
func (r Records) Get() []string {
	s := make([]string, len(r))
	for i, rec := range r {
		s[i] = rec.Query
	}
	return s
}

// Add a new record. If it exists, increase its Rank by 1.
func (r *Records) Add(query string) {
	if r == nil {
		r = &Records{&Record{Rank: 0, Query: query}}
		return
	}

	found := false
	for _, rec := range *r {
		if rec.Query == query {
			found = true
			rec.Rank++
			break
		}
	}
	if !found {
		*r = append(*r, &Record{Rank: 0, Query: query})
	}
}
