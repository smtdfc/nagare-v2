package error_utils

var (
	ErrAgentStateEmpty        = NewError("Agent state empty", 400)
	ErrAgentLoopMaxIterations = NewError("The agent has reached the maximum number of iterations", 400)
)
