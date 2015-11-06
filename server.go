package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"log"
	"seb-blog/utils"
	"os/exec"
	"bytes"
	"time"
	"seb-blog/handler"
)

func main() {
	cnfg := utils.NewConfiguration("config.json")
	ticker := time.NewTicker(time.Second * cnfg.Updateinterval)
	//update repository automatically
	go UpdateRepo(cnfg.Gitfolder, ticker)

	posts := handler.NewPostsHandler(cnfg)

	r := mux.NewRouter()
	r.HandleFunc("/posts/{postslug}", posts.GetPost)

	r.HandleFunc("/listposts/{offset}/{limit}", posts.ListPosts)


	//provide static content ;-)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(cnfg.Htmlfiles)))
	http.Handle("/", r)

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