func (h *[name]Handler) Get[nameUpper]By[itemUpper](c *gin.Context) {
	[item] := c.Param("[item]")
	[convert]

	[name], err := h.[name]Service.Get[nameUpper]By[itemUpper]([itemParam])

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Get [nameUpper] Failed", http.StatusUnprocessableEntity, "Failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.APIResponse("Get [nameUpper] Success", http.StatusOK, "Success", [name])

	c.JSON(http.StatusOK, response)
}