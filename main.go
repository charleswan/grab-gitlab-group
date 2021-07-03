package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"sync"

	"github.com/charleswan/grab-gitlab-group/config"
)

type groupProjectNameUnit struct {
	Name string `json:"name"`
}

type groupProjects struct {
	List []groupProjectNameUnit
}

func main() {
	log.Println("start working...")

	arrs := []string{}
	n := 0
	u := config.Get().GetGroupURL()

	for {
		if n > 0 {
			u = fmt.Sprintf("%s?page=%d", config.Get().GetGroupURL(), n)
		}
		arr := getProjectNameCore(u)
		if len(arr) == 0 {
			break
		}
		arrs = append(arrs, arr...)
		n++
	}

	log.Printf("got from gitlab: %v\n", len(arrs))

	var wg sync.WaitGroup
	wg.Add(len(arrs))

	l := NewLimiter(5)

	for _, v := range arrs {
		v := v
		go func() {
			for {
				if gitCloneProject(l, v) {
					wg.Done()
					break
				}
			}
		}()
	}
	wg.Wait()
}

func gitCloneProject(l *Limiter, name string) bool {
	if !l.Get() {
		return false
	}
	defer l.Release()
	log.Printf("cloning %s...", name)

	// git clone ssh://git@xxx.xxx.xxx:xxx/xxx/aaa.git /Users/xxx.xxx.xxx/xxx/aaa
	src := config.Get().GetGitPrefixURL(name)
	dest := config.Get().GetClonePath(name)
	cmd := exec.Command("git", "clone", src, dest)
	if err := cmd.Run(); err != nil {
		log.Printf("cmd.Run failed: %v\n", err)
	}

	return true
}

func getProjectNameCore(u string) []string {
	arr := []string{}

	client := &http.Client{}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Fatalf("http.NewRequest failed: %v", err)
		return arr
	}
	req.Header.Add("cookie", config.Get().GetCookie())
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("client.Do failed: %v", err)
		return arr
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll failed: %v", err)
		return arr
	}

	// log.Println(string(body))

	list := groupProjects{}
	if err := json.Unmarshal(body, &list.List); err != nil {
		log.Fatalf("json.Unmarshal failed: %v", err)
		return arr
	}

	// log.Printf("len = %d\n", len(list.List))

	for _, v := range list.List {
		arr = append(arr, v.Name)
	}

	return arr
}
