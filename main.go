package main

import (
	"fmt"
	"encoding/json"
	"reflect"
	"io/ioutil"
	"os"
	"flag"
)

func CompareJSON(jsonA, jsonB string) (bool, error) {
	var iJsonA interface{}
	var iJsonB interface{}
	var err error

	err = json.Unmarshal([]byte(jsonA), &iJsonA)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal([]byte(jsonB), &iJsonB)
	if err != nil {
		return false, err
	}

	return Compare(iJsonA, iJsonB), nil
}


func Compare(jsonA, jsonB interface{}) bool {
	if reflect.TypeOf(jsonA) != reflect.TypeOf(jsonB) {
		return false
	}

	switch a := jsonA.(type) {
	case map[string]interface{}:
		b := jsonB.(map[string]interface{})

		if len(a) != len(b) {
			return false
		}

		for k, v := range a {
			val2 := b[k]

			if (v == nil) != (val2 == nil) {
				return false
			}

			if !Compare(v, val2) {
				return false
			}
		}

		return true
	case []interface{}:
		b := jsonB.([]interface{})

		if len(a) != len(b) {
			return false
		}

		var matches int
		flagged := make([]bool, len(b))
		for _, v := range a {
			for i, v2 := range b {
				if Compare(v, v2) && !flagged[i] {
					matches++
					flagged[i] = true

					break
				}
			}
		}

		return matches == len(a)
	default:
		return jsonA == jsonB
	}
}



func JsonFileToStr(path string) (string, error){
	fileJson, err := ioutil.ReadFile(path)
	if err != nil {
		return "",err
	}
	return fmt.Sprintf("%s", fileJson), nil
}


func main(){
	file1 := "a.json"
	file2 := "b.json"
	flag.Parse()
	if len(flag.Args()) == 1 {
		fmt.Printf("Comparing %s with a.json ", flag.Args()[0]);
		file2 = flag.Args()[0] 
	}
	if len(flag.Args()) == 2  {
		fmt.Printf("Comparing %s with %s ", flag.Args()[0], flag.Args()[0]);
		file1 = flag.Args()[0] 
		file2 = flag.Args()[1] 
	}
	if len(flag.Args())<1 {
		fmt.Printf("Comparing a.json with b.json ");
	}
	fmt.Println();
	
	// check if file 1 exists
	if _, err := os.Stat(file1); os.IsNotExist(err) {
		fmt.Println(err);
		os.Exit(0)
	}
	// check if file 2 exists
	if _, err := os.Stat(file2); os.IsNotExist(err) {
		fmt.Println(err);
		os.Exit(0)
	}

	
	file1String, err := JsonFileToStr(file1)
	if err != nil {
		fmt.Println(err)
	}


	file2String, err := JsonFileToStr(file2)
	if err != nil {
		fmt.Println(err)
	}

	testCompare, err := CompareJSON(file1String, file2String)
	if err != nil {
		fmt.Println(err)
	}else{
		fmt.Println(testCompare)
	}
	



}
