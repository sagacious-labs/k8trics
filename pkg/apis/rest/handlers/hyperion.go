package handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sagacious-labs/k8trics/pkg/protos/v1alpha1/api"
	"github.com/sagacious-labs/k8trics/pkg/protos/v1alpha1/base"
	"github.com/sagacious-labs/k8trics/pkg/rpc"
)

func (h *Handlers) Apply(c *gin.Context) {
	req := api.ApplyRequest{}
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "failed to parse request object: " + err.Error()})
		return
	}

	resp, err := rpc.HyperionApply(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *Handlers) Delete(c *gin.Context) {
	moduleName := c.Param("name")
	if moduleName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "module name is required"})
		return
	}

	req := api.DeleteRequest{
		Core: &base.ModuleCore{Name: moduleName},
	}

	resp, err := rpc.HyperionDelete(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handlers) Get(c *gin.Context) {
	moduleName := c.Param("name")
	if moduleName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "module name is required"})
		return
	}

	req := api.GetRequest{
		Core: &base.ModuleCore{Name: moduleName},
	}

	resp, err := rpc.HyperionGet(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handlers) List(c *gin.Context) {
	labels, ok := c.GetQueryMap("labels")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "labels are required as query params"})
		return
	}

	req := api.ListRequest{
		Filter: &api.ListRequest_Label{
			Label: &base.LabelSelector{
				Selector: labels,
			},
		},
	}

	resp, err := rpc.HyperionList(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.Stream(func(w io.Writer) bool {
		item, ok := <-resp
		if !ok {
			return false
		}

		c.SSEvent("module", item)
		return true
	})
}

func (h *Handlers) WatchData(c *gin.Context) {
	moduleName := c.Param("name")
	if moduleName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "module name is a require parameter"})
		return
	}

	req := api.WatchDataRequest{
		Filter: &base.ModuleCore{
			Name: moduleName,
		},
	}

	resp, err := rpc.HyperionWatchData(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.Stream(func(w io.Writer) bool {
		item, ok := <-resp
		if !ok {
			return false
		}

		c.SSEvent("data", item)
		return true
	})
}

func (h *Handlers) WatchLog(c *gin.Context) {
	moduleName := c.Param("name")
	if moduleName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "module name is a require parameter"})
		return
	}

	req := api.WatchLogRequest{
		Filter: &base.ModuleCore{
			Name: moduleName,
		},
	}

	resp, err := rpc.HyperionWatchLog(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.Stream(func(w io.Writer) bool {
		item, ok := <-resp
		if !ok {
			return false
		}

		c.SSEvent("log", item)
		return true
	})
}
