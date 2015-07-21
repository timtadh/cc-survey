package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

import (
	"github.com/timtadh/cc-survey/clones"
	"github.com/timtadh/cc-survey/models"
	"github.com/timtadh/data-structures/set"
	"github.com/timtadh/data-structures/types"
)


type SurveyLogStore struct {
	questions []models.Renderable
	clones []*clones.Clone
	cloneIdxs *set.SortedSet
	answersPath string
}

func NewSurveyStore(dir string, questions []models.Renderable, clones []*clones.Clone) (*SurveyLogStore, error) {
	fi, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		err := os.Mkdir(dir, 0775)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else if !fi.IsDir() {
		return nil, fmt.Errorf("%v is not a directory", dir)
	}
	cloneIdxs := set.NewSortedSet(len(clones))
	for i := 0; i < len(clones); i++ {
		cloneIdxs.Add(types.Int(i))
	}
	st := &SurveyLogStore{
		questions: questions,
		clones: clones,
		cloneIdxs: cloneIdxs,
		answersPath: filepath.Join(dir, "answers"),
	}
	return st, nil
}

func (st *SurveyLogStore) Do(f func(*models.Survey) error) error {
	answersCount, s, err := st.load()
	if err != nil {
		return err
	}
	err = f(s)
	if err != nil {
		return err
	}
	return st.save(answersCount, s)
}

func (st *SurveyLogStore) load() (int, *models.Survey, error) {
	answers := make([]*models.SurveyAnswer, 0, len(st.clones)*2)
	answered := set.NewSortedSet(len(st.clones))
	err := createOrOpen(st.answersPath,
		func(path string) (err error) {
			// create file
			f, err := os.Create(path)
			if err != nil {
				return err
			}
			return f.Close()
		},
		func(path string) (err error) {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			return st.loadFile(f, &answers, answered)
		},
	)
	if err != nil {
		return 0, nil, err
	}
	unanswered := st.cloneIdxs.Subtract(answered)
	s := &models.Survey{
		Questions: st.questions,
		Clones: st.clones,
		Unanswered: unanswered,
		Answers: answers,
	}
	return len(answers), s, nil
}

func (st *SurveyLogStore) loadFile(f io.Reader, answers *[]*models.SurveyAnswer, answered *set.SortedSet) error {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		err := st.loadLine(line, answers, answered)
		if err != nil {
			return err
		}
	}
	return scanner.Err()
}

func (st *SurveyLogStore) loadLine(line []byte, answers *[]*models.SurveyAnswer, answered *set.SortedSet) error {
	var a models.SurveyAnswer
	err := json.Unmarshal(line, &a)
	if err != nil {
		return err
	}
	answered.Add(types.Int(a.CloneID))
	*answers = append(*answers, &a)
	return nil
}

func (st *SurveyLogStore) save(answersCount int, s *models.Survey) error {
	f, err := os.OpenFile(st.answersPath, os.O_APPEND|os.O_SYNC, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	return st.saveFile(f, answersCount, s)
}

func (st *SurveyLogStore) saveFile(f io.Writer, answersCount int, s *models.Survey) error {
	enc := json.NewEncoder(f)
	for i := answersCount; i < len(s.Answers); i++ {
		err := enc.Encode(&s.Answers[i])
		if err != nil {
			return err
		}
		_, err = f.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}
	return nil
}

func createOrOpen(path string, create, open func(string) error) error {
	fi, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		// ok the file does not exist
		return create(path)
	} else if err != nil {
		return err
	} else if fi.IsDir() {
		return fmt.Errorf("%v is a directory", path)
	} else {
		return open(path)
	}
}

