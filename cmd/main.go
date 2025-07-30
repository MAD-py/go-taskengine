package main

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver

	"github.com/MAD-py/go-taskengine/taskengine"
	"github.com/MAD-py/go-taskengine/taskengine/store/postgresql"
)

func Test(ctx *taskengine.Context) error {
	ctx.Logger().Info("Hi, I'm an example task running every 10 seconds!")
	time.Sleep(3 * time.Second)
	return nil
}

func main() {
	db, err := sql.Open(
		"postgres",
		"postgres://user:password@localhost:5432/taskengine?sslmode=disable",
	)
	if err != nil {
		panic(err)
	}

	store := postgresql.NewStore(db)

	// store.DeleteStores(context.Background())

	engine, err := taskengine.New(store)
	if err != nil {
		panic(err)
	}

	task, err := taskengine.NewTask("example_task", Test)
	if err != nil {
		panic(err)
	}

	// trigger, err := taskengine.NewIntervalTrigger(30*time.Second, true)
	trigger, err := taskengine.NewCronTrigger("* * * * *")
	if err != nil {
		panic(err)
	}

	err = engine.RegisterTask(
		task,
		taskengine.WorkerPolicySerial,
		trigger,
		true,
		20,
	)

	if err != nil {
		panic(err)
	}

	engine.Run()

	// engine.Start()
	// time.Sleep(2 * time.Minute)
	// engine.Shutdown()
	// println("Engine shutdown complete")
}
