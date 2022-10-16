package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		html, err := os.ReadFile("index.html")
		if err != nil {
			ctx.Data(http.StatusBadGateway, "text/html; charset=utf-8", []byte(err.Error()))
		}
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", html)
	})
	router.GET("/favicon.ico", func(ctx *gin.Context) {
		png, err := os.ReadFile("favicon.ico")
		if err != nil {
			ctx.Data(http.StatusBadGateway, "text/html; charset=utf-8", []byte(err.Error()))
		}
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", png)
	})
	router.POST("/generate", func(ctx *gin.Context) {

		content, err := ctx.FormFile("file")

		if err != nil {
			ctx.Data(http.StatusBadGateway, "text/html; charset=utf-8", []byte(err.Error()))
			return
		}

		split := strings.Split(content.Filename, ".")
		extention := split[len(split)-1]

		if extention != "json" {
			ctx.Data(http.StatusBadGateway, "text/html; charset=utf-8", []byte("Only .json file"))
			return
		}

		file, err := content.Open()

		if err != nil {
			ctx.Data(http.StatusBadGateway, "text/html; charset=utf-8", []byte(err.Error()))
			return
		}
		defer file.Close()

		data, err := ioutil.ReadAll(file)

		if err != nil {
			ctx.Data(http.StatusBadGateway, "text/html; charset=utf-8", []byte(err.Error()))
			return
		}

		var jsonContent map[string]interface{}
		err = json.Unmarshal(data, &jsonContent)

		if err != nil {
			ctx.Data(http.StatusBadGateway, "text/html; charset=utf-8", []byte(err.Error()))
			return
		}

		objs := map[string][]string{}

		relation := []map[string]string{}

		database := map[string]string{}

		projectObject := jsonContent["projectName"]
		entity := jsonContent["entity"]
		relationObject := jsonContent["relation"]
		databaseObject := jsonContent["database"]
		isUsingWebsocketObject := jsonContent["isUsingWebsocket"]

		project := ""
		if projectObject != nil {
			project = strings.ToLower(fmt.Sprintf("%v", projectObject))
		}

		isUsingWebsocket := false
		if isUsingWebsocketObject != nil {
			isUsingWebsocket = isUsingWebsocketObject.(bool)
		}

		for key, val := range databaseObject.(map[string]interface{}) {
			database[key] = val.(string)
		}

		if relationObject != nil {
			for _, obj := range relationObject.([]interface{}) {
				temp := map[string]string{}
				for key, val := range obj.(map[string]interface{}) {
					temp[key] = val.(string)
				}
				relation = append(relation, temp)
			}
		}

		if entity != nil {
			for key, obj := range entity.(map[string]interface{}) {
				temp := []string{}
				for _, val := range obj.([]interface{}) {
					temp = append(temp, val.(string))
				}
				objs[key] = temp
			}
		}
		result, err := process(objs, project, relation, database, isUsingWebsocket)
		if err != nil {
			ctx.Data(http.StatusBadGateway, "text/html; charset=utf-8", []byte(err.Error()))
			return
		}
		ctx.Data(http.StatusOK, "Application/zip", result)
	})
	router.POST("/template.json", func(ctx *gin.Context) {
		template, err := os.ReadFile("project.json")
		if err != nil {
			ctx.Data(http.StatusBadGateway, "text/html; charset=utf-8", []byte(err.Error()))
			return
		}
		ctx.Data(http.StatusOK, "Application/file", template)

	})
	router.Run()
}

