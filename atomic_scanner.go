package main 

import (
    "flag"
    "fmt"
    "os"
)

const APP_VERSION = "0.1"

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
   
    fmt.Printf(config.User)
   
    
}

