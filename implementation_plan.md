# 修復進階搜尋與開發者留言功能

本計畫將解決進階搜尋中「開發者搜尋」的問題，以及修復開發者無法對自己的遊戲進行留言的問題。

## Open Questions

> [!IMPORTANT]
> **關於「進階搜尋的開發者搜尋」**
> 您提到：「*search 進階 的 author 搜尋不適使用作者搜尋*」。
> 目前系統中，進階搜尋的「開發者」欄位是採用**模糊搜尋 (ILIKE '%keyword%')**。這表示如果輸入 `super_developer_1`，也會配對到 `super_developer_10` 或 `super_developer_120`。
> 請問您的意思是：
> 1. **改為精確配對 (Exact Match)**：輸入 `super_developer_1` 只能完全符合該名稱的開發者？
> 2. 或是希望改成**下拉式選單 (Dropdown)**，讓玩家只能選擇現有的開發者，避免輸入錯誤？
> 
> 目前的計畫中，我暫時將後端改為「精確配對（大小寫不敏感）」。若有其他需求請告訴我！

## Proposed Changes

### Frontend (前端)

#### [MODIFY] [game_detail.html](file:///c:/Users/User/Desktop/workspace/dbms/frontend/pages/store/game_detail.html)
- **問題原因**：前端在判斷「當前使用者是否為該遊戲的開發者」時，使用了 `user.user_id == devId`。但登入時回傳的使用者物件中，ID 欄位名稱是 `id` 而非 `user_id`，導致判斷永遠為 false，開發者因此看不到「發表評論」的區塊與「AUTHOR」身分選項。
- **修改內容**：將 `user.user_id` 的判斷全部替換為 `(user.id || user.user_id)`，確保能正確抓到開發者的 ID。

### Backend (後端)

#### [MODIFY] [game_controller.go](file:///c:/Users/User/Desktop/workspace/dbms/backend/controllers/game_controller.go)
- **修改內容**：在 `GetGames` 函式中，針對 `developer` 的進階搜尋條件，將原本的模糊搜尋 `%+developer+%` 改為精確搜尋（或根據您的回覆進一步調整）。

## Verification Plan

### Manual Verification
1. **開發者留言功能測試**：
   - 使用開發者帳號登入。
   - 進入自己開發的遊戲頁面，確認即使尚未購買該遊戲，也能看見留言區塊。
   - 確認可以選擇 `AUTHOR` 身分發表評論與回覆。
2. **進階搜尋測試**：
   - 在進階搜尋中輸入開發者名稱（例如 `super_developer_1`），確認搜尋結果不會再出現其他名稱包含該字串（例如 `super_developer_10`）的開發者遊戲。
