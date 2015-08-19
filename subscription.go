package agoni

type operationKind int

const (
	newOp operationKind = iota
	updateOp
	deleteOp
)

type Filter interface {
	match(*Operation) bool
}
type FilterOperation struct {
	kind operationKind
}
type FilterExactKey struct {
	k Key
}

// A Subscription provides a channel on which to receive KeyValueStore events.
// A list of Filters may be applied to a Subscription which must all be satisfied
// in order for the event to be dispatched on the channel.
// Multiple Subscriptions may exist for the same channel.
type Subscription struct {
	filters []Filter
	ch      chan *Operation
}

func (s *Subscription) notify(o *Operation) {
	go func() {
		if s.ch != nil {
			select {
			case s.ch <- o:

			}
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
