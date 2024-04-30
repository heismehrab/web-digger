package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web-digger/internal/core/domain/models"
	"web-digger/pkg/logger"
)

type urlToAnalyzeRequest struct {
	Url string `json:"url" binding:"required,url"`
}

type urlToAnalyzeResponse struct {
	Title   string       `json:"title"`
	Version string       `json:"version"`
	HTags   models.HTags `json:"headings"`
}

func (h *Handler) AnalyzeWebPage(c *gin.Context) {
	var reqBody urlToAnalyzeRequest

	if err := c.Bind(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, GetFailedValidationResponse(err))

		return
	}

	// Processing the given URL to obtain its details.
	res, err := h.WebAnalyzerService.Analyze(c, reqBody.Url)

	if err != nil {
		h.logger.LogAPIError(c.Request.RequestURI, err, logger.GetErrorSource(), logger.Detail{
			"URL": reqBody.Url,
		})

		c.JSON(http.StatusInternalServerError, InternalServerError)

		return
	}

	c.JSON(http.StatusOK, GetSuccessResponse(urlToAnalyzeResponse{
		Title:   res.Title,
		Version: res.Version,
		HTags:   res.Hs,
	}))
}
