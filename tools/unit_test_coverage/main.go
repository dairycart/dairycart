package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	functionDeclarationPattern = `func ([a-zA-Z]+)\(`
)

func failIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func findFunctionsInFile(f *os.File) []string {
	funcValidator := regexp.MustCompile(functionDeclarationPattern)

	functionNames := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if funcValidator.MatchString(line) {
			submatches := funcValidator.FindStringSubmatch(line)
			functionName := submatches[1]
			if functionName != `init` {
				functionNames = append(functionNames, functionName)
			}
		}
	}
	return functionNames
}

func findFunctionReferencesInTestFile(f *os.File, functionNames []string) {
	scanner := bufio.NewScanner(f)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		for _, name := range functionNames {
			if strings.Contains(line, name) {
				log.Printf("function %s referenced on line %d of %s\n", name, lineNumber, f.Name())
				// delete that value from the slice here
			}
		}
	}
}

func main() {
	codeFilesToTestFilesMap := map[string]string{
		"api/helpers.go":                  "api/helpers_test.go",
		"api/product_attribute_values.go": "api/product_attribute_values_test.go",
		"api/product_attributes.go":       "api/product_attributes_test.go",
		"api/product_progenitors.go":      "api/product_progenitors_test.go",
		"api/products.go":                 "api/products_test.go",
		"api/queries.go":                  "api/queries_test.go",
	}
	functionsInEachFile := map[string][]string{}

	for codeFile, testFile := range codeFilesToTestFilesMap {
		cf, err := os.Open(codeFile)
		failIfErr(err)

		functionNames := findFunctionsInFile(cf)
		functionsInEachFile[codeFile] = functionNames

		tf, err := os.Open(testFile)
		failIfErr(err)
		findFunctionReferencesInTestFile(tf, functionNames)
	}

	log.Printf(`
	functionsInEachFile: %+v
	`, functionsInEachFile)
}
