package bot

import (
	"FilesToTarBot/config"
	"FilesToTarBot/database"
	"context"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"log"
	"strings"
)

var Dispatcher = tg.NewUpdateDispatcher()
var api *tg.Client
var sender *message.Sender

// RunBot runs the bot to receive the updates
func RunBot(_ context.Context, client *telegram.Client) error {
	api = client.API()
	sender = message.NewSender(tg.NewClient(client))
	Dispatcher.OnNewMessage(func(ctx context.Context, entities tg.Entities, u *tg.UpdateNewMessage) error {
		m, ok := u.Message.(*tg.Message)
		if !ok || m.Out {
			// Outgoing message, not interesting.
			return nil
		}
		// Check if the user is allowed to use the bot
		userID := m.PeerID.(*tg.PeerUser).UserID
		if !config.IsUserAllowed(userID) {
			return nil
		}
		// Check file
		replyText := "Please send a media to bot"
		if m.Media != nil {
			doc, ok := m.Media.(*tg.MessageMediaDocument)
			if ok {
				if doc, ok := doc.Document.AsNotEmpty(); ok {
					var filename string
					for _, attribute := range doc.Attributes {
						if name, ok := attribute.(*tg.DocumentAttributeFilename); ok {
							filename = name.FileName
							break
						}
					}
					err := database.MainDatabase.AddFile(userID, database.File{
						FileReference: doc.FileReference,
						Name:          filename,
						ID:            doc.ID,
						AccessHash:    doc.AccessHash,
						Size:          int64(doc.Size),
					})
					switch err {
					case database.TooBigFileError, database.FileAlreadyExistsError:
						replyText = err.Error()
					case nil:
						replyText = "Added!"
					default:
						replyText = "cannot add file"
						log.Println("cannot add file:", err)
					}
				}
			}
		} else {
			switch m.Message {
			case "/start":
				replyText = "Welcome! Send files to bot in order to add them to the list of files. When done, send /done in order to tar your files."
			case "/reset":
				database.MainDatabase.Reset(userID)
				replyText = "Removed all files"
			}
			if strings.HasPrefix(m.Message, "/done") {
				go uploadFiles(ctx, userID, entities, u, m.Message[:len("/done")])
				return nil
			}
		}

		// Send the link or error
		_, err := sender.Reply(entities, u).Text(ctx, replyText)
		return err
	})
	return nil
}
