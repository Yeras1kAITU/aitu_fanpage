class PostManager {
    constructor() {
        this.currentPost = null;
    }

    async createPost(postData, files = []) {
        const token = localStorage.getItem('token');
        if (!token) {
            throw new Error('Authentication required');
        }

        // If no files, use simple JSON
        if (files.length === 0) {
            const response = await fetch('/api/posts', {
                method: 'POST',
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(postData)
            });

            if (!response.ok) {
                const error = await response.text();
                throw new Error(error || 'Failed to create post');
            }

            const post = await response.json();
            showNotification('Post created successfully!', 'success');
            return post;
        }

        // If files exist, use FormData
        const formData = new FormData();
        formData.append('post', JSON.stringify(postData));

        files.forEach(file => {
            formData.append('files', file);
        });

        const response = await fetch('/api/posts', {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`
            },
            body: formData
        });

        if (!response.ok) {
            const error = await response.text();
            throw new Error(error || 'Failed to create post with files');
        }

        const post = await response.json();
        showNotification('Post created successfully!', 'success');
        return post;
    }

    async getPosts(params = {}) {
        const queryParams = new URLSearchParams({
            limit: params.limit || 10,
            offset: params.offset || 0,
            ...params
        }).toString();

        try {
            const response = await fetchWithAuth(`/api/posts?${queryParams}`);
            return response.json();
        } catch (error) {
            console.error('Get posts error:', error);
            throw error;
        }
    }

    async getPostById(id) {
        try {
            const response = await fetchWithAuth(`/api/posts/${id}`);
            return response.json();
        } catch (error) {
            console.error('Get post error:', error);
            throw error;
        }
    }

    async updatePost(id, postData) {
        try {
            const response = await fetchWithAuth(`/api/posts/${id}`, {
                method: 'PUT',
                body: JSON.stringify(postData)
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.message || 'Failed to update post');
            }

            const post = await response.json();
            showNotification('Post updated successfully!', 'success');
            return post;
        } catch (error) {
            console.error('Update post error:', error);
            throw error;
        }
    }

    async deletePost(id) {
        try {
            const response = await fetchWithAuth(`/api/posts/${id}`, {
                method: 'DELETE'
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.message || 'Failed to delete post');
            }

            showNotification('Post deleted successfully!', 'success');
            return true;
        } catch (error) {
            console.error('Delete post error:', error);
            throw error;
        }
    }

    async likePost(id) {
        try {
            const response = await fetchWithAuth(`/api/posts/${id}/like`, {
                method: 'POST'
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || `Failed to like post: ${response.status}`);
            }

            showNotification('Post liked!', 'success');
            return true;
        } catch (error) {
            console.error('Like post error:', error);

            if (error.message.includes('rate limit exceeded')) {
                showNotification(error.message, 'warning');
            } else {
                showNotification(error.message || 'Failed to like post', 'error');
            }
            throw error;
        }
    }

    async unlikePost(id) {
        try {
            const response = await fetchWithAuth(`/api/posts/${id}/like`, {
                method: 'DELETE'
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || `Failed to unlike post: ${response.status}`);
            }

            showNotification('Post unliked', 'info');
            return true;
        } catch (error) {
            console.error('Unlike post error:', error);
            showNotification(error.message || 'Failed to unlike post', 'error');
            throw error;
        }
    }

    async pinPost(id) {
        try {
            const response = await fetchWithAuth(`/api/posts/${id}/pin`, {
                method: 'POST'
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.message || 'Failed to pin post');
            }

            showNotification('Post pinned!', 'success');
            return true;
        } catch (error) {
            console.error('Pin post error:', error);
            throw error;
        }
    }

    async unpinPost(id) {
        try {
            const response = await fetchWithAuth(`/api/posts/${id}/pin`, {
                method: 'DELETE'
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.message || 'Failed to unpin post');
            }

            showNotification('Post unpinned', 'info');
            return true;
        } catch (error) {
            console.error('Unpin post error:', error);
            throw error;
        }
    }

    async featurePost(id) {
        try {
            const response = await fetchWithAuth(`/api/posts/${id}/feature`, {
                method: 'POST'
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.message || 'Failed to feature post');
            }

            showNotification('Post featured!', 'success');
            return true;
        } catch (error) {
            console.error('Feature post error:', error);
            throw error;
        }
    }

    async unfeaturePost(id) {
        try {
            const response = await fetchWithAuth(`/api/posts/${id}/feature`, {
                method: 'DELETE'
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.message || 'Failed to unfeature post');
            }

            showNotification('Post unfeatured', 'info');
            return true;
        } catch (error) {
            console.error('Unfeature post error:', error);
            throw error;
        }
    }

    async searchPosts(query, limit = 10) {
        try {
            const response = await fetchWithAuth(`/api/posts/search?q=${encodeURIComponent(query)}&limit=${limit}`);
            return response.json();
        } catch (error) {
            console.error('Search posts error:', error);
            throw error;
        }
    }

    async getPostsByCategory(category, limit = 10) {
        try {
            const response = await fetchWithAuth(`/api/posts?category=${category}&limit=${limit}`);
            return response.json();
        } catch (error) {
            console.error('Get posts by category error:', error);
            throw error;
        }
    }

    async getPinnedPosts(limit = 5) {
        try {
            const response = await fetchWithAuth(`/api/posts/pinned?limit=${limit}`);
            return response.json();
        } catch (error) {
            console.error('Get pinned posts error:', error);
            throw error;
        }
    }

    async getFeaturedPosts(limit = 10) {
        try {
            const response = await fetchWithAuth(`/api/posts/featured?limit=${limit}`);
            return response.json();
        } catch (error) {
            console.error('Get featured posts error:', error);
            throw error;
        }
    }

    async getPopularPosts(limit = 10, days = 7) {
        try {
            const response = await fetchWithAuth(`/api/posts/popular?limit=${limit}&days=${days}`);
            return response.json();
        } catch (error) {
            console.error('Get popular posts error:', error);
            throw error;
        }
    }

    async getCategoriesStats() {
        try {
            const response = await fetchWithAuth('/api/posts/categories/stats');
            return response.json();
        } catch (error) {
            console.error('Get categories stats error:', error);
            throw error;
        }
    }

    async getFeed(category, limit = 10, offset = 0) {
        try {
            const url = category
                ? `/api/posts/feed?category=${category}&limit=${limit}&offset=${offset}`
                : `/api/posts/feed?limit=${limit}&offset=${offset}`;

            const response = await fetchWithAuth(url);
            return response.json();
        } catch (error) {
            console.error('Get feed error:', error);
            throw error;
        }
    }

    renderPost(post) {
        return `
            <div class="post-card" data-post-id="${post.id}">
                <div class="post-header">
                    <div class="post-author">
                        <img src="${post.author_profile_image || '/assets/default-avatar.svg'}" 
                             alt="${post.author_name}" class="author-avatar">
                        <div class="author-info">
                            <h4>${post.author_name}</h4>
                            <span class="post-time">${formatTime(post.created_at)}</span>
                        </div>
                    </div>
                    <div class="post-badges">
                        ${post.is_pinned ? '<span class="badge pinned"><i class="fas fa-thumbtack"></i> Pinned</span>' : ''}
                        ${post.is_featured ? '<span class="badge featured"><i class="fas fa-star"></i> Featured</span>' : ''}
                        ${post.category ? `<span class="badge category">${getCategoryLabel(post.category)}</span>` : ''}
                    </div>
                </div>
                
                <div class="post-content">
                    <h3 class="post-title">${post.title || 'Untitled Post'}</h3>
                    ${post.description ? `<p class="post-description">${post.description}</p>` : ''}
                    ${post.content ? `<div class="post-body">${post.content}</div>` : ''}
                </div>
                
                ${this.renderMedia(post.media)}
                
                <div class="post-footer">
                    <div class="post-stats">
                        <span class="stat"><i class="fas fa-heart"></i> ${post.like_count || 0}</span>
                        <span class="stat"><i class="fas fa-comment"></i> ${post.comment_count || 0}</span>
                        <span class="stat"><i class="fas fa-eye"></i> ${post.view_count || 0}</span>
                    </div>
                    
                    <div class="post-actions">
                        <button class="btn-icon like-btn" data-liked="${post.user_liked || false}" data-post-id="${post.id}">
                            <i class="fas ${post.user_liked ? 'fa-heart text-danger' : 'fa-heart'}"></i> 
                            ${post.user_liked ? 'Liked' : 'Like'}
                        </button>
                        <a href="post-detail.html?id=${post.id}" class="btn-icon">
                            <i class="fas fa-comment"></i> Comment
                        </a>
                        ${this.renderModerationButtons(post)}
                        <button class="btn-icon share-btn">
                            <i class="fas fa-share"></i> Share
                        </button>
                    </div>
                </div>
            </div>
        `;
    }

    renderMedia(media) {
        if (!media || media.length === 0) return '';

        if (media.length === 1) {
            const item = media[0];
            return `
                <div class="post-media">
                    <img src="${item.url}" alt="${item.caption || 'Post media'}" class="media-preview single">
                    ${item.caption ? `<p class="media-caption">${item.caption}</p>` : ''}
                </div>
            `;
        }

        return `
            <div class="post-media grid">
                ${media.slice(0, 4).map((item, index) => `
                    <div class="media-item ${index === 3 && media.length > 4 ? 'with-overlay' : ''}">
                        <img src="${item.url}" alt="${item.caption || `Media ${index + 1}`}">
                        ${index === 3 && media.length > 4 ?
            `<div class="media-overlay">+${media.length - 3}</div>` : ''}
                    </div>
                `).join('')}
            </div>
        `;
    }

    renderModerationButtons(post) {
        const user = authManager.getUser();
        if (!user || !authManager.can('manage_posts')) return '';

        return `
            <div class="moderation-dropdown">
                <button class="btn-icon mod-btn">
                    <i class="fas fa-ellipsis-h"></i> Mod
                </button>
                <div class="moderation-menu">
                    ${!post.is_pinned ?
            `<button class="mod-action" data-action="pin" data-post-id="${post.id}">
                            <i class="fas fa-thumbtack"></i> Pin Post
                        </button>` :
            `<button class="mod-action" data-action="unpin" data-post-id="${post.id}">
                            <i class="fas fa-thumbtack"></i> Unpin Post
                        </button>`
        }
                    ${!post.is_featured ?
            `<button class="mod-action" data-action="feature" data-post-id="${post.id}">
                            <i class="fas fa-star"></i> Feature Post
                        </button>` :
            `<button class="mod-action" data-action="unfeature" data-post-id="${post.id}">
                            <i class="fas fa-star"></i> Unfeature Post
                        </button>`
        }
                    <hr>
                    <button class="mod-action text-danger" data-action="delete" data-post-id="${post.id}">
                        <i class="fas fa-trash"></i> Delete Post
                    </button>
                </div>
            </div>
        `;
    }
}

const postManager = new PostManager();

function initPostEventListeners() {
    document.addEventListener('click', async (e) => {
        if (e.target.closest('.like-btn')) {
            const likeBtn = e.target.closest('.like-btn');
            const postId = likeBtn.dataset.postId;
            const isLiked = likeBtn.dataset.liked === 'true';

            try {
                if (!authManager.isAuthenticated()) {
                    showNotification('Please login to like posts', 'warning');
                    window.location.href = 'login.html';
                    return;
                }

                if (isLiked) {
                    await postManager.unlikePost(postId);
                    likeBtn.dataset.liked = 'false';
                    likeBtn.innerHTML = '<i class="fas fa-heart"></i> Like';

                    // Update like count
                    const likeCount = likeBtn.closest('.post-footer').querySelector('.fa-heart').parentElement;
                    const currentCount = parseInt(likeCount.textContent) || 0;
                    if (currentCount > 0) {
                        likeCount.textContent = currentCount - 1;
                    }
                } else {
                    await postManager.likePost(postId);
                    likeBtn.dataset.liked = 'true';
                    likeBtn.innerHTML = '<i class="fas fa-heart text-danger"></i> Liked';

                    // Update like count
                    const likeCount = likeBtn.closest('.post-footer').querySelector('.fa-heart').parentElement;
                    const currentCount = parseInt(likeCount.textContent) || 0;
                    likeCount.textContent = currentCount + 1;
                }
            } catch (error) {
                console.error('Like/unlike error:', error);
            }
        }

        if (e.target.closest('.mod-action')) {
            const actionBtn = e.target.closest('.mod-action');
            const action = actionBtn.dataset.action;
            const postId = actionBtn.dataset.postId;

            if (!authManager.can('manage_posts')) {
                showNotification('Permission denied', 'error');
                return;
            }

            try {
                switch(action) {
                    case 'pin':
                        await postManager.pinPost(postId);
                        actionBtn.innerHTML = '<i class="fas fa-thumbtack"></i> Unpin Post';
                        actionBtn.dataset.action = 'unpin';
                        break;
                    case 'unpin':
                        await postManager.unpinPost(postId);
                        actionBtn.innerHTML = '<i class="fas fa-thumbtack"></i> Pin Post';
                        actionBtn.dataset.action = 'pin';
                        break;
                    case 'feature':
                        await postManager.featurePost(postId);
                        actionBtn.innerHTML = '<i class="fas fa-star"></i> Unfeature Post';
                        actionBtn.dataset.action = 'unfeature';
                        break;
                    case 'unfeature':
                        await postManager.unfeaturePost(postId);
                        actionBtn.innerHTML = '<i class="fas fa-star"></i> Feature Post';
                        actionBtn.dataset.action = 'feature';
                        break;
                    case 'delete':
                        if (confirm('Are you sure you want to delete this post?')) {
                            await postManager.deletePost(postId);
                            const postCard = actionBtn.closest('.post-card');
                            postCard.style.opacity = '0.5';
                            setTimeout(() => postCard.remove(), 300);
                        }
                        break;
                }
            } catch (error) {
                showNotification(error.message || `Failed to ${action} post`, 'error');
            }
        }
    });

    document.addEventListener('click', (e) => {
        if (e.target.closest('.share-btn')) {
            const postCard = e.target.closest('.post-card');
            const postId = postCard.dataset.postId;
            const url = `${window.location.origin}/post-detail.html?id=${postId}`;

            navigator.clipboard.writeText(url)
                .then(() => showNotification('Link copied to clipboard!', 'success'))
                .catch(() => showNotification('Failed to copy link', 'error'));
        }
    });
}

document.addEventListener('DOMContentLoaded', () => {
    initPostEventListeners();
});