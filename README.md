# open-bilibili-live

个人维护的Bilibili直播开放平台SDK For Golang

官方文档请查看 [documentation](https://open-live.bilibili.com/document/bdb1a8e5-a675-5bfe-41a9-7a7163f75dbf)

Bilibili直播开放平台SDK Golang

## Installation

```
go get github.com/aynakeya/open-bilibili-live
```

## Usage

**Create App**

```go
bliveApp := &openblive.BLiveApp{
		AppID:        1692609633998,
		AccessKey:    os.Getenv("openblive_access_key_id"),
		AccessSecret: os.Getenv("openblive_access_key_secret"),
	}
```

**Create & Start & End Client**

```go
client := bliveApp.CreateClient("BPCRO8L9EOE82")
fmt.Println(client.Start())
fmt.Println(client.End())
```

**Establish Long Connection**

```go
conn := client.GetLongConn()
conn.OnDanmu(func(data openblive.DanmakuData) {
	fmt.Println(data.UName, data.FansMedalName, data.Msg)
})
fmt.Println(conn.EstablishConnection(context.Background()))
quit := make(chan os.Signal)
signal.Notify(quit, os.Interrupt)
<-quit
fmt.Println(conn.CloseConnection())
```