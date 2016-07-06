package main

import (
	"os"
	"fmt"
	"log"
	"path"
	"time"
	"strconv"
 	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/aviddiviner/gin-limit"
	"github.com/gin-gonic/gin"
)

var (
	tmplDirectory string = "./tmpl"
	dbFilename    string = "./db/electronicArtArgentina.sqlite3"
	tableName     string = "exhibiciones"
	httpPort      string = "8080"
	apiHostname   string = "http://arte-electronico.cyberpunk.com.ar"
	apiPath       string = "/api"
	perPage       int = 20
	firstPage     string = "1"
	lastPage      string = "80" // (1593 / perPage) // TODO FIX ME

)

var dbmap = initDb()

type (
	Exhibition struct {
		Id         int64 `db:"id" json:"-"`
		Anio       int `db:"anio" json:"year"`
		Nombre     string `db:"nombre_exhibicion" json:"exhibition"`
		Locacion   string `db:"nombre_locacion" json:"location"`
		Curador0   string `db:"curador_0" json:"curator_0,omitempty"`
		Curador1   string `db:"curador_1" json:"curator_1,omitempty"`
		Curador2   string `db:"curador_2" json:"curator_2,omitempty"`
		Curador3   string `db:"curador_3" json:"curator_3,omitempty"`
		FechaIncio time.Time `db:"fecha_inicio" json:"start"`
		FechaFin   time.Time `db:"fecha_finalizacion" json:"end"`
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

	PaginationLinks struct {
		First      string `json:"first"`
		Last       string `json:"last"`
		Prev       string `json:"prev"`
		Next       string `json:"next"`
	}

	Meta struct {
		TotalPages string `json:"total_pages"`
	}

	ResourceObject struct {
		Id         int64 `json:"id"`
		Type       string `json:"type"`
		Attributes Exhibition `json:"attributes"`
	}

	ResourceObjects struct {
		Data  []ResourceObject `json:"data"`
		Links PaginationLinks `json:"links"`
		Meta  Meta `json:"meta"`
	}

	Error struct {
		Id     string `json:"id"`
		Status int `json:"status"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
	}

)

func paginationLink(endpoint string, pg string, which string) string {
	page, err := strconv.ParseInt(pg, 10, 32)
	checkErr(err)
	first, ferr := strconv.ParseInt(firstPage, 10, 32)
	checkErr(ferr)
	last, lerr := strconv.ParseInt(lastPage, 10, 32)
	checkErr(lerr)


	switch which {
	case "prev":
		if page > first {
			page = page - 1
		}
	case "next":
		if page <= last {
			page = page + 1
		}
	}


	link := fmt.Sprintf("%s/%s?page=%d&perPage=%d", apiPath, endpoint, page, perPage)
	if endpoint == "search" && which == "last" { // TODO ugly
		link = ""
	}

	return link
}

func newResourceObject(attr Exhibition) ResourceObject {
	return ResourceObject {
		Id: attr.Id,
		Type: "exhibitions",
		Attributes: attr,
	}
}

func newResourceObjects (exhibitions []Exhibition, endpoint string, pg string) ResourceObjects {
	resources := make([]ResourceObject, 0)
	for e := range exhibitions {
		resources = append(resources, newResourceObject(exhibitions[e]))
	}

	links := PaginationLinks {
		First: paginationLink(endpoint, firstPage, "first"),
		Last: paginationLink(endpoint, lastPage, "last"),
		Prev: paginationLink(endpoint, pg, "prev"),
		Next: paginationLink(endpoint, pg, "next"),
	}

	return ResourceObjects{Data: resources, Links: links, Meta: Meta{TotalPages: "80"}}
}

func pageOffset(pg string) int {
	page, err := strconv.Atoi(pg)
	if err != nil {
		page = 1
	}

	return perPage * (page - 1)
}

func PaginateQuery(query, pg string) string {
	return fmt.Sprintf("%s LIMIT %d OFFSET %d", query, perPage, pageOffset(pg))
}

func init() {
	p := os.Getenv("HTTP_PORT")
	if p != "" {
		httpPort = p
	}

	d := os.Getenv("DB_FILENAME")
	if d != "" {
		dbFilename = d
	}

	t:= os.Getenv("TMPL_DIR")
	if t != "" {
		tmplDirectory = t
	}
}

func main() {
	router := gin.Default()
	router.Use(ApiCors())
	router.Use(limit.MaxAllowed(5))

	// API v1
	v1 := router.Group("api/v1")
	{
		v1.GET("/", ListEndpoints)  // list all exhibitions
		v1.GET("/exhibitions", GetExhibitions)  // list all exhibitions
		v1.GET("/exhibitions/:e_id", GetExhibition)  // show an exhibition
		v1.GET("/search", SearchExhibitions)  // search exhibitions

		v1.OPTIONS("/", EndpointsOptions) // http options
		v1.OPTIONS("/exhibitions", EndpointsOptions) // http options
		v1.OPTIONS("/exhibitions/:e_id", EndpointsOptions) // http options
		v1.OPTIONS("/search", EndpointsOptions) // http options
	}

	// API -- TODO improve this
	current := router.Group("api")
	{
		current.GET("/", ListEndpoints)  // list all exhibitions
		current.GET("/exhibitions", GetExhibitions)  // list all exhibitions
		current.GET("/exhibitions/:e_id", GetExhibition)  // show an exhibition
		current.GET("/search", SearchExhibitions)  // search exhibitions

		current.OPTIONS("/", EndpointsOptions) // http options
		current.OPTIONS("/exhibitions", EndpointsOptions) // http options
		current.OPTIONS("/exhibitions/:e_id", EndpointsOptions) // http options
		current.OPTIONS("/search", EndpointsOptions) // http options
	}

	// Home
	router.LoadHTMLGlob(path.Join(tmplDirectory, "*"))
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.tmpl", gin.H{
			"title": "JSON API for Dataset Arte Electrónico Argentino",
			"version": "v1",
		})
	})

	router.Run(fmt.Sprintf(":%s", httpPort))
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func initDb() *gorp.DbMap {
	db, err := sql.Open("sqlite3", dbFilename)
	checkErr(err)
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbmap.AddTableWithName(Exhibition{}, tableName).SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err)
	return dbmap
}

func EndpointsOptions(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Next()
}

func ApiCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func ListEndpoints(c *gin.Context) {
	api_url := fmt.Sprintf("%s%s/", apiHostname, apiPath)
	c.JSON(200, gin.H{
		"exhibitions_url": fmt.Sprintf("%s%s{/exhibition_id}{?page}", api_url, "exhibitions"),
		"search_url": fmt.Sprintf("%s%s?q={query}{&year,when,technique,curator,artist,work,page}", api_url, "search"),
	})
}

// http http://localhost:8080/api/v1/exhibitions
func GetExhibitions(c *gin.Context) {
	var exhibitions []Exhibition
	page := c.DefaultQuery("page", "1")
	query := PaginateQuery("SELECT * FROM exhibiciones ORDER BY id", page)
	_, err := dbmap.Select(&exhibitions, query)
	checkErr(err)

	if err == nil {
		content := newResourceObjects(exhibitions, "exhibitions", page)
		c.JSON(200, content)
	} else {
		jsonErr := &Error{"not_found", 404, "Not Found Error", "No exhibition found."}
		c.JSON(404, jsonErr)
	}
}

// http http://localhost:8080/api/v1/exhibitions/1
func GetExhibition(c *gin.Context) {
	id := c.Params.ByName("e_id")
	var exhibition Exhibition
	err := dbmap.SelectOne(&exhibition, "SELECT * FROM exhibiciones WHERE id=? LIMIT 1", id)
	checkErr(err)

	if err == nil {
		content := newResourceObject(exhibition)
		c.JSON(200, content)
	} else {
		jsonErr := &Error{"not_found", 404, "Not Found Error", "Your query didn't match any exhibition."}
		c.JSON(404, jsonErr)
	}
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
	page := c.DefaultQuery("page", "1")

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


	query := fmt.Sprintf(`SELECT *
FROM exhibiciones
NATURAL JOIN artistas
NATURAL JOIN curadores
WHERE nombre_exhibicion LIKE '%%%s%%'
%s %s %s %s %s %s`, q,
		and_year, and_where, and_technique, and_work, and_artist, and_curator)

	query = PaginateQuery(query, page)
	_, err := dbmap.Select(&exhibitions, query)
	checkErr(err)

	if err == nil {
		content := newResourceObjects(exhibitions, "search", page)
		c.JSON(200, content)
	} else {
		jsonErr := &Error{"not_found", 404, "Not Found Error", "Your query didn't match any exhibition."}
		c.JSON(404, jsonErr)
	}
}
