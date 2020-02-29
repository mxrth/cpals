package file

import (
	"bufio"
	"os"
)

//TODO: document, write tests

func ReadLines(fileName string) (lines []string) {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	i := 0
	for scanner.Scan() {
		bs := scanner.Bytes()
		lines = append(lines, string(bs))
		i++
	}
	//fmt.Print(len(lines))
	return
}
