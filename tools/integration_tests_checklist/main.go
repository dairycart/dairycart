package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	checklistNamePattern = `(\()(\[)?Test[a-zA-Z]+(\]?)`
	functionNamePattern  = `Test[a-zA-Z]+`

	// This tool should be run from two directories up.
	testsFilePath        = "integration_tests/main_test.go"
	checklistFilePath    = "integration_tests/README.md"
	newChecklistFilePath = "integration_tests/README_linked.md"
)

func failIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func generateFunctionToLineNumberMap(integrationTestsFile *os.File, nameValidator *regexp.Regexp) map[string]int {
	functionNamesToLineNumberMap := map[string]int{}
	validationPattern := fmt.Sprintf(`func %s\(t \*testing\.T\) \{`, functionNamePattern)
	funcValidator := regexp.MustCompile(validationPattern)
	lineNumber := 0

	scanner := bufio.NewScanner(integrationTestsFile)
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
	return functionNamesToLineNumberMap
}

func replaceLinksInChecklistFile(old *os.File, new *os.File, nameValidator *regexp.Regexp, functionNamesToLineNumberMap map[string]int) {
	checklistValidator := regexp.MustCompile(checklistNamePattern)

	lineNumber := 0
	scanner := bufio.NewScanner(old)
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if checklistValidator.MatchString(line) {
			functionName := string(nameValidator.Find([]byte(line)))
			if _, ok := functionNamesToLineNumberMap[functionName]; !ok {
				log.Fatalf("encountered function name %s that doesn't have a corresponding entry in the line number map", functionName)
			}
			link := fmt.Sprintf(`([%s](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L%d))`, functionName, functionNamesToLineNumberMap[functionName])
			checklistPart := strings.Split(line, "(")[0]
			newLine := fmt.Sprintf("%s%s", checklistPart, link)
			_, err := new.WriteString(fmt.Sprintf("%s\n", newLine))
			failIfErr(err)
			delete(functionNamesToLineNumberMap, functionName)
		} else {
			_, err := new.WriteString(fmt.Sprintf("%s\n", line))
			failIfErr(err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	testsFile, err := os.Open(testsFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer testsFile.Close()

	nameValidator := regexp.MustCompile(functionNamePattern)
	functionNamesToLineNumberMap := generateFunctionToLineNumberMap(testsFile, nameValidator)

	checklistFile, err := os.Open(checklistFilePath)
	failIfErr(err)
	defer checklistFile.Close()

	newChecklistFile, err := os.Create(newChecklistFilePath)
	failIfErr(err)
	defer newChecklistFile.Close()

	replaceLinksInChecklistFile(checklistFile, newChecklistFile, nameValidator, functionNamesToLineNumberMap)
	failIfErr(os.Remove(checklistFilePath))
	failIfErr(os.Rename(newChecklistFilePath, checklistFilePath))

	if len(functionNamesToLineNumberMap) != 0 {
		log.Fatalf("Tests exist which are unaccounted for in the README: ")
		for f, ln := range functionNamesToLineNumberMap {
			log.Printf("\t%s (line %d)\n", f, ln)
		}
	}
}
