package main

import (
	//"encoding/json"
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"os"
	"path/filepath"
	//"strings"
	//"time"
	"sync"
)

const (
	APP_VERSION      = "0.1_3.2"
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

	// now find all input directories
	scanDirectories, _ := ioutil.ReadDir(scan_input_path)

	// add a inotify watcher to watch the "statusWriteDir"
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err)
	}
	defer watcher.Close()

	var wg sync.WaitGroup

	//ok all is fine start scan
	// note every input scan dir is scanned parrellel
	for _, dir := range scanDirectories {

		inputPath := filepath.Join(scan_input_path, dir.Name())
		outputPath := filepath.Join(scan_output_path, dir.Name())
		statusDir := filepath.Join(outputPath, "status")
		os.MkdirAll(statusDir, 0775)
		// this is a bid magical somehow iscan creates a directory status to dump the status json
		err = watcher.Add(statusDir)
		if err != nil {
			//fmt.Fatal(err)
			fmt.Println(err)
		}
		wg.Add(1)
		go ScanImage(inputPath, config, outputPath)
	}

	// ok all scans are running now so lets wait for the statusWrite json file
	go func() {
		for {
			select {
			case event := <-watcher.Events:

				if event.Op&fsnotify.Create == fsnotify.Create {
					fmt.Println("create : ", event.Name)
					wg.Done()
					// download report
				}

			case err := <-watcher.Errors:
				fmt.Println("error:", err)
			}
		}
	}()
	fmt.Println("wait")
	wg.Wait()
	fmt.Println("the end")

}
