package app

import (
	"sync"
	"time"
)

type Scheduled struct {
	EveryReboot []func(app *App)
	EverySecond []func(app *App)
	EveryMinute []func(app *App)
	EveryHour   []func(app *App)
	EveryDay    []func(app *App)
	EveryWeek   []func(app *App)
	EveryMonth  []func(app *App)
	EveryYear   []func(app *App)
}

type Task struct {
	Interval time.Duration
	Funcs    []func(app *App)
}

func RunScheduleTasks(app *App, poolSize int, stop <-chan struct{}) {

	for _, f := range app.ScheduledTasks.EveryReboot {
		f(app)
	}
	tasks := []Task{
		{Interval: time.Second, Funcs: app.ScheduledTasks.EverySecond},
		{Interval: time.Minute, Funcs: app.ScheduledTasks.EveryMinute},
		{Interval: time.Hour, Funcs: app.ScheduledTasks.EveryHour},
		{Interval: 24 * time.Hour, Funcs: app.ScheduledTasks.EveryDay},
		{Interval: 7 * 24 * time.Hour, Funcs: app.ScheduledTasks.EveryWeek},
		{Interval: 30 * 24 * time.Hour, Funcs: app.ScheduledTasks.EveryMonth},
		{Interval: 365 * 24 * time.Hour, Funcs: app.ScheduledTasks.EveryYear},
	}

	var wg sync.WaitGroup

	runners := make([]chan bool, len(tasks))

	for i, task := range tasks {
		runner := make(chan bool, poolSize)
		runners[i] = runner
		wg.Add(i)

		go func(task Task, runner chan bool) {
			defer wg.Done()
			ticker := time.NewTicker(task.Interval)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					for _, f := range task.Funcs {
						runner <- true
						go func(f func(app *App)) {
							defer func() { <-runner }()
							f(app)
						}(f)
					}
				case <-stop:
					return
				}
			}

		}(task, runner)
	}
	wg.Wait()

	for _, runner := range runners {
		close(runner)
	}

}
