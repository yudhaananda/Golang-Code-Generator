package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		html, err := os.ReadFile("index.html")
		if err != nil {
			ctx.Data(http.StatusBadGateway, "text/html; charset=utf-8", []byte(""))
		}
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", html)
	})
	router.POST("/generate", func(ctx *gin.Context) {

		content, err := ctx.FormFile("file")

		if err != nil {
			ctx.Data(http.StatusBadGateway, "text/html; charset=utf-8", []byte(""))
			return
		}

		file, err := content.Open()

		if err != nil {
			ctx.Data(http.StatusBadGateway, "text/html; charset=utf-8", []byte(""))
			return
		}
		defer file.Close()

		data, err := ioutil.ReadAll(file)

		if err != nil {
			ctx.Data(http.StatusBadGateway, "text/html; charset=utf-8", []byte(""))
			return
		}

		var jsonContent map[string]interface{}
		err = json.Unmarshal(data, &jsonContent)

		if err != nil {
			ctx.Data(http.StatusBadGateway, "text/html; charset=utf-8", []byte(""))
			return
		}

		objs := map[string][]string{}

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
		result, err := process(objs, project)
		if err != nil {
			ctx.Data(http.StatusBadGateway, "text/html; charset=utf-8", []byte(""))
			return
		}
		ctx.Data(http.StatusOK, "Application/zip", result)
	})
	router.Run()
}

func process(objs map[string][]string, project string) ([]byte, error) {

	for key, obj := range objs {
		if key == "user" {
			err := createAuthService(obj, project)
			if err != nil {
				return nil, err
			}
		}
		err := createEntity(obj, key, project)
		if err != nil {
			return nil, err
		}
		err = createRepository(obj, key, project)
		if err != nil {
			return nil, err
		}
		err = createService(obj, key, project)
		if err != nil {
			return nil, err
		}
		err = createInput(obj, key, project)
		if err != nil {
			return nil, err
		}
		err = createHandler(obj, key, project)
		if err != nil {
			return nil, err
		}
	}
	err := createHelper(project)
	if err != nil {
		return nil, err
	}
	err = createAuthHandler(project)
	if err != nil {
		return nil, err
	}
	err = createFormatter(project)
	if err != nil {
		return nil, err
	}
	err = createJwtService(project)
	if err != nil {
		return nil, err
	}
	err = createMain(objs, project)
	if err != nil {
		return nil, err
	}
	result, err := zipping(project)
	if err != nil {
		return nil, err
	}
	err = delete(project)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func delete(project string) error {
	err := os.RemoveAll(project)
	if err != nil {
		return err
	}
	err = os.Remove(project + ".zip")
	if err != nil {
		return err
	}
	return nil
}

func zipping(project string) ([]byte, error) {
	baseFolder := project + "/"

	zipName, err := os.Create(project + ".zip")

	if err != nil {
		return nil, err
	}

	defer zipName.Close()

	w := zip.NewWriter(zipName)

	addFiles(w, baseFolder)

	err = w.Close()
	if err != nil {
		return nil, err
	}
	zipFile, err := os.ReadFile(project + ".zip")

	if err != nil {
		return nil, err
	}

	return zipFile, nil
}

func addFiles(w *zip.Writer, basePath string) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		// fmt.Println(basePath + file.Name())
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(basePath + file.Name())
			if err != nil {
				fmt.Println(err)
			}

			// Add some files to the archive.
			var f io.Writer
			f, err = w.Create(basePath + file.Name())

			if err != nil {
				fmt.Println(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(basePath + file.Name())
		} else if file.IsDir() {

			newBase := basePath + file.Name() + "/"

			addFiles(w, newBase)
		}
	}
}

func createJwtService(project string) error {
	err := os.MkdirAll(project+"\\service\\", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "\\service\\jwtService.go")

	if err != nil {
		return err
	}
	defer file.Close()

	fileTemplate, err := os.ReadFile("template\\jwtService.txt")

	if err != nil {
		return err
	}

	template := string(fileTemplate)

	template = strings.Replace(template, "[project]", project, -1)

	_, err = fmt.Fprintln(file, template)
	if err != nil {
		return err
	}
	return nil
}

func createFormatter(project string) error {
	err := os.MkdirAll(project+"\\formatter\\", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "\\formatter\\authFormatter.go")

	if err != nil {
		return err
	}
	defer file.Close()

	fileTemplate, err := os.ReadFile("template\\authFormatter.txt")

	if err != nil {
		return err
	}

	template := string(fileTemplate)

	template = strings.Replace(template, "[project]", project, -1)

	_, err = fmt.Fprintln(file, template)
	if err != nil {
		return err
	}
	return nil
}

