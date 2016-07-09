package main

import (
	"time"
	_ "github.com/mattn/go-sqlite3"
	"github.com/gin-gonic/gin"
)

type (
	Exhibition struct {
		Id         int64 `db:"id" json:"-"`
		Anio       int `db:"anio" json:"year,omitempty"`
		Nombre     string `db:"nombre_exhibicion" json:"exhibition,omitempty"`
		Locacion   string `db:"nombre_locacion" json:"location,omitempty"`
		Curador0   string `db:"curador_0" json:"curator_0,omitempty"`
		Curador1   string `db:"curador_1" json:"curator_1,omitempty"`
		Curador2   string `db:"curador_2" json:"curator_2,omitempty"`
		Curador3   string `db:"curador_3" json:"curator_3,omitempty"`
		FechaIncio time.Time `db:"fecha_inicio" json:"start,omitempty"`
		FechaFin   time.Time `db:"fecha_finalizacion" json:"end,omitempty"`
		Obra       string `db:"nombre_obra" json:"work,omitempty"`
		Artista0   string `db:"nombre_artista_0" json:"artist_0,omitempty"`
		Artista1   string `db:"nombre_artista_1" json:"artist_1,omitempty"`
		Artista2   string `db:"nombre_artista_2" json:"artist_2,omitempty"`
		Artista3   string `db:"nombre_artista_3" json:"artist_3,omitempty"`
		Artista4   string `db:"nombre_artista_4" json:"artist_3,omitempty"`
		Artista5   string `db:"nombre_artista_5" json:"artist_5,omitempty"`
		Artista6   string `db:"nombre_artista_6" json:"artist_6,omitempty"`
		Artista7   string `db:"nombre_artista_7" json:"artist_7,omitempty"`
		Artista8   string `db:"nombre_artista_8" json:"artist_8,omitempty"`
		Tecnica    string `db:"tecnica" json:"technique,omitempty"`
	}

	ResourceExhibition struct {
		Id         int64 `json:"id"`
		Type       string `json:"type"`
		Attributes Exhibition `json:"attributes"`
	}

	ResourceExhibitions struct {
		Data  []ResourceExhibition `json:"data"`
		Links PaginationLinks `json:"links"`
		Meta  Meta `json:"meta"`
	}
)

func newResourceExhibition (attr Exhibition) ResourceExhibition {
	return ResourceExhibition {
		Id: attr.Id,
		Type: "exhibitions",
		Attributes: attr,
	}
}

func newResourceExhibitions (exhibitions []Exhibition, endpoint string, pg string, totalPages string) ResourceExhibitions {
	resources := make([]ResourceExhibition, 0)
	for e := range exhibitions {
		resources = append(resources, newResourceExhibition(exhibitions[e]))
	}

	links := PaginationLinks {First: "",Last: "",Prev: "",Next: ""}
	if pg != "" {
		links = PaginationLinks {
			First: paginationLink(endpoint, firstPage, "first"),
			Last: paginationLink(endpoint, lastPage, "last"),
			Prev: paginationLink(endpoint, pg, "prev"),
			Next: paginationLink(endpoint, pg, "next"),
		}
	}

	if totalPages == "" {
		totalPages = "1"
	}

	return ResourceExhibitions{Data: resources, Links: links, Meta: Meta{TotalPages: totalPages}}
}


// http http://localhost:8080/api/exhibitions
func GetExhibitions(c *gin.Context) {
	var exhibitions []Exhibition
	page := c.DefaultQuery("page", "1")

	query := PaginateQuery("SELECT * FROM exhibiciones ORDER BY id", page)
	_, err := dbmap.Select(&exhibitions, query)
	checkErr(err)

	if err == nil {
		content := newResourceExhibitions(exhibitions, "exhibitions", page, "80")
		c.JSON(200, content)
	} else {
		jsonErr := &Error{"not_found", 404, "Not Found Error", "No exhibition found."}
		c.JSON(404, jsonErr)
	}
}

// http http://localhost:8080/api/exhibitions/1
func GetExhibition(c *gin.Context) {
	id := c.Params.ByName("e_id")
	var exhibition Exhibition
	err := dbmap.SelectOne(&exhibition, "SELECT * FROM exhibiciones WHERE id=? LIMIT 1", id)
	checkErr(err)

	if err == nil {
		content := newResourceExhibition(exhibition)
		c.JSON(200, content)
	} else {
		jsonErr := &Error{"not_found", 404, "Not Found Error", "Your query didn't match any exhibition."}
		c.JSON(404, jsonErr)
	}
}
