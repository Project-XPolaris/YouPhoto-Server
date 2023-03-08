package service

type Signal struct {
}

const (
	TaskTypeScanLibrary = "ScanLibrary"
	TaskTypeRemove      = "RemoveLibrary"
)
const (
	TaskStatusRunning = "Running"
	TaskStatusDone    = "Done"
	TaskStatusError   = "Error"
	TaskStatusStop    = "Stop"
)
