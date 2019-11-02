package engine

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/clients"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	chat "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

/*
А зачем нам рум гарбаж коллекток, если можно ходить по слайсу playing?
*/
func (lobby *Lobby) launchGarbageCollector(timeout float64) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_handle.go Join()")
		lobby.wGroup.Done()
	}()
	//utils.Debug(false, "launchGarbageCollector", lobby.Waiting.len())
	it := NewConnectionsIterator(lobby.Waiting)
	for it.Next() {
		waiter := it.Value()
		if waiter == nil {
			utils.Debug(true, "nil waiter")
		}
		t := waiter.Time()
		if time.Since(t).Seconds() > timeout {
			//utils.Debug(false, waiter.User.Name, " - bad")
			lobby.Leave(waiter, "")
		} else {
			//utils.Debug(false, waiter.User.Name, " - good", waiter.Disconnected(), time.Since(t).Seconds())
		}
	}
}

/*
sendMessagesToDB send unsent actions on messages(from lobby and
 rooms) to the database.

 If the chat ID associated with the lobby is not specified, an
attempt is made to contact the chat microservice and get the chat ID
and all messages in it. If successful, messages from the database are
added to the beginning of the lobby messages. In this case, the lobby
sends old messages marked 'deleted' and re-sends all lobby messages
to all 'waiting' connections. This is necessary for the correct order
of messages(first taken from the database, then written during the
absence of connection to the database). If the chat ID cannot be
retrieved, the function shuts down.

After receiving the identifier( or if the identifier has already been
set), an attempt is made to record unsent actions on messages.

single - the channel by which it is guaranteed that only 1 goroutine
 is trying to get the ID and the messages of the chat lobby
*/
func (lobby *Lobby) sendMessagesToDB(single chan interface{}) {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_handle.go Join()")
		lobby.wGroup.Done()
	}()

	var (
		err         error
		chatID      int32
		messages    []*MessageWithAction
		msgID       *chat.MessageID
		newMessages []*models.Message
	)
	utils.Debug(false, "sendMessagesToDB:")
	single <- nil
	// TODO change it if implement several service instances(i am about 0 )
	if lobby.dbChatID() == 0 {
		chatID, newMessages, err = GetChatIDAndMessages(lobby.location(),
			chat.ChatType_LOBBY, 0, lobby.SetImage)
		if err != nil {
			utils.Debug(false, "sendMessagesToDB error:", err.Error())
			return
		}
		oldMessages := lobby.insertMessages(newMessages)
		sendMessagesTodelete(lobby.send, All, oldMessages...)
		sendMessages(lobby.send, All, newMessages...)
		sendMessages(lobby.send, All, oldMessages...)

		lobby.setDBChatID(chatID)
	}
	<-single
	messages = lobby.NotSavedMessagesGetAndClear()
	utils.Debug(false, "lost messages:", len(messages))
	for _, messageAction := range messages {
		if messageAction.origin == nil {
			continue
		}
		if messageAction.message.ChatId == 0 {

			chatID, err = messageAction.getChatID()
			if err != nil {
				lobby.AddNotSavedMessage(messageAction)
				continue
			}
			messageAction.message.ChatId = chatID
		}
		switch messageAction.action {
		case models.Write:
			msgID, err = clients.ALL.Chat().AppendMessage(context.Background(), messageAction.message)
			if err == nil {
				messageAction.origin.ID = msgID.Value
			}
		case models.Update:
			_, err = clients.ALL.Chat().UpdateMessage(context.Background(), messageAction.message)

		case models.Delete:
			_, err = clients.ALL.Chat().DeleteMessage(context.Background(), messageAction.message)
		}
		if err != nil {
			utils.Debug(false, "again error:", err.Error())
			lobby.AddNotSavedMessage(messageAction)
		}
	}
}

func (lobby *Lobby) sendGamesToDB() {
	if lobby.done() {
		return
	}
	lobby.wGroup.Add(1)
	defer func() {
		utils.CatchPanic("lobby_handle.go Join()")
		lobby.wGroup.Done()
	}()

	var (
		err   error
		games = lobby.NotSavedGamesGetAndClear()
	)
	for _, game := range games {
		if err = lobby.db().SaveGame(*game); err != nil {
			lobby.AddNotSavedGame(game)
		}
	}
}

/*
Run accepts connections and messages from them
Goroutine. When it is finished, the lobby will be cleared
*/
func (lobby *Lobby) Run() {
	defer func() {
		utils.CatchPanic("lobby_handle.go Run()")
		lobby.Free()
	}()

	s10 := time.NewTicker(time.Second * 5)
	defer s10.Stop()
	var timeout float64
	timeout = 10

	single := make(chan interface{}, 1)
	defer close(single)

	go lobby.sendMessagesToDB(single)
	for {
		select {
		case <-s10.C:
			go lobby.launchGarbageCollector(timeout)
			go lobby.sendMessagesToDB(single)
			go lobby.sendGamesToDB()
		case connection := <-lobby.chanJoin:
			go lobby.Join(connection)
		case message := <-lobby.chanBroadcast:
			lobby.Analize(message)
		case <-lobby.chanBreak:
			return
		}
	}
}
