class MediaManager {
    constructor() {
        this.maxFileSize = 10 * 1024 * 1024; // 10MB
        this.allowedTypes = ['.jpg', '.jpeg', '.png', '.gif', '.mp4', '.mov', '.pdf'];
        this.maxFiles = 10;
    }

    async uploadMedia(files) {
        if (files.length > this.maxFiles) {
            throw new Error(`Maximum ${this.maxFiles} files allowed`);
        }

        const formData = new FormData();

        for (const file of files) {
            this.validateFile(file);
            formData.append('files', file);
        }

        try {
            const response = await fetchWithAuth('/media/upload', {
                method: 'POST',
                body: formData
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.message || 'Upload failed');
            }

            const uploadedFiles = await response.json();
            showNotification(`${files.length} file(s) uploaded successfully`, 'success');
            return uploadedFiles;
        } catch (error) {
            console.error('Media upload error:', error);
            throw error;
        }
    }

    validateFile(file) {
        if (file.size > this.maxFileSize) {
            throw new Error(`${file.name} is too large (max ${formatFileSize(this.maxFileSize)})`);
        }

        const ext = '.' + file.name.split('.').pop().toLowerCase();
        if (!this.allowedTypes.includes(ext)) {
            throw new Error(`${file.name} is not a supported file type`);
        }

        return true;
    }

    async deleteMedia(fileUrl) {
        try {
            const encodedUrl = encodeURIComponent(fileUrl);
            const response = await fetchWithAuth(`/media/${encodedUrl}`, {
                method: 'DELETE'
            });

            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.message || 'Delete failed');
            }

            showNotification('Media deleted successfully', 'success');
            return true;
        } catch (error) {
            console.error('Media delete error:', error);
            throw error;
        }
    }

    async getMediaInfo(fileUrl) {
        try {
            const encodedUrl = encodeURIComponent(fileUrl);
            const response = await fetchWithAuth(`/media/info/${encodedUrl}`);

            if (!response.ok) {
                throw new Error('Failed to get media info');
            }

            return response.json();
        } catch (error) {
            console.error('Get media info error:', error);
            throw error;
        }
    }

    renderMediaPreview(file, index) {
        const preview = document.createElement('div');
        preview.className = 'file-preview-item';
        preview.dataset.index = index;

        if (file.type.startsWith('image/')) {
            const reader = new FileReader();
            reader.onload = function(e) {
                const img = document.createElement('img');
                img.src = e.target.result;
                preview.appendChild(img);
            };
            reader.readAsDataURL(file);
        } else if (file.type.startsWith('video/')) {
            const video = document.createElement('video');
            video.src = URL.createObjectURL(file);
            video.controls = true;
            preview.appendChild(video);
        } else {
            preview.innerHTML = `
                <div style="display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100%;">
                    <i class="fas fa-file fa-3x" style="color: var(--gray-color);"></i>
                    <span style="font-size: 0.8rem; margin-top: 10px;">${file.name}</span>
                </div>
            `;
        }

        const removeBtn = document.createElement('button');
        removeBtn.className = 'remove-file';
        removeBtn.innerHTML = '<i class="fas fa-times"></i>';
        removeBtn.addEventListener('click', () => this.removeFile(index));

        preview.appendChild(removeBtn);
        return preview;
    }

    createUploadProgressBar() {
        const progressBar = document.createElement('div');
        progressBar.className = 'progress-bar';
        progressBar.innerHTML = '<div class="progress-fill"></div>';
        return progressBar;
    }

    updateProgressBar(progressBar, percentage) {
        const fill = progressBar.querySelector('.progress-fill');
        fill.style.width = `${percentage}%`;
    }
}

const mediaManager = new MediaManager();

// Media upload functionality
function initMediaUpload(uploadElementId, previewElementId) {
    const uploadElement = document.getElementById(uploadElementId);
    const previewElement = document.getElementById(previewElementId);

    if (!uploadElement || !previewElement) return;

    let selectedFiles = [];

    // Drag and drop
    uploadElement.addEventListener('dragover', (e) => {
        e.preventDefault();
        uploadElement.style.borderColor = 'var(--primary-color)';
        uploadElement.style.backgroundColor = 'rgba(74, 111, 165, 0.1)';
    });

    uploadElement.addEventListener('dragleave', () => {
        uploadElement.style.borderColor = '';
        uploadElement.style.backgroundColor = '';
    });

    uploadElement.addEventListener('drop', (e) => {
        e.preventDefault();
        uploadElement.style.borderColor = '';
        uploadElement.style.backgroundColor = '';

        if (e.dataTransfer.files.length) {
            handleFiles(e.dataTransfer.files);
        }
    });

    // File input
    const fileInput = uploadElement.querySelector('input[type="file"]');
    if (fileInput) {
        uploadElement.addEventListener('click', () => fileInput.click());

        fileInput.addEventListener('change', (e) => {
            if (e.target.files.length) {
                handleFiles(e.target.files);
            }
        });
    }

    function handleFiles(files) {
        const validFiles = Array.from(files).filter(file => {
            try {
                mediaManager.validateFile(file);
                return true;
            } catch (error) {
                showNotification(error.message, 'error');
                return false;
            }
        });

        if (selectedFiles.length + validFiles.length > mediaManager.maxFiles) {
            showNotification(`Maximum ${mediaManager.maxFiles} files allowed`, 'error');
            validFiles.length = mediaManager.maxFiles - selectedFiles.length;
        }

        selectedFiles.push(...validFiles);
        renderPreview();
    }

    function renderPreview() {
        previewElement.innerHTML = '';
        selectedFiles.forEach((file, index) => {
            const preview = mediaManager.renderMediaPreview(file, index);
            previewElement.appendChild(preview);
        });

        // Update file input
        if (fileInput) {
            const dataTransfer = new DataTransfer();
            selectedFiles.forEach(file => dataTransfer.items.add(file));
            fileInput.files = dataTransfer.files;
        }
    }

    function removeFile(index) {
        selectedFiles.splice(index, 1);
        renderPreview();
    }

    return {
        getFiles: () => selectedFiles,
        clearFiles: () => {
            selectedFiles = [];
            renderPreview();
        },
        upload: async () => {
            if (selectedFiles.length === 0) {
                throw new Error('No files selected');
            }

            try {
                const uploadedFiles = await mediaManager.uploadMedia(selectedFiles);
                selectedFiles = [];
                renderPreview();
                return uploadedFiles;
            } catch (error) {
                throw error;
            }
        }
    };
}