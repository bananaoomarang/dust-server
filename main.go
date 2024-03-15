package main

import (
	"database/sql"
	"strings"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type LevelReq struct {
	Name string `json:"name" binding:"required"`
	Data string `json:"data" binding:"required"`
}

type LevelCreated struct {
	Id int `json:"id" binding:"required"`
}

type LevelRes struct {
	Id int `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
	Data string `json:"data"`
}

type SearchResults struct {
	Results []LevelRes `json:"results"`
	HasMore bool `json:"hasMore"`
}

type SearchQuery struct {
	Name string `form:"name"`
	Limit int `form:"limit,default=10"`
	Offset int `form:"offset,default=0"`
}

func main() {
	ConnectDB()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	router.GET("/", rootGET)
	router.GET("/levels", levelsGET)
	router.GET("/levels/:id", levelGET)
	router.POST("/levels", levelPOST)

	router.Run("localhost:8080")
}

func queryDb(query SearchQuery) (*sql.Rows, error) {
	if (query.Name != "") {
		q := `
SELECT id,
	name,
	data
FROM levels

WHERE ts @@ to_tsquery('english', $1)

LIMIT $2
OFFSET $3
`
		return Db.Query(q, strings.ReplaceAll(query.Name, " ", "+") + ":*", query.Limit + 1, query.Offset)
	} else {
		q := `
SELECT id,
	name,
	data
FROM levels

ORDER BY id DESC
LIMIT $1
OFFSET $2
`
		return Db.Query(q, query.Limit + 1, query.Offset)
	}
}

func rootGET(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"api": "hello",
	})
}

func levelsGET(c *gin.Context) {
	result := SearchResults{
		Results: []LevelRes{},
		HasMore: false,
	}

	var query SearchQuery

	if err := c.BindQuery(&query); err != nil {
		c.AbortWithError(500, err)
		return
	}

	rows, err := queryDb(query)

	if err != nil{
		c.AbortWithError(500, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var level LevelRes
        if err := rows.Scan(&level.Id, &level.Name, &level.Data); err != nil {
			c.AbortWithError(500, err)
        }
		result.Results = append(result.Results, level)
	}

	if err = rows.Err(); err != nil {
		c.AbortWithError(500, err)
	}

	if len(result.Results) > query.Limit {
		result.HasMore = true
		result.Results = result.Results[0:query.Limit]
	}

	c.JSON(http.StatusOK, result)
}


func levelGET(c *gin.Context) {
	var result LevelRes

	levelId := c.Param("id")

	q := `SELECT id,name,data FROM levels WHERE id=$1`
	if err := Db.QueryRow(q, levelId).Scan(&result.Id, &result.Name, &result.Data); err != nil {
		c.AbortWithError(500, err)
	}
	c.JSON(http.StatusOK, result)
}

func levelPOST(c *gin.Context) {
	var newLevel LevelReq
	var result LevelCreated

	if err := c.ShouldBindJSON(&newLevel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	q := `INSERT INTO levels(name, data) VALUES($1,$2) RETURNING id`
	if err := Db.QueryRow(q, newLevel.Name, newLevel.Data).Scan(&result.Id); err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(http.StatusCreated, result)
}
