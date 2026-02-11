const API_BASE = window.location.origin;

function formatTime(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diff = Math.floor((now - date) / 1000);

    if (diff < 60) return 'Just now';
    if (diff < 3600) return `${Math.floor(diff / 60)} minutes ago`;
    if (diff < 86400) return `${Math.floor(diff / 3600)} hours ago`;
    if (diff < 604800) return `${Math.floor(diff / 86400)} days ago`;

    return date.toLocaleDateString();
}

function truncateText(text, maxLength = 100) {
    if (text.length <= maxLength) return text;
    return text.substring(0, maxLength) + '...';
}

function showNotification(message, type = 'info') {
    const existing = document.querySelector('.notification');
    if (existing) existing.remove();

    const notification = document.createElement('div');
    notification.className = `notification ${type}`;
    notification.innerHTML = `
        <i class="fas fa-${getNotificationIcon(type)}"></i>
        <span>${message}</span>
    `;

    document.body.appendChild(notification);

    setTimeout(() => {
        notification.style.opacity = '0';
        notification.style.transform = 'translateX(100%)';
        setTimeout(() => notification.remove(), 300);
    }, 3000);
}

function getNotificationIcon(type) {
    switch(type) {
        case 'success': return 'check-circle';
        case 'error': return 'exclamation-circle';
        case 'warning': return 'exclamation-triangle';
        default: return 'info-circle';
    }
}

function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

function validateEmail(email) {
    const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return re.test(email);
}

function validatePassword(password) {
    return password.length >= 8;
}

function getAuthToken() {
    return localStorage.getItem('token');
}

function setAuthToken(token) {
    localStorage.setItem('token', token);
}

function removeAuthToken() {
    localStorage.removeItem('token');
}

function getUserRole() {
    const user = JSON.parse(localStorage.getItem('user') || '{}');
    return user.role || null;
}

function isAdmin() {
    return getUserRole() === 'admin';
}

function isModerator() {
    const role = getUserRole();
    return role === 'admin' || role === 'moderator';
}

async function fetchWithAuth(url, options = {}) {
    const token = getAuthToken();
    const headers = {
        'Content-Type': 'application/json',
        ...options.headers
    };

    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE}${url}`, {
        ...options,
        headers
    });

    if (response.status === 401) {
        removeAuthToken();
        localStorage.removeItem('user');
        window.location.href = '/login.html';
        throw new Error('Authentication required');
    }

    if (response.status === 403) {
        showNotification('Permission denied', 'error');
        throw new Error('Permission denied');
    }

    return response;
}

async function fetchJSON(url, options = {}) {
    const response = await fetchWithAuth(url, options);
    return response.json();
}

function createLoadingSpinner() {
    const spinner = document.createElement('div');
    spinner.className = 'spinner';
    return spinner;
}

function toggleLoading(element, isLoading) {
    if (isLoading) {
        element.classList.add('loading');
        element.appendChild(createLoadingSpinner());
    } else {
        element.classList.remove('loading');
        const spinner = element.querySelector('.spinner');
        if (spinner) spinner.remove();
    }
}

function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

function getFileTypeIcon(fileType) {
    if (fileType.includes('image')) return 'fa-image';
    if (fileType.includes('video')) return 'fa-video';
    if (fileType.includes('pdf')) return 'fa-file-pdf';
    if (fileType.includes('word') || fileType.includes('document')) return 'fa-file-word';
    return 'fa-file';
}

const CATEGORY_LABELS = {
    'meme': 'Memes',
    'event': 'Events',
    'news': 'News',
    'question': 'Questions',
    'lost_found': 'Lost & Found',
    'academic': 'Academic',
    'social': 'Social',
    'sports': 'Sports'
};

function getCategoryLabel(category) {
    return CATEGORY_LABELS[category] || category;
}

const ROLE_LABELS = {
    'student': 'Student',
    'admin': 'Administrator',
    'alumni': 'Alumni',
    'moderator': 'Moderator'
};

function getRoleLabel(role) {
    return ROLE_LABELS[role] || role;
}

// Color mapping for roles
const ROLE_COLORS = {
    'student': '#4a6fa5',
    'admin': '#e74c3c',
    'alumni': '#2ecc71',
    'moderator': '#f39c12'
};

function getRoleColor(role) {
    return ROLE_COLORS[role] || '#6c757d';
}