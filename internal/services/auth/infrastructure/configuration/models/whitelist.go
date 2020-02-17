package models

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/infrastructure/configuration"
	"gopkg.in/oauth2.v3/models"
)

type Client struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
	Domain string `json:"domain"`
	UserID string `json:"user_id"`
}

func (c *Client) Get() *models.Client {
	return &models.Client{
		ID:     c.ID,
		Secret: c.Secret,
		Domain: c.Domain,
		UserID: c.UserID,
	}
}

func (c *Client) Set(m *models.Client) {
	c.ID = m.ID
	c.Secret = m.Secret
	c.Domain = m.Domain
	c.UserID = m.UserID
}

type Whitelist []*Client

func (w *Whitelist) Get() configuration.Whitelist {
	var list configuration.Whitelist
	for _, client := range *w {
		list = append(list, client.Get())
	}
	return list
}

func (w *Whitelist) Set(list configuration.Whitelist) {
	*w = []*Client{}
	for _, client := range list {
		var c = &Client{}
		c.Set(client)
		*w = append(*w, c)
	}
}
