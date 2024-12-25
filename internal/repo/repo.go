package repo

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/runsong2018/auto_generate_repo/internal/format"
	"gopkg.in/yaml.v3"
	"io"
	"reflect"
	"regexp"
	"strings"
)

type (
	Repo struct {
		PackageName string `yaml:"packageName,omitempty"`
		Name        string `yaml:"name,omitempty"`
		TableName   string `yaml:"tableName,omitempty"`

		Fields  []Field `yaml:"fields,omitempty"`
		Indexes []Index `yaml:"indexes,omitempty"`
	}

	Field struct {
		Name    string      `yaml:"name,omitempty"`    // 名称
		Type    string      `yaml:"type,omitempty"`    // 类型
		Size    *int64      `yaml:"size,omitempty"`    // 大小
		Default interface{} `yaml:"default,omitempty"` // 默认
		Comment string      `yaml:"comment,omitempty"` // 备注

		Gorm   string
		Option string //fuzzy 支持模糊匹配;
	}

	Index struct {
		Name   []string `yaml:"name,omitempty"`   // 索引
		Unique bool     `yaml:"unique,omitempty"` // 唯一
	}
)

func (r *Repo) Parse(reader io.Reader) (err error) {
	return r.parse(reader)
}

func (r *Repo) parse(reader io.Reader) (err error) {
	buf := bytes.NewBuffer(nil)
	if _, err = buf.ReadFrom(reader); err != nil {
		return
	}

	body, fieldToCommentMap := r.parseComment(buf)
	if err = yaml.Unmarshal(body, r); err != nil {
		return
	}
	for idx, field := range r.Fields {
		var ok bool
		var comment string

		if comment, ok = fieldToCommentMap[field.Name]; ok {
			r.Fields[idx].Comment = comment
		}
	}

	r.fillDefaultColumn()

	err = r.parseGorm()
	if err != nil {
		err = errors.New("GORM配置解析失败")
		return
	}

	return
}

func (r *Repo) parseComment(buf *bytes.Buffer) (body []byte, m map[string]string) {
	prefixToken := `- name:`
	suffixToken := "#"

	m = make(map[string]string)
	lines := strings.Split(buf.String(), "\n")
	for _, line := range lines {
		if strings.Contains(line, prefixToken) && strings.Contains(line, suffixToken) {

			prefixIdx := strings.Index(line, prefixToken)
			suffixIdx := strings.Index(line, suffixToken)

			field := line[prefixIdx+len(prefixToken) : suffixIdx]
			comment := line[suffixIdx+1:]

			m[strings.ReplaceAll(field, " ", "")] = strings.ReplaceAll(comment, " ", "")
		}
	}

	reg := regexp.MustCompile(`#.*`)
	body = reg.ReplaceAll(buf.Bytes(), []byte(""))
	return
}

func (r *Repo) parseGorm() (err error) {
	fieldToIndexesMap := make(map[string][]string)
	for _, index := range r.Indexes {
		var indexName string

		if index.Unique {
			indexName = "unique_index:idx"
		} else {
			indexName = "index:idx"
		}

		for _, field := range index.Name {
			indexName += "_" + format.SnakeString(field)
		}

		for _, field := range index.Name {
			fieldToIndexesMap[field] = append(fieldToIndexesMap[field], indexName)
		}
	}

	for idx, field := range r.Fields {
		var blocks []string

		if field.Size != nil {
			blocks = append(blocks, fmt.Sprintf(`size:%d`, *field.Size))
		}

		if field.Default != nil && !reflect.ValueOf(field.Default).IsZero() {
			if str, ok := field.Default.(string); ok {
				blocks = append(blocks, fmt.Sprintf(`default:'%s'`, str))
			} else {
				blocks = append(blocks, fmt.Sprintf(`default:%v`, field.Default))
			}
		}

		if indexes, ok := fieldToIndexesMap[field.Name]; ok {
			blocks = append(blocks, strings.Join(indexes, ";"))
		}

		if len(r.Fields[idx].Gorm) > 0 {
			blocks = append(blocks, r.Fields[idx].Gorm)
		}

		if len(blocks) != 0 {
			r.Fields[idx].Gorm = fmt.Sprintf(`gorm:"%s"`, strings.Join(blocks, ";"))
		}
	}

	return
}

func (r *Repo) fillDefaultColumn() {
	size := int64(100)
	r.Fields = append(r.Fields, Field{
		Name: "CreatedBy",
		Type: "string",
		Size: &size,
		Gorm: "comment:创建者",
	})
	r.Fields = append(r.Fields, Field{
		Name: "CreatedAt",
		Type: "time.Time",
		Gorm: "comment:创建时间",
	})
	r.Fields = append(r.Fields, Field{
		Name: "ModifiedBy",
		Type: "string",
		Size: &size,
		Gorm: "comment:修改者",
	})
	r.Fields = append(r.Fields, Field{
		Name: "ModifiedAt",
		Type: "time.Time",
		Gorm: "comment:修改时间",
	})
	r.Fields = append(r.Fields, Field{
		Name: "DeletedAt",
		Type: "gorm.DeletedAt",
		Gorm: "comment:删除时间",
	})
}
