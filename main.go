package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	tgclient "github.com/featt/tg/tg"
	middleware "github.com/featt/tg/tg/middlewares"
	"github.com/go-faster/errors"
	"github.com/gotd/contrib/storage"
	"github.com/gotd/td/tg"
)



func main() {
	tgClient := tgclient.NewTG()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	tgClient.Dispatcher.OnNewMessage(func(ctx context.Context, e tg.Entities, u *tg.UpdateNewMessage) error {
		msg, ok := u.Message.(*tg.Message)
		if !ok {
			return nil
		}
		if msg.Out {
			return nil
		}
		p, err := storage.FindPeer(ctx, tgClient.PeerDB, msg.GetPeerID())
		if err != nil {
			return err
		}
		var chId int64
		for k, _ := range e.Channels {
			chId = k
		}
		
		if chId != 0 {
			ch := e.Channels[chId]
			tgClient.Api.ChannelsJoinChannel(ctx, ch.AsInput())
			req := &tg.MessagesGetHistoryRequest{
				Peer: ch.AsInputPeer(),
				Limit: 2,
			}			
			msgs, err := tgClient.Api.MessagesGetHistory(ctx, req)
			if err != nil {
				fmt.Println(err.Error())
			}
			data, _ := msgs.AsModified()
			for _, v := range data.GetMessages() {
				fmt.Println(v)
			}
			//fmt.Println(data.GetMessages())
		}
		fmt.Printf("%s: %s\n", p, msg.Message)
		return nil
	})
	err := middleware.Info(ctx, tgClient.Waiter, tgClient.Client, tgClient.Api, tgClient.UpdatesRecovery)
	if err != nil {
		if errors.Is(err, context.Canceled) && ctx.Err() == context.Canceled {
			fmt.Println("\rClosed")
			os.Exit(0)
		}
		_, _ = fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Done")
		os.Exit(0)
	}

}