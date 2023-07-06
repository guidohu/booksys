package database

type SessionByTime []Session

func (a SessionByTime) Len() int           { return len(a) }
func (a SessionByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SessionByTime) Less(i, j int) bool { return a[i].StartTime.Before(a[j].StartTime) }
