package main

import (
	

)


type Report struct {
	UUID     string   `json:"uuid"`
	ScannerName string  `json:"scanname"`
	Time     string  `json:"time"`
	ReportUrl string `json:"reporturl"`
	Vulnerabilities []vulnerability `json:"vulnerabilities"`
}

