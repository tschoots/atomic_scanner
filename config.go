package main

import (
	"archive/zip"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"path/filepath"
)

const (
	config_file = "conf/config.json"
	cli_url     = "/download/scan.cli.zip"
	cli_path    = "scanner"
)

var key = []byte("pronktjiesparktr535afasdf asdnvr")

type config struct {
	Url      string `json:"url"`
	Host     string `json:"hubhost"`
	Port     string `json:"port"`
	Scheme   string `json:"scheme"`
	User     string `json:"user"`
	Password string `json:"password"`
	ScannerDir  string `json:"scannerdir"`
}

func GetConfig() (*config, bool) {
	// check if the json file in conf/config.json exist
	if _, err := os.Stat(config_file); err != nil {
		if os.IsNotExist(err) {
			return nil, false
		} else {
			fmt.Errorf("file %s has problems : \n%s\n", config_file, err)
			return nil, false
		}

	} else {
		// config_file exists read it and unmarshall
		conf := &config{}
		
		file, err := ioutil.ReadFile(config_file)
		if err != nil {
			fmt.Errorf("ERROR readfile %s : \n%s\n", config_file, err)
			return nil, false
		}
		
		json.Unmarshal(file, conf)
		conf.Password = decrypt(key, conf.Password)
		return conf, true

	}

}

func CreateConfig() (*config , bool) {
	conf := config{ScannerDir: cli_path}
	conf.init()
	
	jsonString , err := json.Marshal(conf)
	if err != nil {
		fmt.Errorf("Error marshall configuration: %s\n\n", err)
		return nil, false
	}
	
	if err := os.MkdirAll(filepath.Dir(config_file), 0755); err != nil {
		fmt.Errorf("Error creating directory %s : \n%s\n\n", filepath.Dir(config_file), err)
		return nil, false
	}
	
	if err := ioutil.WriteFile(config_file, jsonString, 0755); err != nil {
		fmt.Errorf("error writing config file : %s\n\n", err)
		return nil, false
	}
	
	conf.Password = decrypt(key, conf.Password)
	
	if err := os.MkdirAll(cli_path, 0755); err != nil {
		fmt.Errorf("Error creating directory %s: \n%s\n\n", cli_path, err)
		return nil, false
	}
	
	downloadUrl := fmt.Sprintf("%s/%s", conf.Url, cli_url)
	downloadFromUrl(downloadUrl)
	return &conf, true
}

func (c *config) init() {

	// config file doesn't exist so get data from the user
	var hubUrl string

	fmt.Println("Enter url formate http|https://<server>[:port]: ")
	fmt.Scanln(&hubUrl)

	fmt.Printf("hub user: \n")
	fmt.Scanln(&c.User)

	fmt.Printf("password:")
	fmt.Println("\033[8m")
	fmt.Scanln(&c.Password)
	fmt.Println("\033[28m")

	u, err := url.Parse(hubUrl)
	if err != nil {
		panic(err)
	}

	c.Url = fmt.Sprintf("%s://%s", u.Scheme, u.Host)

	c.Scheme = u.Scheme

	if strings.Compare(c.Scheme, "https") == 0 {
		c.Port = "443"
	} else {
		c.Port = "80"
	}
	c.Host = u.Host
	if strings.Contains(c.Host, ":") {
		host, port, _ := net.SplitHostPort(c.Host)
		c.Host = host
		c.Port = port
	}

	if valid := c.validUseridPassword(); !valid {
		fmt.Printf("ERROR user id password combination not valid.")
		os.Exit(2)
	}
	c.Password = encrypt(key, c.Password)

}

// encrypt string to base64 crypto using AES
func encrypt(key []byte, text string) string {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// decrypt from base64 to decrypted string
func decrypt(key []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}

func (c *config) validUseridPassword() bool {

	fmt.Println(c.Url)
	u, err := url.ParseRequestURI(c.Url)
	if err != nil {
		fmt.Printf("ERROR : url.ParseRequestURI\n%s\n", err)
		return false
	}
	resource := "/j_spring_security_check"
	u.Path = resource
	data := url.Values{}
	data.Add("j_username", c.User)
	data.Add("j_password", c.Password)

	client := &http.Client{}

	jar := &myjar{}
	jar.jar = make(map[string][]*http.Cookie)
	client.Jar = jar

	urlStr := fmt.Sprintf("%v", u)
	req, err := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	if err != nil {
		fmt.Printf("ERROR NewRequest:\n%s\n", err)
		return false
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded") // needed this the prevend 401 Unauthorized

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("ERROR client.do\n%s\n", err)
		return false
	}
	resp.Body.Close()
	if resp.StatusCode != 204 {
		fmt.Printf("ERROR : resp status : %s\n%d\n", resp.Status, resp.StatusCode)
		return false
	}
	return true
}

type myjar struct {
	jar map[string][]*http.Cookie
}

func (p *myjar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	//	fmt.Printf("in SetCookies The URL is : %s\n", u.String())
	//	fmt.Printf("u.Host : %s\n", u.Host)
	//	fmt.Printf("The cookie being set is : %s\n", cookies)
	p.jar[u.Host] = cookies
}

func (p *myjar) Cookies(u *url.URL) []*http.Cookie {
	//	fmt.Printf("in cookies The URL is : %s\n", u.String())
	//	fmt.Printf("u.Host : %s\n", u.Host)
	//	fmt.Printf("Cookie being returned is : %s\n", p.jar[u.Host])
	return p.jar[u.Host]
}

func downloadFromUrl(url string) {
	tokens := strings.Split(url, "/")
	dirName := cli_path
	filename := tokens[len(tokens)-1]
	fullFileName := fmt.Sprintf("%s/%s", dirName, filename)
	fmt.Printf("downloading %s to %s\n", url, fullFileName)

	
	os.MkdirAll(cli_path, 0755)
	output, err := os.Create(fullFileName)
	if err != nil {
		fmt.Printf("ERROR creating file %s \n%s\n", fullFileName, err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("ERROR while downloading %s\n%s\n", url, err)
		return
	}

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Printf("ERROR copy to file %s\n%s\n", fullFileName, err)
		return
	}
	fmt.Sprintf("%d bytes downloaded.", n)

	r, err := zip.OpenReader(fullFileName)
	if err != nil {
		fmt.Printf("ERROR zip.OpenReader : %s\n%s\n", fullFileName, err)
		return
	}

	//iterate through the archive create the files and directories
	for _, f := range r.File {
		name := fmt.Sprintf("%s/%s", dirName, f.Name)
		if f.FileInfo().IsDir() {
			fmt.Printf("dir : %s\n", f.Name)
			os.Mkdir(name, 0775)
		} else {
			fmt.Printf("file : %s\n", f.Name)
			rc, err := f.Open()
			if err != nil {
				fmt.Printf("ERROR can't open file : %s\n", f.Name)
				continue
			}
			fo, err := os.Create(name)
			if err != nil {
				fmt.Printf("ERROR creating file %s\n%s\n", name, err)
				continue
			}
			defer fo.Close()

			//buffer the writing to file
			buf := make([]byte, 1024)
			for {
				n, err := rc.Read(buf)
				if n == 0 {
					break
				}
				if err != nil {
					fmt.Printf("ERROR reading buffer of %s\n%s\n", name, err)
					break
				}

				// write a chunk
				if _, err := fo.Write(buf[:n]); err != nil {
					fmt.Printf("ERROR writing buffer to file %s\n%s\n", name, err)
					break
				}
			}

		}
	}
	os.Remove(fullFileName)
}
