# 非同期処理の共通ルール

## 冪等性の担保

Workerがメッセージを重複して受信した場合や、処理途中で失敗して再試行された場合でも、ポイントが二重に付与されないことを保証します。

```sql
-- point_transactions テーブルで重複を防止
INSERT IGNORE INTO point_transactions 
  (id, idempotency_key, user_id, point_amount, created_at)
VALUES (?, ?, ?, ?, ?)
```

- **仕組み**: テーブルの `idempotency_key` (UNIQUE制約) を利用。
- **挙動**: `INSERT IGNORE` (または `ON CONFLICT DO NOTHING`) が成功した（＝新規処理）場合のみ、ポイント残高テーブル `point_root` を更新します。失敗した（＝処理済み）場合は、何もせずに正常終了として扱います。



## 状態遷移

```
PENDING (待機中)
   ↓ スケジューラがPick
PROCESSING (処理中 / Queue送信済み)
   ↓ Workerが処理
DONE (完了)  or  FAILED (失敗)
```