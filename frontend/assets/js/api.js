// ============================================================
// api.js — VaporAuror 集中式 API 模組
// 所有後端呼叫都在此統一管理
// ============================================================

const API_BASE = '';

// ============================================================
// 核心工具：帶 JWT Token 的 fetch 包裝
// ============================================================
async function authFetch(url, options = {}) {
    const token = localStorage.getItem('token');
    const headers = {
        'Content-Type': 'application/json',
        'Cache-Control': 'no-cache',
        ...(options.headers || {})
    };
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }
    const fetchOptions = { ...options, headers, cache: 'no-store' };
    const res = await fetch(`${API_BASE}${url}`, fetchOptions);

    // 401 → token 失效，踢回登入頁
    if (res.status === 401) {
        localStorage.removeItem('token');
        localStorage.removeItem('currentUser');
        localStorage.removeItem('userRole');
        alert('登入已過期，請重新登入。');
        window.location.href = '/pages/auth/login.html';
        throw new Error('Unauthorized');
    }

    return res;
}

// 解析回應，若非 2xx 則拋出含訊息的錯誤
async function parseResponse(res) {
    let data;
    const contentType = res.headers.get('content-type') || '';
    if (contentType.includes('application/json')) {
        data = await res.json();
    } else {
        data = { message: await res.text() };
    }
    if (!res.ok) {
        throw new Error(data.message || data.error || `HTTP ${res.status}`);
    }
    return data;
}

// ============================================================
// 1. 認證 (Auth)
// ============================================================

// POST /api/auth/login
async function apiLogin(email, password) {
    const res = await fetch(`${API_BASE}/api/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
    });
    return parseResponse(res);
}

// POST /api/auth/register
async function apiRegister(username, email, password, isDeveloper = false) {
    const res = await fetch(`${API_BASE}/api/auth/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, email, password, is_developer: isDeveloper })
    });
    return parseResponse(res);
}

// POST /api/auth/logout
async function apiLogout() {
    try {
        const res = await authFetch('/api/auth/logout', { method: 'POST' });
        await parseResponse(res);
    } finally {
        localStorage.removeItem('token');
        localStorage.removeItem('currentUser');
        localStorage.removeItem('userRole');
    }
}

// PUT /api/users/profile
async function apiUpdateProfile(data) {
    const res = await authFetch('/api/users/profile', {
        method: 'PUT',
        body: JSON.stringify(data)
    });
    const result = await parseResponse(res);
    return result.user || result;
}

// ============================================================
// 2. 商店與遊戲 (Games)
// ============================================================

// GET /api/games?q=keyword&tag=tag&developer=name&min_price=0&max_price=1000&sort=popular
async function apiGetGames(query = '') {
    const params = new URLSearchParams();
    if (typeof query === 'string') {
        if (query) params.set('q', query);
    } else if (query && typeof query === 'object') {
        Object.entries(query).forEach(([key, value]) => {
            if (value !== undefined && value !== null && String(value).trim() !== '') {
                params.set(key, value);
            }
        });
    }
    const qs = params.toString();
    const url = qs ? `/api/games?${qs}` : '/api/games';
    const options = {};
    const token = localStorage.getItem('token');
    if (token && params.get('hide_owned') === 'true') {
        options.headers = { 'Authorization': `Bearer ${token}` };
    }
    const res = await fetch(`${API_BASE}${url}`, options);
    return parseResponse(res);
}

async function apiGetTags() {
    const res = await fetch(`${API_BASE}/api/tags`);
    return parseResponse(res);
}

// GET /api/games/{id}
async function apiGetGame(id) {
    const res = await fetch(`${API_BASE}/api/games/${id}`);
    return parseResponse(res);
}

// GET /api/games/{id}/reviews  (透過 game detail 一起取得或單獨端點)
async function apiGetReviews(gameId) {
    const res = await fetch(`${API_BASE}/api/games/${gameId}/reviews`);
    return parseResponse(res);
}

// POST /api/social/games/{id}/reviews
async function apiPostReview(gameId, attitude, content, postAsRole = 'USERS') {
    const res = await authFetch(`/api/social/games/${gameId}/reviews`, {
        method: 'POST',
        body: JSON.stringify({ attitude: attitude.toUpperCase(), content, post_as_role: postAsRole })
    });
    return parseResponse(res);
}

// POST /api/social/reviews/{id}/replies
async function apiPostReply(reviewId, content, postAsRole = 'USERS') {
    const res = await authFetch(`/api/social/reviews/${reviewId}/replies`, {
        method: 'POST',
        body: JSON.stringify({ content, post_as_role: postAsRole })
    });
    return parseResponse(res);
}

