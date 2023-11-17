package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Nmaprun struct {
	XMLName          xml.Name `xml:"nmaprun"`
	Text             string   `xml:",chardata"`
	Scanner          string   `xml:"scanner,attr"`
	Start            string   `xml:"start,attr"`
	Version          string   `xml:"version,attr"`
	Xmloutputversion string   `xml:"xmloutputversion,attr"`
	Scaninfo         struct {
		Text     string `xml:",chardata"`
		Type     string `xml:"type,attr"`
		Protocol string `xml:"protocol,attr"`
	} `xml:"scaninfo"`
	Host     []Host `xml:"host"`
	Runstats struct {
		Text     string `xml:",chardata"`
		Finished struct {
			Text    string `xml:",chardata"`
			Time    string `xml:"time,attr"`
			Timestr string `xml:"timestr,attr"`
			Elapsed string `xml:"elapsed,attr"`
		} `xml:"finished"`
		Hosts struct {
			Text  string `xml:",chardata"`
			Up    string `xml:"up,attr"`
			Down  string `xml:"down,attr"`
			Total string `xml:"total,attr"`
		} `xml:"hosts"`
	} `xml:"runstats"`
}

type Host struct {
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
			_, err = writer.WriteString(ip + ":" + port + "\n")
			if err != nil {
				log.Fatal(err)
			}
		}
	} else if os.Args[1] == "xml" {
		d := xml.NewDecoder(infile)
		for {
			var host Host
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
					_, err = writer.WriteString(host.Address.Addr + ":" + host.Ports.Port.Portid)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	} else if os.Args[1] == "json" {

	} else {
		fmt.Printf("unknown input file type: '%s'\n", os.Args[1])
		fmt.Println("available input types: xml, json, list")
	}
	writer.Flush()
}