func createAuthHandler(project string) error {
	err := os.MkdirAll(project+"\\handler\\", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "\\handler\\authHandler.go")

	if err != nil {
		return err
	}
	defer file.Close()

	fileTemplate, err := os.ReadFile("template\\authHandler.txt")

	if err != nil {
		return err
	}

	template := string(fileTemplate)

	template = strings.Replace(template, "[project]", project, -1)

	_, err = fmt.Fprintln(file, template)
	if err != nil {
		return err
	}
	return nil
}

func createAuthService(items []string, project string) error {
	err := os.MkdirAll(project+"\\service\\", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "\\service\\authService.go")

	if err != nil {
		return err
	}
	defer file.Close()

	fileTemplate, err := os.ReadFile("template\\authService.txt")

	if err != nil {
		return err
	}

	template := string(fileTemplate)

	registerItem := ""

	for i := 1; i < len(items)-6; i++ {
		if strings.Split(items[i], " ")[0] != "Password" {
			registerItem += strings.Split(items[i], " ")[0] + ": input." + strings.Split(items[i], " ")[0] + ",\n"
		}
	}
	registerItem += "Password: string(password),\n"
	registerItem += "CreatedBy: input.UserName,\n"
	registerItem += "CreatedDate: time.Now(),"

	template = strings.Replace(template, "[project]", project, -1)
	template = strings.Replace(template, "[registerItem]", registerItem, -1)

	_, err = fmt.Fprintln(file, template)
	if err != nil {
		return err
	}
	return nil
}

func createMain(objs map[string][]string, project string) error {
	err := os.MkdirAll(project, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "\\main.go")

	if err != nil {
		return err
	}

	defer file.Close()

	fileTemplate, err := os.ReadFile("template\\main.txt")

	if err != nil {
		return err
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
	for key, items := range objs {
		keys := strings.Split(key, "")
		keys[0] = strings.ToUpper(keys[0])
		keyUpper := strings.Join(keys, "")
		varRepo := key + "Repository"
		repoArea += varRepo + " := repository.New" + keyUpper + "Repository(db)\n"
		varService := key + "Service"
		serviceArea += varService + " := service.New" + keyUpper + "Service(" + varRepo + ")\n"
		handlerArea += key + "Handler := handler.New" + keyUpper + "Handler(" + varService + ")\n"
		apiArea += "api.POST(\"/create" + key + "\", authMiddleware(jwtService, userService), " + key + "Handler.Create" + keyUpper + ")\n"
		apiArea += "api.POST(\"/edit" + key + "\", authMiddleware(jwtService, userService), " + key + "Handler.Edit" + keyUpper + ")\n"
		apiArea += "api.GET(\"/getall" + key + "s\", authMiddleware(jwtService, userService), " + key + "Handler.GetAll" + keyUpper + "s)\n"
		apiArea += "api.GET(\"/delete" + key + "/:id\", authMiddleware(jwtService, userService), " + key + "Handler.Delete" + keyUpper + ")\n"
		if key == "user" {
			serviceArea += "authService := service.NewAuthService(userRepository)\n"
			serviceArea += "jwtService := service.NewJwtService()\n"
			handlerArea += "authHandler := handler.NewAuthHandler(authService, jwtService)\n"
			apiArea += "api.POST(\"/register\", authHandler.RegisterUser)\n"
			apiArea += "api.POST(\"/login\", authHandler.Login)\n"
		}
		for i := 0; i < len(items)-6; i++ {
			itemSplit := strings.Split(strings.Split(items[i], " ")[0], "")
			if strings.Split(items[i], " ")[0] != "Password" {
				itemSplit[0] = strings.ToLower(itemSplit[0])
				itemLower := strings.Join(itemSplit, "")
				apiArea += "api.GET(\"/get" + key + "by" + strings.ToLower(itemLower) + "/:" + itemLower + "\", authMiddleware(jwtService, userService), " + key + "Handler.Get" + keyUpper + "By" + strings.Split(items[i], " ")[0] + ")\n"
			}
		}
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
		return err
	}

	fmt.Println("Main Created")
	return nil
}

func createHelper(project string) error {
	err := os.MkdirAll(project+"\\helper", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "\\helper\\helper.go")

	if err != nil {
		return err
	}

	defer file.Close()
	copy, err := os.Open("template\\helper.txt")
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, copy)
	if err != nil {
		return err
	}

	err = file.Sync()
	if err != nil {
		return err
	}
	return nil
}