// POST /api/developer/games  [DEVELOPER]
async function apiGetDeveloperGames() {
    const res = await authFetch('/api/developer/games');
    return parseResponse(res);
}

async function apiCreateGame(title, price, desc) {
    const res = await authFetch('/api/developer/games', {
        method: 'POST',
        body: JSON.stringify({ title, price, desc })
    });
    return parseResponse(res);
}

async function apiUpdateGame(id, price, desc) {
    const res = await authFetch(`/api/developer/games/${id}`, {
        method: 'PUT',
        body: JSON.stringify({ price, desc })
    });
    return parseResponse(res);
}

async function apiCreateTag(tagName) {
    const res = await authFetch('/api/developer/tags', {
        method: 'POST',
        body: JSON.stringify({ tag_name: tagName })
    });
    return parseResponse(res);
}

async function apiAddTagToGame(gameId, tagId) {
    const res = await authFetch(`/api/developer/games/${gameId}/tags`, {
        method: 'POST',
        body: JSON.stringify({ tag_id: tagId })
    });
    return parseResponse(res);
}

async function apiRemoveTagFromGame(gameId, tagId) {
    const res = await authFetch(`/api/developer/games/${gameId}/tags/${tagId}`, { method: 'DELETE' });
    return parseResponse(res);
}

// DELETE /api/developer/games/{id}  [DEVELOPER]
async function apiDeleteGame(id) {
    const res = await authFetch(`/api/developer/games/${id}`, { method: 'DELETE' });
    return parseResponse(res);
}

// DELETE /api/admin/games/{id} [ADMIN]
async function apiAdminDeleteGame(id) {
    const res = await authFetch(`/api/admin/games/${id}`, { method: 'DELETE' });
    return parseResponse(res);
}

// GET /api/developer/games/{id}/stats  [DEVELOPER]
async function apiGetGameStats(id) {
    const res = await authFetch(`/api/developer/games/${id}/stats`);
    const result = await parseResponse(res);
    return result.stats || result;
}

// ============================================================
// 3. 購物車 (Cart)
// ============================================================

// GET /api/protected/cart
async function apiGetCart() {
    const res = await authFetch('/api/protected/cart');
    return parseResponse(res);
}

// POST /api/protected/cart  { game_id }
async function apiAddToCart(gameId) {
    const res = await authFetch('/api/protected/cart', {
        method: 'POST',
        body: JSON.stringify({ game_id: gameId })
    });
    return parseResponse(res);
}

// DELETE /api/protected/cart/{game_id}
async function apiRemoveFromCart(gameId) {
    const res = await authFetch(`/api/protected/cart/${gameId}`, { method: 'DELETE' });
    return parseResponse(res);
}

// POST /api/protected/checkout
async function apiCheckout() {
    const res = await authFetch('/api/protected/checkout', { method: 'POST' });
    return parseResponse(res);
}

// ============================================================
// 4. 購買紀錄與退款 (Transactions & Refunds)
// ============================================================

// GET /api/protected/transactions
async function apiGetTransactions() {
    const res = await authFetch('/api/protected/transactions');
    return parseResponse(res);
}

// POST /api/social/refunds  { transaction_item_id, reason }
async function apiRequestRefund(transactionItemId, reason) {
    const res = await authFetch('/api/social/refunds', {
        method: 'POST',
        body: JSON.stringify({ transaction_item_id: transactionItemId, reason: reason })
    });
    return parseResponse(res);
}

// GET /api/protected/refunds (取得個人退款歷史)
async function apiGetMyRefunds() {
    const res = await authFetch('/api/protected/refunds');
    return parseResponse(res);
}

// GET /api/csr/refunds  [CSR]
async function apiGetRefunds() {
    const res = await authFetch('/api/csr/refunds');
    return parseResponse(res);
}

// PUT /api/csr/refunds/{id}  [CSR]
async function apiProcessRefund(refundId, status, rejectReason = "") {
    const res = await authFetch(`/api/csr/refunds/${refundId}`, {
        method: 'PUT',
        body: JSON.stringify({ status, reject_reason: rejectReason })
    });
    return parseResponse(res);
}

// ============================================================
// 5. 遊戲庫與願望清單 (Library & Wishlist)
// ============================================================

// GET /api/protected/library
async function apiGetLibrary() {
    const res = await authFetch('/api/protected/library');
    return parseResponse(res);
}

