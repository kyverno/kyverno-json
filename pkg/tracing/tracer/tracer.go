package tracer

type Tracer struct {
	traces []string
}

func (t *Tracer) Trace(trace string) {
	t.traces = append(t.traces, trace)
}

func (t *Tracer) Traces() []string {
	return t.traces
}
