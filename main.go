package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var inputs = [][]string{
	{"cl", "class", "infrastructure/utils"},
	{"co", "constant", "domain/constants"},
	{"c", "component", "ui/components"},
	{"l", "component", "ui/layouts"},
	{"lc", "component", "ui/layouts/components"},
	{"v", "component", "ui/views"},
	{"vc", "component", "ui/views"},
	{"d", "directive", "infrastructure/directives"},
	{"e", "enum", "domain/enums"},
	{"g", "guard", "infrastructure/guards"},
	{"ic", "interceptor", "infrastructure/interceptors"},
	{"i", "interface", "domain/interfaces"},
	{"p", "pipe", "infrastructure/pipes"},
	{"s", "service", "infrastructure/services"},
}

var skipIndexed = []int{3, 5}
var updateRouter = 5

func main() {
	cmd := os.Args

	for i, e := range inputs {

		if cmd[1] == e[0] || cmd[1] == e[1] {

			if len(cmd[2]) <= 3 {
				fmt.Println(e[1] + " valid name is required")
				break
			}

			p := e[2] + "/" + cmd[2]

			ex := exec.Command("ng", "g", e[1], p)

			err := ex.Run()

			if err != nil {
				fmt.Println(err)
				break
			}

			if i == updateRouter {
				fmt.Println("update router")
			}

			if i != skipIndexed[0] && i != skipIndexed[1] {
				root := "./src/app/"
				path := strings.Split(p, "/")
				path[len(path)-1] = ""

				indexedDirectory(root + strings.Join(path, "/"))
			}

		}
	}
}

func indexedDirectory(p string) {
	idx := "index.ts"

	var paths []string

	err := filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Name() == idx || d.Name() == "." {
			return nil
		}

		ts := ".ts"

		if strings.Contains(d.Name(), ts) {
			root := strings.Replace(p, "./", "", 1)
			host := strings.Replace(path, root, "", 1)
			host = strings.Replace(host, ts, "", 1)
			line :=
				"export * from './" +
					host +
					"';"

			paths = append(paths, line)
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	paths = append(paths, "")
	bytes := []byte(strings.Join(paths, "\n"))

	filePath := p + idx

	idxFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
	}

	defer idxFile.Close()

	_, err = idxFile.Seek(0, 0)
	if err != nil {
		fmt.Println(err)
	}

	err = idxFile.Truncate(0)
	if err != nil {
		fmt.Println(err)
	}

	writer := bufio.NewWriter(idxFile)
	_, err = writer.Write(bytes)
	if err != nil {
		fmt.Println(err)
	}

	err = writer.Flush()
	if err != nil {
		fmt.Println(err)
	}
}
