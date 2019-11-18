package main

import (
	"flag"
	"fmt"
	"github.com/kataras/iris"
	"io/ioutil"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

const (
	ver = "0.01"
)

var (
	rootDir = flag.String("rootdir", "", "Specify the root directory for share")
	port    = flag.String("port", "", "Specify the port for listening")
	version = flag.Bool("version", false, "Show version")
	help    = flag.Bool("help", false, "Show this help message")
)

type Config struct {
	RootDir string
	Port    string
}

func (c *Config) Read(configFile string) {
	c.RootDir = "~/SharedFiles"
	c.Port = ":3000"
	if _, err := os.Stat(configFile); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Config file doesn't exist.") //Debug
		} else {
			fmt.Println(err)
		}
	} else {
		file, err := ioutil.ReadFile(configFile)
		if err != nil {
			fmt.Println(err)
		} else {
			for _, line := range strings.Split(string(file), "\n") {
				line = strings.Trim(line, " \t")
				ReplaceRept(&line, " ")
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
	if *rootDir != "" {
		c.RootDir = *rootDir
	}
	if *port != "" {
		c.Port = ":" + *port
	}
	ReplaceRept(&c.RootDir, "/")
}

func ReplaceRept(str *string, target string) { //Replace specified Repeated string to Single
	*str = strings.ReplaceAll(*str, target+target, target)
	if strings.Contains(*str, target+target) {
		ReplaceRept(str, target)
	}
}

type File struct {
	Name  string
	Size  string
	Date  string
	URL   string
	IsDir bool
}

type Files []File

func (f Files) Len() int { return len(f) }

func (f Files) Less(i, j int) bool {
	/*length := func() int {
		if len(f[i].Name) > len(f[j].Name) {
			return len(f[j].Name)
		} else {
			return len(f[i].Name)
		}
	}()
	for ii := 0; ii < length; ii++ {
		if f[i].Name[ii] > f[j].Name[ii] {
			return true
		} else if f[i].Name[ii] < f[j].Name[ii] {
			return false
		} else {
			continue
		}
	}*/
	return f[i].Name < f[j].Name
}

func (f Files) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

func (f *Files) Sort() {
	var files, dirs Files
	for _, v := range *f {
		if v.IsDir {
			dirs = append(dirs, v)
		} else {
			files = append(files, v)
		}
	}
	sort.Sort(files)
	sort.Sort(dirs)
	*f = append(dirs, files...)
}

func main() {
	flag.Parse()
	if *version {
		fmt.Println("Software Version:", ver)
		fmt.Println("Go Compiler Version:", strings.ToUpper(runtime.Version()))
		fmt.Println("Arch:", strings.ToUpper(runtime.GOARCH))
		fmt.Println("System:", strings.ToUpper(runtime.GOOS))
		return
	}
	if *help {
		flag.Usage()
		return
	}

	var config Config
	config.Read("./config")

	app := iris.New()
	//	app.Logger().SetLevel("debug")
	app.RegisterView(iris.HTML("./views", ".html").Layout("templates/layout.html"))
	app.HandleDir("/public", "./public")
	app.HandleDir("/", config.RootDir)
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().GetStringDefault("message", "The page you're looking for doesn't exist"))
		ctx.View("universal/404.html")
	})
	app.Favicon("./public/favicon.ico")

	app.Get("/{path:path}", func(ctx iris.Context) {
		reqPath := ctx.Path()
		//	fmt.Println("Request Path:", reqPath) //Debug
		fullPath := config.RootDir + reqPath
		ReplaceRept(&fullPath, "/")
		//	fmt.Println("Full Path:", fullPath) //Debug
		path, err := os.Stat(fullPath)
		//	fmt.Println("Path:", &path) //Debug
		if os.IsNotExist(err) || path == nil {
			ctx.ViewData("Message", "Not Found.")
			ctx.View("universal/404.html")
		}
		if path.IsDir() {
			var (
				file        File
				files       Files
				filesX, err = ioutil.ReadDir(fullPath)
			)
			if err != nil {
				ctx.ViewData("Message", err)
				ctx.View("universal/404.html")
			}
			for _, fileX := range filesX {
				if fileX.IsDir() {
					file.IsDir = true
					file.Name = fileX.Name() + "/"
					file.Size = "-"
					file.Date = fileX.ModTime().Format("2006-01-02 15:04:05")
					/*file.URL = file.Name
					ReplaceRept(&file.URL, "/")*/
					files = append(files, file)
				} else {
					file.IsDir = false
					file.Name = fileX.Name()
					file.Size = strconv.FormatInt(fileX.Size(), 10)
					file.Date = fileX.ModTime().Format("2006-01-02 15:04:05")
					/*file.URL = file.Name
					ReplaceRept(&file.URL, "/")*/
					files = append(files, file)
				}
			}
			files.Sort()
			if reqPath != "/" {
				var (
					updir File
					X     Files
				)
				updir.IsDir = true
				updir.Name = "../"
				updir.Size = "-"
				updir.Date = path.ModTime().Format("2006-01-02 15:04:05")
				X = append(X, updir)
				files = append(X, files...)
			}
			ctx.ViewData("Location", reqPath)
			ctx.ViewData("Date", path.ModTime().Format("2006-01-02 15:04:05"))
			ctx.ViewData("Files", files)
			ctx.View("filelist.html")
		} else if !path.IsDir() {
			//User Download
			ctx.SendFile(fullPath, path.Name())
		}
	})
	fmt.Println("Root Directory:", config.RootDir) //Debug
	app.Run(iris.Addr(config.Port))
}
