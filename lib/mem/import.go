package mem

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"time"
)

func Import(data io.Reader) (Mem, error) {
	reader := bufio.NewReader(data)

	header, err := parseHeader(reader)
	if err != nil {
		return Mem{}, err
	}

	return Mem{MemHeader: header}, nil
}

var headerRegex = regexp.MustCompile(`^\s+([^:]+):\s+(.*)`)

func parseHeader(reader *bufio.Reader) (MemHeader, error) {
	mp := map[string]string{}

	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			return MemHeader{}, err
		}

		if len(s) == 1 {
			// Done with header section; break!
			break
		}
		result := headerRegex.FindStringSubmatch(s)

		mp[result[1]] = result[2]
	}

	var nc bool
	if mp["NC/disease"] == "NC" {
		nc = true
	}

	var sex Sex
	switch mp["Sex"] {
	case "M":
		sex = MaleSex
	case "F":
		sex = FemaleSex
	default:
		sex = OtherSex
	}

	temp, err := strconv.ParseFloat(mp["Temperature"], 32)
	if err != nil {
		return MemHeader{}, err
	}

	age, err := strconv.Atoi(mp["Age"])
	if err != nil {
		return MemHeader{}, err
	}

	layout := "2/1/06 15:04:05"
	date, err := time.Parse(layout, mp["Date"]+" "+mp["Start time"])
	if err != nil {
		return MemHeader{}, err
	}

	header := MemHeader{
		File:          mp["File"],
		Name:          mp["Name"],
		Protocol:      mp["Protocol"],
		Date:          date,
		Age:           age,
		Sex:           sex,
		Temperature:   float32(temp),
		SRSites:       mp["S/R sites"],
		NormalControl: nc,
		Operator:      mp["Operator"],
		Comment:       mp["Comments"],
	}

	return header, nil
}
