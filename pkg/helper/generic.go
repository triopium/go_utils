package helper

import "fmt"

type Iterator[T any] interface {
	ForEachRemaining(action func(T) error) error
	// other methods
}

// Struct version
type MyStruct[T any] struct{}

var _ Iterator[uint64] = &MyStruct[uint64]{} // Validate this type conforms to interface

func (s *MyStruct[T]) ForEachRemaining(action func(T) error) error {
	return nil
}

// Slice version
type MySlice[T any] []T

var _ Iterator[uint64] = MySlice[uint64]{} // Validate this type conforms to interface

func (s MySlice[T]) ForEachRemaining(action func(T) error) error {
	for _, e := range s {
		if err := action(e); err != nil {
			return err
		}
	}
	return nil
}

func CheckInt(in int) (bool, error) {
	if in == 0 {
		return false, fmt.Errorf("0")
	}
	if in > 1 && in < 10 {
		return true, nil
	}
	return false, nil
}

// Struct version
type MyConfig[T any] struct {
}

func (s *MyConfig[T]) Check(
	action func(T) (bool, error), input T) error {
	ok, err := action(input)
	if err != nil {
		return err
	}
	if ok {
		fmt.Println("input is ok")
		return nil
	}
	fmt.Println("input not ok")
	return nil
}

type Fopt struct {
	// MyConfig
	My any
}

// func (f *Fopt) AddMethod(fu MyConfig[int]) {
func (f *Fopt) AddMethod(fu any) {
	f.My = fu
}

type Opt[T any] struct {
	Description string
	FuncMatch   func(v T) (bool, error)
}

func (o *Opt[T]) CheckujAlloved(input T, aFunc func(T) (bool, error)) {
	ok, err := aFunc(input)
	if err != nil {
		panic(fmt.Errorf("error checking value: %s", err))
	}
	if !ok {
		panic("value not alloved")
	}
	if ok {
		fmt.Println("param ok")
	}
}

// type OpstMap struct {

// }
