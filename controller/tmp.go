package controller

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/utils"
	"html/template"
	"io"
	"io/ioutil"
	"regexp"

	//"html/template"
	"os"
	"path/filepath"
	"strings"
)

type templateFile struct {
	root  string
	files map[string][]string
}

var beeViewPathTemplates map[string]*template.Template

//　加载模板文件
func init() {

	beeViewPathTemplates = make(map[string]*template.Template)

	viewPath := "templates"

	beeViewPathTemplates = make(map[string]*template.Template)
	BuildBeegoTemplate(viewPath)

}

// visit will make the paths into two part,the first is subDir (without tf.root),the second is full path(without tf.root).
// if tf.root="views" and
// paths is "views/errors/404.html",the subDir will be "errors",the file will be "errors/404.html"
// paths is "views/admin/errors/404.html",the subDir will be "admin/errors",the file will be "admin/errors/404.html"
func (tf *templateFile) visit(paths string, f os.FileInfo, err error) error {
	if f == nil {
		return err
	}
	if f.IsDir() || (f.Mode()&os.ModeSymlink) > 0 {
		return nil
	}

	replace := strings.NewReplacer("\\", "/")
	file := strings.TrimLeft(replace.Replace(paths[len(tf.root):]), "/")
	subDir := filepath.Dir(file)

	tf.files[subDir] = append(tf.files[subDir], file)
	return nil
}

// BuildTemplate will build all template files in a directory.
// it makes beego can render any template file in view directory.
func BuildBeegoTemplate(dir string, files ...string) error {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return errors.New("dir open err")
	}

	self := &templateFile{
		root:  dir,
		files: make(map[string][]string),
	}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		return self.visit(path, f, err)
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
		return err
	}

	for _, v := range self.files {
		for _, file := range v {
			t, _ := getTemplate(self.root, file, v...)
			beeViewPathTemplates[file] = t
		}
	}
	return nil
}

func getTplDeep(root, file, parent string, t *template.Template) (*template.Template, error) {
	var fileAbsPath string
	var rParent string
	if filepath.HasPrefix(file, "../") {
		rParent = filepath.Join(filepath.Dir(parent), file)
		fileAbsPath = filepath.Join(root, filepath.Dir(parent), file)
	} else {
		rParent = file
		fileAbsPath = filepath.Join(root, file)
	}
	if e := utils.FileExists(fileAbsPath); !e {
		panic("can't find template file:" + file)
	}
	data, err := ioutil.ReadFile(fileAbsPath)
	if err != nil {
		return nil, err
	}
	t, err = t.New(file).Parse(string(data))
	if err != nil {
		return nil, err
	}
	reg := regexp.MustCompile("{{" + "[ ]*template[ ]+\"([^\"]+)\"")
	allSub := reg.FindAllStringSubmatch(string(data), -1)
	for _, m := range allSub {
		if len(m) == 2 {
			tl := t.Lookup(m[1])
			if tl != nil {
				continue
			}

			_, err = getTplDeep(root, m[1], rParent, t)
			if err != nil {
				return nil, err
			}
		}
	}
	return t, nil
}

func getTemplate(root, file string, others ...string) (t *template.Template, err error) {
	t = template.New(file)
	t, err = getTplDeep(root, file, "", t)
	if err != nil {
		return nil, err
	}
	return
}

//　执行模板
func ExecuteViewPathTemplate(wr io.Writer, name string, viewPath string, data interface{}) {

	var err error
	if t, ok := beeViewPathTemplates[viewPath]; ok {
		if t.Lookup(name) != nil {
			err = t.ExecuteTemplate(wr, name, data)
		} else {
			err = t.Execute(wr, data)
		}
		if err != nil {
			fmt.Fprintf(wr, "Unknow error： %s", err)
		}
		return
	}
	fmt.Fprintf(wr, "Unknown view path: %s", viewPath)
	return

}
