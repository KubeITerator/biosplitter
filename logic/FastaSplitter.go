package logic

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type FastaSplitter struct {
}

type FastaParams struct {
	// Maximum number of records per split
	MaxRecords int `json:"maxrecord,omitempty"`
	// Maximum Size in Bytes
	ByteSize int `json:"bytesize,omitempty"`
}

func (f FastaSplitter) GetRanges() []Range {

	url := os.Getenv("DATASOURCE")
	if url == "" {
		fmt.Println("empty envvar: DATASOURCE, abort")
		os.Exit(2)
	}
	//url = "https://s3.computational.bio.uni-giessen.de/swift/v1/testdata/small_vir.faa"
	params := f.ParseParams()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error in getting file")
		os.Exit(2)
	}
	defer resp.Body.Close()

	offset := 0
	var splitSites []int
	buf := make([]byte, 16384)
	for {
		n, err := resp.Body.Read(buf)
		if err == io.EOF {
			splitSites = append(splitSites, int(resp.ContentLength))
			break
		}
		for index, data := range buf[:n] {
			// Byte is an alias for uint8, 62 == '>'
			if data == 62 {
				splitSites = append(splitSites, index+offset)
			}
		}
		offset += n
	}

	return f.GetRangesWithSiteList(splitSites, params)

}
func (f FastaSplitter) GetRangesWithSiteList(sitelist []int, params FastaParams) (returnRange []Range) {
	counter := 0
	if params.MaxRecords != 0 {
		currentRecordNum := -1
		byteRange := Range{0, -1, counter}
		for _, site := range sitelist {
			currentRecordNum++
			if currentRecordNum == params.MaxRecords {
				currentRecordNum = 0
				byteRange.StopByte = site - 1
				returnRange = append(returnRange, byteRange)
				counter++
				byteRange = Range{site, -1, counter}
			}
		}
		byteRange.StopByte = sitelist[len(sitelist)-1]
		returnRange = append(returnRange, byteRange)

	} else if params.ByteSize != 0 {
		lastSplitSite := 0
		byteRange := Range{0, -1, counter}
		for index, site := range sitelist {
			if index < len(sitelist)-1 {
				if sitelist[index+1]-lastSplitSite > params.ByteSize {
					lastSplitSite = site
					byteRange.StopByte = site - 1
					if byteRange.StopByte != -1 {
						returnRange = append(returnRange, byteRange)
						counter++
					}
					byteRange = Range{site, -1, counter}
				}
			}
		}
		byteRange.StopByte = sitelist[len(sitelist)-1]
		returnRange = append(returnRange, byteRange)
	}

	return returnRange
}
func (f FastaSplitter) ParseParams() FastaParams {
	param := os.Getenv("PARAMS")
	if param == "" {
		fmt.Println("empty envvar: PARAMS, abort")
		os.Exit(2)
	}
	fparams := FastaParams{}
	err := json.Unmarshal([]byte(param), &fparams)
	if err != nil {
		panic(err)
	}
	return fparams
}
