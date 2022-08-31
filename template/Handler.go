package handler

import (
	"[project]/helper"
	"[project]/input"
	"[project]/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type [name]Handler struct {
	[name]Service service.[nameUpper]Service
	jwtService  service.JwtService
}

func New[nameUpper]Handler([name]Service service.[nameUpper]Service) *[name]Handler {
	return &[name]Handler{[name]Service}
}

func (h *[name]Handler) Create[nameUpper](c *gin.Context) {
	var input input.[nameUpper]Input

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Create [nameUpper] Failed", http.StatusUnprocessableEntity, "Failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	[name], err := h.[name]Service.Create[nameUpper](input)

	if err != nil {

		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Create [nameUpper] Failed", http.StatusUnprocessableEntity, "Failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.APIResponse("Create [nameUpper] Success", http.StatusOK, "Success", [name])

	c.JSON(http.StatusOK, response)
}

func (h *[name]Handler) Edit[nameUpper](c *gin.Context) {
	var input input.[nameUpper]EditInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Edit [nameUpper] Failed", http.StatusUnprocessableEntity, "Failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	[name], err := h.[name]Service.Edit[nameUpper](input)

	if err != nil {

		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Edit [nameUpper] Failed", http.StatusUnprocessableEntity, "Failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.APIResponse("Edit [nameUpper] Success", http.StatusOK, "Success", [name])

	c.JSON(http.StatusOK, response)
}

func (h *[name]Handler) GetAll[nameUpper]s(c *gin.Context) {
	[name]s, err := h.[name]Service.GetAll[nameUpper]()

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Get All [nameUpper] Failed", http.StatusUnprocessableEntity, "Failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.APIResponse("Get All [nameUpper] Success", http.StatusOK, "Success", [name]s)

	c.JSON(http.StatusOK, response)
}

func (h *[name]Handler) Get[nameUpper]ById(c *gin.Context) {
	id := c.Param("id")

	[name], err := h.[name]Service.Get[nameUpper]ById(id)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Get [nameUpper] Failed", http.StatusUnprocessableEntity, "Failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.APIResponse("Get [nameUpper] Success", http.StatusOK, "Success", [name])

	c.JSON(http.StatusOK, response)
}
