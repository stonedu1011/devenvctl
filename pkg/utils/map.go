package utils

// OrderedMap map with a track or entry order.
// Not concurrent safe
// Not designed for high performance
// Not for frequent set and delete
type OrderedMap[K comparable, V any] struct {
	store map[K]V
	keys  []K
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return NewOrderedMapWithCap[K, V](10)
}

func NewOrderedMapWithCap[K comparable, V any](cap int) *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		store: map[K]V{},
		keys:  make([]K, 0, cap),
	}
}

func (m *OrderedMap[K, V]) Len() int {
	return len(m.store)
}

func (m *OrderedMap[K, V]) Get(k K) V {
	v, _ := m.store[k]
	return v
}

func (m *OrderedMap[K, V]) Lookup(k K) (v V, ok bool) {
	v, ok = m.store[k]
	return
}

func (m *OrderedMap[K, V]) Set(k K, v V) {
	if _, ok := m.store[k]; !ok {
		m.keys = append(m.keys, k)
	}
	m.store[k] = v
}

func (m *OrderedMap[K, V]) Delete(k K) bool {
	if _, ok := m.store[k]; !ok {
		return false
	}
	var i int
	for i = range m.keys {
		if m.keys[i] == k {
			break
		}
	}
	if i >= len(m.keys) {
		return false
	}
	m.keys = append(m.keys[:i], m.keys[i+1:]...)
	delete(m.store, k)
	return true
}

func (m *OrderedMap[K, V]) Keys() []K {
	return m.keys
}
