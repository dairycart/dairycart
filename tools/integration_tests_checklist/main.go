package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	checklistNamePattern = `(\()(\[)?Test[a-zA-Z]+(\]?)`
	functionNamePattern  = `Test[a-zA-Z]+`

	// This tool should be run from two directories up.
	testsFolder          = "integration_tests"
	checklistFilePath    = "integration_tests/README.md"
	newChecklistFilePath = "integration_tests/README_linked.md"
)

type Result struct {
	Filename string
	Line     int
}

func failIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func generateFunctionToLineNumberMap(file *os.File, filename string, nameValidator *regexp.Regexp) map[string]Result {
	functionNamesToLineNumberMap := map[string]Result{}
	validationPattern := fmt.Sprintf(`func %s\(t \*testing\.T\) \{`, functionNamePattern)
	funcValidator := regexp.MustCompile(validationPattern)
	lineNumber := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if funcValidator.MatchString(line) {
			functionName := string(nameValidator.Find([]byte(line)))
			functionNamesToLineNumberMap[functionName] = Result{
				Filename: filename,
				Line:     lineNumber,
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return functionNamesToLineNumberMap
}

func replaceLinksInChecklistFile(old *os.File, new *os.File, nameValidator *regexp.Regexp, functionNamesToLineNumberMap map[string]Result) {
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
			link := fmt.Sprintf(`([%s](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/%s.go#L%d))`, functionName, functionNamesToLineNumberMap[functionName].Filename, functionNamesToLineNumberMap[functionName].Line)
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
	files, err := ioutil.ReadDir(testsFolder)
	if err != nil {
		log.Fatal(err)
	}

	functionNamesToLineNumberMap := map[string]Result{}
	nameValidator := regexp.MustCompile(functionNamePattern)
	for _, f := range files {
		fileName := f.Name()
		if f.IsDir() || !strings.HasSuffix(fileName, "test.go") {
			continue
		}
		actualPath := fmt.Sprintf("%s/%s", testsFolder, fileName)
		log.Printf("working on %s", fileName)
		file, err := os.Open(actualPath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		for k, v := range generateFunctionToLineNumberMap(file, fileName, nameValidator) {
			functionNamesToLineNumberMap[k] = v
		}
	}

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
		for f, ln := range functionNamesToLineNumberMap {
			log.Printf("\t%s (line %d)\n", f, ln.Line)
		}
		log.Fatalf("There are %d tests which are unaccounted for in the README", len(functionNamesToLineNumberMap))
	}
}
