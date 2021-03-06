package server

import (
	"log"
	"net/http"

	"github.com/fabbricadigitale/scimd/api/patch"

	"github.com/gin-gonic/gin/binding"

	"github.com/fabbricadigitale/scimd/api/delete"
	"github.com/globalsign/mgo"

	"github.com/fabbricadigitale/scimd/api"
	"github.com/fabbricadigitale/scimd/api/create"
	"github.com/fabbricadigitale/scimd/api/messages"
	"github.com/fabbricadigitale/scimd/api/query"
	"github.com/fabbricadigitale/scimd/api/update"
	"github.com/fabbricadigitale/scimd/schemas/core"
	"github.com/fabbricadigitale/scimd/schemas/resource"
	"github.com/fabbricadigitale/scimd/storage"
	"github.com/gin-gonic/gin"
)

// ResourceService describes ...
type ResourceService struct {
	rt core.ResourceType
	Service
	Lister
	Getter
	Putter
	Patcher
	Deleter
	Searcher
}

// NewResourceService creates a new `ResourceService` for the given `core.ResourceTyper`
func NewResourceService(rt core.ResourceType) *ResourceService {
	return &ResourceService{
		rt: rt,
	}
}

// Path returns the endpoint of the `ResourceService`
func (rs *ResourceService) Path() string {
	return rs.rt.Endpoint
}

// List ...
func (rs *ResourceService) List(c *gin.Context) {
	params := api.NewSearch()
	// Using the form binding engine (query)
	if err := c.ShouldBindWith(params, binding.Form); err != nil {
		err := messages.NewError(err)
		c.JSON(err.Status, err)
	}

	// Go ahead ...
	params.Attributes.Explode()
	log.Printf("%+v\n", params)

	// Retrieve the storage adapter
	store, ok := c.Get("storage")
	if !ok {
		err := messages.NewError(&api.InternalServerError{
			Detail: "Missing storage setup ...",
		})
		c.JSON(err.Status, err)
	}

	rtArr := make([]*core.ResourceType, 0)
	rtArr = append(rtArr, &rs.rt)

	list, err := query.Resources(store.(storage.Storer), rtArr, params)
	if err != nil {
		err := messages.NewError(err)
		c.JSON(err.Status, err)
	}

	c.JSON(http.StatusOK, list)
}

// Get ...
func (rs *ResourceService) Get(c *gin.Context) {
	var attrs api.Attributes
	// Using the form binding engine (query)
	if err := c.ShouldBindWith(&attrs, binding.Form); err != nil {
		err := messages.NewError(err)
		c.JSON(err.Status, err)
	}

	// Explode the attributes
	attrs.Explode()

	// Retrieve the storage adapter
	store, ok := c.Get("storage")
	if !ok {
		err := messages.NewError(&api.InternalServerError{
			Detail: "Missing storage setup ...",
		})
		c.JSON(err.Status, err)
	}
	// Retrieve the id segment
	id := c.Param("id")

	res, err := query.Resource(store.(storage.Storer), &rs.rt, id, &attrs)
	if err != nil {
		err := messages.NewError(err)
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, res.(*resource.Resource))
}

// Post ...
func (rs *ResourceService) Post(c *gin.Context) {
	var contents resource.Resource
	if err := c.ShouldBindWith(&contents, binding.JSON); err != nil {
		err := messages.NewError(err)
		c.JSON(err.Status, err)
	} else {
		// Retrieve the storage adapter
		store, ok := c.Get("storage")
		if !ok {
			err := messages.NewError(&api.InternalServerError{
				Detail: "Missing storage setup ...",
			})
			c.JSON(err.Status, err)
		}

		res, err := create.Resource(store.(storage.Storer), &rs.rt, &contents)
		if err != nil {
			we := handlingError(err)
			c.JSON(we.Status, we)
			return
		}
		c.JSON(http.StatusOK, res.(*resource.Resource))
	}
}

