package models

import (
)

import (
	"github.com/timtadh/cc-survey/clones"
	"github.com/timtadh/data-structures/set"
	"github.com/timtadh/data-structures/types"
)


type Survey struct {
	Questions []Renderable
	Clones []*clones.Clone
	Unanswered *set.SortedSet
	Answers []*SurveyAnswer
}

type SurveyAnswer struct {
	UserEmail string
	CloneID int
	CloneExtID int
	CloneDir string
	SelectionPr float64
	ConditionalPr float64
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

func (s *Survey) CountAnswers(email string) int {
	var count int
	for _, a := range s.Answers {
		if a.UserEmail == email {
			count++
		}
	}
	return count
}

func (s *Survey) Next() (cid int, c *clones.Clone) {
	cidT, err := s.Unanswered.Random()
	if err != nil {
		return -1, nil
	}
	cid = int(cidT.(types.Int))
	c = s.Clones[cid]
	return cid, c
}

func (s *Survey) Answer(answer *SurveyAnswer) {
	s.Unanswered.Remove(types.Int(answer.CloneID))
	s.Answers = append(s.Answers, answer)
}

