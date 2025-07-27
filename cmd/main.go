package main

import (
	"database/sql"
	"time"

	"github.com/MAD-py/go-taskengine/taskengine"
	"github.com/MAD-py/go-taskengine/taskengine/store/postgresql"
)

func main() {
	db, err := sql.Open("postgres", "")
	if err != nil {
		panic(err)
	}

	store := postgresql.NewStore(db)
	engine := taskengine.New(store)

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

	// trigger, err := taskengine.NewIntervalTrigger(30*time.Second, true)
	trigger, err := taskengine.NewCronTrigger("* * * * *")
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
	time.Sleep(2 * time.Minute)
	// engine.Shutdown()
	// println("Engine shutdown complete")
}
