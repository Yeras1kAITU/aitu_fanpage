class CommentManager {
    constructor() {
        this.currentPostId = null;
    }

    async createComment(postId, content) {
        try {
            const response = await fetchWithAuth(`/api/posts/${postId}/comments`, {
                method: 'POST',
                body: JSON.stringify({ content })
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.message || 'Failed to create comment');
            }

            const comment = await response.json();
            showNotification('Comment added!', 'success');
            return comment;
        } catch (error) {
            console.error('Create comment error:', error);
            throw error;
        }
    }

    async getComments(postId, limit = 50, offset = 0) {
        try {
            const response = await fetchWithAuth(`/api/posts/${postId}/comments?limit=${limit}&offset=${offset}`);
            return response.json();
        } catch (error) {
            console.error('Get comments error:', error);
            throw error;
        }
    }

    async updateComment(commentId, content) {
        try {
            const response = await fetchWithAuth(`/api/comments/${commentId}`, {
                method: 'PUT',
                body: JSON.stringify({ content })
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.message || 'Failed to update comment');
            }

            const comment = await response.json();
            showNotification('Comment updated!', 'success');
            return comment;
        } catch (error) {
            console.error('Update comment error:', error);
            throw error;
        }
    }

    async deleteComment(commentId) {
        try {
            const response = await fetchWithAuth(`/api/comments/${commentId}`, {
                method: 'DELETE'
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.message || 'Failed to delete comment');
            }

            showNotification('Comment deleted!', 'success');
            return true;
        } catch (error) {
            console.error('Delete comment error:', error);
            throw error;
        }
    }

    async getCommentCount(postId) {
        try {
            const response = await fetchWithAuth(`/api/posts/${postId}/comments/count`);
            const data = await response.json();
            return data.comment_count || 0;
        } catch (error) {
            console.error('Get comment count error:', error);
            return 0;
        }
    }

    renderComment(comment) {
        const user = authManager.getUser();
        const canEdit = user && (user.id === comment.author_id || authManager.can('manage_comments'));

        return `
            <div class="comment" data-comment-id="${comment.id}">
                <div class="comment-header">
                    <div class="comment-author">
                        <img src="${comment.author_profile_image || '/assets/default-avatar.svg'}" 
                             alt="${comment.author_name}" class="comment-avatar">
                        <div class="comment-author-info">
                            <h5>${comment.author_name}</h5>
                            <span class="comment-time">${formatTime(comment.created_at)}</span>
                        </div>
                    </div>
                    ${canEdit ? `
                        <div class="comment-actions">
                            <button class="btn-icon edit-comment-btn" data-comment-id="${comment.id}">
                                <i class="fas fa-edit"></i>
                            </button>
                            <button class="btn-icon delete-comment-btn" data-comment-id="${comment.id}">
                                <i class="fas fa-trash"></i>
                            </button>
                        </div>
                    ` : ''}
                </div>
                <div class="comment-content">
                    <p>${comment.content}</p>
                </div>
                ${comment.updated_at !== comment.created_at ? `
                    <div class="comment-edited">
                        <small>Edited ${formatTime(comment.updated_at)}</small>
                    </div>
                ` : ''}
            </div>
        `;
    }

    renderCommentForm() {
        return `
            <div class="comment-form">
                <div class="form-group">
                    <textarea class="form-control" id="comment-text" 
                              placeholder="Write your comment here..." 
                              rows="3" required></textarea>
                </div>
                <div class="form-actions">
                    <button type="submit" class="btn btn-primary">
                        <i class="fas fa-paper-plane"></i> Post Comment
                    </button>
                </div>
            </div>
        `;
    }

    renderEditForm(comment) {
        return `
            <div class="comment-edit-form" data-comment-id="${comment.id}">
                <div class="form-group">
                    <textarea class="form-control edit-comment-text">${comment.content}</textarea>
                </div>
                <div class="form-actions">
                    <button type="button" class="btn btn-primary save-edit-btn">
                        <i class="fas fa-save"></i> Save
                    </button>
                    <button type="button" class="btn btn-secondary cancel-edit-btn">
                        <i class="fas fa-times"></i> Cancel
                    </button>
                </div>
            </div>
        `;
    }
}

const commentManager = new CommentManager();

