package taskengine

import (
	"context"
	"time"
)

type contextKey string

const ContextKey contextKey = "go-taskengine/contextKey"

var _ context.Context = (*Context)(nil)

type Context struct {
	ctx context.Context

	tick *Tick

	taskName string

	logger Logger
}

// ====================================
// ======== [ Execution Info ] ========
// ====================================

func (c *Context) Logger() Logger { return c.logger }

func (c *Context) TaskName() string { return c.taskName }

func (c *Context) LastTick() time.Time { return c.tick.lastTick }

func (c *Context) CurrentTick() time.Time { return c.tick.currentTick }

// ====================================
// ====== [ Context Interface ] =======
// ====================================

func (c *Context) Deadline() (deadline time.Time, ok bool) { return c.ctx.Deadline() }

func (c *Context) Done() <-chan struct{} { return c.ctx.Done() }

func (c *Context) Err() error { return c.ctx.Err() }

func (c *Context) Value(key any) any {
	if key == ContextKey {
		return c
	}

	if c.ctx == nil {
		return nil
	}
	return c.ctx.Value(key)
}
