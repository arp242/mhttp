package ztpl

import (
	"bytes"
	"html/template"
	"io"
	"path/filepath"
	"strings"

	"zgo.at/zhttp/ztpl/internal"
	"zgo.at/zhttp/ztpl/tplfunc"
	"zgo.at/zlog"
)

// Init sets up the templates.
func Init(path string, pack map[string][]byte) {
	if pack == nil {
		internal.Templates.Set(path, nil)
		Reload()
		return
	}

	t := New()
	for k, v := range pack {
		k = strings.Trim(strings.TrimPrefix(k, path), "/")
		t = template.Must(t.New(k).Parse(string(v)))
	}
	internal.Templates.Set(path, t)
}

// IsLoaded reports if templates have been loaded.
func IsLoaded() bool {
	return internal.Templates != nil
}

// New creates a new empty template instance.
func New() *template.Template {
	return template.New("").Option("missingkey=error").Funcs(tplfunc.FuncMap)
}

// Reload the templates from the filesystem; errors are logged but not fatal!
func Reload() {
	hp := internal.Templates.Path + "/*.gohtml"
	html, err := filepath.Glob(hp)
	if err != nil {
		zlog.Printf("ztpl.Reload: reading templates from %q: %s", hp, err)
	}
	tp := internal.Templates.Path + "/*.gotxt"
	txt, err := filepath.Glob(tp)
	if err != nil {
		zlog.Printf("ztpl.Reload: reading templates from %q: %s", tp, err)
	}

	t, err := New().ParseFiles(append(html, txt...)...)
	if err != nil {
		zlog.Errorf("ztpl.Reload: parsing files: %s (from: %q and %q)", err, hp, tp)
	}
	internal.Templates.Set(internal.Templates.Path, t)
}

// Execute a named template.
func Execute(w io.Writer, name string, data interface{}) error {
	return internal.Templates.ExecuteTemplate(w, name, data)
}

// ExecuteBytes a named template and return the data as bytes.
func ExecuteBytes(name string, data interface{}) ([]byte, error) {
	w := new(bytes.Buffer)
	err := internal.Templates.ExecuteTemplate(w, name, data)
	return w.Bytes(), err
}

// ExecuteString a named template and return the data as bytes.
func ExecuteString(name string, data interface{}) (string, error) {
	w := new(bytes.Buffer)
	err := internal.Templates.ExecuteTemplate(w, name, data)
	return w.String(), err
}
