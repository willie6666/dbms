# API Spec 差異對照表

比對對象：`add_api_spec.md`（較新版本）vs `bananaapi_spec.md`（較舊版本）

> **閱讀說明**
> - ✅ `add_api_spec` 有、`bananaapi_spec` 缺少或不同
> - ❌ `bananaapi_spec` 有、`add_api_spec` 缺少或不同
> - 🔴 **影響前端整合**：欄位名稱或回傳格式不同，需確認哪個版本正確

---

## 1. 使用者與權限 (Users & Auth)

---

### `[POST] /api/auth/register` (註冊新帳號)

| 項目 | add_api_spec | bananaapi_spec |
|------|-------------|----------------|
| 201 回傳 user 物件 | 包含 `email` 欄位 | **缺少** `email` 欄位 |

```diff
// add_api_spec (較完整)
"user": {
  "id": 1,
  "username": "PlayerOne",
+ "email": "player1@test.com",
  "role": "USERS"
}
```

---

### `[POST] /api/auth/login` (使用者登入)

| 項目 | add_api_spec | bananaapi_spec |
|------|-------------|----------------|
| 額外錯誤碼 | ✅ `403 Forbidden: "This account is not active"` (帳號被停權時) | ❌ 無此錯誤碼 |

---

### `[PUT] /api/users/profile` (修改個人資料)

| 項目 | add_api_spec | bananaapi_spec |
|------|-------------|----------------|
| 🔴 成功回傳格式 | 包含 `user` 物件（有 id/username/email/role） | 僅回傳 `{"message": "Profile updated successfully"}` |

```diff
// add_api_spec 的 200 回應
{
  "message": "Profile updated successfully",
+ "user": { "id": 1, "username": "NewName", "email": "new@test.com", "role": "USERS" }
}
```

---

### `[PUT] /api/admin/users/{id}/suspend` (停權)

| 項目 | add_api_spec | bananaapi_spec |
|------|-------------|----------------|
| 行為說明 | ✅ 明確說明為 **Toggle** 行為（ACTIVE ↔ DEACTIVE） | ❌ 無 Toggle 說明，給人印象是單向停權 |

---

## 2. 商店與遊戲 (Store & Games)

---

### `[GET] /api/developer/games` (查看自己的遊戲列表)

| 項目 | add_api_spec | bananaapi_spec |
|------|-------------|----------------|
| 回傳格式說明 | ✅ 明確說明包含 media：`...game_objects_with_media...` | 僅 `[ { "game_id": 1, "title": "...", "price": 350 } ]` |
| ADMIN 行為說明 | ✅ 說明 ADMIN 可查看全部遊戲 | ❌ 無此說明 |

---

### `[POST] /api/developer/games` (上架新遊戲)

| 項目 | add_api_spec | bananaapi_spec |
|------|-------------|----------------|
| `desc` 備註 | ✅ 補充說明「支援 Markdown」 | 無此備註 |

> 兩者 Request Body 欄位相同（title, price, desc），無功能性差異。

---

### `[PUT] /api/developer/games/{id}` (編輯遊戲資訊)

| 項目 | add_api_spec | bananaapi_spec |
|------|-------------|----------------|
| Request Body 欄位 | `price`, `desc` | `price`, `desc` |
| 🔴 注意事項 | 僅說明「ADMIN 可編輯任何遊戲」 | ✅ 有額外說明「`title` 無法透過此 API 修改；`price` 若傳 `0` 仍會寫入」 |
| ADMIN 說明 | ✅ 說明 ADMIN 可編輯任何遊戲 | ❌ 無此說明 |

---

### `[POST] /api/developer/games/{id}/media` (上傳遊戲素材)

| 項目 | add_api_spec | bananaapi_spec |
|------|-------------|----------------|
| 表格格式 | ✅ 更完整（4 欄：欄位名稱/類型/必填/說明） | 3 欄（欄位/必填/說明），無類型欄 |
| 儲存路徑說明 | ✅ 列出兩種路徑的完整 assets 目錄格式 | 僅說明對外 URL（`/media/images/...` 和 `/downloads/...`） |
| `400 Bad Request` | ❌ 無此錯誤碼 | ✅ `{"error": "Missing file field"}` |

> 注意：`bananaapi_spec` 補了 `400` 錯誤碼但 `add_api_spec` 沒有，建議兩者應合併。

---

### `[DELETE] /api/developer/games/{id}/media/{media_id}` (刪除遊戲素材)

