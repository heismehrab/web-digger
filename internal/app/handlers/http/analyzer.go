package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web-digger/pkg/logger"
)

type urlToAnalyzeRequest struct {
	Url string `json:"url" binding:"required,url"`
}

type urlToAnalyzeResponse struct {
}

func (h *Handler) AnalyzeWebPage(c *gin.Context) {
	var reqBody urlToAnalyzeRequest

	if err := c.Bind(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, GetFailedValidationResponse(err))

		return
	}

	// Processing the given URL to obtain its details.
	_, err := h.WebAnalyzerService.Analyze(c, reqBody.Url)

	if err != nil {
		h.logger.LogAPIError(c.Request.RequestURI, err, logger.GetErrorSource(), logger.Detail{
			"URL": reqBody.Url,
		})

		c.JSON(http.StatusInternalServerError, InternalServerError)

		return
	}

	c.JSON(http.StatusOK, GetSuccessResponse(urlToAnalyzeResponse{}))
}
