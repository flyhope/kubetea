package view

import (
	"html/template"
	"time"
)

// 模板函数定义
var defaultTemplate = template.New("kubetea").Funcs(template.FuncMap{
	"PodPhaseView": PodPhaseView,
	"PodReadyView": PodReadyView,
	"FormatTime":   FormatTime,
})

// FormatTime 格式化模板时间
func FormatTime(t time.Time) string {
	return t.Format(time.DateTime)
}

// NewTmplParse 创建模板
func NewTmplParse(text string) (*template.Template, error) {
	newTmpl, err := defaultTemplate.Clone()
	if err != nil {
		return nil, err
	}
	return newTmpl.Parse(text)
}

// NewTmplParseSlice 创建一组模板
func NewTmplParseSlice(texts []string) ([]*template.Template, error) {
	templates := make([]*template.Template, 0, len(texts))
	for _, text := range texts {
		newTmpl, err := NewTmplParse(text)
		if err != nil {
			return nil, err
		}
		templates = append(templates, newTmpl)
	}
	return templates, nil
}
