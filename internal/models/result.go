package models

import "fmt"

type Stats struct {
	Value float64
	Name  string
	Id    int
}

func (s *Stats) ToString() string {
	return s.Name + ": " + fmt.Sprintf("[%d] %f", s.Id, s.Value)
}
