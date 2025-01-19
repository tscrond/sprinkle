package lib

import (
	"encoding/json"
	"fmt"
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
