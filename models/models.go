package models

type Coordinates struct {
	IsFollowMouse bool
	Interval      int
}

type TempWait struct {
	SleepVal int
	IsKill   bool
}
