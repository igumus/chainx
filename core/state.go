package core

import "fmt"

type State struct {
	data map[string][]byte
}

func NewState() *State {
	return &State{
		data: make(map[string][]byte),
	}
}

func (s *State) Put(k, v []byte) error {
	s.data[string(k)] = v
	return nil
}

func (s *State) Delete(k []byte) error {
	delete(s.data, string(k))
	return nil
}

func (s *State) Get(k []byte) ([]byte, error) {
	key := string(k)
	v, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("key not found in state: %s", key)
	}
	return v, nil
}

func (s *State) Merge(other *State) {
	for otk, otv := range other.data {
		s.data[otk] = otv
	}
}
