package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/ini.v1"
)

func main() {
	flag.Parse()
	a := flag.Args()
	if len(a) == 0 {
		must("please specify at least 1 argument")
	}

	switch a[0] {
	case "open":
		if err := open(); err != nil {
			must(err.Error())
		}
	}
}

func must(msg string, more ...interface{}) {
	fmt.Printf(msg, more...)
	fmt.Println()
	os.Exit(1)
}

func open() error {
	url, err := buildURL()
	if err != nil {
		return err
	}

	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func buildURL() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	root, err := getGitRoot(dir)
	if err != nil {
		return "", err
	}

	t, err := readConfig(root)
	if err != nil {
		return "", err
	}

	giturl, err := getURL(t)
	if err != nil {
		return "", err
	}
	return parseGitURL(giturl)
}

func getGitRoot(path string) (string, error) {
	var y bool
	var err error

	p := path

	for !y {
		if path == "/" {
			break
		}
		y, err = isGitRoot(p)
		if err != nil {
			return "", err
		}
		if y {
			return p, nil
		}
		p = filepath.Dir(p)
	}

	return "", fmt.Errorf("no git repository found in %s", p)
}

func isGitRoot(p string) (bool, error) {
	i, err := ioutil.ReadDir(p)
	if err != nil {
		return false, err
	}

	for _, f := range i {
		if f.Name() == ".git" {
			return true, nil
		}
	}

	return false, nil
}

func readConfig(path string) (*ini.File, error) {
	path = fmt.Sprintf("%s/.git/config", path)
	return readConfigFile(path)
}

func readConfigFile(path string) (*ini.File, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func getURL(f *ini.File) (string, error) {
	r, err := f.GetSection(`remote "origin"`)
	if err != nil {
		return "", errors.New("no remote found in git config")
	}
	u := r.Key("url").String()
	if len(u) == 0 {
		return "", errors.New(`remote "origin" does not seem to have a URL`)
	}
	return u, nil
}

func parseGitURL(u string) (string, error) {
	u = u[strings.Index(u, "@")+1:]
	u = strings.Replace(u, ":", "/", -1)
	u = u[0:strings.Index(u, ".git")]
	return fmt.Sprintf("https://%s", u), nil
}
