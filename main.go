package main

import (
	"bio-splitter/logic"
	"encoding/json"
	"fmt"
)

type MarshallStruct struct {
	Index int    `json:"index"`
	Range string `json:"range"`
}

func marshalAndReturn(s logic.Splitter) {
	ranges := s.GetRanges()
	rnstrings := convertRangeToJsonRange(ranges)
	va, err := json.Marshal(rnstrings)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(va))
}

func convertRangeToJsonRange(rngs []logic.Range) (retstr []MarshallStruct) {
	for _, rn := range rngs {

		retstr = append(retstr, MarshallStruct{
			Index: rn.Index,
			Range: fmt.Sprintf("Range:bytes=%v-%v", rn.StartByte, rn.StopByte),
		})
	}
	return retstr
}

func main() {

	// Created Splitter logic, must react to ENV-VARS INPUTFILE AND PARAMS
	// INPUTFILE = String with URL
	// PARAMS = JSON obj. with list of parameters
	// Must implement logic.Splitter interface
	splitter := logic.FastaSplitter{}
	marshalAndReturn(splitter)
}
