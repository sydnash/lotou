package core

func StartService(name string, d int, m Module) uint {
	s := newService(name)
	s.m = m
	id := registerService(s)
	m.SetService(s)
	m.OnInit()
	if d > 0 {
		s.runWithLoop(d)
	} else {
		s.run()
	}
	return id
}

func Start(isMulti, isMaster bool) {
}
