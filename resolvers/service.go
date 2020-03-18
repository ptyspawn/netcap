package resolvers

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// https://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.csv

// excerpt:
//Service Name,Port Number,Transport Protocol,Description,Assignee,Contact,Registration Date,Modification Date,Reference,Service Code,Unauthorized Use Reported,Assignment Notes
//,0,tcp,Reserved,[Jon_Postel],[Jon_Postel],,,,,,
//,0,udp,Reserved,[Jon_Postel],[Jon_Postel],,,,,,
//tcpmux,1,tcp,TCP Port Service Multiplexer,[Mark_Lottor],[Mark_Lottor],,,,,,
//tcpmux,1,udp,TCP Port Service Multiplexer,[Mark_Lottor],[Mark_Lottor],,,,,,
//compressnet,2,tcp,Management Utility,,,,,,,,
//compressnet,2,udp,Management Utility,,,,,,,,
//compressnet,3,tcp,Compression Process,[Bernie_Volz],[Bernie_Volz],,,,,,
//compressnet,3,udp,Compression Process,[Bernie_Volz],[Bernie_Volz],,,,,,

var (
	udpPortMap = make(map[int]port)
	tcpPortMap = make(map[int]port)
)

type port struct {
	service string
	num int
}

func InitServiceDB() {

	var (
		f, err = os.Open(filepath.Join(dataBaseSource, "service-names-port-numbers.csv"))
		csvReader = csv.NewReader(f)
	)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for i, r := range records {

		// skip CSV header
		if i == 0 {
			continue
		}

		// skip empty or broken records
		if r[1] == "" {
			continue
		}

		// handle ranges
		if strings.Contains(r[1], "-") {
			parts := strings.Split(r[1], "-")
			if len(parts) != 2 {
				fmt.Println("invalid parts length", parts)
				continue
			}
			start, err :=  strconv.Atoi(parts[0])
			if err != nil {
				fmt.Println(err)
				continue
			}
			end, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println(err)
				continue
			}
			if end < start {
				fmt.Println("invalid range", parts)
				continue
			}
			for i := start; i <= end; i++ {
				p := port{
					service: r[3],
					num: i,
				}
				if r[2] == "tcp" {
					tcpPortMap[i] = p
				} else if r[2] == "udp" {
					udpPortMap[i] = p
				} else {
					//fmt.Println("ignoring:", r)
				}
			}
		} else {
			// add port
			num, err := strconv.Atoi(r[1])
			if err != nil {
				fmt.Println(err)
				continue
			}
			p := port{
				service: r[3],
				num: num,
			}
			if r[2] == "tcp" {
				tcpPortMap[num] = p
			} else if r[2] == "udp" {
				udpPortMap[num] = p
			} else {
				//fmt.Println("ignoring:", r)
			}
		}
	}

	if !Quiet {
		fmt.Println("loaded", len(tcpPortMap), "TCP service records")
		fmt.Println("loaded", len(udpPortMap), "UDP service records")
	}
}

func LookupServiceByPort(port int, typ string) string {
	if typ == "TCP" {
		if res, ok := tcpPortMap[port]; ok {
			return res.service
		}
	} else {
		if res, ok := udpPortMap[port]; ok {
			return res.service
		}
	}
	if !Quiet {
		log.Println("no result for", port, typ)
	}
	return ""
}