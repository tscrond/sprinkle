package lib

import (
	"encoding/json"
	"fmt"
	"strings"
)

// PrettyPrintStruct takes any struct and prints it in a human-readable format.
func PrettyPrintStruct(v interface{}) {
	prettyJSON, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("Error pretty-printing struct: %v\n", err)
		return
	}
	fmt.Println(string(prettyJSON))
}

func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

// trimSuffixAfter removes everything after (and including) the first occurrence of the separator
func TrimSuffixAfter(s, sep string) string {
	if idx := strings.Index(s, sep); idx != -1 {
		return s[:idx]
	}
	return s
}

// trimLastSuffixAfter removes everything after (and including) the last occurrence of the separator
func TrimLastSuffixAfter(s, sep string) string {
	if idx := strings.LastIndex(s, sep); idx != -1 {
		return s[:idx]
	}
	return s
}
