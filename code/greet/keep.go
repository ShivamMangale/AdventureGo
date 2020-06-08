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

//NewHandler  a html handler
func NewHandler(s Story) http.Handler {
	return handler{s}
}

type handler struct {
	s Story
}

func createTemplate(filename string) *template.Template {
	return template.Must(template.New("").Parse(filename))
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

var sampleTemplate = createTemplate(defaultTmplt)

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tpl, err := sampleTemplate.Clone()
	if err != nil {
		log.Printf("%v\n", err)
		http.Error(w, "Something went wrong!", http.StatusInternalServerError)
	}
	curpath := strings.TrimSpace(r.URL.Path)
	if curpath == "" || curpath == "/" {
		curpath = "/intro"
	}
	curpath = curpath[1:]

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
