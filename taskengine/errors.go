package taskengine

import "errors"

var (
	ErrorPolicyMismatch        = errors.New("policy mismatch")
	ErrorJobNameMismatch       = errors.New("job name mismatch")
	ErrorTriggerMismatch       = errors.New("trigger mismatch")
	ErrorTaskAlreadyRegistered = errors.New("task is already registered")
)
