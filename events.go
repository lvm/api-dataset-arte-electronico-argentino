package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/gin-gonic/gin"
)

type (
	Event struct {
		Id    int64 `db:"id" json:"-"`
		Event string `db:"nombre_exhibicion" json:"event,omitempty"`
		Total string `db:"count(nombre_exhibicion)" json:"total,omitempty"`
	}

	ResourceEvent struct {
		Id         int64 `json:"id"`
		Type       string `json:"type"`
		Attributes Event `json:"attributes"`
	}

	ResourceEvents struct {
		Data  []ResourceEvent `json:"data"`
		Links PaginationLinks `json:"links"`
		Meta  Meta `json:"meta"`
	}
)

func newResourceEvent (attr Event) ResourceEvent {
	return ResourceEvent {
		Id: attr.Id,
		Type: "events",
		Attributes: attr,
	}
}

func newResourceEvents (events []Event, endpoint string, pg string) ResourceEvents {
	resources := make([]ResourceEvent, 0)
	for l := range events {
		resources = append(resources, newResourceEvent(events[l]))
	}

	links := PaginationLinks {
		First: paginationLink(endpoint, firstPage, "first"),
		Last: paginationLink(endpoint, lastPage, "last"),
		Prev: paginationLink(endpoint, pg, "prev"),
		Next: paginationLink(endpoint, pg, "next"),
	}

	return ResourceEvents{Data: resources, Links: links, Meta: Meta{TotalPages: "3"}}
}


// http http://localhost:8080/api/events
func GetEvents(c *gin.Context) {
	var events []Event
	page := c.DefaultQuery("page", "1")

	query := PaginateQuery(groupByQuery("SELECT id, %s, count(%s) FROM exhibiciones", "event"), page)
	_, err := dbmap.Select(&events, query)
	checkErr(err)

	if err == nil {
		content := newResourceEvents(events, "events", page)
		c.JSON(200, content)
	} else {
		jsonErr := &Error{"not_found", 404, "Not Found Error", "No event found."}
		c.JSON(404, jsonErr)
	}
}
