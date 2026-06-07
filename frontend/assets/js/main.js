// ============================================================
// VaporAuror — main.js
// 全域狀態、假資料、共用元件
// 所有後端 API 呼叫預留點皆以 // [API] 標記
// ============================================================

// 移除 mockGames，全數改為 API 介接

// ============================================================
// 頁面初始化
// ============================================================
document.addEventListener("DOMContentLoaded", () => {
    renderHeader();

    // 搜尋框事件（首頁）
    const searchForm = document.getElementById('search-form');
    if (searchForm) {
        searchForm.addEventListener('submit', (e) => {
            e.preventDefault();
            const keyword = document.getElementById('search-input').value.trim();
            if (keyword) {
                // [API] GET /api/games?q={keyword}
                window.location.href = `/pages/store/search.html?q=${encodeURIComponent(keyword)}`;
            }
        });
    }
});

// ============================================================
// 渲染共用導覽列
// ============================================================
function renderHeader() {
    const header = document.getElementById('global-header');
    if (!header) return;

    // 登入頁與註冊頁只顯示 Logo
    const path = window.location.pathname;
    if (path.includes('login.html') || path.includes('register.html')) {
        header.innerHTML = `<a href="/" class="logo">VaporAuror</a>`;
        return;
    }

    // [API] 未來由後端 session/token 判斷，目前由 localStorage 模擬
    const currentRole = localStorage.getItem('userRole') || 'GUEST';
    const userDataStr = localStorage.getItem('currentUser');
    const username = userDataStr ? JSON.parse(userDataStr).username : '玩家';

    // 左側選單
    let leftHtml = `
        <div class="nav-left">
            <a href="/" class="logo">VaporAuror</a>
            <ul class="main-menu">
                <li><a href="/">商店首頁</a></li>
    `;

    if (currentRole !== 'GUEST') {
        leftHtml += `
                <li><a href="/pages/user/library">遊戲庫</a></li>
                <li><a href="/pages/user/social">社群</a></li>
        `;
        if (currentRole === 'DEVELOPER') {
            leftHtml += `<li><a href="/pages/dashboard/dev_dashboard" style="color:#e5a93d;">開發者中心</a></li>`;
        } else if (currentRole === 'CSR') {
            leftHtml += `<li><a href="/pages/dashboard/csr_dashboard" style="color:#66c0f4;">客服中心</a></li>`;
        } else if (currentRole === 'ADMIN') {
            leftHtml += `<li><a href="/pages/dashboard/csr_dashboard" style="color:#66c0f4;">客服中心</a></li>`;
            leftHtml += `<li><a href="/pages/dashboard/admin_dashboard" style="color:#ff5959;">管理後台</a></li>`;
        }
    }
    leftHtml += `</ul></div>`;

    // 右側
    let rightHtml = `<div class="nav-right">`;
    if (currentRole === 'GUEST') {
        rightHtml += `<a href="/pages/auth/login" class="btn-primary" style="margin:0;padding:8px 22px;width:auto;font-size:14px;">登入 / 註冊</a>`;
    } else {
        rightHtml += `
            <a href="/pages/user/cart" class="nav-cart">購物車</a>
            <div class="user-dropdown">
                <button class="user-btn">${escapeHtml(username)}</button>
                <div class="dropdown-content">
                    <a href="/pages/user/profile">基本資料</a>
                    <a href="/pages/user/history">購買紀錄</a>
                    <a href="/pages/user/wishlist">願望清單</a>
                    <a href="#" id="logout-btn">登出</a>
                </div>
            </div>
        `;
    }
    rightHtml += `</div>`;

    header.innerHTML = leftHtml + rightHtml;

    // 登出事件
    const logoutBtn = document.getElementById('logout-btn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', async (e) => {
            e.preventDefault();
            try {
                await apiLogout();
            } catch (_) {
                // 即使 API 失敗也清除本地狀態
                localStorage.removeItem('token');
                localStorage.removeItem('currentUser');
                localStorage.removeItem('userRole');
            }
            window.location.href = '/';
        });
    }
}

