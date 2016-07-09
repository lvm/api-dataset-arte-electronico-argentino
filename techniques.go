package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/gin-gonic/gin"
)

type (
	Technique struct {
		Id        int64 `db:"id" json:"-"`
		Technique string `db:"tecnica" json:"technique,omitempty"`
		Total     string `db:"count(tecnica)" json:"total,omitempty"`
	}

	ResourceTechnique struct {
		Id         int64 `json:"id"`
		Type       string `json:"type"`
		Attributes Technique `json:"attributes"`
	}

	ResourceTechniques struct {
		Data  []ResourceTechnique `json:"data"`
		Links PaginationLinks `json:"links"`
		Meta  Meta `json:"meta"`
	}
)

func newResourceTechnique (attr Technique) ResourceTechnique {
	return ResourceTechnique {
		Id: attr.Id,
		Type: "techniques",
		Attributes: attr,
	}
}

func newResourceTechniques (techniques []Technique, endpoint string, pg string) ResourceTechniques {
	resources := make([]ResourceTechnique, 0)
	for l := range techniques {
		if techniques[l].Technique != "" {
			resources = append(resources, newResourceTechnique(techniques[l]))
		}
	}

	links := PaginationLinks {
		First: paginationLink(endpoint, firstPage, "first"),
		Last: paginationLink(endpoint, lastPage, "last"),
		Prev: paginationLink(endpoint, pg, "prev"),
		Next: paginationLink(endpoint, pg, "next"),
	}

	return ResourceTechniques{Data: resources, Links: links, Meta: Meta{TotalPages: "6"}}
}


// http http://localhost:8080/api/techniques
func GetTechniques(c *gin.Context) {
	var techniques []Technique
	page := c.DefaultQuery("page", "1")

	query := PaginateQuery(groupByQuery("SELECT id, %s, count(%s) FROM exhibiciones", "technique"), page)
	_, err := dbmap.Select(&techniques, query)
	checkErr(err)

	if err == nil {
		content := newResourceTechniques(techniques, "techniques", page)
		c.JSON(200, content)
	} else {
		jsonErr := &Error{"not_found", 404, "Not Found Error", "No technique found."}
		c.JSON(404, jsonErr)
	}
}
