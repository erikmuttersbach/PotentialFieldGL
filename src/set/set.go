package set

type Set struct {
	data map[interface{}]bool
}

func NewSet() *Set {
	return &Set{
		make(map[interface{}]bool),
	}
}

func (s *Set) Len() int {
	return len(s.data)
}

func (s *Set) Put(obj interface{}) {
	s.data[obj] = true
}

func (s *Set) Remove(obj interface{}) {
	delete(s.data, obj)
}

func (s *Set) Contains(obj interface{}) bool {
	return s.data[obj]
}

func (s *Set) Elements() []interface{} {
	elms := make([]interface{}, len(s.data))
	i := 0
	for key, _ := range s.data {
		elms[i] = key
		i++
	} 
	return elms
}




