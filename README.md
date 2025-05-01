# tcmrsv

練習室自動予約ツール

### Getting Started

```go
package main

import (
  "time"

  "github.com/ekkx/tcmrsv"
)

func main() {
  rsv := tcmrsv.New()

  if err := rsv.Login(&tcmrsv.LoginParams{
    UserID:   "-----",
    Password: "-----",
  }); err != nil {
    panic(err)
  }

  if err = rsv.Reserve(&tcmrsv.ReserveParams{
    Campus:     tcmrsv.CampusIkebukuro,
    RoomID:     "42d1eacc-60d5-428b-8c64-aef11a512c30",
    Date:       time.Now(),
    FromHour:   12,
    FromMinute: 0,
    ToHour:     14,
    ToMinute:   0,
  }); err != nil {
    panic(err)
  }
}
```

### Roadmap

- [x] ログイン
- [x] 予約
- [x] 予約のキャンセル
- [x] 自分の予約一覧取得
- [x] 利用可能な練習室の取得
- [x] 操作が正常に終了したか判定
- [x] 練習室一覧を固定化（ID, ピアノの種類）
- [x] 練習室とキャンパスを紐付け
- [ ] できる範囲で入力値のバリデーション