| 項目 | add_api_spec | bananaapi_spec |
|------|-------------|----------------|
| 🔴 404 錯誤訊息 | `"Media not found"` 或 `"Game not found"` (兩種可能) | 僅 `"Game not found"` |

---

### `[POST] /api/developer/tags` (建立新標籤)

| 項目 | add_api_spec | bananaapi_spec |
|------|-------------|----------------|
| 🔴 201 成功回傳 | 包含 `data` 物件：`{ "tag_id": 1, "tag_name": "Action" }` | 僅 `{"message": "Tag created successfully"}` |

---

## 3. 訂單、購物車與客服

> 此區段兩份文件**完全一致**，無差異。

---

## 4. 遊戲庫與願望清單

---

### `[GET] /api/protected/library/{game_id}/download` (下載遊戲)

| 項目 | add_api_spec | bananaapi_spec |
|------|-------------|----------------|
| 🔴 成功回傳格式 | **直接回傳檔案串流 (binary)**，附帶 `Content-Disposition` header | 回傳 JSON：`{"message": "Download link generated", "download_url": "..."}` |
| 404 Not Found | ✅ `{"error": "No downloadable game file is available"}` | ❌ 無此錯誤碼 |

> 🔴 **這是最大的不相容差異**：前端整合方式完全不同。
> - `add_api_spec`：前端應觸發瀏覽器直接下載（`window.location.href` 或 `<a>` 標籤）
> - `bananaapi_spec`：前端應解析 JSON 取得 `download_url` 再跳轉

---

## 5. 社交、評論與通訊

---

### `[POST] /api/social/friends/request` (發送好友邀請)

| 項目 | add_api_spec | bananaapi_spec |
|------|-------------|----------------|
| 🔴 Request Body | 支援兩種格式：`{ "receiver_id": 2 }` 或 `{ "username": "PlayerTwo" }` | 僅支援 `{ "receiver_id": 2 }` |
| 額外錯誤碼 | ✅ `400: "Friend request already exists"` | ❌ 無 |
| 額外錯誤碼 | ✅ `400: "Cannot send a friend request to yourself"` | ❌ 無 |
| 額外錯誤碼 | ✅ `404: "User not found"` (使用 username 時) | ❌ 無 |

---

### `[POST] /api/social/blacklist` (加入黑名單)

| 項目 | add_api_spec | bananaapi_spec |
|------|-------------|----------------|
| 🔴 Request Body | 支援 `{ "blocked_id": 5 }` 或 `{ "user_id": 5 }` 兩種欄位名 | 僅支援 `{ "blocked_id": 5 }` |
| 額外錯誤碼 | ✅ `400: "Cannot blacklist yourself"` | ❌ 無 |
| 額外錯誤碼 | ✅ `400: "blocked_id is required"` | ❌ 無 |

---

## 差異摘要總覽

| # | API 端點 | 嚴重程度 | 說明 |
|---|---------|---------|------|
| 1 | `POST /api/auth/register` | 低 | add 的 201 回傳多了 `email` 欄位 |
| 2 | `POST /api/auth/login` | 低 | add 多了帳號停權的 `403` 錯誤碼 |
| 3 | `PUT /api/users/profile` | 中 | add 的 200 回傳包含 `user` 物件；banana 只有 message |
| 4 | `PUT /api/admin/users/{id}/suspend` | 低 | add 說明 toggle 行為；banana 無此說明 |
| 5 | `GET /api/developer/games` | 低 | add 說明回傳包含 media 且 ADMIN 可見全部 |
| 6 | `PUT /api/developer/games/{id}` | 低 | banana 多了 price=0 的說明；add 多了 ADMIN 可編輯的說明 |
| 7 | `POST /api/developer/games/{id}/media` | 低 | add 有完整儲存路徑；banana 多了 `400` 錯誤碼 |
| 8 | `DELETE /api/developer/games/{id}/media` | 低 | add 的 404 有兩種可能訊息 |
| 9 | `POST /api/developer/tags` | 中 | add 的 201 回傳包含 `data` 物件（tag 資料）；banana 只有 message |
| 10 | **`GET /api/protected/library/{id}/download`** | 🔴 **高** | add 回傳 binary stream；banana 回傳含 URL 的 JSON，前端整合方式完全不同 |
| 11 | `POST /api/social/friends/request` | 中 | add 支援 username 查找，多了多種 400/404 錯誤碼 |
| 12 | `POST /api/social/blacklist` | 中 | add 支援 `user_id` 替代欄位，多了 400 錯誤碼 |
