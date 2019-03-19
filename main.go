package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "0000"
	DB_NAME     = "go_backend"
)

type DataSnap struct {
	Sno string `json:"sno"`
	Ac1 string `json:"ax"`
	Ac2 string `json:"ay"`
	Ac3 string `json:"az"`
	Gy1 string `json:"gx"`
	Gy2 string `json:"gy"`
	Gy3 string `json:"gz"`
	Or1 string `json:"ox"`
	Or2 string `json:"oy"`
	Or3 string `json:"oz"`
	Srn int    `json:"srn"`
}

type DataSnaps struct {
	DataSnaps []DataSnap `json:"data_snaps"`
}

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	} else {
		fmt.Println("DB Connected...")
	}
	defer db.Close()

	e := echo.New()

	e.GET("/test/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	e.GET("/all/", func(c echo.Context) error {
		sqlStatement := "SELECT * FROM data_dump"
		rows, err := db.Query(sqlStatement)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		result := DataSnaps{}
		for rows.Next() {
			snap := DataSnap{}
			err2 := rows.Scan(&snap.Sno, &snap.Ac1, &snap.Ac2, &snap.Ac3, &snap.Gy1, &snap.Gy2, &snap.Gy3, &snap.Or1, &snap.Or2, &snap.Or3, &snap.Srn)
			if err2 != nil {
				return err2
			}
			result.DataSnaps = append(result.DataSnaps, snap)
		}
		return c.JSON(http.StatusCreated, result)
	})

	e.GET("/after/:id/", func(c echo.Context) error {
		id := c.Param("id")
		sqlStatement := "SELECT * FROM data_dump WHERE srn > " + id
		rows, err := db.Query(sqlStatement)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		result := DataSnaps{}
		for rows.Next() {
			snap := DataSnap{}
			err2 := rows.Scan(&snap.Sno, &snap.Ac1, &snap.Ac2, &snap.Ac3, &snap.Gy1, &snap.Gy2, &snap.Gy3, &snap.Or1, &snap.Or2, &snap.Or3, &snap.Srn)
			if err2 != nil {
				return err2
			}
			result.DataSnaps = append(result.DataSnaps, snap)
		}
		return c.JSON(http.StatusCreated, result)
	})

	e.GET("/at/:id/", func(c echo.Context) error {
		id := c.Param("id")
		sqlStatement := "SELECT * FROM data_dump WHERE srn = " + id
		rows, err := db.Query(sqlStatement)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		result := DataSnaps{}
		for rows.Next() {
			snap := DataSnap{}
			err2 := rows.Scan(&snap.Sno, &snap.Ac1, &snap.Ac2, &snap.Ac3, &snap.Gy1, &snap.Gy2, &snap.Gy3, &snap.Or1, &snap.Or2, &snap.Or3, &snap.Srn)
			if err2 != nil {
				return err2
			}
			result.DataSnaps = append(result.DataSnaps, snap)
		}
		return c.JSON(http.StatusCreated, result)
	})

	e.GET("/last/:id/", func(c echo.Context) error {
		id := c.Param("id")
		sqlStatement := "SELECT * FROM data_dump WHERE srn > ((SELECT COUNT(srn) FROM data_dump) - " + id + ")"
		rows, err := db.Query(sqlStatement)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		result := DataSnaps{}
		for rows.Next() {
			snap := DataSnap{}
			err2 := rows.Scan(&snap.Sno, &snap.Ac1, &snap.Ac2, &snap.Ac3, &snap.Gy1, &snap.Gy2, &snap.Gy3, &snap.Or1, &snap.Or2, &snap.Or3, &snap.Srn)
			if err2 != nil {
				return err2
			}
			result.DataSnaps = append(result.DataSnaps, snap)
		}
		return c.JSON(http.StatusCreated, result)
	})

	e.POST("/send/", func(c echo.Context) error {
		u := new(DataSnap)
		if err := c.Bind(u); err != nil {
			return err
		}
		sqlStatement := "INSERT INTO data_dump (sno, ax, ay, az, gx, gy, gz) VALUES ($1, $2, $3, $4, $5, $6, $7)"
		res, err := db.Exec(sqlStatement, u.Sno, u.Ac1, u.Ac2, u.Ac3, u.Gy1, u.Gy2, u.Gy3)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
			return c.JSON(http.StatusCreated, u)
		}
		return c.String(http.StatusOK, "ok")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
