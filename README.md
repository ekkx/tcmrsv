# tcmrsv

東京音楽大学 - 練習室予約ライブラリ

### Getting Started

```go
package main

import (
  "time"

  "github.com/ekkx/tcmrsv"
)

func main() {
  client := tcmrsv.New()

  if err := client.Login(&tcmrsv.LoginParams{
    UserID:   "your_user_id",
    Password: "your_password",
  }); err != nil {
    panic(err)
  }

  if err = client.Reserve(&tcmrsv.ReserveParams{
    Campus:     tcmrsv.CampusIkebukuro,
    RoomID:     "42d1eacc-60d5-428b-8c64-aef11a512c30",
    Date:       tcmrsv.Today().AddDays(1),
    FromHour:   12,
    FromMinute: 0,
    ToHour:     14,
    ToMinute:   0,
  }); err != nil {
    panic(err)
  }
}
```

詳しい使い方は [cmd/example](https://github.com/ekkx/tcmrsv/tree/master/cmd/example) を参照してください。

### Roadmap

- [x] ログイン
- [x] 予約
- [x] 予約のキャンセル
- [x] 自分の予約一覧取得
- [x] 利用可能な練習室の取得
- [x] 操作が正常に終了したか判定
- [x] 練習室一覧を固定化（ID, ピアノの種類）
- [x] 練習室とキャンパスを紐付け
- [x] 独自日付パッケージの作成
- [x] できる範囲で入力値のバリデーション
- [x] 3 日以内の予約しか受け付けない（12:00AM~12:00PM までは 2 日以内）
- [x] 利用可能な練習室として、◯ か disabled ではない `<input>` のみ取得
- [x] レスポンスのモックテスト
