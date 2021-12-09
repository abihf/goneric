package collection

type Set[V any] struct {
	m Map[V, bool]
}

func (s *Set[V]) Delete(value V) {
	s.m.Delete(value)
}
func (s *Set[V]) Has(value V) (ok bool) {
	_, ok = s.m.Load(value)
	return
}

func (s *Set[V]) Range(f func(value V) bool) {
	s.m.Range(func(v V, _ bool) bool {
		return f(v)
	})
}
func (s *Set[V]) Add(value V) {
	s.m.Store(value, true)
}
