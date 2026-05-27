package generator

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type xmlNode struct {
	Name     string
	Attrs    map[string]string
	Text     string
	Children []*xmlNode
}

func JSONFromXML(input string) (string, error) {
	dec := xml.NewDecoder(strings.NewReader(input))
	var stack []*xmlNode
	var root *xmlNode
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			node := &xmlNode{Name: t.Name.Local, Attrs: map[string]string{}}
			for _, attr := range t.Attr {
				node.Attrs["@"+attr.Name.Local] = attr.Value
			}
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				parent.Children = append(parent.Children, node)
			}
			stack = append(stack, node)
			if root == nil {
				root = node
			}
		case xml.CharData:
			if len(stack) > 0 {
				text := strings.TrimSpace(string(t))
				if text != "" {
					stack[len(stack)-1].Text += text
				}
			}
		case xml.EndElement:
			if len(stack) == 0 {
				return "", fmt.Errorf("unexpected XML close tag %s", t.Name.Local)
			}
			stack = stack[:len(stack)-1]
		}
	}
	if root == nil {
		return "", fmt.Errorf("no XML root found")
	}
	wrapped := map[string]any{root.Name: nodeValue(root)}
	data, err := json.MarshalIndent(wrapped, "", "  ")
	return string(data), err
}

func nodeValue(node *xmlNode) any {
	if len(node.Attrs) == 0 && len(node.Children) == 0 {
		return node.Text
	}
	obj := map[string]any{}
	for key, value := range node.Attrs {
		obj[key] = value
	}
	grouped := map[string][]any{}
	for _, child := range node.Children {
		grouped[child.Name] = append(grouped[child.Name], nodeValue(child))
	}
	for key, values := range grouped {
		if len(values) == 1 {
			obj[key] = values[0]
		} else {
			obj[key] = values
		}
	}
	if node.Text != "" {
		obj["#text"] = node.Text
	}
	return obj
}

func (n *xmlNode) String() string {
	var b bytes.Buffer
	b.WriteString(n.Name)
	return b.String()
}
