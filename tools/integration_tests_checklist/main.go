package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	cmdName := "git"
	cmdArgs := []string{"rev-parse", "--verify", "HEAD"}
	cmdOut, err := exec.Command(cmdName, cmdArgs...).Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running git rev-parse command: ", err)
		os.Exit(1)
	}
	commit := strings.TrimSpace(string(cmdOut))

	testsFile, err := os.Open("../../integration_tests/main_test.go")
	if err != nil {
		log.Fatal(err)
	}
	defer testsFile.Close()

	functionNamePattern := `Test[a-zA-Z]+`
	checklistNamePattern := `\(Test[a-zA-Z]+\)`
	validationPattern := fmt.Sprintf(`func %s\(t \*testing\.T\) \{`, functionNamePattern)
	commitPattern := `[0-9a-f]{40}`
	funcValidator := regexp.MustCompile(validationPattern)
	checklistValidator := regexp.MustCompile(checklistNamePattern)
	nameValidator := regexp.MustCompile(functionNamePattern)
	commitValidator := regexp.MustCompile(commitPattern)

	functionNamesToLineNumberMap := map[string]int{}

	lineNumber := 0
	scanner := bufio.NewScanner(testsFile)
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if funcValidator.MatchString(line) {
			functionName := string(nameValidator.Find([]byte(line)))
			functionNamesToLineNumberMap[functionName] = lineNumber
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	checklistFilePath := "../../integration_tests/README.md"
	newChecklistFilePath := "../../integration_tests/README_linked.md"

	checklistFile, err := os.Open(checklistFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer checklistFile.Close()

	newChecklistFile, err := os.Create(newChecklistFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer newChecklistFile.Close()

	lineNumber = 0
	scanner = bufio.NewScanner(checklistFile)
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if checklistValidator.MatchString(line) {
			functionName := string(nameValidator.Find([]byte(line)))
			if _, ok := functionNamesToLineNumberMap[functionName]; !ok {
				log.Fatalf("encountered function name %s that doesn't have a corresponding entry in the line number map", functionName)
			}
			link := fmt.Sprintf(`[%s](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/%s/integration_tests/main_test.go#L%d)`, functionName, commit, functionNamesToLineNumberMap[functionName])
			newLine := strings.Replace(line, functionName, link, 1)
			_, err := newChecklistFile.WriteString(fmt.Sprintf("%s\n", newLine))
			if err != nil {
				log.Fatal(err)
			}
		} else if commitValidator.MatchString(line) {
			oldCommit := string(commitValidator.Find([]byte(line)))
			newLine := strings.Replace(line, oldCommit, commit, 1)
			_, err := newChecklistFile.WriteString(fmt.Sprintf("%s\n", newLine))
			if err != nil {
				log.Fatal(err)
			}
		} else {
			n, err := newChecklistFile.WriteString(fmt.Sprintf("%s\n", line))
			if err != nil {
				log.Fatalf("%v, %v", n, err)
			}
		}
	}

	err = os.Remove(checklistFilePath)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Rename(newChecklistFilePath, checklistFilePath)
	if err != nil {
		log.Fatal(err)
	}
}
