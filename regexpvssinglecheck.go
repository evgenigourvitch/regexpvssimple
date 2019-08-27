package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"
)

var (
	ifaRegex      = regexp.MustCompile("^(?i)\\{?[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}\\}?$")
	cMaxIFALength = 36
	gFactor       = 100
	cZerosIfa     = "00000000-0000-0000-0000-000000000000"
)

func main() {
	ifas, err := loadIFAs()
	if err != nil {
		fmt.Printf("got error: %+v\n", err)
		return
	}
	size := int64(len(ifas))
	start := time.Now().UnixNano()
	checkUsingRegexp(ifas)
	doneIn := time.Now().UnixNano() - start
	singleCheckTimeRegExp := float64(doneIn) / float64(size)
	fmt.Printf("RegExp check results: checked %d ifas in %d nannosecs, avg time for 1 check is %.2f\n", size, doneIn, singleCheckTimeRegExp)
	start = time.Now().UnixNano()
	checkSimple(ifas)
	doneIn = time.Now().UnixNano() - start
	singleCheckTimeSingle := float64(doneIn) / float64(size)
	fmt.Printf("Simple check results: checked %d ifas in %d nannosecs, avg time for 1 check is %.2f\n", size, doneIn, singleCheckTimeSingle)
	fmt.Printf("Diff: %.2f\n", singleCheckTimeSingle*100/singleCheckTimeRegExp)
}

func checkUsingRegexp(ifas []string) {
	for _, ifa := range ifas {
		isValidIFA(ifa)
	}
}

func checkSimple(ifas []string) {
	for _, ifa := range ifas {
		validateIFA(ifa)
	}
}

func loadIFAs() ([]string, error) {
	file, err := os.Open("ifas/ifas.small")
	if err != nil {
		return nil, err
	}
	res := []string{}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for i := 0; i < gFactor; i++ {
			res = append(res, scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func isValidIFA(ifa string) bool {
	if len(ifa) == 0 {
		return false
	}
	if ifa == cZerosIfa {
		return false
	}
	return ifaRegex.MatchString(ifa)
}

func validateIFA(ifa string) bool {
	ifaLength := len(ifa)
	if ifaLength != cMaxIFALength {
		return false
	}
	bytes := []byte(ifa)
	for i := 0; i < cMaxIFALength; i++ {
		ch := bytes[i]
		if (ch >= 48 && ch <= 57) || (ch >= 65 && ch <= 70) || (ch >= 97 && ch <= 102) {
			continue
		}
		if ch == 45 && (i == 8 || i == 13 || i == 18 || i == 23) {
			continue
		}
		return false
	}
	return true
}
