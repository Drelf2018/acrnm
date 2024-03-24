package acrnm

import (
	mapset "github.com/deckarep/golang-set/v2"
)

type Value[K comparable, V any] interface {
	MapKey() K
	Equal(V) bool
}

type Map[K comparable, V Value[K, V]] struct {
	m      map[K]V
	alive  mapset.Set[K]
	victim mapset.Set[K]
}

func (m *Map[K, V]) Len() int {
	return len(m.m)
}

func (m *Map[K, V]) Updates(vals []V) (newVals, diffVals, delVals []V) {
	if 2*len(vals) <= m.Len() {
		return
	}
	keys := make([]K, 0, len(vals))
	for _, val := range vals {
		k := val.MapKey()
		keys = append(keys, k)
		if m.victim.ContainsOne(k) {
			m.victim.Remove(k)
			m.alive.Add(k)
		} else if !m.alive.ContainsOne(k) {
			m.alive.Add(k)
			newVals = append(newVals, val)
			m.m[k] = val
			continue
		}
		if !m.m[k].Equal(val) {
			diffVals = append(diffVals, val)
			m.m[k] = val
		}
	}
	m.victim.Each(func(k K) bool {
		delVals = append(delVals, m.m[k])
		delete(m.m, k)
		return false
	})
	newSet := mapset.NewSet(keys...)
	m.victim = m.alive.Difference(newSet)
	m.alive = newSet
	return
}

func NewMap[K comparable, V Value[K, V]](vals ...V) (m *Map[K, V]) {
	m = &Map[K, V]{
		m:      make(map[K]V),
		alive:  mapset.NewSet[K](),
		victim: mapset.NewSet[K](),
	}
	for _, val := range vals {
		k := val.MapKey()
		m.m[k] = val
		m.alive.Add(k)
	}
	return
}
