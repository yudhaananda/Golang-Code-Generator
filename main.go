package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	objs := map[string][]string{}

	content, err := os.ReadFile("project.json")
	if err != nil {
		log.Fatal(err)
	}

	var jsonContent map[string]interface{}
	err = json.Unmarshal(content, &jsonContent)

	if err != nil {
		log.Fatal(err)
	}

	projectObject := jsonContent["projectName"]
	entity := jsonContent["entity"]

	project := strings.ToLower(fmt.Sprintf("%v", projectObject))

	for key, obj := range entity.(map[string]interface{}) {
		temp := []string{}
		for _, val := range obj.([]interface{}) {
			temp = append(temp, val.(string))
		}
		objs[key] = temp
	}

	for key, obj := range objs {
		createEntity(obj, key, project)
		createRepository(obj, key, project)
		createService(obj, key, project)
		createInput(obj, key, project)
		createHandler(obj, key, project)
	}
	createHelper(project)
	createMain(objs, project)
}

func createMain(objs map[string][]string, project string) {
	err := os.MkdirAll(project, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(project + "\\main.go")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fileTemplate, err := os.ReadFile("template\\main.go")

	if err != nil {
		log.Fatal(err)
	}

	migrateArea := ""
	for key := range objs {
		keys := strings.Split(key, "")
		keys[0] = strings.ToUpper(keys[0])
		keyUpper := strings.Join(keys, "")
		migrateArea += "entity." + keyUpper + "{},"
	}

	repoArea := ""
	serviceArea := ""
	handlerArea := ""
	apiArea := ""
	for key := range objs {
		keys := strings.Split(key, "")
		keys[0] = strings.ToUpper(keys[0])
		keyUpper := strings.Join(keys, "")
		varRepo := key + "Repository"
		repoArea += varRepo + " := repository.New" + keyUpper + "Repository(db)\n"
		varService := key + "Service"
		serviceArea += varService + " := service.New" + keyUpper + "Service(" + varRepo + ")\n"
		handlerArea += key + "Handler := handler.New" + keyUpper + "Handler(" + varService + ")\n"
		apiArea += "api.POST(\"/create" + key + "\", " + key + "Handler.Create" + keyUpper + ")\n"
		apiArea += "api.POST(\"/edit" + key + "\", " + key + "Handler.Edit" + keyUpper + ")\n"
		apiArea += "api.POST(\"/getall" + key + "s\", " + key + "Handler.GetAll" + keyUpper + "s)\n"
		apiArea += "api.POST(\"/get" + key + "byid/:id\", " + key + "Handler.Get" + keyUpper + "ById)\n"
	}

	template := string(fileTemplate)

	template = strings.Replace(template, "[project]", project, -1)
	template = strings.Replace(template, "[migrateArea]", migrateArea, -1)
	template = strings.Replace(template, "[repoArea]", repoArea, -1)
	template = strings.Replace(template, "[serviceArea]", serviceArea, -1)
	template = strings.Replace(template, "[handlerArea]", handlerArea, -1)
	template = strings.Replace(template, "[apiArea]", apiArea, -1)

	_, err = fmt.Fprintln(file, template)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Main Created")
}

func createHelper(project string) {
	err := os.MkdirAll(project+"\\helper", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(project + "\\helper\\helper.go")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	copy, err := os.Open("template\\helper.go")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	_, err = io.Copy(file, copy)
	if err != nil {
		log.Fatal(err)
	}

	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}
}

func createHandler(items []string, name string, project string) {
	err := os.MkdirAll(project+"\\handler", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(project + "\\handler\\" + name + "Handler.go")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fileTemplate, err := os.ReadFile("template\\Handler.go")

	if err != nil {
		log.Fatal(err)
	}

	template := string(fileTemplate)

	names := strings.Split(name, "")
	names[0] = strings.ToUpper(names[0])
	nameUpper := strings.Join(names, "")

	template = strings.Replace(template, "[name]", name, -1)
	template = strings.Replace(template, "[nameUpper]", nameUpper, -1)
	template = strings.Replace(template, "[project]", project, -1)
	_, err = fmt.Fprintln(file, template)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(name + " Handler Created")
}

func createInput(items []string, name string, project string) {
	err := os.MkdirAll(project+"\\input", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(project + "\\input\\" + name + "Input.go")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	names := strings.Split(name, "")
	names[0] = strings.ToUpper(names[0])
	nameUpper := strings.Join(names, "")

	var codes = []string{
		"package input",
		"",
		"type " + nameUpper + "Input struct {",
	}
	for i := 1; i < len(items)-5; i++ {
		codes = append(codes, items[i]+" `json:\""+strings.ToLower(strings.Split(items[i], " ")[0])+"\" binding:\"required\"`")
	}
	codes = append(codes, []string{
		"}",
		"",
		"type " + nameUpper + "EditInput struct {",
	}...)
	for i := 0; i < len(items)-6; i++ {
		codes = append(codes, items[i]+" `json:\""+strings.ToLower(strings.Split(items[i], " ")[0])+"\" binding:\"required\"`")
	}
	codes = append(codes, items[len(items)-4]+" `json:\""+strings.ToLower(strings.Split(items[len(items)-4], " ")[0])+"\" binding:\"required\"`")
	codes = append(codes, "}")

	for _, code := range codes {
		_, err := fmt.Fprintln(file, code)
		if err != nil {
			log.Fatal(err)
		}

	}

	fmt.Println(name + " Input Created")
}

func createService(items []string, name string, project string) {
	err := os.MkdirAll(project+"\\service", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(project + "\\service\\" + name + "Service.go")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fileTemplate, err := os.ReadFile("template\\Service.go")

	if err != nil {
		log.Fatal(err)
	}

	template := string(fileTemplate)

	names := strings.Split(name, "")
	names[0] = strings.ToUpper(names[0])
	nameUpper := strings.Join(names, "")

	createItem := ""

	for i := 1; i < len(items)-6; i++ {
		createItem += strings.Split(items[i], " ")[0] + ": input." + strings.Split(items[i], " ")[0] + ",\n"
	}

	createItem += "\nCreatedDate: time.Now(),"

	editItem := ""
	for i := 0; i < len(items)-6; i++ {
		editItem += strings.Split(items[i], " ")[0] + ": input." + strings.Split(items[i], " ")[0] + ",\n"
	}
	editItem += "CreatedBy: old" + nameUpper + ".CreatedBy,\n"
	editItem += "CreatedDate: old" + nameUpper + ".CreatedDate,\n"
	editItem += "UpdatedBy: input.UpdatedBy,\n"

	template = strings.Replace(template, "[name]", name, -1)
	template = strings.Replace(template, "[nameUpper]", nameUpper, -1)
	template = strings.Replace(template, "[project]", project, -1)
	template = strings.Replace(template, "[createItem]", createItem, -1)
	template = strings.Replace(template, "[editItem]", editItem, -1)

	_, err = fmt.Fprintln(file, template)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(name + " Service Created")
}

func createRepository(items []string, name string, project string) {
	err := os.MkdirAll(project+"\\repository", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(project + "\\repository\\" + name + "Repository.go")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	fileTemplate, err := os.ReadFile("template\\Repository.go")

	if err != nil {
		log.Fatal(err)
	}

	template := string(fileTemplate)

	names := strings.Split(name, "")
	names[0] = strings.ToUpper(names[0])
	nameUpper := strings.Join(names, "")

	template = strings.Replace(template, "[name]", name, -1)
	template = strings.Replace(template, "[nameUpper]", nameUpper, -1)
	template = strings.Replace(template, "[project]", project, -1)
	_, err = fmt.Fprintln(file, template)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(name + " repository Created")
}

func createEntity(items []string, name string, project string) {
	err := os.MkdirAll(project+"\\entity", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(project + "\\entity\\" + name + "Entity.go")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	names := strings.Split(name, "")
	names[0] = strings.ToUpper(names[0])
	name = strings.Join(names, "")
	var codes = []string{
		"package entity",
		"",
		"type " + name + " struct{",
	}

	for _, value := range items {
		if value == "Id int" {
			codes = append(codes, value+"`gorm:\"primarykey;autoIncrement:true\"`")
		} else {
			codes = append(codes, value)
		}
	}

	codes = append(codes, "}")

	for _, code := range codes {
		_, err := fmt.Fprintln(file, code)
		if err != nil {
			log.Fatal(err)
		}

	}
	fmt.Println(name + " entity Created")
}
