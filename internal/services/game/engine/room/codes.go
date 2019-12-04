package room

// Codes of events
const (
	UpdateStatus     = 0
	UpdatePeople     = 1
	UpdateConnection = 2
)

// Game status
const (
	StatusRecruitment = 0
	StatusFlagPlacing = 1
	StatusRunning     = 2
	StatusFinished    = 3
	StatusAborted     = 4
	StatusHistory     = 5
)

// People status
const (
	AllExit       = 0
	Full          = 1
	AllDied       = 2
	PlayerEnter   = 3
	ObserverEnter = 4
)
