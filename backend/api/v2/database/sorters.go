package database

type SessionByTime []Session

func (a SessionByTime) Len() int           { return len(a) }
func (a SessionByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SessionByTime) Less(i, j int) bool { return a[i].StartTime.Before(a[j].StartTime) }

type SessionByTimeDesc []Session

func (a SessionByTimeDesc) Len() int           { return len(a) }
func (a SessionByTimeDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SessionByTimeDesc) Less(i, j int) bool { return a[i].StartTime.After(a[j].StartTime) }
