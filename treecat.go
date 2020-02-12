package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
)

var (
	Info = Teal
	Warn = Yellow
	Fata = Red
)

var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
)

type Counter struct {
	dirs  int
	files int
}

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

func (counter *Counter) index(path string) bool {
	stat, _ := os.Stat(path)
	if stat.IsDir() {
		counter.dirs += 1
		return false
	} else {
		counter.files += 1
		return true
	}
}

func (counter *Counter) output() string {
	return fmt.Sprintf("\n%d directories, %d files", counter.dirs, counter.files)
}

func dirnamesFrom(base string) []string {
	file, err := os.Open(base)
	if err != nil {
		fmt.Println(err)
	}

	names, _ := file.Readdirnames(0)
	file.Close()

	sort.Strings(names)
	return names
}

func fsPrint(isFile bool, name string) string {
	if isFile {
		return Teal(name)
	} else {
		return name
	}
}

func catPrint(isFile bool, path string, prefix string) {
	if isFile {
		contents := cat(path)
		if len(contents) > 0 {
			fmt.Print(prefix + Yellow("  >> "+cat(path)))
		}
	}
}

func tree(counter *Counter, base string, prefix string) {
	names := dirnamesFrom(base)

	for index, name := range names {
		if name[0] == '.' {
			continue
		}

		subpath := path.Join(base, name)
		isFile := counter.index(subpath)

		if index == len(names)-1 {
			fmt.Println(prefix+"└──", fsPrint(isFile, name))
			catPrint(isFile, subpath, prefix+"    ")
			tree(counter, subpath, prefix+"    ")
		} else {
			fmt.Println(prefix+"├──", fsPrint(isFile, name))
			catPrint(isFile, subpath, prefix+"│   ")
			tree(counter, subpath, prefix+"│   ")
		}
	}
}

func cat(file string) string {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Print(err)
	}

	return string(bytes)
}

func main() {
	var directory string
	if len(os.Args) > 1 {
		directory = os.Args[1]
	} else {
		directory = "."
	}

	counter := new(Counter)
	fmt.Println(directory)

	tree(counter, directory, "")
	fmt.Println(counter.output())
}
