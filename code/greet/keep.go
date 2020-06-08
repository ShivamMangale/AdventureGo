package greet

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

//Story  A structure for story
type Story map[string]Chapter

//Chapter  A structure for Chapter
type Chapter struct {
	Title     string   `json:"title,omitempty"`
	StoryText []string `json:"story,omitempty"`
	Paths     []Path   `json:"options,omitempty"`
}

//Path  A structure for Path
type Path struct {
	Text    string `json:"text,omitempty"`
	Chapter string `json:"arc,omitempty"`
}

//SourceToJSON  a function that takes in the source file and returns the json format
func SourceToJSON(fileToRead io.Reader) (Story, error) {
	fileDecoder := json.NewDecoder(fileToRead)
	var story Story
	if err := fileDecoder.Decode(&story); err != nil {
		return nil, err
	}

	return story, nil
}

//Options  to structue the possible input from user
type Options func(h *handler)

//ChooseTempl  is a function that will assign the template to be used
//based on what the user wants
func ChooseTempl(t int) Options {
	return func(h *handler) {
		if t == 2 {
			h.t = createTemplate(secondtemplate)
		}
	}
}

//ChoosePathFn  is a function that will assign the appropriate function
//to be used on what the user wants
func ChoosePathFn(fn func(r *http.Request) string) Options {
	return func(h *handler) {
		h.pathFunction = fn
	}
}

//NewHandler  a html handler
func NewHandler(s Story, options ...Options) http.Handler {
	var h handler
	h.s = s
	h.t = sampleTemplate
	h.pathFunction = defaultPathFn

	for _, opt := range options {
		opt(&h)
	}

	return h
}

type handler struct {
	s            Story
	t            *template.Template
	pathFunction func(r *http.Request) string
}

func createTemplate(filename string) *template.Template {
	return template.Must(template.New("").Parse(filename))
}

var sampleTemplate = createTemplate(defaultTmplt)

func defaultPathFn(r *http.Request) string {
	pat := r.URL.Path
	curpath := strings.TrimSpace(pat)
	if curpath == "" || curpath == "/" {
		curpath = "/intro"
	}
	curpath = curpath[1:]

	return curpath
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tpl, err := h.t.Clone()
	if err != nil {
		log.Printf("%v\n", err)
		http.Error(w, "Something went wrong!", http.StatusInternalServerError)
	}

	curpath := h.pathFunction(r)

	if chapter, ok := h.s[curpath]; ok {
		err := tpl.Execute(w, chapter)
		if err != nil {
			log.Printf("%v\n", err)
			http.Error(w, "Something went wrong!", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Chapter Not Found!", http.StatusNotFound)
	}
}

var defaultTmplt = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="X-UA-Compatible" content="ie=edge">
  <title>Adventure Go - {{ .Title }}</title>
  <link href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
  <div class="container">
    <h1>Adventure Go</h1>
    {{template "storyArc" .}}
  </div>
</body>
<style>
  a {
    display: block;
    padding: 0.25em;
    margin: 1em;
  }
</style>
<!-- jquery & Bootstrap JS -->
<script src="//ajax.googleapis.com/ajax/libs/jquery/1.11.3/jquery.min.js">
</script>
<script src="//maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js">
</script>
</html>

{{define "storyArc"}}
  <h2>{{.Title}}</h2>
  {{range .StoryText}}
    <p>{{.}}</p>
  {{end}}
  {{if .Paths}}
    {{range .Paths}}
      <a href="{{.Chapter}}">{{.Text}}</a>
    {{end}}
  {{else}}
    <p>The end</p>
  {{end}}
{{end}}`

var secondtemplate = `
<!DOCTYPE html>
<html>
  <head>
    <title>{{.Title}}</title>
    <link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/materialize/1.0.0-rc.2/css/materialize.min.css">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
  </head>
  <body class="grey lighten-4">
    <div class="row">
        <div class="col m12">
          <div class="card large horizontal z-depth-4 brown white-text">
            <div class="card-stacked">
              <h2 class="card-title">{{.Title}}</h2>
              <div class="card-content">
                {{range .StoryText}}
                  <p>{{.}}</p>
                {{end}}
              </div>
              <div class="card-action">
                {{range .Paths}}
                  <a href="{{.Chapter}}">{{.Text}}</p>
                {{end}}
              </div>
            </div>
          </div>
        </div>
      </div>
  </body>
</html>`
