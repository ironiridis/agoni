package agoni

type operationKind int

const (
	newOp operationKind = iota
	updateOp
	deleteOp
)

// A Filter is a condition the consumer of a Subscription defines. All of the
// defined Filters for a Subscription must be met for a notification to fire.
type Filter interface {
	match(*Operation) bool
}

// A FilterOperation is a Filter that matches on whether the operation is a
// NewOperation, an UpdateOperation, or a DeleteOperation.
type FilterOperation struct {
	kind operationKind
}

// A FilterExactKey is a Filter that matches on the key being operated on.
type FilterExactKey struct {
	k Key
}

// A Subscription provides a channel on which to receive KeyValueStore events.
// A list of Filters may be applied to a Subscription which must all be
// satisfied in order for the event to be dispatched on the channel.
// Multiple Subscriptions may exist for the same channel.
type Subscription struct {
	filters []Filter
	ch      chan *Operation
	ctl     chan int
}

// Destroy will signal all pending notify goroutines to exit.
func (s *Subscription) Destroy() {
	close(s.ctl)
}

func (s *Subscription) notify(o *Operation) {
	go func() {
		// If we destroy a Subscription, s.ch may not be reading its end anymore.
		// We can close s.ctl to make it always readable so goroutines don't leak.
		select {
		case <-s.ctl: // if readable, Subscription has been destroyed
		case s.ch <- o:
		}
	}()
	return
}

func (s *Subscription) match(o *Operation) bool {
	for i := range s.filters {
		if !s.filters[i].match(o) {
			return false
		}
	}
	return true
}

func (f *FilterOperation) match(o *Operation) bool {
	switch (*o).(type) {
	case *NewOperation:
		return (f.kind == newOp)
	case *UpdateOperation:
		return (f.kind == updateOp)
	case *DeleteOperation:
		return (f.kind == deleteOp)
	}
	return false
}

func (f *FilterExactKey) match(o *Operation) bool {
	return f.k.Compare((*o).Key())
}