func process(objs map[string][]string, project string, relation []map[string]string, database map[string]string, isUsingWebsocket bool) ([]byte, error) {

	for key, obj := range objs {

		objNew, err := createEntity(obj, key, project, relation)
		if err != nil {
			return nil, err
		}
		if key == "user" {
			err := createAuthService(objNew, project)
			if err != nil {
				return nil, err
			}
		}
		err = createRepository(objNew, key, project, obj)
		if err != nil {
			return nil, err
		}
		err = createService(objNew, key, project)
		if err != nil {
			return nil, err
		}
		err = createInput(objNew, key, project)
		if err != nil {
			return nil, err
		}
		err = createHandler(objNew, key, project, isUsingWebsocket)
		if err != nil {
			return nil, err
		}
		objs[key] = objNew
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
	err = createMain(objs, project, database, isUsingWebsocket)
	if err != nil {
		return nil, err
	}
	err = createEnvFile(project, database)
	if err != nil {
		return nil, err
	}
	err = createMakeFile(project, database, isUsingWebsocket)
	if err != nil {
		return nil, err
	}
	err = createAuthMiddleware(project)
	if err != nil {
		return nil, err
	}
	err = createEnvEntity(project)
	if err != nil {
		return nil, err
	}
	err = createApiListHtml(objs, project)
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

func createApiListHtml(objs map[string][]string, project string) error {

	err := os.MkdirAll(project, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "/index.html")

	if err != nil {
		return err
	}

	defer file.Close()

	fileTemplate, err := os.ReadFile("template/index.txt")

	if err != nil {
		return err
	}

	templateLoop, err := os.ReadFile("template/indexLoop.txt")

	if err != nil {
		return err
	}

	template := string(fileTemplate)
	loop := string(templateLoop)

	apiArea := ""
	count := 0

	for key, items := range objs {

		apiArea += strings.Replace(loop, "[item]", "api.POST(\"/create"+key+"\")\n", -1)
		count += 1
		apiArea = strings.Replace(apiArea, "[id]", strconv.FormatInt(int64(count), 10), -1)
		apiArea = strings.Replace(apiArea, "[jsonArea]", "[jsonAreaCreate]", -1)
		apiArea += strings.Replace(loop, "[item]", "api.PUT(\"/edit"+key+"\")\n", -1)
		count += 1
		apiArea = strings.Replace(apiArea, "[id]", strconv.FormatInt(int64(count), 10), -1)
		apiArea = strings.Replace(apiArea, "[jsonArea]", "[jsonAreaEdit]", -1)
		apiArea += strings.Replace(loop, "[item]", "api.GET(\"/getall"+key+"s\")\n", -1)
		count += 1
		apiArea = strings.Replace(apiArea, "[id]", strconv.FormatInt(int64(count), 10), -1)
		apiArea = strings.Replace(apiArea, "[jsonArea]", "[jsonAreaGet]", -1)
		apiArea += strings.Replace(loop, "[item]", "api.DELETE(\"/delete"+key+"/:id\")\n", -1)
		count += 1
		apiArea = strings.Replace(apiArea, "[id]", strconv.FormatInt(int64(count), 10), -1)
		apiArea = strings.Replace(apiArea, "[jsonArea]", "", -1)
		if key == "user" {
			apiArea += strings.Replace(loop, "[item]", "api.POST(\"/register\")\n", -1)
			count += 1
			apiArea = strings.Replace(apiArea, "[id]", strconv.FormatInt(int64(count), 10), -1)
			apiArea = strings.Replace(apiArea, "[jsonArea]", "[jsonAreaRegister]", -1)
			apiArea += strings.Replace(loop, "[item]", "api.POST(\"/login\")\n", -1)
			count += 1
			apiArea = strings.Replace(apiArea, "[id]", strconv.FormatInt(int64(count), 10), -1)
			apiArea = strings.Replace(apiArea, "[jsonArea]", "<p><b>Json Request</b></p>\n{\"username\" : string\n <br> \"password\" : string}", -1)
		}
		jsonAreaCreate := "<p><b>Json Request</b></p>\n{"
		jsonAreaEdit := "<p><b>Json Request</b></p>\n{"
		jsonAreaGet := "<p><b>Json Request</b></p>\n{\"Page\": int <br>\n \"Take\": int <br>\n \"OrderBy\": string <br>\n}"
		for i := 0; i < len(items)-6; i++ {
			if len(strings.Split(items[i], " ")) > 1 {
				itemSplit := strings.Split(strings.Split(items[i], " ")[0], "")
				if strings.Split(items[i], " ")[0] != "Password" && !strings.Contains(strings.Split(items[i], " ")[1], "time.Time") {
					itemSplit[0] = strings.ToLower(itemSplit[0])
					itemLower := strings.Join(itemSplit, "")
					apiArea += strings.Replace(loop, "[item]", "api.GET(\"/get"+key+"by"+strings.ToLower(itemLower)+"/:"+itemLower+"\")\n", -1)
					count += 1
					apiArea = strings.Replace(apiArea, "[id]", strconv.FormatInt(int64(count), 10), -1)
					apiArea = strings.Replace(apiArea, "[jsonArea]", "[jsonAreaGet]", -1)
				}
				if i != 0 {
					jsonAreaCreate += "\"" + strings.ToLower(strings.Split(items[i], " ")[0]) + "\" : " + strings.Split(items[i], " ")[1] + "<br>\n"
				}
				jsonAreaEdit += "\"" + strings.ToLower(strings.Split(items[i], " ")[0]) + "\" : " + strings.Split(items[i], " ")[1] + "<br>\n"
			}
		}
		jsonAreaCreate += "}"
		jsonAreaEdit += "}"
		if key == "user" {
			apiArea = strings.Replace(apiArea, "[jsonAreaRegister]", jsonAreaCreate, -1)
		}
		apiArea = strings.Replace(apiArea, "[jsonAreaEdit]", jsonAreaEdit, -1)
		apiArea = strings.Replace(apiArea, "[jsonAreaCreate]", jsonAreaCreate, -1)
		apiArea = strings.Replace(apiArea, "[jsonAreaGet]", jsonAreaGet, -1)
	}
	template = strings.Replace(template, "[itemLoop]", apiArea, -1)

	_, err = fmt.Fprintln(file, template)
	if err != nil {
		return err
	}
	return nil
}

func createAuthMiddleware(project string) error {
	err := os.MkdirAll(project+"/middleware/", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "/middleware/authMiddleware.go")

	if err != nil {
		return err
	}
	defer file.Close()

	fileTemplate, err := os.ReadFile("template/authMiddleware.txt")

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

func createEnvFile(project string, database map[string]string) error {
	err := os.MkdirAll(project, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "/.env")

	if err != nil {
		return err
	}

	defer file.Close()

	code := "DB_USER = \"" + database["user"] + "\"\nDB_PASS = \"" + database["pass"] + "\"\nDB_PORT = \"" + database["port"] + "\"\nDB_HOST = \"" + database["host"] + "\"\nDB_NAME = \"" + project + "\"\nJWT_SECRET_TOKEN=\"" + database["jwttoken"] + "\""

	_, err = fmt.Fprintln(file, code)
	if err != nil {
		return err
	}

	return nil
}

func createMakeFile(project string, database map[string]string, isUsingWebsocket bool) error {
	err := os.MkdirAll(project, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "/Makefile")

	if err != nil {
		return err
	}

	defer file.Close()

	code := "initialize:\n	go mod init " + project + "\n	go get github.com/gin-gonic/gin\n	go get github.com/go-playground/validator/v10\n	go get gorm.io/gorm\n	go get golang.org/x/crypto/bcrypt\n	go get gorm.io/driver/" + database["type"] + "\n	go get github.com/joho/godotenv\n	go get github.com/gin-contrib/cors\n	go get github.com/dgrijalva/jwt-go"

	if isUsingWebsocket {
		code += "\n	go get github.com/gorilla/websocket"
	}

	_, err = fmt.Fprintln(file, code)
	if err != nil {
		return err
	}

	return nil
}

func createEnvEntity(project string) error {
	err := os.MkdirAll(project+"/entity", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "/entity/envEntity.go")

	if err != nil {
		return err
	}

	defer file.Close()
	copy, err := os.Open("template/envEntity.txt")
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

func createJwtService(project string) error {
	err := os.MkdirAll(project+"/service/", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "/service/jwtService.go")

	if err != nil {
		return err
	}
	defer file.Close()

	fileTemplate, err := os.ReadFile("template/jwtService.txt")

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
	err := os.MkdirAll(project+"/formatter/", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "/formatter/authFormatter.go")

	if err != nil {
		return err
	}
	defer file.Close()

	fileTemplate, err := os.ReadFile("template/authFormatter.txt")

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
	err := os.MkdirAll(project+"/handler/", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "/handler/authHandler.go")

	if err != nil {
		return err
	}
	defer file.Close()

	fileTemplate, err := os.ReadFile("template/authHandler.txt")

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
	err := os.MkdirAll(project+"/service/", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "/service/authService.go")

	if err != nil {
		return err
	}
	defer file.Close()

	fileTemplate, err := os.ReadFile("template/authService.txt")

	if err != nil {
		return err
	}

	template := string(fileTemplate)

	registerItem := ""

	for i := 1; i < len(items)-6; i++ {
		if len(strings.Split(items[i], " ")) > 1 {
			if strings.Split(items[i], " ")[0] != "Password" {
				registerItem += "		" + strings.Split(items[i], " ")[0] + ": input." + strings.Split(items[i], " ")[0] + ",\n"
			}
		}
	}
	registerItem += "		Password: string(password),\n"
	registerItem += "		CreatedBy: input.UserName,\n"
	registerItem += "		CreatedDate: time.Now(),"

	template = strings.Replace(template, "[project]", project, -1)
	template = strings.Replace(template, "[registerItem]", registerItem, -1)

	_, err = fmt.Fprintln(file, template)
	if err != nil {
		return err
	}
	return nil
}

func createMain(objs map[string][]string, project string, database map[string]string, isUsingWebsocket bool) error {
	err := os.MkdirAll(project, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "/main.go")

	if err != nil {
		return err
	}

	defer file.Close()

	fileTemplate, err := os.ReadFile("template/main.txt")

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
		if isUsingWebsocket {
			handlerArea += key + "Handler := handler.New" + keyUpper + "Handler(" + varService + ", upgrader)\n"
		} else {
			handlerArea += key + "Handler := handler.New" + keyUpper + "Handler(" + varService + ")\n"
		}
		apiArea += "api.POST(\"/create" + strings.ToLower(key) + "\", authMiddleware.AuthMiddleware, " + key + "Handler.Create" + keyUpper + ")\n"
		apiArea += "api.PUT(\"/edit" + strings.ToLower(key) + "\", authMiddleware.AuthMiddleware, " + key + "Handler.Edit" + keyUpper + ")\n"
		apiArea += "api.GET(\"/getall" + strings.ToLower(key) + "s\", authMiddleware.AuthMiddleware, " + key + "Handler.GetAll" + keyUpper + "s)\n"
		apiArea += "api.DELETE(\"/delete" + strings.ToLower(key) + "/:id\", authMiddleware.AuthMiddleware, " + key + "Handler.Delete" + keyUpper + ")\n"
		if key == "user" {
			serviceArea += "authService := service.NewAuthService(userRepository)\n"
			serviceArea += "jwtService := service.NewJwtService()\n"
			handlerArea += "authHandler := handler.NewAuthHandler(authService, jwtService)\n"
			apiArea += "api.POST(\"/register\", authHandler.RegisterUser)\n"
			apiArea += "api.POST(\"/login\", authHandler.Login)\n"
		}
		for i := 0; i < len(items)-6; i++ {
			if len(strings.Split(items[i], " ")) > 1 {
				itemSplit := strings.Split(strings.Split(items[i], " ")[0], "")
				if strings.Split(items[i], " ")[0] != "Password" && !strings.Contains(strings.Split(items[i], " ")[1], "time.Time") {
					itemSplit[0] = strings.ToLower(itemSplit[0])
					itemLower := strings.Join(itemSplit, "")
					apiArea += "api.GET(\"/get" + key + "by" + strings.ToLower(itemLower) + "/:" + itemLower + "\", authMiddleware.AuthMiddleware, " + key + "Handler.Get" + keyUpper + "By" + strings.Split(items[i], " ")[0] + ")\n"
				}
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
	template = strings.Replace(template, "[dbType]", database["type"], -1)

	if isUsingWebsocket {
		template = strings.Replace(template, "[isWebsocket]", "\"github.com/gorilla/websocket\"", -1)
		template = strings.Replace(template, "[upgrader]", "var upgrader = websocket.Upgrader{}", -1)
	}

	template = strings.Replace(template, "[isWebsocket]", "", -1)
	template = strings.Replace(template, "[upgrader]", "", -1)

	if database["type"] == "mysql" {
		template = strings.Replace(template, "[dsnString]", "\"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local\", env.DB_USER, env.DB_PASS, env.DB_HOST, env.DB_PORT, env.DB_NAME", -1)
	} else if database["type"] == "postgres" {
		template = strings.Replace(template, "[dsnString]", "\"host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Shanghai\", env.DB_HOST, env.DB_USER, env.DB_PASS, env.DB_NAME, env.DB_PORT", -1)
	} else if database["type"] == "sqlite" {
		template = strings.Replace(template, "[dsnString]", "gorm.db", -1)
	} else if database["type"] == "sqlserver" {
		template = strings.Replace(template, "[dsnString]", "\"sqlserver://%s:%s@%s:%s?database=%s\", env.DB_USER, DB_PASS, env.DB_HOST, env.DB_PORT, env.DB_NAME", -1)
	}

	_, err = fmt.Fprintln(file, template)
	if err != nil {
		return err
	}

	fmt.Println("Main Created")
	return nil
}

func createHelper(project string) error {
	err := os.MkdirAll(project+"/helper", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "/helper/helper.go")

	if err != nil {
		return err
	}

	defer file.Close()
	copy, err := os.Open("template/helper.txt")
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

func createHandler(items []string, name string, project string, isUsingWebsocket bool) error {
	err := os.MkdirAll(project+"/handler", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "/handler/" + name + "Handler.go")

	if err != nil {
		return err
	}

	defer file.Close()

	fileTemplate, err := os.ReadFile("template/Handler.txt")

	if err != nil {
		return err
	}

	fileTemplateGetByHandler, err := os.ReadFile("template/GetByHandler.txt")

	if err != nil {
		return err
	}

	fileTemplateHandlerConvert, err := os.ReadFile("template/HandlerConvert.txt")

	if err != nil {
		return err
	}

	fileTemplateHandlerPaging, err := os.ReadFile("template/GetByPagingHandler.txt")
	if err != nil {
		return err
	}

	fileTemplateGetAllHandler, err := os.ReadFile("template/GetAllHandler.txt")
	if err != nil {
		return err
	}

	if isUsingWebsocket {
		fileTemplateGetAllHandler, err = os.ReadFile("template/GetAllHandlerWS.txt")
		if err != nil {
			return err
		}

		fileTemplateGetByHandler, err = os.ReadFile("template/GetByHandlerWS.txt")

		if err != nil {
			return err
		}

		fileTemplateHandlerPaging, err = os.ReadFile("template/GetByPagingHandlerWS.txt")
		if err != nil {
			return err
		}

	}

	template := string(fileTemplate)
	templateGetByHandler := string(fileTemplateGetByHandler)
	templateHandlerConvert := string(fileTemplateHandlerConvert)
	templatePaging := string(fileTemplateHandlerPaging)
	templateGetallhandler := string(fileTemplateGetAllHandler)
	getByHandler := ""

	names := strings.Split(name, "")
	names[0] = strings.ToUpper(names[0])
	nameUpper := strings.Join(names, "")

	templateGetByHandler = strings.Replace(templateGetByHandler, "[name]", name, -1)
	templateGetByHandler = strings.Replace(templateGetByHandler, "[nameUpper]", nameUpper, -1)
	templateHandlerConvert = strings.Replace(templateHandlerConvert, "[nameUpper]", nameUpper, -1)
	templatePaging = strings.Replace(templatePaging, "[nameUpper]", nameUpper, -1)
	templateGetallhandler = strings.Replace(templateGetallhandler, "[nameUpper]", nameUpper, -1)
	templateGetallhandler = strings.Replace(templateGetallhandler, "[name]", name, -1)

	for i := 0; i < len(items)-6; i++ {
		if len(strings.Split(items[i], " ")) > 1 {

			itemSplit := strings.Split(strings.Split(items[i], " ")[0], "")
			itemSplit[0] = strings.ToLower(itemSplit[0])
			itemLower := strings.Join(itemSplit, "")
			typeSplit := strings.Split(strings.Split(items[i], " ")[1], "")
			typeSplit[0] = strings.ToUpper(typeSplit[0])
			type_ := strings.Join(typeSplit, "")
			if strings.Contains(type_, "Float") {
				type_ = "Float"
			}
			paging := ""
			pagingItem := ""
			tempGetByHandler := ""
			if itemLower != "password" && !strings.Contains(strings.Split(items[i], " ")[1], "time.Time") {
				if itemLower != "id" {
					pagingItem = ", helper.SetPagingDefault(paging)"
					paging = templatePaging

				}
				tempGetByHandler = strings.Replace(templateGetByHandler, "[item]", itemLower, -1)
				tempGetByHandler = strings.Replace(tempGetByHandler, "[itemUpper]", strings.Split(items[i], " ")[0], -1)
				tempGetByHandler = strings.Replace(tempGetByHandler, "[type]", type_, -1)
				tempGetByHandler = strings.Replace(tempGetByHandler, "[paging]", paging, -1)
				if type_ == "String" {
					tempGetByHandler = strings.Replace(tempGetByHandler, "[itemParam]", itemLower+pagingItem, -1)
					tempGetByHandler = strings.Replace(tempGetByHandler, "[convert]", "", -1)
				} else {
					tempHandlerConvert := strings.Replace(templateHandlerConvert, "[itemParam]", itemLower+type_, -1)
					tempHandlerConvert = strings.Replace(tempHandlerConvert, "[item]", itemLower, -1)
					tempHandlerConvert = strings.Replace(tempHandlerConvert, "[type]", type_, -1)
					if strings.Contains(strings.Split(items[i], " ")[1], "int") {
						tempHandlerConvert = strings.Replace(tempHandlerConvert, "[parseType]", "Atoi", -1)
						tempHandlerConvert = strings.Replace(tempHandlerConvert, "[param]", "", -1)
					}
					if strings.Contains(strings.Split(items[i], " ")[1], "float64") {
						tempHandlerConvert = strings.Replace(tempHandlerConvert, "[parseType]", "ParseFloat", -1)
						tempHandlerConvert = strings.Replace(tempHandlerConvert, "[param]", ", 64", -1)
					}
					if strings.Contains(strings.Split(items[i], " ")[1], "float32") {
						tempHandlerConvert = strings.Replace(tempHandlerConvert, "[parseType]", "ParseFloat", -1)
						tempHandlerConvert = strings.Replace(tempHandlerConvert, "[param]", ", 32", -1)
					}
					if strings.Contains(strings.Split(items[i], " ")[1], "bool") {
						tempHandlerConvert = strings.Replace(tempHandlerConvert, "[parseType]", "ParseBool", -1)
						tempHandlerConvert = strings.Replace(tempHandlerConvert, "[param]", "", -1)
					}
					tempGetByHandler = strings.Replace(tempGetByHandler, "[itemParam]", itemLower+type_+pagingItem, -1)
					tempGetByHandler = strings.Replace(tempGetByHandler, "[convert]", tempHandlerConvert, -1)
				}
				getByHandler += tempGetByHandler + "\n"
			}
		}
	}

	template = strings.Replace(template, "[name]", name, -1)
	template = strings.Replace(template, "[nameUpper]", nameUpper, -1)
	template = strings.Replace(template, "[project]", project, -1)
	template = strings.Replace(template, "[getByHandler]", getByHandler, -1)
	template = strings.Replace(template, "[getAllHandler]", templateGetallhandler, -1)
	if isUsingWebsocket {
		template = strings.Replace(template, "[isWebsocket]", "\"github.com/gorilla/websocket\"", -1)
		template = strings.Replace(template, "[webSocketItem1]", "upgrader websocket.Upgrader", -1)
		template = strings.Replace(template, "[webSocketItem2]", ", upgrader websocket.Upgrader", -1)
		template = strings.Replace(template, "[webSocketItem3]", ", upgrader", -1)
	}
	template = strings.Replace(template, "[isWebsocket]", "", -1)
	template = strings.Replace(template, "[webSocketItem1]", "", -1)
	template = strings.Replace(template, "[webSocketItem2]", "", -1)
	template = strings.Replace(template, "[webSocketItem3]", "", -1)

	_, err = fmt.Fprintln(file, template)
	if err != nil {
		return err
	}

	fmt.Println(name + " Handler Created")
	return nil
}

func createInput(items []string, name string, project string) error {
	err := os.MkdirAll(project+"/input", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "/input/" + name + "Input.go")

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
	for i := 1; i < len(items)-6; i++ {
		if len(strings.Split(items[i], " ")) > 1 {
			temp := strings.Split(items[i], " ")
			codes = append(codes, "	"+temp[0]+" "+temp[1]+" `json:\""+strings.ToLower(strings.Split(items[i], " ")[0])+"\" binding:\"required\"`")
		}
	}
	codes = append(codes, []string{
		"}",
		"",
		"type " + nameUpper + "EditInput struct {",
	}...)
	for i := 0; i < len(items)-6; i++ {
		if len(strings.Split(items[i], " ")) > 1 {
			temp := strings.Split(items[i], " ")
			codes = append(codes, "	"+temp[0]+" "+temp[1]+" `json:\""+strings.ToLower(strings.Split(items[i], " ")[0])+"\" binding:\"required\"`")
		}
	}
	codes = append(codes, "}")
	if name == "user" {
		codes = append(codes, []string{
			"",
			"type LoginInput struct {",
		}...)
		for i := 1; i < 3; i++ {
			codes = append(codes, "	"+items[i]+" `json:\""+strings.ToLower(strings.Split(items[i], " ")[0])+"\" binding:\"required\"`")
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
	err := os.MkdirAll(project+"/service", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "/service/" + name + "Service.go")

	if err != nil {
		return err
	}

	defer file.Close()

	fileTemplate, err := os.ReadFile("template/Service.txt")

	if err != nil {
		return err
	}

	fileTemplateGetByServiceMethod, err := os.ReadFile("template/GetByServiceMethod.txt")

	if err != nil {
		return err
	}

	fileTemplateGetByService, err := os.ReadFile("template/GetByService.txt")

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
		if len(strings.Split(items[i], " ")) > 1 {
			createItem += "		" + strings.Split(items[i], " ")[0] + ": input." + strings.Split(items[i], " ")[0] + ",\n"
		}
	}

	createItem += "		CreatedBy: userLogin.UserName,\n"
	createItem += "		CreatedDate: time.Now(),"

	editItem := ""
	for i := 0; i < len(items)-6; i++ {
		if len(strings.Split(items[i], " ")) > 1 {

			editItem += "		" + strings.Split(items[i], " ")[0] + ": input." + strings.Split(items[i], " ")[0] + ",\n"
			itemSplit := strings.Split(strings.Split(items[i], " ")[0], "")
			itemSplit[0] = strings.ToLower(itemSplit[0])
			itemLower := strings.Join(itemSplit, "")
			tempGetByServiceMethod := ""
			tempGetByService := ""
			pagingParam := ""
			pagingItem := ""

			if itemLower != "password" {

				if itemLower != "id" {
					pagingParam = ", paging helper.Paging"
					pagingItem = ", paging"
				}
				tempGetByServiceMethod = strings.Replace(templateGetByServiceMethod, "[itemUpper]", strings.Split(items[i], " ")[0], -1)
				tempGetByServiceMethod = strings.Replace(tempGetByServiceMethod, "[itemParam]", itemLower+" "+strings.Split(items[i], " ")[1]+pagingParam, -1)
				tempGetByServiceMethod = strings.Replace(tempGetByServiceMethod, "[item]", itemLower+pagingItem, -1)
				getByServiceMethod += tempGetByServiceMethod + "\n"

				tempGetByService = strings.Replace(templateGetByService, "[itemUpper]", strings.Split(items[i], " ")[0], -1)
				tempGetByService = strings.Replace(tempGetByService, "[itemParam]", itemLower+" "+strings.Split(items[i], " ")[1]+pagingParam, -1)
				getByService += tempGetByService + "\n"
			}
		}
	}
	editItem += "		CreatedBy: old" + nameUpper + ".CreatedBy,\n"
	editItem += "		CreatedDate: old" + nameUpper + ".CreatedDate,\n"
	editItem += "		UpdatedBy: userLogin.UserName,\n"
	editItem += "		UpdatedDate: time.Now(),"

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

func createRepository(items []string, name string, project string, itemCompares []string) error {
	err := os.MkdirAll(project+"/repository", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(project + "/repository/" + name + "Repository.go")

	if err != nil {
		return err
	}

	defer file.Close()

	fileTemplate, err := os.ReadFile("template/Repository.txt")

	if err != nil {
		return err
	}

	fileTemplateFindByRepoMethod, err := os.ReadFile("template/FindByRepoMethod.txt")

	if err != nil {
		return err
	}

	fileTemplateFindByRepo, err := os.ReadFile("template/FindByRepo.txt")

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
	diference := len(items) - len(itemCompares)
	preload := ""

	for i := len(items) - 6 - diference; i < len(items)-6; i++ {
		if len(strings.Split(items[i], " ")) == 1 {
			preload += ".Preload(\"" + items[i] + "\")"
		}
	}

	templateFindByMethod = strings.Replace(templateFindByMethod, "[name]", name, -1)
	templateFindByMethod = strings.Replace(templateFindByMethod, "[nameUpper]", nameUpper, -1)
	templateFindByMethod = strings.Replace(templateFindByMethod, "[preload]", preload, -1)
	templateFindBy = strings.Replace(templateFindBy, "[nameUpper]", nameUpper, -1)

	for i := 0; i < len(items)-6; i++ {
		if len(strings.Split(items[i], " ")) > 1 {

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
			pagingParam := ""
			paging := ""

			if itemLower != "password" {
				if itemLower != "id" {
					pagingParam = ", paging helper.Paging"
					paging = ".Offset((paging.Page - 1) * paging.Take).Limit(paging.Take).Order(paging.OrderBy)"
				}
				tempFindByMethod = strings.Replace(templateFindByMethod, "[item]", itemLower, -1)
				tempFindByMethod = strings.Replace(tempFindByMethod, "[item_]", item_, -1)
				tempFindByMethod = strings.Replace(tempFindByMethod, "[itemUpper]", strings.Split(items[i], " ")[0], -1)
				tempFindByMethod = strings.Replace(tempFindByMethod, "[itemParam]", itemLower+" "+strings.Split(items[i], " ")[1]+pagingParam, -1)
				tempFindByMethod = strings.Replace(tempFindByMethod, "[paging]", paging, -1)
				findByMethod += tempFindByMethod

				tempFindBy = strings.Replace(templateFindBy, "[itemUpper]", strings.Split(items[i], " ")[0], -1)
				tempFindBy = strings.Replace(tempFindBy, "[item]", itemLower, -1)
				tempFindBy = strings.Replace(tempFindBy, "[itemParam]", itemLower+" "+strings.Split(items[i], " ")[1]+pagingParam, -1)

				findBy += tempFindBy + "\n"
			}
		}
	}

	template = strings.Replace(template, "[name]", name, -1)
	template = strings.Replace(template, "[nameUpper]", nameUpper, -1)
	template = strings.Replace(template, "[project]", project, -1)
	template = strings.Replace(template, "[findByMethod]", findByMethod, -1)
	template = strings.Replace(template, "[findBy]", findBy, -1)
	template = strings.Replace(template, "[preload]", preload, -1)

	_, err = fmt.Fprintln(file, template)
	if err != nil {
		return err
	}

	fmt.Println(name + " repository Created")
	return nil
}

func createEntity(items []string, name string, project string, relation []map[string]string) ([]string, error) {
	itemReturn := []string{}
	err := os.MkdirAll(project+"/entity", os.ModePerm)
	if err != nil {
		return items, err
	}

	file, err := os.Create(project + "/entity/" + name + "Entity.go")

	if err != nil {
		return items, err
	}

	defer file.Close()

	names := strings.Split(name, "")
	names[0] = strings.ToUpper(names[0])
	nameUpper := strings.Join(names, "")
	var codes = []string{
		"package entity",
		"",
		"import \"time\"",
		"",
		"type " + nameUpper + " struct{",
	}

	var status []string
	var relationPartner []string

	for _, value := range relation {
		if value["table1"] == name {
			if value["relationName"] == "1M" {
				status = append(status, "table1 1M")
				relationPartner = append(relationPartner, strings.Split(value["table2"], " ")[0])
			}
			if value["relationName"] == "MM" {
				status = append(status, "table1 MM")
				relationPartner = append(relationPartner, strings.Split(value["table2"], " ")[0])
			}
			if value["relationName"] == "11" {
				status = append(status, "table1 11")
				relationPartner = append(relationPartner, strings.Split(value["table2"], " ")[0])
			}
		}

		if value["table2"] == name {
			if strings.Contains(value["relationName"], "1M") {
				status = append(status, "table2 1M")
				relationPartner = append(relationPartner, strings.Split(value["table1"], " ")[0])
			}
			if value["relationName"] == "MM" {
				status = append(status, "table2 MM")
				relationPartner = append(relationPartner, strings.Split(value["table1"], " ")[0])
			}
			if value["relationName"] == "11" {
				status = append(status, "table2 11")
				relationPartner = append(relationPartner, strings.Split(value["table1"], " ")[0])
			}
		}

	}

	var relationPartnerUpper []string
	for _, value := range relationPartner {
		if value != "" {
			temp := strings.Split(value, "")
			temp[0] = strings.ToUpper(temp[0])
			relationPartnerUpper = append(relationPartnerUpper, strings.Join(temp, ""))
		}
	}

	for _, value := range items {
		for key, vStatus := range status {

			if vStatus == "table1 1M" {
				if value == "CreatedBy string" {
					codes = append(codes, "	"+relationPartnerUpper[key]+"s []"+relationPartnerUpper[key]+"`gorm:\"ForeignKey:"+nameUpper+"Id\"`")
					itemReturn = append(itemReturn, relationPartnerUpper[key]+"s")
				}
			}
			if vStatus == "table2 1M" {
				if value == "CreatedBy string" {
					codes = append(codes, "	"+relationPartnerUpper[key]+"Id int")
					itemReturn = append(itemReturn, relationPartnerUpper[key]+"Id int "+relationPartnerUpper[key])
					itemReturn = append(itemReturn, relationPartnerUpper[key])
					codes = append(codes, "	"+relationPartnerUpper[key]+" "+relationPartnerUpper[key]+" `gorm:\"ForeignKey:"+relationPartnerUpper[key]+"Id\"`")
				}
			}
			if vStatus == "table1 11" {
				if value == "CreatedBy string" {
					codes = append(codes, "	"+relationPartnerUpper[key]+"Id int")
					itemReturn = append(itemReturn, relationPartnerUpper[key]+"Id int "+relationPartnerUpper[key])
					itemReturn = append(itemReturn, relationPartnerUpper[key])
					codes = append(codes, "	"+relationPartnerUpper[key]+" "+relationPartnerUpper[key]+" `gorm:\"ForeignKey:"+relationPartnerUpper[key]+"Id\"`")
				}
			}
		}
		itemReturn = append(itemReturn, value)
		if value == "Id int" {
			codes = append(codes, "	"+value+"`gorm:\"primarykey;autoIncrement:true\"`")
		} else {
			codes = append(codes, "	"+value)
		}
	}

	codes = append(codes, "}")

	for _, code := range codes {
		_, err := fmt.Fprintln(file, code)
		if err != nil {
			return itemReturn, err
		}

	}
	fmt.Println(name + " entity Created")
	return itemReturn, nil
}
