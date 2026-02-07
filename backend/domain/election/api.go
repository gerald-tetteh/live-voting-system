package election

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ElectionAPI struct {
	service *ElectionService
	log *zap.Logger
}

func NewElectionAPI(service *ElectionService, logger *zap.Logger) *ElectionAPI {
	return &ElectionAPI{service: service, log: logger}
}

func (api *ElectionAPI) RegisterRoutes(server *gin.Engine) {
	server.GET("/elections", api.getElections)
	server.GET("/elections/:id", api.getElection)
	server.POST("/elections", api.createElection)
	server.PATCH("/elections/:id", api.updateElection)
}

func (api *ElectionAPI) createElection(ctx *gin.Context) {
	requestId := ctx.GetString("requestId")
	var election Election
	err := ctx.ShouldBindJSON(&election)
	if err != nil {
		api.log.Error(err.Error())
		api.log.Error("could not parse election", zap.String("request_id", requestId))
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "could not parse election"})
		return
	}
	err = api.service.CreateElection(ctx, &election)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "created election"})
}

func (api *ElectionAPI) getElections(ctx *gin.Context) {
	requestId := ctx.GetString("requestId")
	rawStatus := ctx.DefaultQuery("status", "draft")
	status := ElectionStatus(rawStatus)
	if !status.IsValid() {
		api.log.Error("status is invalid", zap.String("request_id", requestId))
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "status is invalid"})
		return
	}
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		api.log.Error(err.Error())
		api.log.Error("could not parse page number", zap.String("request_id", requestId))
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "could not parse page number"})
		return
	}
	pageSize, err := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if err != nil {
		api.log.Error(err.Error())
		api.log.Error("could not parse page size", zap.String("request_id", requestId))
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "could not parse page size"})
		return
	}
	offset := (page - 1) * pageSize
	queryParams := ElectionQueryParams{
		Status: status,
		Limit: pageSize,
		Offset: offset,
	}
	elections, err := api.service.GetElections(ctx, queryParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, elections)
}

func (api *ElectionAPI) getElection(ctx *gin.Context) {
	electionId := ctx.Param("id")
	election, err := api.service.GetElection(ctx, electionId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, election)
}

func (api *ElectionAPI) updateElection(ctx *gin.Context) {
	requestId := ctx.GetString("requestId")
	electionId := ctx.Param("id")
	var updatedElection Election
	err := ctx.ShouldBindJSON(&updatedElection)
	if err != nil {
		api.log.Error(err.Error())
		api.log.Error("could not parse update information", zap.String("request_id", requestId))
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "could not parse update information"})
		return
	}
	err = api.service.UpdateElection(ctx, electionId, &updatedElection)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "updated election"})
}