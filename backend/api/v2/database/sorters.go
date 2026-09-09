package database

// SessionByTime sorts sessions by their start time, earliest first.
// It implements sort.Interface.
type SessionByTime []Session

// Len implements sort.Interface.
func (a SessionByTime) Len() int { return len(a) }

// Swap implements sort.Interface.
func (a SessionByTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less implements sort.Interface.
func (a SessionByTime) Less(i, j int) bool { return a[i].StartTime.Before(a[j].StartTime) }

// SessionByTimeDesc sorts sessions by their start time, latest first.
// It implements sort.Interface.
type SessionByTimeDesc []Session

// Len implements sort.Interface.
func (a SessionByTimeDesc) Len() int { return len(a) }

// Swap implements sort.Interface.
func (a SessionByTimeDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less implements sort.Interface.
func (a SessionByTimeDesc) Less(i, j int) bool { return a[i].StartTime.After(a[j].StartTime) }
