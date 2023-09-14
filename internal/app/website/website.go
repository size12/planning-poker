package website

import "html/template"

type Website struct {
	files *template.Template
	url   string
}

func New(url string) (*Website, error) {
	files, err := template.ParseFiles(htmlFiles...)
	if err != nil {
		return nil, err
	}

	return &Website{files: files, url: url}, err
}
