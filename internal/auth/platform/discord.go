package platform

import (
	"fmt"

	"study-planner/internal/auth"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/oauth2"
)

type DiscordUserSupplier struct{}

func NewDiscordUserSupplier() *DiscordUserSupplier {
	return &DiscordUserSupplier{}
}

func (s *DiscordUserSupplier) GetUserInfo(token *oauth2.Token) (*auth.UserInfo, error) {
	session, err := discordgo.New(fmt.Sprintf("%s %s", token.TokenType, token.AccessToken))
	if err != nil {
		return nil, fmt.Errorf("ini")
	}

	user, err := session.User("@me")
	if err != nil {
		return nil, err
	}

	return &auth.UserInfo{
		ExternalID: user.ID,
		Platform:   "discord",
		Name:       user.Username,
		AvatarURL:  user.AvatarURL(""),
	}, nil
}
