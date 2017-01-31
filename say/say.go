package thingsay

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Say interface {
	This(string) string
	One([]string) string
}

type Speaker struct {
	Art     string
	Wrapper string
}

func (s *Speaker) One(things []string) string {
	return s.This(things[rand.Intn(len(things))])
}

func (s *Speaker) This(this string) string {
	return fmt.Sprint("%s\n%s\n%s\n%s\n",
		s.Wrapper,
		this,
		s.Wrapper,
		s.Art)
}
