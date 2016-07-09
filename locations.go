package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/gin-gonic/gin"
)

type (
	Location struct {
		Id         int64 `db:"id" json:"-"`
		Location   string `db:"nombre_locacion" json:"location,omitempty"`
		Total      string `db:"count(nombre_locacion)" json:"total,omitempty"`
	}

	ResourceLocation struct {
		Id         int64 `json:"id"`
		Type       string `json:"type"`
		Attributes Location `json:"attributes"`
	}

	ResourceLocations struct {
		Data  []ResourceLocation `json:"data"`
		Links PaginationLinks `json:"links"`
		Meta  Meta `json:"meta"`
	}
)

func newResourceLocation (attr Location) ResourceLocation {
	return ResourceLocation {
		Id: attr.Id,
		Type: "locations",
		Attributes: attr,
	}
}

func newResourceLocations (locations []Location, endpoint string, pg string) ResourceLocations {
	resources := make([]ResourceLocation, 0)
	for l := range locations {
		resources = append(resources, newResourceLocation(locations[l]))
	}

	links := PaginationLinks {
		First: paginationLink(endpoint, firstPage, "first"),
		Last: paginationLink(endpoint, lastPage, "last"),
		Prev: paginationLink(endpoint, pg, "prev"),
		Next: paginationLink(endpoint, pg, "next"),
	}

	return ResourceLocations{Data: resources, Links: links, Meta: Meta{TotalPages: "1"}}
}


// http http://localhost:8080/api/locations
func GetLocations(c *gin.Context) {
	var locations []Location
	page := c.DefaultQuery("page", "1")

	query := PaginateQuery(groupByQuery("SELECT id, %s, count(%s) FROM exhibiciones", "location"), page)
	_, err := dbmap.Select(&locations, query)
	checkErr(err)

	if err == nil {
		content := newResourceLocations(locations, "locations", page)
		c.JSON(200, content)
	} else {
		jsonErr := &Error{"not_found", 404, "Not Found Error", "No location found."}
		c.JSON(404, jsonErr)
	}
}
