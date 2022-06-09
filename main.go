package main

import (
	"fmt"
	"log"
	"time"

	"rock_example/config"

	"github.com/go-rock/rock"
)

var _ = [][]int{{1}}

func main() {
	foo := "bar"
	fmt.Println(foo)
	app := rock.New()
	app.Use(rock.Recovery())
	app.Use(Logger())
	app.NoRoute(func(c rock.Context) {
		// c.JSON(404, rock.M{"msg": "404 not found"})
		c.Status(404)
		c.HTML("404")
	})

	// app.NoMethod(func(c rock.Context) {
	// 	// c.JSON(404, rock.M{"msg": "404 not found"})
	// 	c.Status(405)
	// 	c.String(405, format string, values ...interface{})
	// })

	config.Setup(app)
	app.Static("/assets", "./static")
	app.Static("/themes", "./themes/assets")
	// base.HTMLRender(render.Default())
	// app.RegisterView(render.New(render.ViewConfig{
	// 	ViewDir:   "./templates/pg2/",
	// 	Extension: ".html",
	// }))
	// app.RegisterView(rock.NewHtmlEngine("html engine"))
	app.Get("/", Home)

	// app.Get("/posts/:id", Post)
	// api := app.Group("/api")
	// api.Use(onlyForApi())
	// {
	// 	api.Get("/home", ApiIndex)
	// 	v1 := api.Group("/v1")
	// 	{
	// 		v1.Get("/home", ApiIndex)
	// 	}
	// }
	admin := app.Group("/admin")
	admin.Use(auth())

	// admin.RegisterView(rock.NewPgEngine("pg engine"))
	// admin.RegisterView(render.New(render.ViewConfig{
	// 	ViewDir:   "./views/",
	// 	Extension: ".html",
	// }))

	admin.NoRoute(func(c rock.Context) {
		c.JSON(404, rock.M{"msg": "404 not found"})
	})

	// app.GetHTMLRender().SetViewDir("./tem/")
	// app.GetHTMLRender().SetViewDir("./template/")
	{
		// app.GetHTMLRender().SetViewDir("./tem/")
		admin.Get("/login", AdminLogin)
		admin.Post("/login", AdminLogin)
	}

	app.Get("/panic", func(c rock.Context) {
		names := []string{"geektutu"}
		c.String(200, names[100])
	})

	err := app.Run()
	if err != nil {
		panic(err)
	}
}

func Post(c rock.Context) {
	log.Printf("query from is %s %d", c.Query("from"), c.QueryInt("cid"))
	c.String(200, "post id is %s", c.Param("id"))
}

func Home(c rock.Context) {
	// c.JSON(200, rock.H{"msg": "ok"})
	log.Println(c.GetView(), "HOME")
	c.HTML("home")
}

type Error struct {
	Name string `json:"name"`
	Msg  string `json:"msg"`
}

// admin
func AdminLogin(c rock.Context) {
	log.Println("admin auth action")
	log.Println(c.GetView().Engine.Ext())

	error := &Error{Name: "render error", Msg: "Error msg"}
	// error := rock.M{"Msg": "xiao"}
	// c.JSON(http.StatusOK, rock.H{"msg": "admin login"})
	// c.Status(422)
	c.HTML("admin/login", rock.M{"data": error})
}

// Api
func ApiIndex(c rock.Context) {
	c.JSON(200, rock.H{"msg": "api v1 index"})
}

// middlewares
func onlyForApi() rock.HandlerFunc {
	return func(c rock.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("Api only code [%d] %s in %v for group api", c.StatusCode(), c.Request().RequestURI, time.Since(t))
	}
}

func auth() rock.HandlerFunc {
	return func(c rock.Context) {
		log.Println("auth before")
		c.Next()
		log.Println("auth after")
	}
}

func Logger() rock.HandlerFunc {
	return func(c rock.Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s %s in %v", c.StatusCode(), c.Request().Method, c.Request().RequestURI, time.Since(t))
	}
}