// Search ...
func (rs *ResourceService) Search(c *gin.Context) {
	contents := &messages.SearchRequest{}
	if err := c.ShouldBindWith(contents, binding.JSON); err != nil {
		err := messages.NewError(err)
		c.JSON(err.Status, err)
	}

	// Retrieve the storage adapter
	store, ok := c.Get("storage")
	if !ok {
		err := messages.NewError(&api.InternalServerError{
			Detail: "Missing storage setup ...",
		})
		c.JSON(err.Status, err)
	}

	rtArr := make([]*core.ResourceType, 0)
	rtArr = append(rtArr, &rs.rt)

	list, err := query.SearchRequest(store.(storage.Storer), rtArr, contents)
	if err != nil {
		err := messages.NewError(err)
		c.JSON(err.Status, err)
	}
	c.JSON(http.StatusOK, list)

}

// Put ...
func (rs *ResourceService) Put(c *gin.Context) {
	var attrs api.Attributes
	// Using the form binding engine (query)
	if err := c.ShouldBindWith(&attrs, binding.Form); err != nil {
		err := messages.NewError(err)
		c.JSON(err.Status, err)
	}

	// Explode the attributes
	attrs.Explode()
	// Retrieve the id segment
	id := c.Param("id")

	contents := &resource.Resource{}
	if err := c.ShouldBindWith(contents, binding.JSON); err != nil {
		err := messages.NewError(err)
		c.JSON(err.Status, err)
	}

	// Retrieve the storage adapter
	store, ok := c.Get("storage")
	if !ok {
		err := messages.NewError(&api.InternalServerError{
			Detail: "Missing storage setup ...",
		})
		c.JSON(err.Status, err)
	}
	res, err := update.Resource(store.(storage.Storer), &rs.rt, id, contents)
	if err != nil {
		we := handlingError(err)
		c.JSON(we.Status, we)
		return
	}
	c.JSON(http.StatusOK, res.(*resource.Resource))

}

// Patch ...
func (rs *ResourceService) Patch(c *gin.Context) {

	contents := &messages.PatchOp{}
	if err := c.ShouldBindWith(contents, binding.JSON); err != nil {
		err := messages.NewError(err)
		c.JSON(err.Status, err)
		return
	}

	// Using the form binding engine (query)
	var attrs api.Attributes
	if err := c.ShouldBindWith(&attrs, binding.Form); err != nil {
		err := messages.NewError(err)
		c.JSON(err.Status, err)
		return
	}

	// Explode the attributes
	attrs.Explode()
	// Retrieve the id segment
	id := c.Param("id")

	// Retrieve the storage adapter
	store, ok := c.Get("storage")
	if !ok {
		err := messages.NewError(&api.InternalServerError{
			Detail: "Missing storage setup ...",
		})
		c.JSON(err.Status, err)
		return
	}

	res, err := patch.Resource(store.(storage.Storer), &rs.rt, id, contents.Operations[0].Op, contents.Operations[0].Path, contents.Operations[0].Value)
	if err != nil {
		we := handlingError(err)
		c.JSON(we.Status, we)
		return
	}
	c.JSON(http.StatusOK, res.(*resource.Resource))

}

// Delete ...
func (rs *ResourceService) Delete(c *gin.Context) {

	id := c.Param("id")
	// Retrieve the storage adapter
	store, ok := c.Get("storage")
	if !ok {
		err := messages.NewError(&api.InternalServerError{
			Detail: "Missing storage setup ...",
		})
		c.JSON(err.Status, err)
	}

	err := delete.Resource(store.(storage.Storer), &rs.rt, id)
	if err != nil {
		err := messages.NewError(err)
		c.JSON(err.Status, err)
	} else {
		c.JSON(http.StatusOK, nil)
	}
}

func handlingError(err error) messages.Error {

	detail := err.Error()

	var newError error

	switch err.(type) {
	case *mgo.LastError:
		newError = &api.UniquenessError{
			Detail: detail,
		}
	default:
		newError = &api.InternalServerError{
			Detail: detail,
		}
	}

	return messages.NewError(newError)
}
