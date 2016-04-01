package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
	"encoding/json"
)

const (
	APP_VERSION      = "0.1"
	scan_input_path  = "scanin"
	scan_output_path = "scanout"
)

type input struct {
	apiUrl   string
	username string
	password string
}

var in input

func init() {
	flag.StringVar(&in.apiUrl, "a", "REQUIRED", "api url for example https://saleshub.blackducksoftware.com")
	flag.StringVar(&in.username, "u", "REQUIRED", "the username that can be used to access rest endpoints")
	flag.StringVar(&in.password, "p", "REQUIRED", "the password")
}

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "Print the version number.")

func main() {
	flag.Parse() // Scan the arguments list

	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
	}

	config, ok := GetConfig()
	if !ok {
		// create config and exit
		fmt.Printf("Config has to be created.")
		CreateConfig()
		os.Exit(1)
	}

	// check if the input path exists
	if _, err := os.Stat(scan_input_path); os.IsNotExist(err) {
		fmt.Errorf("ERROR input directory %s doesn't exist\n%s\n\n", scan_input_path, err)
		fmt.Printf("ERROR input directory %s doesn't exist\n%s\n\n", scan_input_path, err)
		os.Exit(1)
	} else if _, err := os.Stat(scan_output_path); os.IsNotExist(err) {
		fmt.Printf("ERROR ouput directory %s doesn't exist\n%s\n\n", scan_output_path, err)
		os.Exit(1)
	} else {
		// check if the directory is empty
		files, err := ioutil.ReadDir(scan_input_path)
		if err != nil {
			fmt.Printf("ERROR reading dir %s\n%s\n\n", scan_input_path, err)
			os.Exit(1)
		}
		if len(files) == 0 {
			fmt.Printf("WARNING dir %s is empty nothing to scan\n", scan_input_path)
			os.Exit(1)
		}
	}

	//ok all is fine start scan
	ScanImage(scan_input_path, config)
	
	
	// now the scan is finished find the BOM report
	p, _ := filepath.Abs(scan_input_path)
	p = strings.Replace(p, "\\", "/", -1)
	h, _ := os.Hostname()
	searchString := fmt.Sprintf("%s/%s", h, p)
	searchString = strings.Replace(searchString, "//", "/", -1)
	hub := HubServer{Config: config}
	if ok := hub.login(); !ok {
		fmt.Printf("ERROR login into the hub.\n")
		os.Exit(1)
	}
	
	// check if the scan was completed
	codelocations := hub.findCodeLocations(searchString)
	if len(codelocations.Items) != 1 {
		fmt.Printf("ERROR no code locations for search string : \n%s\n\n", searchString)
		os.Exit(1)
	}
	for strings.Compare(codelocations.Items[0].Status, "COMPLETE") != 0 {
		time.Sleep( 1 * time.Minute )
		codelocations = hub.findCodeLocations(searchString)
		fmt.Printf("Scan status : %s\n", codelocations.Items[0].Status)
	}
	
	// check if the BOM creation is completed
	bomRows := hub.getBomRows(codelocations.Items[0].Version.Id, 1)
	count := bomRows.TotalCount
	for {
		time.Sleep(1 * time.Minute)
		bomRows = hub.getBomRows(codelocations.Items[0].Version.Id, 1)
		if count == bomRows.TotalCount {
			break
		}
		count = bomRows.TotalCount
		fmt.Printf("Building BOM : %d rows\n", count)
	}
	
	
	// now get the components with vulnerabilities
	vulnBom := hub.getVulnerabilityBom(codelocations.Items[0].Version.Id, 5000)
	
	
	// get the vulnerabilities per component
	var totalVulnerabilitiesList []vulnerability
	for _, v := range vulnBom.Items {
		vulnList := hub.getVulnerabilities(codelocations.Items[0].Version.Id, v.ChannelRelease.Id, v.ProducerReleases[0].Id, 5000)
		totalVulnerabilitiesList = append(totalVulnerabilitiesList, vulnList.Items...)
	}
	
	fmt.Printf("total nr of vulnerabilities : %d\n", len(totalVulnerabilitiesList))
	//timeStamp := time.Now().Format(time.RFC3339)
	t := time.Now()
	timeStamp := fmt.Sprintf("%d-%02d-%02dT%02d_%02d_%02d-00_00",
        t.Year(), t.Month(), t.Day(),
        t.Hour(), t.Minute(), t.Second())
	reportUrl := fmt.Sprintf("%s/#versions/id:%s/vier:bom", config.Url, codelocations.Items[0].Version.Id)
	report := &Report{UUID: "uuid", ScannerName: "blackduck", Time: timeStamp, Vulnerabilities: totalVulnerabilitiesList, ReportUrl: reportUrl}
	
	jsonReport , err := json.Marshal(report)
	if err != nil {
		fmt.Printf("ERROR marshall of report went wrong : \n%s\n", err)
		os.Exit(1)
	}
	
	fileName := fmt.Sprintf("%s.json", timeStamp)
	if err := ioutil.WriteFile(filepath.Join(scan_output_path, fileName), jsonReport, 0755); err != nil {
		fmt.Printf("ERROR writing file %s \n%s\n\n", fileName, err)
		os.Exit(1)
	}

}
