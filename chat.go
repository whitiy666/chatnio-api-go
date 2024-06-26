package chatnio

import "github.com/whitiy666/chatnio-api-go/utils"

type Chat struct {
	Id    int
	Uri   string
	Token string
	Conn  *utils.WebSocket
}

type ChatAuthForm struct {
	Id    int    `json:"id"`
	Token string `json:"token"`
}

type ChatRequestForm struct {
	Message string `json:"message"`
	Model   string `json:"model"`
	Web     bool   `json:"web"`
}

type ChatPartialResponse struct {
	Message      string  `json:"message"`
	Conversation int64   `json:"conversation"`
	Keyword      string  `json:"keyword"`
	Quota        float32 `json:"quota"`
	End          bool    `json:"end"`
}

func (i *Instance) NewChat(id int) (*Chat, error) {

	conn, err := utils.NewWebsocket(i.GetChatEndpoint())
	if err != nil {
		return nil, err
	}
	return &Chat{
		Id:    id,
		Uri:   i.GetChatEndpoint(),
		Token: i.GetApiKey(),
		Conn:  conn,
	}, nil
}

func (c *Chat) Send(v interface{}) bool {
	return c.Conn.Send(v)
}

func (c *Chat) Close() error {
	return c.Conn.Close()
}

func (c *Chat) DeferClose() {
	c.Conn.DeferClose()
}

func (c *Chat) SendAuthRequest() bool {
	return c.Send(ChatAuthForm{
		Id:    c.Id,
		Token: c.Token,
	})
}

func (c *Chat) AskStream(form *ChatRequestForm, callback func(ChatPartialResponse)) {
	// for authentication
	if c.Conn.IsEmpty() {
		c.SendAuthRequest()
	}

	c.Send(map[string]interface{}{
		"type":    "chat",
		"message": form.Message,
		"model":   form.Model,
		"web":     form.Web,
	})

	for {
		form := utils.ReadForm[ChatPartialResponse](c.Conn)
		if form == nil {
			continue
		}

		callback(*form)
		if form.End {
			break
		}
	}
}

func (c *Chat) Ask(form *ChatRequestForm, channel chan ChatPartialResponse) {
	worker := func() {
		c.AskStream(form, func(res ChatPartialResponse) {
			channel <- res
		})

		close(channel)
	}

	go worker()
}
