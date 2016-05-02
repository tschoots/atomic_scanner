package main

import (
	

)


type Report struct {
	ScanType  string  `json:"Scan Type"`
	Scanner   string  `json:"Scanner"`
	FinishedTime string `json:"Finished Time"`
	UUID     string   `json:"UUID"`
	ScannerName string  `json:"scanname"`
	Time     string  `json:"Time"`
	Successful  bool `json:"Successful"`
	CVEFeedLastUpdated  string `json:"CVE Feed Last Updated"`
	ReportUrl string `json:"reporturl"`
	Vulnerabilities []vulnerability `json:"vulnerabilities"`
}

