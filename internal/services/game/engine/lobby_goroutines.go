package engine

import (
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/synced"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"

	chat "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/proto"
	ctypes "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/database"
)

/*
А зачем нам рум гарбаж коллекток, если можно ходить по слайсу playing?
Чтобы комната корректно обрабатывала уход игроков
*/
func (lobby *Lobby) launchGarbageCollector() {
	lobby.s.Do(func() {
		//utils.Debug(false, "launchGarbageCollector", lobby.Waiting.len())
		it := NewConnectionsIterator(lobby.Waiting)
		timeout := lobby.config().Lobby.ConnectionTimeout.Duration
		for it.Next() {
			waiter := it.Value()
			if waiter == nil {
				utils.Debug(true, "nil waiter")
			}
			t := waiter.Time()
			if time.Since(t) > timeout {
				//utils.Debug(false, waiter.User.Name, " - bad")
				lobby.Leave(waiter, "")
			} else {
				//utils.Debug(false, waiter.User.Name, " - good", waiter.Disconnected(), time.Since(t).Seconds())
			}
		}
	})
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
func (lobby *Lobby) sendMessagesToDB() {
	lobby.s.Do(func() {
		var (
			err         error
			chatID      int32
			messages    []*MessageWithAction
			msgID       *chat.MessageID
			newMessages []*models.Message
		)
		utils.Debug(false, "sendMessagesToDB:")
		// TODO change it if implement several service instances(i am about 0 )
		if lobby.dbChatID() == 0 {
			chatID, newMessages, err = GetChatIDAndMessages(lobby.ChatService, lobby.location(),
			ctypes.LobbyType, 0, lobby.SetImage)
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
				msgID, err = lobby.ChatService.AppendMessage(messageAction.message)
				if err == nil {
					messageAction.origin.ID = msgID.Value
				}
			case models.Update:
				_, err = lobby.ChatService.UpdateMessage(messageAction.message)

			case models.Delete:
				_, err = lobby.ChatService.DeleteMessage(messageAction.message)
			}
			if err != nil {
				utils.Debug(false, "again error:", err.Error())
				lobby.AddNotSavedMessage(messageAction)
			}
		}
	})
}

func (lobby *Lobby) sendGamesToDB() {
	lobby.s.Do(func() {
		var (
			err   error
			games = lobby.NotSavedGamesGetAndClear()
		)
		for _, game := range games {
			if err = lobby.db().Save(*game); err != nil {
				lobby.AddNotSavedGame(game)
			}
		}
	})
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

	intervals := lobby.config().Lobby.Intervals

	var gc = synced.SingleGoroutine{}
	gc.Init(intervals.GarbageCollector.Duration, lobby.launchGarbageCollector)
	defer gc.Close()

	var m2db = synced.SingleGoroutine{}
	m2db.Init(intervals.MessagesToDB.Duration, lobby.sendMessagesToDB)
	defer m2db.Close()

	var g2db = synced.SingleGoroutine{}
	g2db.Init(intervals.GamesToDB.Duration, lobby.sendGamesToDB)
	defer g2db.Close()

	fmt.Println("we run!")
	lobby.sendGamesToDB()
	for {
		select {
		case <-gc.C():
			go gc.Do()
		case <-m2db.C():
			go m2db.Do()
		case <-g2db.C():
			go g2db.Do()
		case connection := <-lobby.chanJoin:
			go lobby.Join(connection)
		case message := <-lobby.chanBroadcast:
			utils.Debug(false, "analize me")
			lobby.Analize(message)
		case <-lobby.chanBreak:
			return
		}
	}
}
