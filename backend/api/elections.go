package api

import (
	"net/http"
	"strconv"

	"geraldaddo.com/live-voting-system/models"
	"geraldaddo.com/live-voting-system/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ElectionAPI struct {
	Service *services.ElectionService
	Logger *zap.Logger
}

func (api *ElectionAPI) RegisterRoutes(server *gin.Engine) {
	server.POST("/elections", api.createElection)
	server.GET("/elections", api.getElections)
}

func (api *ElectionAPI) createElection(ctx *gin.Context) {
	requestId := ctx.GetString("requestId")
	var election models.Election
	err := ctx.ShouldBindJSON(&election)
	if err != nil {
		api.Logger.Error(err.Error())
		api.Logger.Error("could not parse election", zap.String("request_id", requestId))
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "could not parse election"})
		return
	}
	err = api.Service.Save(&election)
	if err != nil {
		api.Logger.Error(err.Error())
		api.Logger.Error("could not create election", zap.String("request_id", requestId))
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not create election"})
		return
	}
	api.Logger.Info("Created election", zap.String("request_id", requestId))
	ctx.JSON(http.StatusOK, gin.H{"message": "created election"})
}
func (api * ElectionAPI) getElections(ctx *gin.Context) {
	requestId := ctx.GetString("requestId")

	rawStatus := ctx.Query("status")
	status := models.ElectionStatus(rawStatus)
	if !status.IsValid() {
		api.Logger.Error("status is invalid", zap.String("request_id", requestId))
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "status is invalid"})
		return
	}
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		api.Logger.Error(err.Error())
		api.Logger.Error("could not parse page number", zap.String("request_id", requestId))
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "could not parse page number"})
		return
	}
	pageSize, err := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if err != nil {
		api.Logger.Error(err.Error())
		api.Logger.Error("could not parse page size", zap.String("request_id", requestId))
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "could not parse page size"})
		return
	}

	offset := (page - 1) * pageSize
	
	queryParams := services.ElectionQueryParams{
		Status: status,
		Limit: pageSize,
		Offset: offset,
	}

	elections, err := api.Service.GetAll(queryParams)
	if err != nil {
		api.Logger.Error(err.Error())
		api.Logger.Error("could not get elections", zap.String("request_id", requestId))
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not get elections"})
	}
	ctx.JSON(http.StatusOK, elections)
}