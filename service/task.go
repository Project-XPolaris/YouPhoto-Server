package service

import (
	"fmt"
	. "github.com/ahmetb/go-linq/v3"
	youlogcore "github.com/project-xpolaris/youplustoolkit/youlog"
	"github.com/projectxpolaris/youphoto/youlog"
	"github.com/rs/xid"
	"sync"
	"time"
)

type Signal struct {
}

const (
	TaskTypeScanLibrary = iota + 1
	TaskTypeRemove
)
const (
	TaskStatusRunning = iota + 1
	TaskStatusDone
	TaskStatusError
)

var TaskLogger = youlog.DefaultYouLogPlugin.Logger.NewScope("Task")
var TaskTypeNameMapping map[int]string = map[int]string{
	TaskTypeScanLibrary: "ScanLibrary",
	TaskTypeRemove:      "RemoveLibrary",
}

var TaskStatusNameMapping map[int]string = map[int]string{
	TaskStatusRunning: "Running",
	TaskStatusDone:    "Done",
	TaskStatusError:   "Error",
}
var DefaultTaskPool TaskPool = TaskPool{
	Tasks: []Task{},
}

type Task interface {
	GetId() string
	GetType() int
	GetStatus() int
	GetOutput() interface{}
	GetStartTime() time.Time
}
type BaseTask struct {
	Id        string
	Type      int
	Error     error
	Status    int
	Logger    *youlogcore.Scope
	StartTime time.Time
}

func (t *BaseTask) GetId() string {
	return t.Id
}

func (t *BaseTask) GetType() int {
	return t.Type
}
func (t *BaseTask) GetStatus() int {
	return t.Status
}
func (t *BaseTask) GetStartTime() time.Time {
	return t.StartTime
}
func (t *BaseTask) AbortError(err error) {
	t.Logger.Error(err.Error())
	if err != nil {
		t.Status = TaskStatusError
	}
	t.Error = err
}

func (t *BaseTask) UpdateDoneStatus() {
	t.Status = TaskStatusDone
}
func NewBaseTask(taskType int) BaseTask {
	id := xid.New().String()
	t := BaseTask{
		Id:        id,
		Type:      taskType,
		Logger:    youlog.DefaultYouLogPlugin.Logger.NewScope(fmt.Sprintf("%s-%s", TaskStatusNameMapping[taskType], id)),
		StartTime: time.Now(),
		Status:    TaskStatusRunning,
	}
	return t
}

type TaskPool struct {
	sync.Mutex
	Tasks []Task
}

func (p *TaskPool) RemoveTaskById(id string) {
	p.Lock()
	defer p.Unlock()
	var newTask []Task
	From(p.Tasks).WhereT(func(task Task) bool {
		return task.GetId() != id
	}).ToSlice(&newTask)
	p.Tasks = newTask
}

func (p *TaskPool) AddTask(task Task) {
	p.Lock()
	defer p.Unlock()
	newTasks := make([]Task, 0)
	newTasks = append(newTasks, task)
	newTasks = append(newTasks, p.Tasks...)
	p.Tasks = newTasks
}
func GetTaskList() []Task {
	return DefaultTaskPool.Tasks
}
