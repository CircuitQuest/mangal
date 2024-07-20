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

type UserHistory []string

// Size gets the amount of users in the history.
func (u *UserHistory) Size() int {
	if u == nil {
		return 0
	}
	return len(*u)
}

// Last gets the last user in the history.
func (u *UserHistory) Last() string {
	if u == nil {
		return ""
	}
	if len(*u) != 0 {
		return (*u)[len(*u)-1]
	}
	return ""
}

// Get the user history as a slice of strings.
func (u *UserHistory) Get() []string {
	if u == nil {
		return []string{}
	}
	return []string(*u)
}

// Add will add the new username to the user history,
// if it exists, it will move it to the end of the history.
func (u *UserHistory) Add(username string) {
	if u == nil {
		u = &UserHistory{username}
		return
	}
	var newHistory UserHistory
	for _, user := range *u {
		if user != username {
			newHistory = append(newHistory, user)
		}
	}
	newHistory = append(newHistory, username)
	*u = newHistory
}

// Delete will delete the username from the user history, if existent.
func (u *UserHistory) Delete(username string) {
	if u == nil {
		return
	}
	var newHistory UserHistory
	for _, user := range *u {
		if user != username {
			newHistory = append(newHistory, user)
		}
	}
	*u = newHistory
}
