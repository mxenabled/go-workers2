package workers

type Action interface {
	Call(queue string, message interface{}, next func())
}

type middleware struct {
	actions []Action
}

func (m *middleware) Append(action Action) {
	m.actions = append(m.actions, action)
}

func (m *middleware) Prepend(action Action) {
	actions := make([]Action, len(m.actions)+1)
	actions[0] = action
	copy(actions[1:], m.actions)
	m.actions = actions
}

func (m *middleware) call(queue string, message interface{}, final func()) {
	newContinuation(m.actions, queue, message, final)()
}

func newContinuation(actions []Action, queue string, message interface{}, final func()) func() {
	return func() {
		if len(actions) > 0 {
			actions[0].Call(
				queue,
				message,
				newContinuation(actions[1:], queue, message, final),
			)
		} else {
			final()
		}
	}
}

func newMiddleware(actions ...Action) *middleware {
	return &middleware{
		actions,
	}
}