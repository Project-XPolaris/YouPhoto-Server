package module

import "github.com/allentom/harukap/module/task"

var Task *task.TaskModule

func CreateTaskModule() {
	Task = task.NewTaskModule()
}
