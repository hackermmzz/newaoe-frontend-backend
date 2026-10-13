package service

import (
	"strings"

	"newaoe/Src/forum/dao"
	"newaoe/Src/forum/model"
	"newaoe/Src/util"

	nethtml "golang.org/x/net/html"
)

const forumSummaryMaxLen = 120

func ForumAnnouncementAdd(id, title, contentPath, content string) error {
	_, err := dao.ForumAnnouncementInsert(nil, &model.ForumAnnouncement{
		ID:          id,
		Title:       title,
		ContentPath: contentPath,
		Summary:     generateSummary(content),
		SubmitTime:  util.UTC_Time(),
	})
	return err
}

func generateSummary(content string) string {
	if strings.TrimSpace(content) == "" {
		return ""
	}
	doc, err := nethtml.Parse(strings.NewReader(content))
	if err != nil {
		return ""
	}
	var buf strings.Builder
	collectForumSummaryText(&buf, doc)
	text := strings.Join(strings.Fields(buf.String()), " ")
	return truncateForumSummary(text, forumSummaryMaxLen)
}

func collectForumSummaryText(buf *strings.Builder, node *nethtml.Node) {
	if node.Type == nethtml.ElementNode {
		switch node.Data {
		case "img", "video", "source", "script", "style":
			return
		}
	}
	if node.Type == nethtml.TextNode {
		buf.WriteString(node.Data)
	}
	if isForumSummaryBlock(node) {
		buf.WriteByte(' ')
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		collectForumSummaryText(buf, child)
	}
	if isForumSummaryBlock(node) {
		buf.WriteByte(' ')
	}
}

func isForumSummaryBlock(node *nethtml.Node) bool {
	if node.Type != nethtml.ElementNode {
		return false
	}
	switch node.Data {
	case "p", "div", "li", "br", "blockquote", "pre", "tr", "h1", "h2", "h3", "h4", "h5", "h6", "ul", "ol", "section", "article", "hr":
		return true
	default:
		return false
	}
}

func truncateForumSummary(text string, maxLen int) string {
	runes := []rune(text)
	if len(runes) <= maxLen {
		return text
	}
	cut := maxLen
	for i := maxLen; i > 0; i-- {
		if runes[i-1] == ' ' {
			cut = i - 1
			break
		}
	}
	if cut <= 0 {
		cut = maxLen
	}
	return string(runes[:cut]) + "..."
}
