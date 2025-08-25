package timewheel

import (
	"github.com/Cooooing/cutil/common/logger"
	"testing"
	"time"
)

func TestTimewheel(t *testing.T) {
	tw := NewTimewheel(time.Millisecond*100, 100, 10, 100)
	tw.Start()
	defer tw.Stop()
	task := &Task{
		Key:      "task1",
		Interval: time.Millisecond * 500,
		Times:    -1,
		Job: func(task *Task) {
			logger.Info("%s run... ", task.Key)
		},
	}
	err := tw.AddTask(task)

	if err != nil {
		t.Errorf("add task err: %v", err)
	}
	time.Sleep(time.Second * 10)
}
