package stats

import "github.com/ruscalworld/study-planner/internal/task"

type DisciplineTask struct {
	task.Task
	Discipline Discipline `json:"discipline" db:"discipline"`
}

type Discipline struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
