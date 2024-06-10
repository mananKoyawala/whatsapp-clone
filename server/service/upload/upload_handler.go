package upload

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	api "github.com/mananKoyawala/whatsapp-clone/internal"
)

type AwsHandler struct {
	AwsService
}

func NewAwsHandler(a AwsService) AwsHandler {
	return AwsHandler{
		AwsService: a,
	}
}

func (h *AwsHandler) UploaFile(c *gin.Context) (int, error) {

	// accept the file
	form, err := c.MultipartForm()
	if err != nil {
		return http.StatusBadRequest, err
	}

	files := form.File["files"]

	numbersOfFile := len(files)
	if numbersOfFile != 1 {
		return http.StatusBadRequest, errors.New("only one file is allowed")
	}

	res, err := h.AwsService.UploaFile(files)

	if err != nil {
		return api.WriteError(c, http.StatusInternalServerError, res)
	}

	return api.WriteData(c, http.StatusOK, res)
}

func (h *AwsHandler) DeleteFile(c *gin.Context) (int, error) {

	url := c.Request.FormValue("url")

	// extract the name of the image
	parts := strings.Split(url, "/")
	imageName := parts[len(parts)-1]

	err := h.AwsService.deleteFile(imageName)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	return api.WriteMessage(c, http.StatusOK, "file delelted")
}