func createHandler(items []string, name string, project string) error {
	err := os.MkdirAll(project+"\\handler", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "\\handler\\" + name + "Handler.go")

	if err != nil {
		return err
	}

	defer file.Close()

	fileTemplate, err := os.ReadFile("template\\Handler.txt")

	if err != nil {
		return err
	}

	fileTemplateGetByHandler, err := os.ReadFile("template\\GetByHandler.txt")

	if err != nil {
		return err
	}

	fileTemplateHandlerConvert, err := os.ReadFile("template\\HandlerConvert.txt")

	if err != nil {
		return err
	}

	template := string(fileTemplate)
	templateGetByHandler := string(fileTemplateGetByHandler)
	templateHandlerConvert := string(fileTemplateHandlerConvert)
	getByHandler := ""

	names := strings.Split(name, "")
	names[0] = strings.ToUpper(names[0])
	nameUpper := strings.Join(names, "")

	templateGetByHandler = strings.Replace(templateGetByHandler, "[name]", name, -1)
	templateGetByHandler = strings.Replace(templateGetByHandler, "[nameUpper]", nameUpper, -1)
	templateHandlerConvert = strings.Replace(templateHandlerConvert, "[nameUpper]", nameUpper, -1)

	for i := 0; i < len(items)-6; i++ {
		itemSplit := strings.Split(strings.Split(items[i], " ")[0], "")
		itemSplit[0] = strings.ToLower(itemSplit[0])
		itemLower := strings.Join(itemSplit, "")
		typeSplit := strings.Split(strings.Split(items[i], " ")[1], "")
		typeSplit[0] = strings.ToUpper(typeSplit[0])
		type_ := strings.Join(typeSplit, "")
		if strings.Contains(type_, "Float") {
			type_ = "Float"
		}

		tempGetByHandler := ""
		if itemLower != "password" {
			tempGetByHandler = strings.Replace(templateGetByHandler, "[item]", itemLower, -1)
			tempGetByHandler = strings.Replace(tempGetByHandler, "[itemUpper]", strings.Split(items[i], " ")[0], -1)
			tempGetByHandler = strings.Replace(tempGetByHandler, "[type]", type_, -1)
			if type_ == "String" {
				tempGetByHandler = strings.Replace(tempGetByHandler, "[itemParam]", itemLower, -1)
				tempGetByHandler = strings.Replace(tempGetByHandler, "[convert]", "", -1)
			} else {
				tempHandlerConvert := strings.Replace(templateHandlerConvert, "[itemParam]", itemLower+type_, -1)
				tempHandlerConvert = strings.Replace(tempHandlerConvert, "[item]", itemLower, -1)
				tempHandlerConvert = strings.Replace(tempHandlerConvert, "[type]", type_, -1)
				if type_ == "Int" {
					tempHandlerConvert = strings.Replace(tempHandlerConvert, "[parseType]", "Atoi", -1)
					tempHandlerConvert = strings.Replace(tempHandlerConvert, "[param]", "", -1)
				}
				if type_ == "Float64" {
					tempHandlerConvert = strings.Replace(tempHandlerConvert, "[parseType]", "ParseFloat", -1)
					tempHandlerConvert = strings.Replace(tempHandlerConvert, "[param]", ", 64", -1)
				}
				if type_ == "Float32" {
					tempHandlerConvert = strings.Replace(tempHandlerConvert, "[parseType]", "ParseFloat", -1)
					tempHandlerConvert = strings.Replace(tempHandlerConvert, "[param]", ", 32", -1)
				}
				if type_ == "Bool" {
					tempHandlerConvert = strings.Replace(tempHandlerConvert, "[parseType]", "ParseBool", -1)
					tempHandlerConvert = strings.Replace(tempHandlerConvert, "[param]", "", -1)
				}
				tempGetByHandler = strings.Replace(tempGetByHandler, "[itemParam]", itemLower+type_, -1)
				tempGetByHandler = strings.Replace(tempGetByHandler, "[convert]", tempHandlerConvert, -1)
			}
			getByHandler += tempGetByHandler + "\n"
		}
	}

	template = strings.Replace(template, "[name]", name, -1)
	template = strings.Replace(template, "[nameUpper]", nameUpper, -1)
	template = strings.Replace(template, "[project]", project, -1)
	template = strings.Replace(template, "[getByHandler]", getByHandler, -1)
	_, err = fmt.Fprintln(file, template)
	if err != nil {
		return err
	}

	fmt.Println(name + " Handler Created")
	return nil
}

