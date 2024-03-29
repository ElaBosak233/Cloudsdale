package controller

import (
	"github.com/elabosak233/cloudsdale/internal/model/request"
	"github.com/elabosak233/cloudsdale/internal/model/response"
	"github.com/elabosak233/cloudsdale/internal/service"
	"github.com/elabosak233/cloudsdale/internal/utils/convertor"
	"github.com/elabosak233/cloudsdale/internal/utils/validator"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IPodController interface {
	Create(ctx *gin.Context)
	Remove(ctx *gin.Context)
	Renew(ctx *gin.Context)
	Find(ctx *gin.Context)
	FindById(ctx *gin.Context)
}

type PodController struct {
	podService service.IPodService
}

func NewInstanceController(appService *service.Service) IPodController {
	return &PodController{
		podService: appService.PodService,
	}
}

// Create
// @Summary 创建实例
// @Description	创建实例
// @Tags Pod
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input	body	request.PodCreateRequest	true	"PodCreateRequest"
// @Router /pods/ [post]
func (c *PodController) Create(ctx *gin.Context) {
	instanceCreateRequest := request.PodCreateRequest{}
	err := ctx.ShouldBindJSON(&instanceCreateRequest)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  validator.GetValidMsg(err, &instanceCreateRequest),
		})
		return
	}
	user, _ := ctx.Get("user")
	instanceCreateRequest.UserID = user.(*response.UserResponse).ID
	res, err := c.podService.Create(instanceCreateRequest)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":       http.StatusOK,
		"id":         res.ID,
		"instances":  res.Instances,
		"removed_at": res.RemovedAt,
	})
}

// Remove
// @Summary 停止并删除容器
// @Description	停止并删除容器
// @Tags Pod
// @Produce json
// @Security ApiKeyAuth
// @Param input	body request.PodRemoveRequest true "PodRemoveRequest"
// @Router /pods/{id} [delete]
func (c *PodController) Remove(ctx *gin.Context) {
	instanceRemoveRequest := request.PodRemoveRequest{}
	err := ctx.ShouldBindJSON(&instanceRemoveRequest)
	instanceRemoveRequest.ID = convertor.ToUintD(ctx.Param("id"), 0)
	instanceRemoveRequest.UserID = ctx.GetUint("GameID")
	err = c.podService.Remove(instanceRemoveRequest)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}

// Renew
// @Summary 容器续期
// @Description	容器续期
// @Tags Pod
// @Produce json
// @Security ApiKeyAuth
// @Param input	body request.PodRenewRequest true "PodRenewRequest"
// @Router /pods/{id} [put]
func (c *PodController) Renew(ctx *gin.Context) {
	instanceRenewRequest := request.PodRenewRequest{}
	err := ctx.ShouldBindJSON(&instanceRenewRequest)
	instanceRenewRequest.ID = convertor.ToUintD(ctx.Param("id"), 0)
	instanceRenewRequest.UserID = ctx.GetUint("GameID")
	removedAt, err := c.podService.Renew(instanceRenewRequest)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":       http.StatusOK,
		"removed_at": removedAt,
	})
}

// FindById
// @Summary 实例查询
// @Description	实例查询
// @Tags Pod
// @Produce json
// @Param id path string true "id"
// @Router /pods/{id} [get]
func (c *PodController) FindById(ctx *gin.Context) {
	id := ctx.Param("id")
	rep, err := c.podService.FindById(convertor.ToUintD(id, 0))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": rep,
	})
}

// Find
// @Summary 实例查询
// @Description	实例查询
// @Tags Pod
// @Produce json
// @Security ApiKeyAuth
// @Param input	query request.PodFindRequest false "PodFindRequest"
// @Router /pods/ [get]
func (c *PodController) Find(ctx *gin.Context) {
	podFindRequest := request.PodFindRequest{}
	err := ctx.ShouldBindQuery(&podFindRequest)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  validator.GetValidMsg(err, &podFindRequest),
		})
		return
	}
	user, _ := ctx.Get("user")
	podFindRequest.UserID = user.(*response.UserResponse).ID
	pods, _ := c.podService.Find(podFindRequest)
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": pods,
	})
}
