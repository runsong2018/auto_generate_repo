package mysql

import (
	"bytes"
	"context"
	_ "embed"
	"github.com/runsong2018/auto_generate_repo/internal/format"
	"github.com/runsong2018/auto_generate_repo/internal/repo"
	"io"
	"strings"
	"text/template"
)

//go:embed mysql.tmpl
var mysqlTmpData string

func Handle(ctx context.Context, repo *repo.Repo, writer io.Writer) (err error) {
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"ToLower":     strings.ToLower,
		"SnakeString": format.SnakeString,
	}).Parse(mysqlTmpData)
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(nil)
	err = tmpl.Execute(buf, repo)
	if err != nil {
		return
	}

	_, err = writer.Write(buf.Bytes())
	if err != nil {
		return
	}

	return
}
