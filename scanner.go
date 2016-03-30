package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"strings"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"log"
	"path/filepath"
	"regexp"
	"flag"

)

//func ScanImage(path string,jar string, jarpath string, image string, conf config) {
func ScanImage(conf config) {
	//iscanPath := "/home/blackduck/j_scanner/bin/scan.cli.sh"
	//hostname, _ := os.Hostname()
	img_arr := strings.Split(image, ":")
	img_name := img_arr[0]
	tag := img_arr[1]
	project := fmt.Sprintf("{%s}%s", conf.ScanHost, img_name)
	
	jar, jarPath := getJarFiles(conf.ScannerDir)
	
	onejarpath := fmt.Sprintf("-Done-jar.jar.path=%s", jarpath)

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
		path)
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



