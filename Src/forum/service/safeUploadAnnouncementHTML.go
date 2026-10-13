package service

import (
	"bytes"
	"fmt"
	stdhtml "html"
	"net/url"
	"newaoe/Src/common/download/service"
	"newaoe/Src/common/upload"
	"newaoe/Src/util"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/microcosm-cc/bluemonday"
	nethtml "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func SafeUploadContent(contentPath string, content string, fileMap, imageMap, videoMap map[string]string) error {
	//检查是否为空
	if strings.TrimSpace(content) == "" {
		return fmt.Errorf("内容不能为空")
	}
	//安全过滤前端上传的html，并把图片、视频、文件链接换成真实路径
	safeHTML, err := sanitizeForumContent(content, fileMap, imageMap, videoMap)
	if err != nil {
		return err
	}
	if strings.TrimSpace(safeHTML) == "" {
		return fmt.Errorf("内容不能为空")
	}
	if !upload.UploadData(contentPath, []byte(safeHTML), "text/html; charset=utf-8") {
		return fmt.Errorf("内容上传失败")
	}
	return nil
}

var (
	forumContentPolicyOnce sync.Once
	forumContentPolicy     *bluemonday.Policy
)

func getForumContentPolicy() *bluemonday.Policy {
	forumContentPolicyOnce.Do(func() {
		p := bluemonday.UGCPolicy()
		p.AllowElements("video", "source")
		p.AllowAttrs("controls").OnElements("video")
		p.AllowAttrs("src").OnElements("video", "source")
		p.AllowAttrs("poster").OnElements("video")
		p.AllowAttrs("width", "height").Matching(bluemonday.NumberOrPercent).OnElements("video")
		p.AllowAttrs("type").Matching(regexp.MustCompile(`(?i)^[a-z0-9.+-]+/[a-z0-9.+-]+$`)).OnElements("source")
		forumContentPolicy = p
	})
	return forumContentPolicy
}

func sanitizeForumContent(content string, fileMap, imageMap, videoMap map[string]string) (string, error) {
	nodes, err := nethtml.ParseFragment(strings.NewReader(content), &nethtml.Node{
		Type:     nethtml.ElementNode,
		Data:     "body",
		DataAtom: atom.Body,
	})
	if err != nil {
		return "", err
	}
	tokens := make(map[string]string)
	for _, node := range nodes {
		rewriteForumAssetURLs(node, tokens, fileMap, imageMap, videoMap)
	}
	var buf bytes.Buffer
	for _, node := range nodes {
		if err := nethtml.Render(&buf, node); err != nil {
			return "", err
		}
	}
	safe := getForumContentPolicy().Sanitize(buf.String())
	for token, realPath := range tokens {
		safe = strings.ReplaceAll(safe, token, stdhtml.EscapeString(realPath))
	}
	return safe, nil
}

func rewriteForumAssetURLs(node *nethtml.Node, tokens map[string]string, fileMap, imageMap, videoMap map[string]string) {
	if node.Type == nethtml.ElementNode {
		switch node.Data {
		case "img":
			rewriteOrDropAssetAttr(node, "src", tokens, false, imageMap)
		case "video":
			rewriteOrDropAssetAttr(node, "src", tokens, false, videoMap)
			rewriteOrDropAssetAttr(node, "poster", tokens, false, imageMap)
		case "source":
			rewriteOrDropAssetAttr(node, "src", tokens, false, videoMap)
		case "a":
			//文件走附件下载；图片和视频链接仍直接打开
			if !rewriteAssetAttr(node, "href", tokens, true, fileMap) &&
				!rewriteAssetAttr(node, "href", tokens, false, imageMap, videoMap) {
				removeNodeAttr(node, "href")
			}
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		rewriteForumAssetURLs(child, tokens, fileMap, imageMap, videoMap)
	}
}

func rewriteOrDropAssetAttr(node *nethtml.Node, key string, tokens map[string]string, attachment bool, maps ...map[string]string) {
	if !rewriteAssetAttr(node, key, tokens, attachment, maps...) {
		removeNodeAttr(node, key)
	}
}

func rewriteAssetAttr(node *nethtml.Node, key string, tokens map[string]string, attachment bool, maps ...map[string]string) bool {
	index := -1
	for i, attr := range node.Attr {
		if strings.EqualFold(attr.Key, key) {
			index = i
			break
		}
	}
	if index < 0 {
		return false
	}
	downloadLink, ok := lookupAssetPath(node.Attr[index].Val, attachment, maps...)
	if !ok {
		return false
	}
	token := "forum-asset-" + util.UUID()
	node.Attr[index].Val = token
	tokens[token] = downloadLink
	return true
}

func removeNodeAttr(node *nethtml.Node, key string) {
	for i, attr := range node.Attr {
		if strings.EqualFold(attr.Key, key) {
			node.Attr = append(node.Attr[:i], node.Attr[i+1:]...)
			return
		}
	}
}

func lookupAssetPath(raw string, attachment bool, maps ...map[string]string) (string, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", false
	}
	candidates := []string{raw}
	if unescaped, err := url.PathUnescape(raw); err == nil && unescaped != raw {
		candidates = append(candidates, unescaped)
	}
	if parsed, err := url.Parse(raw); err == nil {
		if parsed.Path != "" {
			candidates = append(candidates, parsed.Path, path.Base(parsed.Path))
		}
		if parsed.Opaque != "" {
			candidates = append(candidates, parsed.Opaque, path.Base(parsed.Opaque))
		}
	}
	slashPath := strings.ReplaceAll(raw, "\\", "/")
	candidates = append(candidates, path.Base(slashPath))

	seen := make(map[string]struct{}, len(candidates))
	for _, candidate := range candidates {
		candidate = strings.TrimSpace(candidate)
		if candidate == "" || candidate == "." || candidate == "/" || candidate == "\\" {
			continue
		}
		if _, ok := seen[candidate]; ok {
			continue
		}
		seen[candidate] = struct{}{}
		for _, assetMap := range maps {
			realPath, ok := assetMap[candidate]
			if !ok || realPath == "" {
				continue
			}
			//根据realPath获取下载链接(强制走CDN)。文件传attachment=true，按附件下载。
			downloadLink, downloadType := service.PublicFileDownload(realPath, attachment, true)
			if downloadLink == "" || downloadType == service.DownloadURLType_UnKnown {
				return "", false
			}
			if attachment && downloadType == service.DownloadURLType_CDN {
				parsedLink, err := url.Parse(downloadLink)
				if err != nil {
					return "", false
				}
				query := parsedLink.Query()
				query.Set("attachment", "true")
				parsedLink.RawQuery = query.Encode()
				downloadLink = parsedLink.String()
			}
			return downloadLink, true
		}
	}
	return "", false
}
