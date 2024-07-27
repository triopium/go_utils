package helper

import (
	"fmt"
	"testing"
)

func TestGeneric(t *testing.T) {
	iters := []Iterator[uint64]{
		&MyStruct[uint64]{},
		MySlice[uint64]{0, 123, 12345678901234567890},
	}
	for _, iter := range iters {
		err := iter.ForEachRemaining(func(e uint64) error {
			fmt.Println(e)
			return nil
		})
		if err != nil {
			panic(err)
		}
	}
}

func TestGeneric2(t *testing.T) {
	iters := []Iterator[string]{
		&MyStruct[string]{},
		MySlice[string]{"kek", "lek"},
	}
	for _, iter := range iters {
		err := iter.ForEachRemaining(func(e string) error {
			fmt.Println(e)
			return nil
		})
		if err != nil {
			panic(err)
		}
	}
}

func TestGenericFunc(t *testing.T) {
	mc := MyConfig[string]{}
	err := mc.Check(DirectoryExists, "/tmp")
	fmt.Println("kek", err)
	mc1 := MyConfig[int]{}
	err = mc1.Check(CheckInt, 2)
	fmt.Println("kek", err)
}

func TestGenericFunc2(t *testing.T) {
	fopt := Fopt{}
	fopt.AddMethod(MyConfig[int]{})
	fopt.AddMethod(MyConfig[string]{})
	kek := fopt.My.(MyConfig[string])
	err := kek.Check(DirectoryExists, "/tmp")
	if err != nil {
		t.Error(err)
	}
	// fopt.My = MyConfig{}
}

func TestGenericFunc3(t *testing.T) {
	o := Opt[string]{}
	o.FuncMatch = DirectoryExists
	ok, err := o.FuncMatch("/tmp")
	if err != nil {
		t.Error(err)
	}
	if ok {
		fmt.Println("kek param ok")
	}
}

func TestGenericFunc4(t *testing.T) {
	sl := []any{Opt[int]{}}
	sl = append(sl, Opt[string]{}) //nolint:all
}
