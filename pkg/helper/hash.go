package helper

import (
	"hash/fnv"
)

func StringToHash(input string) (uint32, error) {
	hash := fnv.New32()
	_, err := hash.Write([]byte(input))
	if err != nil {
		return 0, err
	}
	return hash.Sum32(), nil
}
