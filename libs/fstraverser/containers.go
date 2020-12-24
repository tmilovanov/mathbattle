package fstraverser

type stringStack struct {
	impl []string
}

func (s *stringStack) push(elem string) {
	s.impl = append(s.impl, elem)
}

func (s *stringStack) isEmpty() bool {
	if len(s.impl) == 0 {
		return true
	}
	return false
}

func (s *stringStack) pop() (string, bool) {
	if s.isEmpty() {
		return "", false
	}

	r := s.impl[len(s.impl)-1]

	s.impl = s.impl[0 : len(s.impl)-1]

	return r, true
}

type stringQueue struct {
	impl []string
}

func (s *stringQueue) isEmpty() bool {
	if len(s.impl) == 0 {
		return true
	}

	return false
}

func (s *stringQueue) push(elem string) {
	s.impl = append(s.impl, elem)
}

func (s *stringQueue) pop() (string, bool) {
	if s.isEmpty() {
		return "", false
	}

	r := s.impl[0]

	s.impl = s.impl[1:]

	return r, true
}
