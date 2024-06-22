package helper

import (
	"encoding/json"
	"fmt"
)

type ErrMap map[string][]string

func (em ErrMap) Add(errMain error, errorsPartial ...error) bool {
	if em == nil {
		panic("empty receiver")
	}
	if errMain == nil {
		return false
	}
	_, ok := em[errMain.Error()]
	if !ok {
		em[errMain.Error()] = make([]string, 0)
	}
	if len(errorsPartial) == 0 {
		em[errMain.Error()] = append(
			em[errMain.Error()], "")
		fmt.Println("during", em)
		return true
	}
	if len(errorsPartial) > 0 {
		for _, e := range errorsPartial {
			em[errMain.Error()] = append(
				em[errMain.Error()], e.Error())
		}
	}
	return true
}

func (em ErrMap) Aggregate(errMain error, msgs ...string) bool {
	if em == nil {
		panic("empty receiver")
	}
	if errMain == nil {
		return false
	}
	_, ok := em[errMain.Error()]
	if !ok {
		em[errMain.Error()] = make([]string, 0)
	}
	if len(msgs) == 0 {
		em[errMain.Error()] = append(
			em[errMain.Error()], "")
		return true
	}
	if len(msgs) > 0 {
		for _, e := range msgs {
			em[errMain.Error()] = append(
				em[errMain.Error()], e)
		}
	}
	return true
}

func (em ErrMap) Marshal() string {
	b, err := json.Marshal(em)
	if err != nil {
		return "cannot marshal errors"
	}
	return string(b)
}

func (em ErrMap) MarshalError(errMsg string) error {
	if len(em) == 0 {
		return nil
	}
	errs := em.Marshal()
	return fmt.Errorf("%s: %s", errMsg, errs)
}
