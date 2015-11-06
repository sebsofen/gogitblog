package handler
import (
	"seb-blog/utils"
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
)

type ByModDate []os.FileInfo
func (a ByModDate) Len() int           { return len(a) }
func (a ByModDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByModDate) Less(i, j int) bool { return a[i].ModTime().Nanosecond() < a[j].ModTime().Nanosecond() }


type PostsHandler struct {
	Config *utils.Configuration
}

type Post struct {
	Post string
}

func NewPostsHandler(config *utils.Configuration) (*PostsHandler) {
	posts := &PostsHandler{
		Config: config,
	}

	return posts
}

func (psts* PostsHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	files, _ := ioutil.ReadDir(psts.Config.Gitfolder)
	sort.Sort(ByModDate(files))
	for _, f := range files {
		fmt.Fprintln(w, string(f.Name()))
	}
}


func (psts* PostsHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postslug := vars["postslug"]
	dat, err := ioutil.ReadFile(psts.Config.Gitfolder + "/" + postslug + "/Post.md")

	if (err != nil){
		panic(err)
	}

	fmt.Print(string(dat))

	post := &Post {
		Post: string(dat),
	}

	//Fprintln will print to webpage
	jsonPost, _ := json.Marshal(post)
	fmt.Fprintln(w, string(jsonPost))
	fmt.Println(psts.Config.Gitfolder)
}