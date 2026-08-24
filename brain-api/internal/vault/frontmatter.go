package vault

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

const fence = "---"

func splitFrontmatter(raw []byte) (map[string]any, string, error) {
	text := string(raw)
	if strings.HasPrefix(text, "\uFEFF") {
		text = strings.TrimPrefix(text, "\uFEFF")
	}

	fm := map[string]any{}
	if !strings.HasPrefix(text, fence+"\n") && !strings.HasPrefix(text, fence+"\r\n") {
		return fm, strings.TrimSuffix(text, "\n"), nil
	}

	rest := text[len(fence):]
	rest = strings.TrimPrefix(rest, "\r\n")
	rest = strings.TrimPrefix(rest, "\n")

	end := indexFenceLine(rest)
	if end < 0 {
		return nil, "", fmt.Errorf("unclosed YAML frontmatter")
	}

	yamlBlock := rest[:end]
	body := rest[end:]
	body = strings.TrimPrefix(body, fence)
	body = strings.TrimPrefix(body, "\r\n")
	body = strings.TrimPrefix(body, "\n")
	body = strings.TrimSuffix(body, "\n")

	if strings.TrimSpace(yamlBlock) != "" {
		if err := yaml.Unmarshal([]byte(yamlBlock), &fm); err != nil {
			return nil, "", fmt.Errorf("frontmatter: %w", err)
		}
		if fm == nil {
			fm = map[string]any{}
		}
	}
	return fm, body, nil
}

func indexFenceLine(s string) int {
	offset := 0
	for {
		i := strings.Index(s[offset:], fence)
		if i < 0 {
			return -1
		}
		i += offset
		lineStart := 0
		if i > 0 {
			if s[i-1] != '\n' {
				offset = i + len(fence)
				continue
			}
			lineStart = i
		}
		after := i + len(fence)
		if after == len(s) || s[after] == '\n' || (s[after] == '\r' && after+1 < len(s) && s[after+1] == '\n') {
			return lineStart
		}
		offset = after
	}
}

func encodeNote(n Note) ([]byte, error) {
	var buf bytes.Buffer
	if len(n.Frontmatter) > 0 {
		enc := yaml.NewEncoder(&buf)
		enc.SetIndent(2)
		if err := enc.Encode(n.Frontmatter); err != nil {
			return nil, fmt.Errorf("frontmatter: %w", err)
		}
		if err := enc.Close(); err != nil {
			return nil, err
		}
		yamlBlock := strings.TrimSuffix(buf.String(), "\n")
		buf.Reset()
		fmt.Fprintf(&buf, "%s\n%s\n%s\n", fence, yamlBlock, fence)
	}
	body := n.Body
	if body != "" && !strings.HasSuffix(body, "\n") {
		body += "\n"
	}
	buf.WriteString(body)
	return buf.Bytes(), nil
}
