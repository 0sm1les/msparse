package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func main() {
	infile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer infile.Close()

	outfile, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()

	scanner := bufio.NewScanner(infile)
	writer := bufio.NewWriter(outfile)

	for scanner.Scan() {
		values := strings.Split(scanner.Text(), " ")
		if values[0] == "#masscan" || values[0] == "#end" {
			continue
		}
		ip := values[3]
		port := values[2]
		writer.WriteString(ip + ":" + port + "\n")
	}
	writer.Flush()
}
