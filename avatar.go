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

type AuthAvatar struct{}
var UseAuthAvatar AuthAvatar
func(_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok{
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}