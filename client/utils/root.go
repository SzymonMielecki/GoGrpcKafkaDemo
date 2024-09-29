package utils

import "hash/fnv"

func GetColorForUser(username string) int {
	hash := fnv.New32a()
	hash.Write([]byte(username))
	colors := []int{31, 32, 33, 34, 35, 36} // ANSI color codes
	return colors[hash.Sum32()%uint32(len(colors))]
}
