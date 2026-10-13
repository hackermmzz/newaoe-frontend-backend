package service

import (
	"newaoe/Src/forum/dao"
	"newaoe/Src/forum/model"
)

func CommentAdd(comment *model.ForumComment) error {
	_, err := dao.ForumCommentInsert(nil, comment)
	return err
}
