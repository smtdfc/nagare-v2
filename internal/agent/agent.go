package agent

import "github.com/smtdfc/nagare/pkgs/logger"

type Agent struct {
	State  *State
	logger *logger.Logger
}

func (a *Agent) WithState(state *State) *Agent {
	a.State = state
	return a
}

func NewAgent(l *logger.Logger) *Agent {
	if l == nil {
		l = logger.NewLogger()
	}

	return &Agent{
		State: nil,

		logger: l,
	}
}
