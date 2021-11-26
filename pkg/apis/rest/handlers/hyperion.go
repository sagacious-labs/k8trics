package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

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

	resp, err := h.performRequest(func(ep string) (interface{}, error) {
		return rpc.HyperionApply(c.Request.Context(), &req, ep)
	})
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

	resp, err := h.performRequest(func(ep string) (interface{}, error) {
		return rpc.HyperionDelete(c.Request.Context(), &req, ep)
	})
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

	resp, err := h.performRequest(func(ep string) (interface{}, error) {
		return rpc.HyperionGet(c.Request.Context(), &req, ep)
	})
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

	resp, err := h.performRequestWithChannel(func(ep string) (chan interface{}, error) {
		resp, err := rpc.HyperionList(c.Request.Context(), &req, ep)
		ch := make(chan interface{}, 8)

		go func() {
			for data := range resp {
				ch <- data
			}
		}()

		return ch, err
	})
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

	resp, err := h.performRequestWithChannel(func(ep string) (chan interface{}, error) {
		ctx := context.WithValue(c.Request.Context(), "pod_store", h.store)
		resp, err := rpc.HyperionWatchData(ctx, &req, ep)
		ch := make(chan interface{}, 8)

		go func() {
			for data := range resp {
				ch <- data
			}
		}()

		return ch, err
	})
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

	resp, err := h.performRequestWithChannel(func(ep string) (chan interface{}, error) {
		resp, err := rpc.HyperionWatchLog(c.Request.Context(), &req, ep)
		ch := make(chan interface{}, 8)

		go func() {
			for data := range resp {
				ch <- data
			}
		}()

		return ch, err
	})
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

func (h *Handlers) performRequest(fn func(ep string) (interface{}, error)) (interface{}, error) {
	errs := []error{}
	ress := []interface{}{}

	for _, pod := range h.store.GetByLabels(map[string]string{
		"core.hyperion.io/master": "true",
	}) {
		endpoint, err := pod.Endpoint()
		if err != nil {
			errs = append(errs, err)
			continue
		}

		res, err := fn(endpoint)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		ress = append(ress, res)
	}

	return ress, mergeErrors(errs)
}

func (h *Handlers) performRequestWithChannel(fn func(ep string) (chan interface{}, error)) (chan interface{}, error) {
	errs := []error{}
	centralCh := make(chan interface{}, 8)

	for _, pod := range h.store.GetByLabels(map[string]string{
		"core.hyperion.io/master": "true",
	}) {
		endpoint, err := pod.Endpoint()
		if err != nil {
			errs = append(errs, err)
			continue
		}

		res, err := fn(endpoint)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		go func() {
			for data := range res {
				centralCh <- data
			}
		}()
	}

	return centralCh, mergeErrors(errs)
}
func mergeErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	errStrs := []string{}

	for _, err := range errs {
		errStrs = append(errStrs, err.Error())
	}

	return errors.New(strings.Join(errStrs, "\n"))
}
