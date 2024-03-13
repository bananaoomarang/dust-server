package main

import (
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

type PaginationQuery struct {
	Offset int `form:"offset"`
}

func main() {
	ConnectDB()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	router.GET("/levels", levelsGET)
	router.GET("/levels/:id", levelGET)
	router.POST("/levels", levelPOST)

	router.Run("localhost:8080")
}

func levelsGET(c *gin.Context) {
	result := []LevelRes{}

	var pagination_query PaginationQuery

	if err := c.BindQuery(&pagination_query); err != nil {
		c.AbortWithError(500, err)
		return
	}

	q := `SELECT id,name,data FROM levels ORDER BY id ASC LIMIT 10 OFFSET $1`
	rows, err := Db.Query(q, pagination_query.Offset)

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
		result = append(result, level)
	}

	if err = rows.Err(); err != nil {
		c.AbortWithError(500, err)
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
