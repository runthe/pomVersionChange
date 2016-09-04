package main

import (
	"github.com/beevik/etree"
	"fmt"
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

const (
	PomFile string = "pom.xml"
	projectTag string = "project"
	packagingTag string = "packaging"
	parentTag string = "parent"
	versionTag string = "version"
	packagingPom string = "pom"
	NewLine byte = '\n'
)

func main() {
	var inputDir, inputVersion string = readInput()

	if src, err := os.Stat(inputDir); err != nil {
		panic(err)
	} else if !src.IsDir() {
		fmt.Println("Not Directory")
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("inputDir", inputDir)
	fmt.Println("inputVersion", inputVersion)
	fmt.Println()

	for _, pomFilePath := range readPomFiles(inputDir) {
		doc := etree.NewDocument()

		if err := doc.ReadFromFile(pomFilePath); err != nil {
			panic(err)
		}

		//1.packing 확인 있으면 version 거기 찾음 없으면 다른거 찾기
		if root := doc.SelectElement(projectTag); root != nil {
			doc.SetRoot(root)

			if packaging := root.SelectElement(packagingTag); packaging != nil {
				fmt.Println("check::", packaging.Text())

				//main
				if packaging.Text() == packagingPom {
					if version := root.SelectElement(versionTag); version != nil {
						version.SetText(inputVersion)

						if err := doc.WriteToFile(pomFilePath); err != nil {
							panic(err)
						}
					}
					//sub
				} else {
					if parent := root.SelectElement(parentTag); parent != nil {
						fmt.Println(parentTag)
						//fmt.Println(parent.SelectElement(versionTag).Text())
						if version := parent.SelectElement(versionTag); version != nil {
							version.SetText(inputVersion)

							if err := doc.WriteToFile(pomFilePath); err != nil {
								panic(err)
							}
						}
					}
				}
			}
		}
		fmt.Println(pomFilePath)
	}

	fmt.Println()
	fmt.Println("pom.xml version change Success!")
}

/*
func readFile(dir string) []string {

}
*/

func readPomFiles(inputDir string) []string {
	files := []string{}
	err := filepath.Walk(inputDir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && f.Name() == PomFile {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return files
}

func readInput() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Write pom.xml Root Folder Dir And Version \n")
	fmt.Print("EX::/user/local/power 1.2.4.RELEASE \n")
	fmt.Println("")

	return readDirAndVersion(reader)
}

func readDirAndVersion(reader *bufio.Reader) (string, string) {
	fmt.Print("Enter Dir: ")
	inputDir, _ := reader.ReadString(NewLine)

	fmt.Print("Enter Version: ")
	inputVersion, _ := reader.ReadString(NewLine)

	//NewLine Remove
	inputDir = strings.TrimSpace(inputDir)
	inputVersion = strings.TrimSpace(inputVersion)

	//Todo::재귀호출 위치가 잘못된거 같다.
	if inputDir == "" || inputVersion == "" {
		return readDirAndVersion(reader)
	}

	return inputDir, inputVersion
}

/*
func writePomVersion(Document doc, Element element, string inputText) {
	if version := element.SelectElement(""); version != nil {
		version.SetText(inputText)

		if err := doc.WriteToFile(""); err != nil {
			panic(err)
		}
	}
}
*/
