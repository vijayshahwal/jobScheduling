package interfaces

type Schedule interface {
	CalculateNextRun() int64
}

type ScheduleProcessor interface {
	CanProcess(schedule interface{}) bool
	Process(schedule interface{}, jobKey string) error
}
