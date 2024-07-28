package helper

import "strconv"

func SliceStringToMapString(input []string) map[string]bool {
	out := make(map[string]bool, len(input))
	for _, i := range input {
		out[i] = true
	}
	return out
}

func SliceStringToMapInt(input []string) map[int]bool {
	out := make(map[int]bool, len(input))
	for _, i := range input {
		val, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		out[val] = true
	}
	return out
}
