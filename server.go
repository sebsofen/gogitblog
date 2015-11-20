package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"log"
	"github.com/sebsofen/gogitblog/utils"
	"os/exec"
	"bytes"
	"time"
	"github.com/sebsofen/gogitblog/handler"

	"io/ioutil"
)

func main() {
	cnfg := utils.NewConfiguration("config.json")
	fmt.Println("cnfg" + cnfg.Htmlfiles)

	if cnfg.Updateinterval > 0  {
		ticker := time.NewTicker(time.Second * cnfg.Updateinterval)
		//update repository automatically, triggered by update interval
		go UpdateRepo(cnfg.Gitfolder, ticker)
	}

	r := mux.NewRouter()
	handler.NewPostsHandler(cnfg,r)


	//provide more handlers for static content!
	//for all files in posts folders, provide own fileserver
	files, _ := ioutil.ReadDir(cnfg.Gitfolder)
	for _, f := range files {
		//filter folders only
		//TODO IMPLEMENT MORE FILTERS (maybe)
		if (f.IsDir() && f.Name() != ".git") {
			r.PathPrefix("/" + f.Name() + "/").Handler(
				http.StripPrefix("/" + f.Name() + "/",
					http.FileServer(http.Dir(cnfg.Gitfolder + "/" + f.Name() + "/data/"))))
			fmt.Println(cnfg.Gitfolder + "/" + f.Name() + "/data/")
		}
	}

	//provide static website content
	r.PathPrefix("/").Handler(http.StripPrefix("/",http.FileServer(http.Dir(cnfg.Htmlfiles + "/"))))




	//finally, add handler to http listener
	//set custom profex if desired
	http.Handle("/*", r)

	log.Fatal(http.ListenAndServe(":" + cnfg.Port, r))
}



func UpdateRepo(folder string, ticker *time.Ticker) {
	for t := range ticker.C {
		fmt.Println("Tick at", t)
		cmd := exec.Command("git", "pull", "-u", "origin", "master")
		cmd.Dir = folder
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(out.String())

	}
}