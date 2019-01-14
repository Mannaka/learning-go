package main

import (
	"errors"
)

// ErrNoAvatarはAvatarインスタンスがアバターのURLを返すことができない場合に発生するエラー
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません")

// Avatarはユーザーのプロフィール画像を表す型です
type Avatar interface {
	// GetAvatarURLは指定されたクライアントのアバターのURLを返す
	// 問題が生じた場合はエラーを返す。特に、URLを取得できなかった場合は
	// ErrNoAvatarURLを返す
	GetAvatarURL(c *client) (string, error)
}

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "//www/gracatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type FileSystemAvatar struct{}
var UseFileSystemAvatar FileSystemAvatar
func(_ FileSystemAvatar)GetAvatarURL(c *client)(string, error){
	if userid, ok := c.userData["userid"]; ok{
		if useridStr, ok := userid.(string); ok{
			return "/avatars/" + useridStr + ".jpg", nil
		}
	}
	return "", ErrNoAvatarURL
}
