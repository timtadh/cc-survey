package models

import (
	"log"
	"math/rand"
)

import (
	"github.com/timtadh/cc-survey/clones"
)


func init() {
	rand.Seed(int64(randUint64()))
}

type Survey struct {
	Questions []Renderable
	Clones []*clones.Clone
	Unanswered []int
	Answers []*SurveyAnswer
}

type SurveyAnswer struct {
	UserEmail string
	CloneID int
	Responses []Response
}

type Response struct {
	QuestionID int
	Answer int
	Text string
}

type SurveyStore interface {
	Do(func(*Survey) error) error
}


func newSurvey(questions []Renderable, clones []*clones.Clone) *Survey {
	unanswered := make([]int, 0, len(clones))
	for i := 0; i < len(clones); i++ {
		unanswered = append(unanswered, i)
	}
	return &Survey{
		Questions: questions,
		Clones: clones,
		Unanswered: unanswered,
		Answers: make([]*SurveyAnswer, 0, len(questions)*2),
	}
}

func (s *Survey) Next() (cid int, c *clones.Clone) {
	var idx int
	if len(s.Unanswered) == 0 {
		return -1, nil
	} else if len(s.Unanswered) == 1 {
		idx = 0
	} else {
		idx = rand.Intn(len(s.Unanswered))
	}
	cid = s.Unanswered[idx]
	c = s.Clones[cid]
	return cid, c
}

func (s *Survey) Answer(answer *SurveyAnswer) {
	cid := answer.CloneID
	var idx int = -1
	for i, ucid := range s.Unanswered {
		if ucid == cid {
			idx = i
			break
		}
	}
	if idx >= 0 {
		s.Unanswered = remove(s.Unanswered, idx)
	}
	s.Answers = append(s.Answers, answer)
}

func remove(list []int, i int) []int {
	if i < 0 || i >= len(list) {
		log.Panicf("out of range remove len: %d, i: %d", len(list), i)
	}
	dst := list[i:len(list)-1]
	src := list[i+1:len(list)]
	copy(dst, src)
	return list[:len(list)-1]
}
