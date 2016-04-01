package main

import (
	"fmt"
	"os"
	"os/exec"
	"log"
	"path/filepath"
	"regexp"
	"time"

)



//func ScanImage(path string,jar string, jarpath string, image string, conf config) {
func ScanImage(scanDir string, conf *config) {
	hostname, _ := os.Hostname()
	tag := time.Now().Format(time.RFC850)
	project := fmt.Sprintf("{%s}%s", hostname, "atomic_scan" )
	
	jar, jarPath := getJarFiles(conf.ScannerDir)
	
	onejarpath := fmt.Sprintf("-Done-jar.jar.path=%s", jarPath)

	cmd := exec.Command("java",
		"-Xms256m",
		"-Xmx4096m",
		"-Done-jar.silent=true",
		onejarpath,
		"-jar", jar,
		"--host", conf.Host,
		"--port", conf.Port,
		"--scheme", conf.Scheme,
		"--project", project,
		"--release", tag,
		"--username", conf.User,
		"--password", conf.Password,
		"-v",
		scanDir)
	//log.Println(cmd.Args)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	//go io.Copy(log., os.Stderr)
	err := cmd.Run()
	if err != nil {
		log.Println(err.Error())
		return
	}

}

func getJarFiles(path string)(jar string, jarpath string){
	r_jarpath := regexp.MustCompile("^scan\\.cli\\.impl.*-standalone\\.jar$")
	r_jar := regexp.MustCompile("^scan.cli-.*-standalone\\.jar$")
	
	
	
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
			if r_jar.MatchString(f.Name()) {
				jar = path
			}
			if r_jarpath.MatchString(f.Name()) {
				jarpath = path
			}
		
		return nil
	})
	if err != nil {
		fmt.Printf("ERROR : %s\n", err)
	}
	return jar, jarpath
}



