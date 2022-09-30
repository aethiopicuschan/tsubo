# tsubo

[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen?style=flat-square)](/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/aethiopicuschan/tsubo.svg)](https://pkg.go.dev/github.com/aethiopicuschan/tsubo)

5ちゃんねるの読み書き操作をラップしたAPIサーバ+Golangライブラリ

素直に5ちゃんねるを操作しようとすると厳しいので、気持ちモダンな形で間に挟まるAPIサーバとして機能します

- このソフトウェアは専ブラでもプロキシでもありません
- [5ch.net 専用ブラウザの開発者の皆さまへ](https://developer.5ch.net/) にて禁止されている「5ch.netが提供するAPI」は使用しません

## サポート/TODO

- [x] BBSMENUの取得
- [x] スレ一覧の取得
- [ ] スレの読み込み
- [ ] スレへの書き込み
- [ ] スレ立て

## APIサーバとして

```sh
go install github.com/aethiopicuschan/tsubo@latest
tsubo
```

### エンドポイント例

- `GET /bbsmenu`
- `GET /subjects?board=https://greta.5ch.net/poverty/`

## Golangライブラリとして

主に以下パッケージを `go get` することで利用可能です

```
github.com/aethiopicuschan/tsubo/bbsmenu
github.com/aethiopicuschan/tsubo/subject
```

## テスト

```sh
go test ./...
```
