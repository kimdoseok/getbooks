// testhttp project main.go
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

var dirname = "D:/Books"

var urls = []string{"http://www.allitebooks.com/", "http://www.allitebooks.org/"}

func getIndexStr() string {
	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Connection Error!")
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		bodystr := string(body)
		return bodystr
	}
	return ""
}

func getPageNum(bodystr string) int {
	pageptn := "<span class=\"pages\">[0-9]+ / ([0-9]+) Pages</span>"
	r, _ := regexp.Compile(pageptn)
	watched := r.FindStringSubmatch(bodystr)
	//fmt.Printf("%v %T\n", watched[1], watched[1])
	pagenum, err := strconv.Atoi(watched[1])
	if err != nil {
		fmt.Println("String to Integer error!")
		return 0
	}
	return pagenum
}

func getPageList(pagenum int) [][]string {
	resp, err := http.Get(fmt.Sprintf("http://www.allitebooks.org/page/%d/", pagenum))
	if err != nil {
		fmt.Println("Connection Error!")
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	bodystr := string(body)
	pageptn := "<a href=\"(http://www.allitebooks.com|http://www.allitebooks.org)/([a-z0-9-]+)/\" rel=\"bookmark\">"
	r, _ := regexp.Compile(pageptn)
	watchedpages := r.FindAllStringSubmatch(bodystr, -1)
	return watchedpages
}

func getBookPage(bookname string) [][]string {
	resp, err := http.Get(fmt.Sprintf("http://www.allitebooks.org/%s/", bookname))
	if err != nil {
		fmt.Println("Connection Error!")
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	bodystr := string(body)
	//fmt.Println(bodystr)
	pageptn := "<a href=\"(http://file.allitebooks.com|http://www.allitebooks.org)/(\\d+)/([\\w-. ]+)\\.(pdf|zip|chm|rar|epub)\" target=\"_blank\"\\>"
	r, _ := regexp.Compile(pageptn)
	filelink := r.FindAllStringSubmatch(bodystr, -1)
	return filelink
}

func saveFile(fileline []string) bool {
	fname := fmt.Sprintf("%s/%s.%s", dirname, fileline[3], fileline[4])
	if _, err := os.Stat(fname); err == nil {
		fmt.Println("File exists already", fileline[3], fileline[4])
		return false
	}
	urlstr := fmt.Sprintf("%s/%s/%s.%s", fileline[1], fileline[2], fileline[3], fileline[4])
	resp, err := http.Get(urlstr)
	body, err := ioutil.ReadAll(resp.Body)
	if len(body) < 1024*10 {
		return false
	}
	bodystr := string(body)
	f, err := os.Create(fname)
	if err != nil {
		fmt.Println("Writing Error!", err)
		return false
	}
	defer f.Close()
	f.WriteString(bodystr)
	fmt.Println(fname)
	return true

}

func Process() bool {
	if _, err := os.Stat(dirname); err != nil {
		os.Mkdir(dirname, os.ModeDir)
	}
	bodystr := getIndexStr()
	pagenum := getPageNum(bodystr)
	//pagenum = 5

	for i := 1; i <= pagenum; i++ {
		fmt.Println("Page# ", i)
		lastbookname := ""
		bookpages := getPageList(i)
		downcnt := 0
		for j := 0; j < len(bookpages); j++ {
			curbook := bookpages[j][2]
			if lastbookname != curbook {
				filelink := getBookPage(curbook)
				if len(curbook) < 3 {
					continue
				}
				for _, fileline := range filelink {
					if saveFile(fileline) {
						lastbookname = curbook
						downcnt++
					}
				}
			}
		}
		if downcnt == 0 {
			break
		}

	}
	return true
}

func main() {
	Process()
}
