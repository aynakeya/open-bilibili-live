package main

import (
	"context"
	"fmt"
	"github.com/aynakeya/open-bilibili-live"
	"os"
	"os/signal"
)

func main() {
	bliveApp := &openblive.BLiveApp{
		AppID:        1692609633998,
		AccessKey:    os.Getenv("openblive_access_key_id"),
		AccessSecret: os.Getenv("openblive_access_key_secret"),
	}
	client := bliveApp.CreateClient("BPCRO8L9EOE82")
	perr := client.Start()
	fmt.Println(perr)
	conn := client.GetLongConn()
	conn.OnDanmu(func(data openblive.DanmakuData) {
		fmt.Println(data.UName, data.FansMedalName, data.Msg)
	})
	fmt.Println(conn.EstablishConnection(context.Background()))
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println(conn.CloseConnection())
	fmt.Println(client.End())

}
