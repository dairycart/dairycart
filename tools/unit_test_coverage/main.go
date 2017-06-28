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
			if functionName != "init" {
				functionNames = append(functionNames, functionName)
			}
		}
	}
	return functionNames
}

func sliceIndex(item string, slice []string) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	// this should never happen at all
	return -1
}

func findFunctionReferencesInTestFile(f *os.File, functionNames []string) []string {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		for _, name := range functionNames {
			if strings.Contains(line, name) {
				i := sliceIndex(name, functionNames)
				functionNames = append(functionNames[:i], functionNames[i+1:]...)
				break
			}
		}
	}
	return functionNames
}

func main() {
	/*
		Should this be derived from a walk func? yeah probably, but I hate writing them
		and this code doesn't matter anyway.
	*/
	codeFilesToTestFilesMap := map[string]string{
		"api/helpers.go":               "api/helpers_test.go",
		"api/product_option_values.go": "api/product_option_values_test.go",
		"api/product_options.go":       "api/product_options_test.go",
		"api/products.go":              "api/products_test.go",
		"api/queries.go":               "api/queries_test.go",
		"api/discounts.go":             "api/discounts_test.go",
	}

	functionsInEachFile := map[string][]string{}
	for codeFile, testFile := range codeFilesToTestFilesMap {
		cf, err := os.Open(codeFile)
		failIfErr(err)

		functionNames := findFunctionsInFile(cf)
		functionsInEachFile[codeFile] = functionNames

		tf, err := os.Open(testFile)
		failIfErr(err)
		remainingFunctionNames := findFunctionReferencesInTestFile(tf, functionNames)

		if len(remainingFunctionNames) > 0 {
			fmt.Printf("untested functions in %s:\n\t%s\n\n", testFile, strings.Join(remainingFunctionNames, "\n\t"))
		}
	}
}