// GET /api/protected/wishlist
async function apiGetWishlist() {
    const res = await authFetch('/api/protected/wishlist');
    return parseResponse(res);
}

// POST /api/protected/wishlist  { game_id }
async function apiAddWishlist(gameId) {
    const res = await authFetch('/api/protected/wishlist', {
        method: 'POST',
        body: JSON.stringify({ game_id: gameId })
    });
    return parseResponse(res);
}

// DELETE /api/protected/wishlist/{game_id}
async function apiRemoveWishlist(gameId) {
    const res = await authFetch(`/api/protected/wishlist/${gameId}`, { method: 'DELETE' });
    return parseResponse(res);
}

// ============================================================
// 6. 社交 (Friends, Blacklist, Messages)
// ============================================================

// GET /api/social/friends
async function apiGetFriends() {
    const res = await authFetch('/api/social/friends');
    return parseResponse(res);
}

// GET /api/social/friends/requests
async function apiGetFriendRequests() {
    const res = await authFetch('/api/social/friends/requests');
    return parseResponse(res);
}

// GET /api/social/blacklist
async function apiGetBlacklist() {
    const res = await authFetch('/api/social/blacklist');
    return parseResponse(res);
}

// POST /api/social/friends/request  { username }
async function apiFriendRequest(username) {
    const res = await authFetch('/api/social/friends/request', {
        method: 'POST',
        body: JSON.stringify({ username })
    });
    return parseResponse(res);
}

// PUT /api/social/friends/request/{id}/accept
async function apiFriendAccept(requestId) {
    const res = await authFetch(`/api/social/friends/request/${requestId}/accept`, { method: 'PUT' });
    return parseResponse(res);
}

// PUT /api/social/friends/request/{id}/decline
async function apiFriendDecline(requestId) {
    const res = await authFetch(`/api/social/friends/request/${requestId}/decline`, { method: 'PUT' });
    return parseResponse(res);
}

// DELETE /api/social/friends/request/{id}  (收回邀請)
async function apiFriendCancelRequest(requestId) {
    const res = await authFetch(`/api/social/friends/request/${requestId}`, { method: 'DELETE' });
    return parseResponse(res);
}

// POST /api/social/blacklist  { user_id }
async function apiAddBlacklist(userId) {
    const res = await authFetch('/api/social/blacklist', {
        method: 'POST',
        body: JSON.stringify({ blocked_id: userId })
    });
    return parseResponse(res);
}

// DELETE /api/social/blacklist/{user_id}
async function apiRemoveBlacklist(userId) {
    const res = await authFetch(`/api/social/blacklist/${userId}`, { method: 'DELETE' });
    return parseResponse(res);
}

// GET /api/social/messages/{user_id}
async function apiGetMessages(userId) {
    const res = await authFetch(`/api/social/messages/${userId}`);
    return parseResponse(res);
}

// POST /api/social/messages  { receiver_id, content }
async function apiSendMessage(receiverId, content) {
    const res = await authFetch('/api/social/messages', {
        method: 'POST',
        body: JSON.stringify({ receiver_id: receiverId, content })
    });
    return parseResponse(res);
}

// ============================================================
// 7. 管理員 (Admin)
// ============================================================

// GET /api/admin/users  [ADMIN] — API 文件未列但實務上需要
async function apiGetUsers() {
    const res = await authFetch('/api/admin/users');
    return parseResponse(res);
}

// PUT /api/admin/users/{id}/suspend  [ADMIN]
async function apiSuspendUser(userId) {
    const res = await authFetch(`/api/admin/users/${userId}/suspend`, { method: 'PUT' });
    return parseResponse(res);
}

// DELETE /api/admin/users/{id}  [ADMIN]
async function apiDeleteUser(userId) {
    const res = await authFetch(`/api/admin/users/${userId}`, { method: 'DELETE' });
    return parseResponse(res);
}

// PUT /api/admin/users/{id}/role  [ADMIN]  { role }
async function apiChangeRole(userId, role) {
    const res = await authFetch(`/api/admin/users/${userId}/role`, {
        method: 'PUT',
        body: JSON.stringify({ role })
    });
    return parseResponse(res);
}

// ============================================================
// 工具函式：存入登入狀態
// ============================================================
function saveAuthState(data) {
    // 後端預期回傳: { token, user: { id, username, email, role, ... } }
    localStorage.setItem('token', data.token);
    localStorage.setItem('currentUser', JSON.stringify(data.user));
    localStorage.setItem('userRole', data.user.role);
}
