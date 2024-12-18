package main

import (
	"bytes"
	"context"
	"github.com/runsong2018/auto_generate_repo/internal/mysql"
	"github.com/runsong2018/auto_generate_repo/internal/repo"
	"github.com/urfave/cli/v2"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	version string
	build   string
	logger  *slog.Logger
	conf    Config
)

func init() {
	logger = slog.Default()
}

type Config struct {
	YamlPath string
	Include  string
	Dest     string
}

func main() {
	app := cli.NewApp()
	app.Name = "auto_generate_mysql_repo"
	app.Usage = "depend on yaml file generate repo go file"
	app.Version = version + " (" + build + ")"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "yaml",
			Aliases:     []string{"y"},
			Usage:       "yaml file path",
			Value:       "",
			Required:    false,
			Destination: &conf.YamlPath,
		},
		&cli.StringFlag{
			Name:        "include",
			Aliases:     []string{"i"},
			Usage:       "include file name pattern",
			Value:       "*",
			Required:    false,
			Destination: &conf.Include,
		},
		&cli.StringFlag{
			Name:        "dest",
			Aliases:     []string{"d"},
			Usage:       "generate dest path",
			Value:       "",
			Required:    false,
			Destination: &conf.Dest,
		},
	}

	app.Action = func(c *cli.Context) error {
		return action()
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func action() error {
	if conf.YamlPath == "" {
		conf.YamlPath, _ = os.Getwd()
	}
	if conf.Dest == "" {
		conf.Dest, _ = os.Getwd()
	}
	regex := regexp.MustCompile(".*")
	if conf.Include != "*" {
		regex = regexp.MustCompile(conf.Include)
	}
	err := filepath.Walk(conf.YamlPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".yaml") {
			return nil
		}
		if !regex.MatchString(info.Name()) {
			return nil
		}
		rawBody, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		repo := &repo.Repo{}
		err = repo.Parse(bytes.NewBuffer(rawBody))
		if err != nil {
			panic(err)
		}

		output := bytes.NewBuffer(nil)
		err = mysql.Handle(context.Background(), repo, output)
		if err != nil {
			panic(err)
		}

		path = filepath.Join(conf.Dest, strings.Split(info.Name(), ".")[0]+".mysql.go")
		err = os.WriteFile(path, output.Bytes(), 0777)
		if err != nil {
			panic(err)
		}

		if err = exec.Command("gofmt", "-w", path).Run(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return nil
}
