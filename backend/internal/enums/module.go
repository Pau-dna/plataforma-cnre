package enums

type ModuleStatus string

const (
	ModuleStatusLocked    ModuleStatus = "locked"
	ModuleStatusAvailable ModuleStatus = "available"
	ModuleStatusCompleted ModuleStatus = "completed"
)
