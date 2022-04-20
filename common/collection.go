package common

type ListOption func(settings *ListSettings)

func WithListSize(size int) ListOption {
	return func(settings *ListSettings) {
		settings.size = size
	}
}

type ListSettings struct {
	size int
}

type List[T any] interface {
	Add(items ...T)
	Set(i int, item T)
	Get(i int) T
	Len() int
	Entries() []T
	Iterator() Iterator[T]
	Remove(i int)
	RemoveAll()
}

type genericList[T any] []T

func NewList[T any](opts ...ListOption) List[T] {
	settings := new(ListSettings)
	for _, opt := range opts {
		opt(settings)
	}

	l := make(genericList[T], settings.size)
	return &l
}

func NewListFromSlice[T any](slice ...T) List[T] {
	l := make(genericList[T], len(slice))
	copy(l, slice)
	return &l
}

func (l *genericList[T]) Add(items ...T) {
	(*l) = append((*l), items...)
}

func (l *genericList[T]) Set(i int, item T) {
	(*l)[i] = item
}

func (l *genericList[T]) Get(i int) T {
	return (*l)[i]
}

func (l *genericList[T]) Len() int {
	return len(*l)
}

func (l *genericList[T]) Entries() []T {
	return (*l)
}

func (l *genericList[T]) Iterator() Iterator[T] {
	return NewListIterator[T](l)
}

func (l *genericList[T]) Remove(i int) {
	(*l) = append((*l)[:i], (*l)[i+1:]...)
}

func (l *genericList[T]) RemoveAll() {
	(*l) = make(genericList[T], 0)
}

type Map[K comparable, V any] interface {
	Set(key K, value V)
	Get(key K) (V, bool)
	Len() int
	Entries() map[K]V
	Iterator() Iterator[V]
	Remove(key K)
	RemoveAll()
}

func NewMap[K comparable, V any]() Map[K, V] {
	return make(genericMap[K, V])
}

type genericMap[K comparable, V any] map[K]V

func (m genericMap[K, V]) Set(key K, value V) {
	m[key] = value
}

func (m genericMap[K, V]) Get(key K) (V, bool) {
	value, ok := m[key]
	return value, ok
}

func (m genericMap[K, V]) Len() int {
	return len(m)
}

func (m genericMap[K, V]) Entries() map[K]V {
	return m
}

func (m genericMap[K, V]) Iterator() Iterator[V] {
	return NewMapIterator[K, V](m)
}

func (m genericMap[K, V]) Remove(key K) {
	delete(m, key)
}

func (m genericMap[K, V]) RemoveAll() {
	for key, _ := range m {
		m.Remove(key)
	}
}
