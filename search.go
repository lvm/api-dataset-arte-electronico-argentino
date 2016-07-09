package main

import (
	"fmt"
	"strconv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/gin-gonic/gin"
)

func howManyPages(query string) (int, error) {
	count, err := dbmap.SelectInt(query)
	return int(count) / perPage, err
}

// http http://localhost:8080/api/v1/search?q=festival[&year=1998[&when=since[&technique=videoarte[&artist=rytel[&curator=taquini[&work=sheep]]]]]]
func SearchExhibitions(c *gin.Context) {
	var exhibitions []Exhibition
	q := c.Query("q")
	year := c.Query("year")
	when := c.Query("when")
	where := c.Query("where")
	curator := c.Query("curator")
	artist := c.Query("artist")
	technique := c.Query("technique")
	work := c.Query("work")
	page := c.DefaultQuery("page","1")

	and_year := ""
	if year != "" {
		y, yerr := strconv.ParseInt(year, 10, 32)
		checkErr(yerr)

		cmp := "="
		switch when {
		case "until", "UNTIL":
			cmp = "<="
		case "since", "SINCE":
			cmp = ">="
		}

		and_year = fmt.Sprintf("AND anio %s %d", cmp, y)
	}

	and_where := ""
	if where != "" {
		and_where = fmt.Sprintf("AND nombre_locacion LIKE '%%%s%%'", where)
	}

	and_work := ""
	if work != "" {
		and_work = fmt.Sprintf("AND nombre_obra LIKE '%%%s%%'", work)
	}

	and_technique := ""
	if technique != "" {
		and_technique = fmt.Sprintf("AND tecnica LIKE '%%%s%%'", technique)
	}

	and_curator := ""
	if curator != "" {
		and_curator = fmt.Sprintf("AND curadores  MATCH '*%s*'", curator)
	}

	and_artist := ""
	if artist != "" {
		and_artist = fmt.Sprintf("AND artistas MATCH '*%s*'", artist)
	}

	from_where := fmt.Sprintf(`
FROM exhibiciones
NATURAL JOIN artistas
NATURAL JOIN curadores
WHERE nombre_exhibicion LIKE '%%%s%%'
%s %s %s %s %s %s`, q, and_year, and_where, and_technique, and_work, and_artist, and_curator)

	// TODO FIX ME. sqlite doesn't have a quick way to query for # of rows AND the rows, so this:
	totalPages, perr := howManyPages(fmt.Sprintf("SELECT count(id) %s", from_where))
	if perr == nil {
		query := PaginateQuery(fmt.Sprintf("SELECT * %s", from_where), page)

		_, err := dbmap.Select(&exhibitions, query)
		checkErr(err)

		if err == nil {
			content := newResourceExhibitions(exhibitions, "search", page, strconv.Itoa(totalPages))
			c.JSON(200, content)
		} else {
			jsonErr := &Error{"not_found", 404, "Not Found Error", "Your query didn't match any exhibition."}
			c.JSON(404, jsonErr)
		}
	} else {
		jsonErr := &Error{"server_error", 500, "Server Error", "Something wrong has happened."}
		c.JSON(500, jsonErr)
	}
}
