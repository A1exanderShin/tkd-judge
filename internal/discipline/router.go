package discipline

type Router struct {
	active Discipline
}

func (r *Router) Reset() {
	//TODO implement me
	panic("implement me")
}

func NewRouter(defaultDiscipline Discipline) *Router {
	return &Router{
		active: defaultDiscipline,
	}
}

func (r *Router) Switch(d Discipline) {
	r.active = d
}

func (r *Router) HandleEvent(event any) error {
	return r.active.HandleEvent(event)
}

func (r *Router) Snapshot() map[string]any {
	return r.active.Snapshot()
}

func (r *Router) Realtime() <-chan RealtimeEvent {
	return r.active.Realtime()
}
