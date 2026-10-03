package listener

import "github.com/google/uuid"

type Listener struct {
	id      uuid.UUID
	options map[string]string
	proto   string
}

func NewListener(proto string, options map[string]string) Listener {
	return Listener{
		id:      uuid.New(),
		options: options,
		proto:   proto,
	}
}

func (l Listener) ID() uuid.UUID {
	return l.id
}

func (l Listener) Options() map[string]string {
	return l.options
}

func (l Listener) Protocol() string {
	return l.proto
}

// Update replaces the listener's options. Uses a pointer receiver so the
// mutation persists; the previous value receiver silently discarded it.
func (l *Listener) Update(options map[string]string) {
	l.options = options
}
