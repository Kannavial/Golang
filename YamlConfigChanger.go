package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	var yamlFiles []string
	var path string = ""

	fmt.Println("Enter the path of the directory: ")
	_, err := fmt.Scanln(&path)
	if err != nil {
		log.Fatal(err)
	}

	if CheckDirExistOrNot(path){
		CheckDirectoryOfChild(path, yamlFiles)
		SetYamlConfig(yamlFiles)
	}
}

func CheckDirExistOrNot(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
func HasYamlExtentension(filename string) bool {
	yamlExtensionPattern := `(?i)\.yaml$`
	match, _ := regexp.MatchString(yamlExtensionPattern, filename)
	return match
}
func CheckDirectoryOfChild(path string, yamlFiles []string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		var fileName string = file.Name()
		if HasYamlExtentension(fileName) {
			yamlFiles = append(yamlFiles, fileName)
		}
	}
	return nil
}
func SetYamlConfig(yamlFiles []string) error {
	var configDelimited string = "."

	var input string
	var yamlPath string
	fmt.Println("Enter the yaml config key using (.) as delimiter: ")
	_, errY := fmt.Scanln(&yamlPath)
	fmt.Println("Enter the input value for the key: ")
	_, errI := fmt.Scanln(&input)
	if errI != nil || errY != nil {
		log.Fatal(errI, errY)
	}
	parts := strings.Split(yamlPath, configDelimited)

	ChangeConfig(parts, input, yamlFiles)
	return errI
}
func ChangeConfig(configYamlPath []string, input string, yamlFiles []string) {
	for _, yamlFile := range yamlFiles {
		file, err := os.OpenFile(yamlFile, os.O_RDWR, 0644)
		if err != nil {
			log.Printf("Error opening file %s: %v", yamlFile, err)
			continue
		}

		var data map[string]interface{}
		byteData, err := ioutil.ReadAll(file)
		if err != nil {
			log.Printf("Error reading file %s: %v", yamlFile, err)
			file.Close()
			continue
		}

		err = yaml.Unmarshal(byteData, &data)
		if err != nil {
			log.Printf("Error unmarshalling file %s: %v", yamlFile, err)
			file.Close()
			continue
		}

		currentLevel := data
		for _, key := range configYamlPath[:len(configYamlPath)-1] {
			if nextLevel, ok := currentLevel[key].(map[interface{}]interface{}); ok {
				currentLevel = convertToStringMap(nextLevel)
			} else {
				log.Printf("Key %s not found in file %s", key, yamlFile)
				break
			}
		}

		currentLevel[configYamlPath[len(configYamlPath)-1]] = input

		newYamlData, err := yaml.Marshal(data)
		if err != nil {
			log.Printf("Error marshalling data for file %s: %v", yamlFile, err)
			file.Close()
			continue
		}

		err = file.Truncate(0)
		if err != nil {
			log.Printf("Error truncating file %s: %v", yamlFile, err)
			file.Close()
			continue
		}

		_, err = file.Seek(0, 0)
		if err != nil {
			log.Printf("Error seeking to beginning of file %s: %v", yamlFile, err)
			file.Close()
			continue
		}

		_, err = file.Write(newYamlData)
		if err != nil {
			log.Printf("Error writing to file %s: %v", yamlFile, err)
			file.Close()
			continue
		}

		fmt.Printf("Updated file: %s\n", yamlFile)
		file.Close()
	}
}
func convertToStringMap(input map[interface{}]interface{}) map[string]interface{} {
	output := make(map[string]interface{})
	for key, value := range input {
		if strKey, ok := key.(string); ok {
			output[strKey] = value
		}
	}
	return output
}
