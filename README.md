# JSON API for Dataset Arte Electrónico Argentino

A very simple RESTful JSON API to query the dataset curated by [@cristianReynaga](https://github.com/cristianReynaga), you can read more about it [here](https://github.com/lvm/Dataset-Arte-Electronico-Argentino) (in spanish).  
This API is built using Golang and SQLite3. I also provide a sanitised csv file ready to be imported to a sqlite db (with the necessary statements to ease the process).  
  
This is a WIP with the intention to have a quick way to access this information and to keep learning golang :-)

## How to ...?

As usual

```bash
$ git clone https://github.com/lvm/api-dataset-arte-electronico-argentino
$ cd api-dataset-arte-electronico-argentino/
```

Then, create the SQLite database  

```bash
$ cd data
$ cat csv-to-sqlite | sqlite3 ../db/electronicArtArgentina.sqlite3
```

For testing purpose, to build the API (I recommend using the [official golang image](https://hub.docker.com/_/golang/)) you need to execute

```bash
$ export GOPATH=$HOME
$ go get .
$ go run api.go
```

If everything went OK, you should see something like

```
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /api/v1/exhibitions       --> main.GetExhibitions (4 handlers)
[GIN-debug] GET    /api/v1/exhibitions/:e_id --> main.GetExhibition (4 handlers)
[GIN-debug] GET    /api/v1/search            --> main.SearchExhibitions (4 handlers)
[GIN-debug] OPTIONS /api/v1/exhibitions       --> main.EndpointsOptions (4 handlers)
[GIN-debug] OPTIONS /api/v1/exhibitions/:e_id --> main.EndpointsOptions (4 handlers)
[GIN-debug] OPTIONS /api/v1/search            --> main.EndpointsOptions (4 handlers)
[GIN-debug] Listening and serving HTTP on :8080
```

and you can start querying the API.

### Entrypoints

For the moment, there are two endpoints: `exhibitions` and `search`. It accepts only the `GET` and `OPTION` verbs.

#### `exhibitions`

This method allows to query all of them or just one.  

**List all (around 15xx) exhibitions**
```
GET /api/v1/exhibitions
```

**Just one by passing the ID**
```
GET /api/v1/exhibitions/23
```

#### `search`

This method allows to search, primarily by the exhibition's name and then you can pass more parameters to narrow your search. These extra parameters use the `AND` operator when building the query.  

| param       | values                                                                                               |
| ----------- | ---------------------------------------------------------------------------------------------------- |
| `q`         | The full name or a part of the name of an exhibition. Can be any `string`                            |
| `year`      | The *valid* format would be `YYYY`, but you can pass any `int (32)`                                  |
| `when`      | Only takes `since` and `until`. They work as `greater-than-equal` and `less-than-equal` respectively |
| `technique` | The name of a technique, can be any `string`.                                                        |
| `curator`   | The name of a curator, can be any `string`.                                                          |
| `artist`    | The name of an artist, can be any `string`.                                                          |
| `work`      | The name of the work, can be any `string`.                                                           |


**Examples**

**List all exhibitions with `festival` on its name**
```
GET /api/v1/search?q=festival
```

**...during 1998**
```
GET /api/v1/search?q=festival&year=1998
```

**...until 1998**
```
GET /api/v1/search?q=festival&year=1998&when=until
```

**...since 1998**
```
GET /api/v1/search?q=festival&year=1998&when=since
```

**...with technique `videoart`**
```
GET /api/v1/search?q=festival&year=1998&when=since&technique=videoart
```

**...curated by Graciela Taquini**
```
GET /api/v1/search?q=festival&year=1998&when=since&technique=videoart&curator=taquini
```

**...artist Joanna Rytel**
```
GET /api/v1/search?q=festival&year=1998&when=since&technique=videoart&curator=taquini&artist=rytel
```

**...with the word `sheep` on its name**
```
GET /api/v1/search?q=festival&year=1998&when=since&technique=videoart&curator=taquini&artist=rytel&work=sheep
```


## TODO

* Pagination
* Search by person (artists or curators)
* Dockerfile and config files
* Make the API public
* Home with details about the API

## LICENSE

See [LICENSE](LICENSE)
