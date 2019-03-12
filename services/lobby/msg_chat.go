package main

// onChat : 聊天
func (accountObj *Account) onChat(data []byte) {
	Ctx.Log.Infoln("Chat, account:", accountObj.account)

}
