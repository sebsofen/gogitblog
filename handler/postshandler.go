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
	"time"
	"strconv"
)

type ByModDate []os.FileInfo
func (a ByModDate) Len() int           { return len(a) }
func (a ByModDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByModDate) Less(i, j int) bool { return a[i].ModTime().Nanosecond() < a[j].ModTime().Nanosecond() }


type PostsHandler struct {
	Config *utils.Configuration
}

type PostMetadata struct {
	Title string
}

type Post struct {
	Post string
	Title string
	Slug string
	Date time.Time
	Metadata PostMetadata
}

func NewPostsHandler(config *utils.Configuration, r *mux.Router) (*PostsHandler) {


	posts := &PostsHandler{
		Config: config,
	}

	r.HandleFunc("/posts/{postslug}", posts.GetPost)
	r.HandleFunc("/listposts/{offset:[0-9]*}/{limit:[0-9]*}", posts.ListPosts)
	r.HandleFunc("/listposts/count", posts.TotalPosts)

	return posts
}

func (psts* PostsHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	files, _ := ioutil.ReadDir(psts.Config.Gitfolder)
	sort.Sort(ByModDate(files))
	vars := mux.Vars(r)

	offset, _ := strconv.Atoi(vars["offset"])
	if(offset < 0){
		offset = 0
	}

	limit, _ := strconv.Atoi(vars["limit"])
	if limit > 10 {
		limit = 10
	}

	postList := make([]*Post, 0, limit)


	for _, f := range files {

		//filter folders only
		//TODO IMPLEMENT MORE FILTERS (maybe)
		if (f.IsDir() && f.Name() != ".git" ) {
			if(offset == 0 && len(postList) < limit) {
				metadata_json, _ := ioutil.ReadFile(psts.Config.Gitfolder + "/" + f.Name() + "/metadata.json")
				metadata := &PostMetadata{}
				json.Unmarshal(metadata_json, metadata)

				post_md_b, err := ioutil.ReadFile(psts.Config.Gitfolder + "/" + f.Name() + "/Post.md")
				post_md := ""
				if(err == nil){
					post_md = string(post_md_b)
				}
				fmt.Println(post_md)




				postList = append(postList, &Post{
					Post : post_md,
					Date : f.ModTime(),
					Slug : f.Name(),
					Metadata: *metadata,
				})
			}else{
				offset -= 1;
			}
		}
	}

	json, _ := json.Marshal(postList)
	fmt.Fprintln(w, string(json))
}



func (psts* PostsHandler) TotalPosts(w http.ResponseWriter, r *http.Request) {
	files, _ := ioutil.ReadDir(psts.Config.Gitfolder)
	counter := 0;
	for _, f := range files {
		//filter folders only
		//TODO IMPLEMENT MORE FILTERS (maybe)
		if (f.IsDir() && f.Name() != ".git") {
			counter++
		}
	}

	fmt.Fprintln(w,"{'count':" + string(counter) + "}")
}








func (psts* PostsHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postslug := vars["postslug"]
	post_md, err := ioutil.ReadFile(psts.Config.Gitfolder + "/" + postslug + "/Post.md")
	metadata_json, err := ioutil.ReadFile(psts.Config.Gitfolder + "/" + postslug + "/metadata.json")
	metadata := &PostMetadata{}
	json.Unmarshal(metadata_json, metadata)
	if (err != nil){
		panic(err)
	}

	fmt.Print(string(post_md))

	post := &Post {
		Post: string(post_md),
		Metadata: *metadata,
	}

	//Fprintln will print to webpage
	jsonPost, _ := json.Marshal(post)
	fmt.Fprintln(w, string(jsonPost))
	fmt.Println(psts.Config.Gitfolder)
}