package main

import (
	"os"
	"fmt"
	"log"
	//"net/http"
	//"path"
	"strconv"
 	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/aviddiviner/gin-limit"
	"github.com/gin-gonic/gin"
)

var (
	//tmplDirectory string = "./tmpl"
	dbFilename    string = "./db/electronicArtArgentina.sqlite3"
	tableName     string = "exhibiciones"
	httpPort      string = "8080"
	apiHostname   string = "http://api.arte-electronico.cyberpunk.com.ar"
	apiPath       string = ""
	perPage       int = 20
	firstPage     string = "1"
	lastPage      string = "80" // (1593 / perPage) // TODO FIX ME

)

var dbmap = initDb()

type (
	PaginationLinks struct {
		First      string `json:"first"`
		Last       string `json:"last"`
		Prev       string `json:"prev"`
		Next       string `json:"next"`
	}

	Meta struct {
		TotalPages string `json:"total_pages"`
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

func pageOffset(pg string) int {
	page, err := strconv.Atoi(pg)
	if err != nil {
		page = 1
	}

	return perPage * (page - 1)
}

func groupByQuery(query, groupBy string) string {
	field := ""
	order_by := "id"
	switch groupBy {
	case "event":
		field = "nombre_exhibicion"
	case "location":
		field = "nombre_locacion"
		order_by = "nombre_locacion"
	case "year":
		field = "anio"
		order_by = "anio"
	case "technique":
		field = "tecnica"
		order_by = "tecnica"
	}

	if field != "" {
		query = fmt.Sprintf(query, field, field)
		return fmt.Sprintf("%s GROUP BY %s ORDER BY %s", query, field, order_by)
	} else {
		return query
	}
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

	// t:= os.Getenv("TMPL_DIR")
	// if t != "" {
	// 	tmplDirectory = t
	// }
}

func main() {
	router := gin.Default()
	router.Use(ApiCors())
	router.Use(limit.MaxAllowed(5))

	// API v1
	v1 := router.Group("v1")
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

	// API v2
	v2 := router.Group("v2")
	{
		v2.GET("/", ListEndpoints)  // list all exhibitions
		v2.GET("/exhibitions", GetExhibitions)  // list all exhibitions
		v2.GET("/exhibitions/:e_id", GetExhibition)  // show an exhibition
		v2.GET("/search", SearchExhibitions)  // search exhibitions
		v2.GET("/locations", GetLocations)  // list all exhibitions
		v2.GET("/events", GetEvents)  // list all exhibitions
		v2.GET("/techniques", GetTechniques)  // list all exhibitions

		v2.OPTIONS("/", EndpointsOptions) // http options
		v2.OPTIONS("/exhibitions", EndpointsOptions) // http options
		v2.OPTIONS("/exhibitions/:e_id", EndpointsOptions) // http options
		v2.OPTIONS("/search", EndpointsOptions) // http options
		v2.OPTIONS("/locations", EndpointsOptions) // http options
		v2.OPTIONS("/events", EndpointsOptions) // http options
		v2.OPTIONS("/techniques", EndpointsOptions) // http options
	}

	// API -- TODO improve this
	current := router.Group("/")
	{
		current.GET("/", ListEndpoints)  // list all exhibitions
		current.GET("/exhibitions", GetExhibitions)  // list all exhibitions
		current.GET("/exhibitions/:e_id", GetExhibition)  // show an exhibition
		current.GET("/search", SearchExhibitions)  // search exhibitions
		current.GET("/locations", GetLocations)  // list all exhibitions
		current.GET("/events", GetEvents)  // list all exhibitions
		current.GET("/techniques", GetTechniques)  // list all exhibitions

		current.OPTIONS("/", EndpointsOptions) // http options
		current.OPTIONS("/exhibitions", EndpointsOptions) // http options
		current.OPTIONS("/exhibitions/:e_id", EndpointsOptions) // http options
		current.OPTIONS("/search", EndpointsOptions) // http options
		current.OPTIONS("/locations", EndpointsOptions) // http options
		current.OPTIONS("/events", EndpointsOptions) // http options
		current.OPTIONS("/techniques", EndpointsOptions) // http options
	}

	// Home
	/*
	router.LoadHTMLGlob(path.Join(tmplDirectory, "*"))
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.tmpl", gin.H{
			"title": "JSON API for Dataset Arte Electrónico Argentino",
			"version": "v2",
		})
	})
	*/
	// router.GET("/", func(c *gin.Context) {})

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
		"techniques_url": fmt.Sprintf("%s%s{?page}", api_url, "techniques"),
		"events_url": fmt.Sprintf("%s%s{?page}", api_url, "events"),
		"locations_url": fmt.Sprintf("%s%s{?page}", api_url, "locations"),
		"exhibitions_url": fmt.Sprintf("%s%s{/exhibition_id}{?page}", api_url, "exhibitions"),
		"search_url": fmt.Sprintf("%s%s?q={query}{&year,when,technique,curator,artist,work,page}", api_url, "search"),
	})
}
