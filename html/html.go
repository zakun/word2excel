package html

import (
	"bufio"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
	"word2excel.io/common"
	"word2excel.io/question"
	"word2excel.io/question/rules"
)

const ImageDomain = "https://aop-1314572809.cos.ap-beijing.myqcloud.com/questions/zimp/"

type Html struct {
	No       int
	Name     string
	Template string
}

func NewHtml(no int, name, template string) *Html {
	return &Html{
		No:       no,
		Name:     name,
		Template: template,
	}
}

func (h *Html) ParseTemplateContent() ([]question.Question, common.Judge) {
	// 解析html
	htmlBytes := h.ParseHtmlFile()

	// 生成txt文件
	tmpDir := filepath.Dir(h.Name) + "/tmp"
	if !common.IsExistDir(tmpDir) {
		os.Mkdir(tmpDir, 0777)
	}

	txtName, _ := strings.CutSuffix(filepath.Base(h.Name), ".html")
	txtFullname := tmpDir + "/" + txtName + ".txt"

	// common.DD("txt #%v, %v", h.No, txtName)
	os.WriteFile(txtFullname, htmlBytes, 0766)

	return h.ParseTxtFile(txtFullname), common.Success()
}

func (h *Html) ParseHtmlFile() []byte {
	f1, err := os.Open(h.Name)
	common.Throw_panic(err)
	defer f1.Close()

	// 解析html
	doc, err := html.Parse(f1)
	common.Throw_panic(err)

	var sb strings.Builder
	h.ParseNode(doc, &sb)

	return []byte(sb.String())
}

func (h *Html) ParseTxtFile(name string) []question.Question {
	f, err := os.Open(name)
	common.Throw_panic(err)
	defer f.Close()

	rule, err := rules.GetRuleInstance("html", 30)
	common.Throw_panic(err)

	scanner := bufio.NewScanner(f)
	baseName := filepath.Base(name)
	for scanner.Scan() {
		text := scanner.Text()
		if text != "" {
			rule.StartParse(text, baseName)
		}
	}
	common.Throw_panic(scanner.Err())

	return rule.GetAllQuestions()
}

func (h *Html) ParseNode(n *html.Node, sb *strings.Builder) {
	switch n.Type {
	case html.ElementNode:
		if n.Data == "script" || n.Data == "table" || n.Data == "title" || n.Data == "style" || n.Data == "noscript" {
			return
		}

		if n.Data == "div" && n.Attr != nil {
			for _, attr := range n.Attr {
				if attr.Key == "class" && attr.Val == "mainQuestionDiv" {
					sb.WriteString("\n")
				} else if attr.Key == "style" && strings.Contains(attr.Val, "flex-direction: row") {
					sb.WriteString("\n")
				}
			}
		}

		if n.Data == "img" {
			for _, attr := range n.Attr {
				if attr.Key == "src" && strings.Contains(attr.Val, ";base64,") {
					content := attr.Val
					md5 := fmt.Sprintf("%x", md5.Sum([]byte(content)))

					imgInfo := strings.Split(content, ";")
					// imgExt := imgInfo[0][11:]
					bs64 := imgInfo[1][7:]
					reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(bs64))
					iw, ih, ext := common.GetWidthHeightFromImage(reader)
					// img 替换模板
					imgName := md5 + "." + ext
					tmpImageName := "[img=" + ImageDomain + imgName + "(" + common.Sprintf("%d,%d", iw, ih) + ")]"
					sb.WriteString(tmpImageName)
					// 保存图片
					imgFullName := filepath.Dir(h.Name) + "/tmp" + "/" + imgName
					bs64Bytes, _ := base64.StdEncoding.DecodeString(bs64)
					os.WriteFile(imgFullName, bs64Bytes, 0766)
				}
			}
		}

	case html.TextNode:
		text := strings.TrimSpace(n.Data)
		if strings.Contains(n.Data, "学生答案") {
			sb.WriteString("\n" + text)
		} else {
			sb.WriteString(text)
		}

	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		h.ParseNode(c, sb)
	}
}
