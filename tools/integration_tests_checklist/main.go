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

type result struct {
	Filename        string
	FunctionName    string
	DeclarationLine int
	CompletionLine  int
	UseCount        int
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func generateFunctionToLineNumberMap(file *os.File, filename string, nameValidator *regexp.Regexp) map[string]*result {
	functionNamesToLineNumberMap := map[string]*result{}
	validationPattern := fmt.Sprintf(`func %s\(t \*testing\.T\) \{`, functionNamePattern)
	funcValidator := regexp.MustCompile(validationPattern)
	lineNumber := 0

	scanner := bufio.NewScanner(file)
	currentResult := &result{}

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if line == "}" && currentResult.FunctionName != "" {
			// log.Printf("closing function bracket found on line %d of %s", lineNumber, file.Name())
			currentResult.CompletionLine = lineNumber
			functionNamesToLineNumberMap[currentResult.FunctionName] = currentResult
			currentResult = &result{}
		}
		if funcValidator.MatchString(line) {
			functionName := string(nameValidator.Find([]byte(line)))
			currentResult = &result{
				FunctionName:    functionName,
				Filename:        filename,
				DeclarationLine: lineNumber,
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return functionNamesToLineNumberMap
}

func replaceLinksInChecklistFile(old *os.File, new *os.File, nameValidator *regexp.Regexp, functionNamesToLineNumberMap map[string]*result) {
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
			result := functionNamesToLineNumberMap[functionName]
			link := fmt.Sprintf(`([%s](https://github.com/dairycart/dairycart/blob/master/integration_tests/%s#L%d-L%d))`, functionName, result.Filename, result.DeclarationLine, result.CompletionLine)
			checklistPart := strings.Split(line, "(")[0]
			newLine := fmt.Sprintf("%s%s", checklistPart, link)
			_, err := new.WriteString(fmt.Sprintf("%s\n", newLine))
			must(err)
			result.UseCount++
		} else {
			_, err := new.WriteString(fmt.Sprintf("%s\n", line))
			must(err)
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

	functionNamesToLineNumberMap := map[string]*result{}
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
	must(err)
	defer checklistFile.Close()

	newChecklistFile, err := os.Create(newChecklistFilePath)
	must(err)
	defer newChecklistFile.Close()

	replaceLinksInChecklistFile(checklistFile, newChecklistFile, nameValidator, functionNamesToLineNumberMap)
	must(os.Remove(checklistFilePath))
	must(os.Rename(newChecklistFilePath, checklistFilePath))

	if len(functionNamesToLineNumberMap) != 0 {
		missingDeclarations := 0
		for f, ln := range functionNamesToLineNumberMap {
			if ln.UseCount == 0 {
				missingDeclarations++
				log.Printf("\t%s (line %d)\n", f, ln.DeclarationLine)
			}
		}
		if missingDeclarations != 0 {
			log.Fatalf("There are %d tests which are unaccounted for in the README", missingDeclarations)
		}
	}
}
