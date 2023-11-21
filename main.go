package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type XML struct {
	Text    string `xml:",chardata"`
	Endtime string `xml:"endtime,attr"`
	Address struct {
		Text     string `xml:",chardata"`
		Addr     string `xml:"addr,attr"`
		Addrtype string `xml:"addrtype,attr"`
	} `xml:"address"`
	Ports struct {
		Text string `xml:",chardata"`
		Port struct {
			Text     string `xml:",chardata"`
			Protocol string `xml:"protocol,attr"`
			Portid   string `xml:"portid,attr"`
			State    struct {
				Text      string `xml:",chardata"`
				State     string `xml:"state,attr"`
				Reason    string `xml:"reason,attr"`
				ReasonTtl string `xml:"reason_ttl,attr"`
			} `xml:"state"`
		} `xml:"port"`
	} `xml:"ports"`
}

type JSON struct {
	IP        string `json:"ip"`
	Timestamp string `json:"timestamp"`
	Ports     []struct {
		Port   int    `json:"port"`
		Proto  string `json:"proto"`
		Status string `json:"status"`
		Reason string `json:"reason"`
		TTL    int    `json:"ttl"`
	} `json:"ports"`
}

func main() {
	if len(os.Args) != 4 || len(os.Args) == 1 {
		fmt.Println("usage:\n" +
			"	msparse <input type> <input file> <output file>\n" +
			"input types:\n" +
			"	xml, json, list\n" +
			"example:\n" +
			"	msparse list masscan.txt filteredscan.txt")

		os.Exit(0)
	}

	infile, err := os.Open(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	defer infile.Close()

	outfile, err := os.Create(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()

	scanner := bufio.NewScanner(infile)
	writer := bufio.NewWriter(outfile)

	if os.Args[1] == "list" {
		for scanner.Scan() {
			values := strings.Split(scanner.Text(), " ")
			if values[0] == "#masscan" || values[0] == "#end" {
				continue
			}
			ip := values[3]
			port := values[2]
			_, err = writer.WriteString(fmt.Sprintf("%s:%s\n", ip, port))
			if err != nil {
				log.Fatal(err)
			}
		}
	} else if os.Args[1] == "xml" {
		d := xml.NewDecoder(infile)
		for {
			var host XML
			t, err := d.Token()
			if err != nil {
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatal(err)
				}
			}
			switch t := t.(type) {
			case xml.StartElement:
				if t.Name.Local == "host" {
					err := d.DecodeElement(&host, &t)
					if err != nil {
						log.Fatal(err)
					}
					_, err = writer.WriteString(fmt.Sprintf("%s:%s\n", host.Address.Addr, host.Ports.Port.Portid))
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	} else if os.Args[1] == "json" {
		d := json.NewDecoder(infile)

		// read open bracket
		_, err := d.Token()
		if err != nil {
			log.Fatal(err)
		}

		for d.More() {
			var host JSON
			err = d.Decode(&host)
			if err != nil {
				log.Fatal(err)
			}
			_, err = writer.WriteString(fmt.Sprintf("%s:%d\n", host.IP, host.Ports[0].Port))
			if err != nil {
				log.Fatal(err)
			}

		}

		// read closing bracket
		_, err = d.Token()
		if err != nil {
			log.Fatal(err)
		}

	} else {
		fmt.Printf("unknown input file type: '%s'\n", os.Args[1])
		fmt.Println("available input types: xml, json, list")
	}
	writer.Flush()
	filepath, _ := filepath.Abs(os.Args[2])
	fmt.Println("Done!")
	fmt.Printf("Wrote to: %s\n", filepath)
}
