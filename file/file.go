package file

import (
	"bufio"
	"encoding/base64"
	"io/ioutil"
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

func ReadBase64File(f string) []byte {
	b, ok := ioutil.ReadFile(os.Args[1])
	if ok != nil {
		panic("Could not read file")
	}

	c, ok := base64.StdEncoding.DecodeString(string(b))
	if ok != nil {
		panic("Invalid Base64")
	}
	return c
}
