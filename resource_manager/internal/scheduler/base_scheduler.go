package scheduler

type baseScheduler struct {
	state map[string]any
}

func newBaseScheduler() *baseScheduler {
	return &baseScheduler{map[string]any{}}
}

func (bs *baseScheduler) set(key string, value any) {
	if bs.state == nil {
		bs.state = map[string]any{}
	}

	bs.state[key] = value
}

func (bs *baseScheduler) get(key string) any {
	return bs.state[key]
}
