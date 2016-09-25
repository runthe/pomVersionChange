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
	ProjectTag string = "project"
	PackagingTag string = "packaging"
	ParentTag string = "parent"
	VersionTag string = "version"
	PackagingPom string = "pom"
	NewLine byte = '\n'
)

func main() {
	var inputDir, inputVersion string = readInput()

	checkDirectory(inputDir)
	printDirAndVersion(inputDir, inputVersion)

	for _, pomFilePath := range readPomFiles(inputDir) {
		doc := etree.NewDocument()

		if err := doc.ReadFromFile(pomFilePath); err != nil {
			panic(err)
		}

		//1.packing 확인 있으면 version 거기 찾음 없으면 다른거 찾기
		if root := doc.SelectElement(ProjectTag); root != nil {
			doc.SetRoot(root)

			if packaging := root.SelectElement(PackagingTag); packaging != nil {
				//main
				if packaging.Text() == PackagingPom {
					if version := root.SelectElement(VersionTag); version != nil {
						version.SetText(inputVersion)

						if err := doc.WriteToFile(pomFilePath); err != nil {
							panic(err)
						}
					}
					//sub
				} else {
					if parent := root.SelectElement(ParentTag); parent != nil {
						if version := parent.SelectElement(VersionTag); version != nil {
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

func checkDirectory(inputDir string) {
	if src, err := os.Stat(inputDir); err != nil {
		panic(err)
	} else if !src.IsDir() {
		fmt.Println("Not Directory")
		os.Exit(1)
	}
}

func printDirAndVersion(inputDir, inputVersion string) {
	fmt.Println()
	fmt.Println("inputDir", inputDir)
	fmt.Println("inputVersion", inputVersion)
	fmt.Println()
}

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

	inputDir := readDir(reader);
	inputVersion := readVersion(reader);

	return inputDir, inputVersion
}

func readDir(reader *bufio.Reader) (string) {
	fmt.Print("Enter Dir: ")
	inputDir, _ := reader.ReadString(NewLine)
	inputDir = strings.TrimSpace(inputDir)

	if inputDir == "" {
		return readDir(reader)
	}

	return inputDir
}

func readVersion(reader *bufio.Reader) (string) {
	fmt.Print("Enter Version: ")
	inputVersion, _ := reader.ReadString(NewLine)
	inputVersion = strings.TrimSpace(inputVersion)

	if inputVersion == "" {
		return readVersion(reader)
	}

	return inputVersion
}
