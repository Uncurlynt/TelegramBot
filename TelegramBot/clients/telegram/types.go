package telegram

// Представляет ответ от сервера Telegram на запрос обновлений (новых сообщений)
type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

// Представляет отдельное обновление (новое сообщение) от сервера Telegram
type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

// Представляет входящее сообщение от пользователя в Telegram
type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

// Представляет информацию об отправителе сообщения
type From struct {
	Username string `json:"username"`
}

// Представляет информацию о чате (группе/пользователе), в котором было получено сообщение
type Chat struct {
	ID int `json:"id"`
}