function initCommentSystem(postId) {
    const commentsContainer = document.getElementById('comments-container');
    const commentForm = document.getElementById('comment-form');

    if (!commentsContainer) return;

    commentManager.currentPostId = postId;

    // Load comments
    loadComments();

    // Initialize comment form if user is logged in
    if (commentForm && authManager.isAuthenticated()) {
        commentForm.innerHTML = commentManager.renderCommentForm();

        commentForm.addEventListener('submit', async (e) => {
            e.preventDefault();

            const textarea = document.getElementById('comment-text');
            const content = textarea.value.trim();

            if (!content) {
                showNotification('Comment cannot be empty', 'error');
                return;
            }

            const submitBtn = commentForm.querySelector('button[type="submit"]');
            const originalText = submitBtn.innerHTML;

            try {
                submitBtn.disabled = true;
                submitBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Posting...';

                const comment = await commentManager.createComment(postId, content);

                // Add new comment to the top
                const commentEl = createCommentElement(comment);
                const commentsList = document.querySelector('.comments-list');
                commentsList.insertBefore(commentEl, commentsList.firstChild);

                // Update comment count
                updateCommentCount();

                // Clear form
                textarea.value = '';

                showNotification('Comment posted!', 'success');

            } catch (error) {
                showNotification(error.message || 'Failed to post comment', 'error');
            } finally {
                submitBtn.disabled = false;
                submitBtn.innerHTML = originalText;
            }
        });
    } else if (commentForm) {
        commentForm.innerHTML = `
            <div class="auth-required">
                <p>Please <a href="login.html">login</a> to post comments</p>
            </div>
        `;
    }

    // Event listeners for comment actions
    document.addEventListener('click', async (e) => {
        // Edit comment
        if (e.target.closest('.edit-comment-btn')) {
            const commentId = e.target.closest('.edit-comment-btn').dataset.commentId;
            const commentEl = document.querySelector(`[data-comment-id="${commentId}"]`);

            // Get comment data
            const commentContent = commentEl.querySelector('.comment-content p').textContent;

            // Replace comment with edit form
            const editForm = document.createElement('div');
            editForm.innerHTML = commentManager.renderEditForm({
                id: commentId,
                content: commentContent
            });

            commentEl.querySelector('.comment-content').innerHTML = '';
            commentEl.querySelector('.comment-content').appendChild(editForm);

            // Event listeners for edit form
            const saveBtn = editForm.querySelector('.save-edit-btn');
            const cancelBtn = editForm.querySelector('.cancel-edit-btn');

            saveBtn.addEventListener('click', async () => {
                const newContent = editForm.querySelector('.edit-comment-text').value.trim();

                if (!newContent) {
                    showNotification('Comment cannot be empty', 'error');
                    return;
                }

                try {
                    const updatedComment = await commentManager.updateComment(commentId, newContent);

                    // Replace edit form with updated comment
                    commentEl.querySelector('.comment-content').innerHTML = `<p>${updatedComment.content}</p>`;

                    // Update edited time if available
                    if (updatedComment.updated_at !== updatedComment.created_at) {
                        const editedEl = document.createElement('div');
                        editedEl.className = 'comment-edited';
                        editedEl.innerHTML = `<small>Edited ${formatTime(updatedComment.updated_at)}</small>`;
                        commentEl.appendChild(editedEl);
                    }

                } catch (error) {
                    showNotification(error.message || 'Failed to update comment', 'error');
                }
            });

            cancelBtn.addEventListener('click', () => {
                // Restore original comment
                commentEl.querySelector('.comment-content').innerHTML = `<p>${commentContent}</p>`;
            });
        }

        // Delete comment
        if (e.target.closest('.delete-comment-btn')) {
            const commentId = e.target.closest('.delete-comment-btn').dataset.commentId;

            if (confirm('Are you sure you want to delete this comment?')) {
                try {
                    await commentManager.deleteComment(commentId);

                    const commentEl = document.querySelector(`[data-comment-id="${commentId}"]`);
                    commentEl.style.opacity = '0.5';
                    setTimeout(() => commentEl.remove(), 300);

                    updateCommentCount();

                } catch (error) {
                    showNotification(error.message || 'Failed to delete comment', 'error');
                }
            }
        }
    });
}

async function loadComments() {
    const commentsContainer = document.getElementById('comments-container');
    if (!commentsContainer) return;

    const commentsList = commentsContainer.querySelector('.comments-list');
    if (!commentsList) return;

    try {
        const comments = await commentManager.getComments(commentManager.currentPostId);

        if (comments.length === 0) {
            commentsList.innerHTML = '<div class="no-comments"><p>No comments yet. Be the first to comment!</p></div>';
            return;
        }

        commentsList.innerHTML = comments.map(comment =>
            commentManager.renderComment(comment)
        ).join('');

    } catch (error) {
        console.error('Error loading comments:', error);
        commentsList.innerHTML = '<div class="error"><p>Failed to load comments</p></div>';
    }
}

function createCommentElement(comment) {
    const div = document.createElement('div');
    div.className = 'comment';
    div.dataset.commentId = comment.id;
    div.innerHTML = commentManager.renderComment(comment);
    return div;
}

async function updateCommentCount() {
    const countEl = document.getElementById('comment-count');
    if (!countEl || !commentManager.currentPostId) return;

    try {
        const count = await commentManager.getCommentCount(commentManager.currentPostId);
        countEl.textContent = count;
    } catch (error) {
        console.error('Error updating comment count:', error);
    }
}