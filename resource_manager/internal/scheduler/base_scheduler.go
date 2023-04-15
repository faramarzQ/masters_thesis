package scheduler

type baseScheduler struct {
	state map[string]string
}

func newBaseScheduler() *baseScheduler {
	return &baseScheduler{map[string]string{}}
}

func (bs *baseScheduler) set(key string, value string) {
	if bs.state == nil {
		bs.state = map[string]string{}
	}

	bs.state[key] = value
}

func (bs *baseScheduler) get(key string) string {
	return bs.state[key]
}
