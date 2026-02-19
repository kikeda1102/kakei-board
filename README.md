# kakei-board

ゲーミフィケーション要素のある家計簿アプリ。
イベントソーシングを活かした「過去の自分とのスコアアタック」で支出改善を動機づける。

## 技術スタック

| レイヤー | 技術 |
|---------|------|
| Web Frontend | React + TypeScript + RTK Query |
| Mobile | React Native + TypeScript + RTK Query |
| Backend (Command/Query) | Go (CQRS + Event Sourcing) |
| Event Store / Read DB | MySQL |

## アーキテクチャ

CQRS + Event Sourcing をバーティカルスライスで実装。

### MVP スライス

1. 支出を記録する
2. 月次サマリーを見る
3. 予算を設定して進捗を追跡する
4. スコアボード（過去の自分との比較）

## 開発

```bash
# インフラ起動
docker compose up -d

# バックエンド
cd backend && go test ./...

# Web フロントエンド
cd web && npm run dev
```
