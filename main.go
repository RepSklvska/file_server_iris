package main

import (
	"flag"
	"fmt"
	"github.com/kataras/iris"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var (
	rootDir = flag.String("rootdir", "", "Specify the root directory for share")
	port    = flag.String("port", "", "Specify the port for listening")
)

type Config struct {
	RootDir string "~/SharedFiles"
	Port    string ":3000"
}

func (c *Config) Read(configFile string) {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, line := range strings.Split(string(file), "\n") {
			line = strings.Trim(line, " \t")
			if len(line) == 0 || line[0] == 35 {
				continue
			}
			x := strings.Split(line, " ")
			if len(x) < 2 || x[1] == "" {
				continue
			}
			switch strings.ToLower(x[0]) {
			case "rootdir":
				c.RootDir = x[1]
			case "port":
				c.Port = ":" + x[1]
			default:
				continue
			}
		}
	}
}

type File struct {
	Name  string
	Size  string
	Date  string
	URL   string
	IsDir bool
}

func main() {
	var config Config
	config.Read("./config")
	if *rootDir != "" {
		config.RootDir = *rootDir
	}
	if *port != "" {
		config.Port = ":" + *port
	}

	app := iris.New()
	app.Logger().SetLevel("debug")
	app.RegisterView(iris.HTML("./views", ".html").Layout("templates/layout.html"))
	app.HandleDir("/public", "./public")
	app.HandleDir("/", config.RootDir)
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().GetStringDefault("message", "The page you're looking for doesn't exist"))
		ctx.View("universal/404.html")
	})
	app.Favicon("./public/favicon.ico")

	app.Get("/{path:alphabetical}", func(ctx iris.Context) {
		reqPath := ctx.Path()
		fmt.Println("Request Path:", reqPath) //Debug
		if strings.Index(reqPath, "..") != -1 {
			ctx.ViewData("Message", "Access Denied.")
			ctx.View("universal/404.html")
		}
		fullPath := config.RootDir + reqPath + "/"
		fullPath = strings.ReplaceAll(fullPath, "//", "/")
		fmt.Println("Full Path:", fullPath) //Debug
		path, err := os.Stat(fullPath)
		fmt.Println("Path:", path) //Debug
		if os.IsNotExist(err) || path == nil {
			ctx.ViewData("Message", "Not Found.")
			ctx.View("universal/404.html")
		}
		if path.IsDir() {
			var (
				file        File
				files       []File
				filesX, err = ioutil.ReadDir(fullPath)
			)
			if err != nil {
				ctx.ViewData("Message", err)
				ctx.View("universal/404.html")
			}
			for _, fileX := range filesX {
				if fileX.IsDir() {

				} else {
					file.Name = fileX.Name()
					file.Size = strconv.FormatInt(fileX.Size(), 10)
					file.Date = fileX.ModTime().Format("2010-01-20 13:03:02")
					files = append(files, file)
				}
			}
			//Sort needed
			ctx.ViewData("Files", files)
			ctx.View("filelist.html")
		}
	})

	app.Run(iris.Addr(config.Port))
}