func createInput(items []string, name string, project string) error {
	err := os.MkdirAll(project+"\\input", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "\\input\\" + name + "Input.go")

	if err != nil {
		return err
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
	for i := 0; i < len(items)-5; i++ {
		codes = append(codes, items[i]+" `json:\""+strings.ToLower(strings.Split(items[i], " ")[0])+"\" binding:\"required\"`")
	}
	codes = append(codes, items[len(items)-4]+" `json:\""+strings.ToLower(strings.Split(items[len(items)-4], " ")[0])+"\" binding:\"required\"`")
	codes = append(codes, "}")
	if name == "user" {
		codes = append(codes, []string{
			"",
			"type LoginInput struct {",
		}...)
		for i := 1; i < 3; i++ {
			codes = append(codes, items[i]+" `json:\""+strings.ToLower(strings.Split(items[i], " ")[0])+"\" binding:\"required\"`")
		}
		codes = append(codes, "}")
	}

	for _, code := range codes {
		_, err := fmt.Fprintln(file, code)
		if err != nil {
			return err
		}

	}

	fmt.Println(name + " Input Created")
	return nil
}

func createService(items []string, name string, project string) error {
	err := os.MkdirAll(project+"\\service", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "\\service\\" + name + "Service.go")

	if err != nil {
		return err
	}

	defer file.Close()

	fileTemplate, err := os.ReadFile("template\\Service.txt")

	if err != nil {
		return err
	}

	fileTemplateGetByServiceMethod, err := os.ReadFile("template\\GetByServiceMethod.txt")

	if err != nil {
		return err
	}

	fileTemplateGetByService, err := os.ReadFile("template\\GetByService.txt")

	if err != nil {
		return err
	}

	template := string(fileTemplate)
	templateGetByServiceMethod := string(fileTemplateGetByServiceMethod)
	templateGetByService := string(fileTemplateGetByService)
	getByServiceMethod := ""
	getByService := ""

	names := strings.Split(name, "")
	names[0] = strings.ToUpper(names[0])
	nameUpper := strings.Join(names, "")

	createItem := ""

	templateGetByServiceMethod = strings.Replace(templateGetByServiceMethod, "[name]", name, -1)
	templateGetByServiceMethod = strings.Replace(templateGetByServiceMethod, "[nameUpper]", nameUpper, -1)
	templateGetByService = strings.Replace(templateGetByService, "[nameUpper]", nameUpper, -1)

	for i := 1; i < len(items)-6; i++ {
		createItem += strings.Split(items[i], " ")[0] + ": input." + strings.Split(items[i], " ")[0] + ",\n"
	}

	createItem += "CreatedBy: userName,\n"
	createItem += "CreatedDate: time.Now(),"

	editItem := ""
	for i := 0; i < len(items)-6; i++ {
		editItem += strings.Split(items[i], " ")[0] + ": input." + strings.Split(items[i], " ")[0] + ",\n"
		itemSplit := strings.Split(strings.Split(items[i], " ")[0], "")
		itemSplit[0] = strings.ToLower(itemSplit[0])
		itemLower := strings.Join(itemSplit, "")
		tempGetByServiceMethod := ""
		tempGetByService := ""

		if itemLower != "password" {

			tempGetByServiceMethod = strings.Replace(templateGetByServiceMethod, "[itemUpper]", strings.Split(items[i], " ")[0], -1)
			tempGetByServiceMethod = strings.Replace(tempGetByServiceMethod, "[itemParam]", itemLower+" "+strings.Split(items[i], " ")[1], -1)
			tempGetByServiceMethod = strings.Replace(tempGetByServiceMethod, "[item]", itemLower, -1)
			getByServiceMethod += tempGetByServiceMethod + "\n"

			tempGetByService = strings.Replace(templateGetByService, "[itemUpper]", strings.Split(items[i], " ")[0], -1)
			tempGetByService = strings.Replace(tempGetByService, "[itemParam]", itemLower+" "+strings.Split(items[i], " ")[1], -1)
			getByService += tempGetByService + "\n"
		}
	}
	editItem += "CreatedBy: old" + nameUpper + ".CreatedBy,\n"
	editItem += "CreatedDate: old" + nameUpper + ".CreatedDate,\n"
	editItem += "UpdatedBy: userName,\n"

	template = strings.Replace(template, "[name]", name, -1)
	template = strings.Replace(template, "[nameUpper]", nameUpper, -1)
	template = strings.Replace(template, "[project]", project, -1)
	template = strings.Replace(template, "[createItem]", createItem, -1)
	template = strings.Replace(template, "[editItem]", editItem, -1)
	template = strings.Replace(template, "[getBy]", getByService, -1)
	template = strings.Replace(template, "[getByMethod]", getByServiceMethod, -1)

	_, err = fmt.Fprintln(file, template)
	if err != nil {
		return err
	}

	fmt.Println(name + " Service Created")
	return nil
}

func createRepository(items []string, name string, project string) error {
	err := os.MkdirAll(project+"\\repository", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "\\repository\\" + name + "Repository.go")

	if err != nil {
		return err
	}

	defer file.Close()

	fileTemplate, err := os.ReadFile("template\\Repository.txt")

	if err != nil {
		return err
	}

	fileTemplateFindByRepoMethod, err := os.ReadFile("template\\FindByRepoMethod.txt")

	if err != nil {
		return err
	}

	fileTemplateFindByRepo, err := os.ReadFile("template\\FindByRepo.txt")

	if err != nil {
		return err
	}

	template := string(fileTemplate)
	templateFindByMethod := string(fileTemplateFindByRepoMethod)
	templateFindBy := string(fileTemplateFindByRepo)
	findByMethod := ""
	findBy := ""

	names := strings.Split(name, "")
	names[0] = strings.ToUpper(names[0])
	nameUpper := strings.Join(names, "")

	templateFindByMethod = strings.Replace(templateFindByMethod, "[name]", name, -1)
	templateFindByMethod = strings.Replace(templateFindByMethod, "[nameUpper]", nameUpper, -1)
	templateFindBy = strings.Replace(templateFindBy, "[nameUpper]", nameUpper, -1)

	for i := 0; i < len(items)-6; i++ {
		itemSplit := strings.Split(strings.Split(items[i], " ")[0], "")
		itemSplit[0] = strings.ToLower(itemSplit[0])
		itemLower := strings.Join(itemSplit, "")
		for i := 0; i < len(itemSplit); i++ {
			if itemSplit[i] == strings.ToUpper(itemSplit[i]) {
				itemSplit[i] = "_" + strings.ToLower(itemSplit[i])
			}
		}
		item_ := strings.Join(itemSplit, "")
		tempFindByMethod := ""
		tempFindBy := ""

		if itemLower != "password" {
			tempFindByMethod = strings.Replace(templateFindByMethod, "[item]", itemLower, -1)
			tempFindByMethod = strings.Replace(tempFindByMethod, "[item_]", item_, -1)
			tempFindByMethod = strings.Replace(tempFindByMethod, "[itemUpper]", strings.Split(items[i], " ")[0], -1)
			tempFindByMethod = strings.Replace(tempFindByMethod, "[itemParam]", itemLower+" "+strings.Split(items[i], " ")[1], -1)
			findByMethod += tempFindByMethod

			tempFindBy = strings.Replace(templateFindBy, "[itemUpper]", strings.Split(items[i], " ")[0], -1)
			tempFindBy = strings.Replace(tempFindBy, "[item]", itemLower, -1)
			tempFindBy = strings.Replace(tempFindBy, "[itemParam]", itemLower+" "+strings.Split(items[i], " ")[1], -1)
			fmt.Println(itemLower + " " + strings.Split(items[i], " ")[1])
			findBy += tempFindBy + "\n"
		}
	}

	template = strings.Replace(template, "[name]", name, -1)
	template = strings.Replace(template, "[nameUpper]", nameUpper, -1)
	template = strings.Replace(template, "[project]", project, -1)
	template = strings.Replace(template, "[findByMethod]", findByMethod, -1)
	template = strings.Replace(template, "[findBy]", findBy, -1)

	_, err = fmt.Fprintln(file, template)
	if err != nil {
		return err
	}

	fmt.Println(name + " repository Created")
	return nil
}

func createEntity(items []string, name string, project string) error {
	err := os.MkdirAll(project+"\\entity", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "\\entity\\" + name + "Entity.go")

	if err != nil {
		return err
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
			return err
		}

	}
	fmt.Println(name + " entity Created")
	return nil
}