// ============================================================
// 渲染遊戲卡片列表
// ============================================================
function renderGames(games) {
    const container = document.getElementById('game-list');
    if (!container) return;

    if (!games) {
        // Skeleton Loading 骨架屏狀態
        container.innerHTML = '';
        for (let i = 0; i < 4; i++) {
            const skeleton = document.createElement('div');
            skeleton.className = 'game-list-card';
            skeleton.innerHTML = `
                <div class="game-thumbnail skeleton"></div>
                <div class="game-list-info">
                    <div class="skeleton skeleton-title"></div>
                    <div class="skeleton skeleton-text"></div>
                    <div class="skeleton skeleton-text" style="width: 80%;"></div>
                    <div class="game-list-tags" style="margin-top: 10px;">
                        <span class="tag skeleton" style="width: 40px; height: 16px;"></span>
                        <span class="tag skeleton" style="width: 50px; height: 16px;"></span>
                    </div>
                </div>
            `;
            container.appendChild(skeleton);
        }
        return;
    }

    if (games.length === 0) {
        container.innerHTML = `
            <div class="empty-state">
                <div class="empty-icon">🎮</div>
                <p>找不到符合條件的遊戲</p>
                <a href="/" class="btn-secondary" style="display:inline-block;margin-top:10px;">回首頁</a>
            </div>`;
        return;
    }

    container.innerHTML = '';
    games.forEach(game => {
        const card = document.createElement('div');
        card.className = 'game-list-card';

        const tagsHtml = (game.tags || []).map(t => {
            const tagName = typeof t === 'string' ? t : (t.tag_name || t.name || t);
            return `<span class="tag">${escapeHtml(String(tagName))}</span>`;
        }).join('');
        const priceHtml = game.price === 0
            ? `<span class="game-list-price free">免費</span>`
            : `<span class="game-list-price">NT$ ${game.price.toLocaleString()}</span>`;

        let coverHtml = `<span style="font-size:13px;color:#8f98a0;">[${escapeHtml(game.title)} 封面]</span>`;
        if (game.media && game.media.length > 0) {
            const thumb = game.media.find(m => m.media_type === 'thumbnail') || game.media.find(m => m.media_type === 'media');
            if (thumb) {
                coverHtml = `<img src="${thumb.file_url}" alt="cover" style="width:100%; height:100%; object-fit:cover; border-radius: 8px;">`;
            }
        }

        card.innerHTML = `
            <div class="game-thumbnail">
                ${coverHtml}
            </div>
            <div class="game-list-info">
                <div class="game-list-title">${escapeHtml(game.title)}</div>
                <div class="game-list-desc">${escapeHtml(game.desc)}</div>
                <div class="game-list-tags">${tagsHtml}</div>
                ${priceHtml}
            </div>
        `;

        // [API] 使用真實 id 導向
        card.addEventListener('click', () => {
            const gameId = game.game_id || game.id;
            window.location.href = `/pages/store/game_detail?id=${gameId}`;
        });

        container.appendChild(card);
    });
}

// ============================================================
// 工具函式
// ============================================================
function escapeHtml(str) {
    if (typeof str !== 'string') return '';
    return str
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');
}

function getCurrentUser() {
    // [API] 未來由後端 session 提供
    const str = localStorage.getItem('currentUser');
    return str ? JSON.parse(str) : null;
}

function getCurrentRole() {
    // [API] 未來由後端 JWT Token 解碼
    return localStorage.getItem('userRole') || 'GUEST';
}

function requireLogin(redirectTo = '/pages/auth/login.html') {
    if (!localStorage.getItem('token')) {
        alert('請先登入！');
        window.location.href = redirectTo;
        return false;
    }
    return true;
}

// ============================================================
// 全域 Toast Notification 系統
// type: 'info', 'success', 'error'
// ============================================================
function showToast(message, type = 'info') {
    let container = document.getElementById('toast-container');
    if (!container) {
        container = document.createElement('div');
        container.id = 'toast-container';
        document.body.appendChild(container);
    }

    const toast = document.createElement('div');
    toast.className = `toast-message toast-${type}`;
    
    // Icon selection
    let icon = '💬';
    if (type === 'success') icon = '✅';
    if (type === 'error') icon = '❌';

    toast.innerHTML = `<span>${icon}</span> <span>${escapeHtml(message)}</span>`;
    container.appendChild(toast);

    // Animate in
    setTimeout(() => toast.classList.add('show'), 10);

    // Auto remove after 3s
    setTimeout(() => {
        toast.classList.remove('show');
        setTimeout(() => toast.remove(), 400); // Wait for transition
    }, 3000);
}
