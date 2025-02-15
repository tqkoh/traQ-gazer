package traqmessage

import (
	"context"
	"fmt"
	"h23s_15/model"
	"strings"
	"time"

	"github.com/traPtitech/go-traq"
	"golang.org/x/exp/slog"
)

// go routineの中で呼ぶこと
func PollingMessages() {
	pollingInterval := time.Minute * 3

	lastCheckpoint := time.Now()
	ticker := time.Tick(pollingInterval)

	for range ticker {
		now := time.Now()
		messages, err := collectMessages(lastCheckpoint, now)
		if err != nil {
			slog.Error(fmt.Sprintf("Failled to polling messages: %v", err))
			continue
		}

		lastCheckpoint = now

		slog.Info(fmt.Sprintf("Collect %d messages", len(messages.Hits)))
		// TODO: 取得したメッセージを使っての処理の呼び出し
		messageProcessor(messages.Hits)
	}
}

func messageProcessor(messages []traq.Message) {
	messageList, err := ConvertMessageHits(messages)
	if err != nil {
		slog.Error(fmt.Sprintf("Failled to convert messages: %v", err))
		return
	}
	notifyInfoList, err := model.FindMatchingWords(messageList)
	if err != nil {
		slog.Error(fmt.Sprintf("Failled to process messages: %v", err))
		return
	}

	slog.Info(fmt.Sprintf("Sending %d DMs...", len(notifyInfoList)))

	for _, notifyInfo := range notifyInfoList {
		err := sendMessage(notifyInfo.NotifyTargetTraqUuid, genNotifyMessageContent(notifyInfo.MessageId, notifyInfo.Words...))
		if err != nil {
			slog.Error(fmt.Sprintf("Failled to send message: %v", err))
			continue
		}
	}

	slog.Info("End of send DMs")
}

func genNotifyMessageContent(citeMessageId string, words ...string) string {
	list := make([]string, 0)
	for _, word := range words {
		item := fmt.Sprintf("「%s」", word)
		list = append(list, item)
	}

	return fmt.Sprintf("%s\n https://q.trap.jp/messages/%s", strings.Join(list, ""), citeMessageId)
}

func sendMessage(notifyTargetTraqUUID string, messageContent string) error {
	if model.ACCESS_TOKEN == "" {
		slog.Info("Skip sendMessage")
		return nil
	}

	client := traq.NewAPIClient(traq.NewConfiguration())
	auth := context.WithValue(context.Background(), traq.ContextAccessToken, model.ACCESS_TOKEN)
	_, _, err := client.UserApi.PostDirectMessage(auth, notifyTargetTraqUUID).PostMessageRequest(traq.PostMessageRequest{
		Content: messageContent,
	}).Execute()
	if err != nil {
		slog.Info("Error sending message: %v", err)
		return err
	}
	return nil
}

func collectMessages(from time.Time, to time.Time) (*traq.MessageSearchResult, error) {
	if model.ACCESS_TOKEN == "" {
		slog.Info("Skip collectMessage")
		return &traq.MessageSearchResult{}, nil
	}

	client := traq.NewAPIClient(traq.NewConfiguration())
	auth := context.WithValue(context.Background(), traq.ContextAccessToken, model.ACCESS_TOKEN)

	// 1度での取得上限は100まで　それ以上はoffsetを使うこと
	// https://github.com/traPtitech/traQ/blob/47ed2cf94b2209c8444533326dee2a588936d5e0/service/search/engine.go#L51
	result, _, err := client.MessageApi.SearchMessages(auth).After(from).Before(to).Limit(100).Execute()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func ConvertMessageHits(messages []traq.Message) (model.MessageList, error) {
	messageList := model.MessageList{}
	for _, message := range messages {
		messageList = append(messageList, model.MessageItem{
			Id:       message.Id,
			TraqUuid: message.UserId,
			Content:  message.Content,
		})
	}
	return messageList, nil
}
