package main

import (
	"time"

	"github.com/MAD-py/go-taskengine/taskengine"
)

func main() {
	// store, err := filestore.New("./storage")
	// if err != nil {
	// 	panic(err)
	// }

	engine := taskengine.New()

	task, err := taskengine.NewTask(
		"example_task",
		func(ctx *taskengine.Context) error {
			ctx.Logger().Info("Hi, I'm an example task running every 10 seconds!")
			time.Sleep(30 * time.Second)
			return nil
		},
		5*time.Second,
	)

	if err != nil {
		panic(err)
	}

	trigger, err := taskengine.NewIntervalTrigger(30*time.Second, true)
	if err != nil {
		panic(err)
	}

	engine.RegisterTask(
		task,
		taskengine.WorkerPolicyParallel,
		trigger,
		true,
		20,
	)

	engine.Run()

	// engine.Start()
	time.Sleep(1 * time.Minute)
	// engine.Shutdown()
	// println("Engine shutdown complete")
}
