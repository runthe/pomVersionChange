package main

import (
	"github.com/beevik/etree"
	"fmt"
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"strconv"
)

var reader = bufio.NewReader(os.Stdin)

const (
	PomFile string = "pom.xml"
	ProjectTag string = "project"
	PackagingTag string = "packaging"
	ParentTag string = "parent"
	VersionTag string = "version"
	PackagingPom string = "pom"
	NewLine byte = '\n'

	ChoiceOne = "1"
	ChoiceTwo = "2"
	ChoiceThree = "3"

	numbers string = "0123456789"
)

func main() {
	pwd := getLocalDirectory()
	localPomFiles := readPomFiles(pwd)

	fmt.Println("==BySoo==")

	if len(localPomFiles) == 0 {
		changePomversionByManual()
	} else {
		//Root Version만 읽자
		var choice string
		var nextVersion string
		rootVersion := readPomFilesRootVersion(localPomFiles)

		//Version이 없으면 수동으로 넘긴다.
		if len(strings.TrimSpace(rootVersion)) != 0 {
			nextPatchVersion, err := calculateNextPatchVersion(rootVersion)

			if err != nil {
				panic(err)
			}

			showChoices(rootVersion, nextPatchVersion)
			choice = readChoice()
			nextVersion = nextPatchVersion
		} else {
			choice = ChoiceThree
		}

		if choice == ChoiceOne {
			changePomFilesVersion(pwd, nextVersion)

		} else if choice == ChoiceTwo {
			inputVersion := readVersion()
			changePomFilesVersion(pwd, inputVersion)
		} else if choice == ChoiceThree {
			changePomversionByManual()
		}
	}
}

func hasLeadingZeroes(s string) bool {
	return len(s) > 1 && s[0] == '0'
}

func containsOnly(s string, set string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		return !strings.ContainsRune(set, r)
	}) == -1
}

//https://github.com/blang/semver/blob/master/semver.go
func calculateNextPatchVersion(rootVersion string) (nextVersion string, err error) {
	// Split into major.minor.(patch+pr+meta)
	parts := strings.SplitN(rootVersion, ".", 3)

	if len(parts) != 3 {
		return rootVersion, fmt.Errorf("No Major.Minor.Patch elements found")
	}

	// Major
	if !containsOnly(parts[0], numbers) {
		return rootVersion, fmt.Errorf("Invalid character(s) found in major number %q", parts[0])
	}
	if hasLeadingZeroes(parts[0]) {
		return rootVersion, fmt.Errorf("Major number must not contain leading zeroes %q", parts[0])
	}

	major, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return rootVersion, err
	}

	// Minor
	if !containsOnly(parts[1], numbers) {
		return rootVersion, fmt.Errorf("Invalid character(s) found in minor number %q", parts[1])
	}
	if hasLeadingZeroes(parts[1]) {
		return rootVersion, fmt.Errorf("Minor number must not contain leading zeroes %q", parts[1])
	}

	minor, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return rootVersion, err
	}

	var build, prerelease []string
	patchStr := parts[2]

	if buildIndex := strings.IndexRune(patchStr, '+'); buildIndex != -1 {
		build = strings.Split(patchStr[buildIndex + 1:], ".")
		patchStr = patchStr[:buildIndex]
	}

	if preIndex := strings.IndexRune(patchStr, '-'); preIndex != -1 {
		prerelease = strings.Split(patchStr[preIndex + 1:], ".")
		patchStr = patchStr[:preIndex]
	}

	if !containsOnly(patchStr, numbers) {
		return rootVersion, fmt.Errorf("Invalid character(s) found in patch number %q", patchStr)
	}
	if hasLeadingZeroes(patchStr) {
		return rootVersion, fmt.Errorf("Patch number must not contain leading zeroes %q", patchStr)
	}

	patch, err := strconv.ParseUint(patchStr, 10, 64)
	if err != nil {
		return rootVersion, err
	}

	// Version to string
	b := make([]byte, 0, 5)
	b = strconv.AppendUint(b, major, 10)
	b = append(b, '.')
	b = strconv.AppendUint(b, minor, 10)
	b = append(b, '.')
	b = strconv.AppendUint(b, patch + 1, 10)

	if len(prerelease) > 0 {
		b = append(b, '-')
		b = append(b, prerelease[0]...)
	}

	if len(build) > 0 {
		b = append(b, '+')
		b = append(b, build[0]...)

		for _, build := range build[1:] {
			b = append(b, '.')
			b = append(b, build...)
		}
	}

	return string(b), nil
}

func readPomFilesRootVersion(pomFiles []string) (rootVersion string) {
	for _, pomFilePath := range pomFiles {
		doc := etree.NewDocument()

		if err := doc.ReadFromFile(pomFilePath); err != nil {
			panic(err)
		}

		if root := doc.SelectElement(ProjectTag); root != nil {
			doc.SetRoot(root)

			if packaging := root.SelectElement(PackagingTag); packaging != nil {
				//main
				if packaging.Text() == PackagingPom {
					if version := root.SelectElement(VersionTag); version != nil {
						rootVersion = version.Text()
					}
				}
			}

			break;
		}
	}

	return
}


//Todo:: 올바른 메소드명일까...
func changePomversionByManual() {
	var inputDir, inputVersion string = readInput()

	checkDirectory(inputDir)
	printDirAndVersion(inputDir, inputVersion)
	changePomFilesVersion(inputDir, inputVersion)
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
	fmt.Print("Write pom.xml Root Folder Dir And Version \n")
	fmt.Print("EX::/user/local/power 1.2.4.RELEASE \n")
	fmt.Println("")

	inputDir := readDir();
	inputVersion := readVersion();

	return inputDir, inputVersion
}

func readDir() (string) {
	fmt.Print("Enter Dir: ")
	inputDir, _ := reader.ReadString(NewLine)
	inputDir = strings.TrimSpace(inputDir)

	if inputDir == "" {
		return readDir()
	}

	return inputDir
}

func readVersion() (string) {
	fmt.Print("Enter Version: ")
	inputVersion, _ := reader.ReadString(NewLine)
	inputVersion = strings.TrimSpace(inputVersion)

	if inputVersion == "" {
		return readVersion()
	}

	return inputVersion
}

func getLocalDirectory() (pwd string) {
	pwd, err := os.Getwd()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return pwd
}

func showChoices(pomVersion, nextPatchVersion string) {
	fmt.Println(fmt.Sprintf("local pom.xml version: %s", pomVersion))
	fmt.Printf("Choice 1: auto next patch version: %s", nextPatchVersion)
	fmt.Print(", ")
	fmt.Print("Choice 2: input next version")
	fmt.Print(", ")
	fmt.Print("Choice 3: Manual")
	fmt.Println()
}

func readChoice() (string) {
	fmt.Print("Enter Choice 1 OR 2 OR 3: ")

	inputChoice, _ := reader.ReadString(NewLine)
	inputChoice = strings.TrimSpace(inputChoice)

	if inputChoice != ChoiceOne && inputChoice != ChoiceTwo && inputChoice != ChoiceThree {
		return readChoice()
	}

	fmt.Println()
	return inputChoice
}

func changePomFilesVersion(inputDir, inputVersion string) {
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
	fmt.Println("pom.xml files version change Success!")
}
