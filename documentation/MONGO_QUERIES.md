# AITU Fanpage - Railway Deployment CRUD Examples

## Prerequisites

### Get Auth Token (First)
```bash
# Register a user
curl -X POST "https://aitufanpage-production.up.railway.app/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@aitu.edu",
    "password": "test12345",
    "display_name": "Test User"
  }'

# Login to get token
curl -X POST "https://aitufanpage-production.up.railway.app/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "yerasylhello@gmail.com",
    "password": "Qwerty123!"
  }'

# Save the token from response
export TOKEN="some_hash"
```

## User Operations

### 1. Register User
```bash
curl -X POST "https://aitufanpage-production.up.railway.app/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student@aitu.edu",
    "password": "student123",
    "display_name": "AITU Student"
  }'
```

### 2. Login User
```bash
curl -X POST "https://aitufanpage-production.up.railway.app/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student@aitu.edu",
    "password": "student123"
  }'
```

### 3. Get Current User Profile
```bash
curl -X GET "https://aitufanpage-production.up.railway.app/api/users/me" \
  -H "Authorization: Bearer $TOKEN"
```

### 4. Update User Profile
```bash
curl -X PUT "https://aitufanpage-production.up.railway.app/api/users/me" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "display_name": "yeras1k",
    "bio": "SE-2425"
  }'
```

### 5. Change Password
```bash
curl -X PUT "https://aitufanpage-production.up.railway.app/api/users/me/password" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "current_password": "student123",
    "new_password": "newpassword123"
  }'
```

## Post Operations

### 6. Create Post (JSON only)
```bash
curl -X POST "https://aitufanpage-production.up.railway.app/api/posts" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My Post",
    "content": "This is the content of my post",
    "description": "A brief description",
    "category": "news",
    "tags": ["welcome", "first"]
  }'
```

### 7. Create Post with Files (Using PowerShell)
```powershell
# Save this as upload.ps1
$uri = "https://aitufanpage-production.up.railway.app/api/posts"
$token = "YOUR_TOKEN_HERE"
$imagePath = "C:\Users\YourName\Pictures\test.jpg"

$postData = @{
    title = "Post with Image"
    content = "This post has an image"
    category = "meme"
    description = "Check out this image"
}

$boundary = [System.Guid]::NewGuid().ToString()
$LF = "`r`n"

$bodyLines = (
    "--$boundary",
    "Content-Disposition: form-data; name=`"post`"",
    "Content-Type: application/json",
    "",
    ($postData | ConvertTo-Json),
    "--$boundary",
    "Content-Disposition: form-data; name=`"files`"; filename=`"test.jpg`"",
    "Content-Type: image/jpeg",
    "",
    [System.IO.File]::ReadAllText($imagePath, [System.Text.Encoding]::Default),
    "--$boundary--"
) -join $LF

Invoke-RestMethod -Uri $uri -Method Post `
    -Headers @{ "Authorization" = "Bearer $token" } `
    -ContentType "multipart/form-data; boundary=$boundary" `
    -Body $bodyLines
```

### 8. Get All Posts
```bash
# Get first 10 posts
curl -X GET "https://aitufanpage-production.up.railway.app/api/posts?limit=10&offset=0"

# Get posts by category
curl -X GET "https://aitufanpage-production.up.railway.app/api/posts?category=meme&limit=10"

# Get posts by author
curl -X GET "https://aitufanpage-production.up.railway.app/api/posts?author_id=507f1f77bcf86cd799439011&limit=10"
```

### 9. Get Single Post
```bash
curl -X GET "https://aitufanpage-production.up.railway.app/api/posts/6987be909b6a7491efa27dc16e"
```

### 10. Update Post
```bash
curl -X PUT "https://aitufanpage-production.up.railway.app/api/posts/6987be909b6a7491efa27dc16e" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Title",
    "content": "Updated content here",
    "category": "academic"
  }'
```

### 11. Delete Post
```bash
curl -X DELETE "https://aitufanpage-production.up.railway.app/api/posts/698ae85cc8bce76c6ab7be39" \
  -H "Authorization: Bearer $TOKEN"
```

## Post Interaction Operations

### 12. Like a Post
```bash
curl -X POST "https://aitufanpage-production.up.railway.app/api/posts/698ae38dc8bce76c6ab7be36/like" \
  -H "Authorization: Bearer $TOKEN"
```

### 13. Unlike a Post
```bash
curl -X DELETE "https://aitufanpage-production.up.railway.app/api/posts/698ae38dc8bce76c6ab7be36/like" \
  -H "Authorization: Bearer $TOKEN"
```

### 14. Get Post Likes
```bash
curl -X GET "https://aitufanpage-production.up.railway.app/api/posts/6987be909b6a7491efa27dc16e/likes"
```

## Special Post Queries

### 15. Get Pinned Posts
```bash
curl -X GET "https://aitufanpage-production.up.railway.app/api/posts/pinned?limit=5"
```

### 16. Get Featured Posts
```bash
curl -X GET "https://aitufanpage-production.up.railway.app/api/posts/featured?limit=10"
```

### 17. Get Popular Posts
```bash
curl -X GET "https://aitufanpage-production.up.railway.app/api/posts/popular?limit=10&days=7"
```

### 18. Search Posts
```bash
curl -X GET "https://aitufanpage-production.up.railway.app/api/posts/search?q=welcome&limit=10"
```

### 19. Get Category Statistics
```bash
curl -X GET "https://aitufanpage-production.up.railway.app/api/posts/categories/stats"
```

### 20. Get User Feed (Authenticated)
```bash
curl -X GET "https://aitufanpage-production.up.railway.app/api/posts/feed?limit=10" \
  -H "Authorization: Bearer $TOKEN"
```

## Comment Operations

### 21. Create Comment
```bash
curl -X POST "https://aitufanpage-production.up.railway.app/api/posts/698ae38dc8bce76c6ab7be36/comments" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Great post! Thanks for sharing."
  }'
```

### 22. Get Post Comments
```bash
curl -X GET "https://aitufanpage-production.up.railway.app/api/posts/6987be909b6a7491efa27dc16e/comments?limit=50"
```

### 23. Update Comment
```bash
curl -X PUT "https://aitufanpage-production.up.railway.app/api/comments/698aea56c8bce76c6ab7be3c" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Updated comment text"
  }'
```

### 24. Delete Comment
```bash
curl -X DELETE "https://aitufanpage-production.up.railway.app/api/comments/698aea56c8bce76c6ab7be3c" \
  -H "Authorization: Bearer $TOKEN"
```

### 25. Get Comment Count
```bash
curl -X GET "https://aitufanpage-production.up.railway.app/api/posts/6987be909b6a7491efa27dc16e/comments/count"
```

## Media Operations

### 26. Upload Media (PowerShell)
```powershell
$uri = "https://aitufanpage-production.up.railway.app/api/media/upload"
$token = "your_token"
$imagePath = "C:\Path\To\Image.jpg"

$boundary = [System.Guid]::NewGuid().ToString()
$LF = "`r`n"

$bodyLines = (
    "--$boundary",
    "Content-Disposition: form-data; name=`"files`"; filename=`"upload.jpg`"",
    "Content-Type: image/jpeg",
    "",
    [System.IO.File]::ReadAllText($imagePath, [System.Text.Encoding]::Default),
    "--$boundary--"
) -join $LF

Invoke-RestMethod -Uri $uri -Method Post `
    -Headers @{ "Authorization" = "Bearer $token" } `
    -ContentType "multipart/form-data; boundary=$boundary" `
    -Body $bodyLines
```

### 27. Get Media Info
```bash
curl -X GET "https://aitufanpage-production.up.railway.app/api/media/info/uploads%2Fimages%2F7d0b2891-b090-492b-bc4c-3302904a68f8.jpg" \
  -H "Authorization: Bearer $TOKEN"
```

### 28. Serve Media File
```bash
# Direct access to uploaded files
curl -X GET "https://aitufanpage-production.up.railway.app/uploads/images/7d0b2891-b090-492b-bc4c-3302904a68f8.jpg"
```

## Admin Operations

### 29. Get System Statistics (Admin only)
```bash
curl -X GET "https://aitufanpage-production.up.railway.app/api/admin/stats" \
  -H "Authorization: Bearer $TOKEN"
```

### 30. Get All Users (Admin only)
```bash
curl -X GET "https://aitufanpage-production.up.railway.app/api/admin/users?limit=50&offset=0" \
  -H "Authorization: Bearer $TOKEN"
```

### 31. Search Users (Admin only)
```bash
curl -X GET "https://aitufanpage-production.up.railway.app/api/admin/users/search?q=name&limit=20" \
  -H "Authorization: Bearer $TOKEN"
```

### 32. Update User Role (Admin only)
```bash
curl -X PUT "https://aitufanpage-production.up.railway.app/api/admin/users/6987be909b6a7491efa27dc16e/role" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "moderator"
  }'
```

### 33. Deactivate User (Admin only)
```bash
curl -X PUT "https://aitufanpage-production.up.railway.app/api/admin/users/6987be909b6a7491efa27dc16e/status/deactivate" \
  -H "Authorization: Bearer $TOKEN"
```

### 34. Delete User (Admin only)
```bash
curl -X DELETE "https://aitufanpage-production.up.railway.app/api/admin/users/6987be909b6a7491efa27dc16e" \
  -H "Authorization: Bearer $TOKEN"
```

## Post Moderation Operations

### 35. Pin Post (Moderator/Admin)
```bash
curl -X POST "https://aitufanpage-production.up.railway.app/api/posts/6987be909b6a7491efa27dc16e/pin" \
  -H "Authorization: Bearer $TOKEN"
```

### 36. Unpin Post (Moderator/Admin)
```bash
curl -X DELETE "https://aitufanpage-production.up.railway.app/api/posts/6987be909b6a7491efa27dc16e/pin" \
  -H "Authorization: Bearer $TOKEN"
```

### 37. Feature Post (Moderator/Admin)
```bash
curl -X POST "https://aitufanpage-production.up.railway.app/api/posts/6987be909b6a7491efa27dc16e/feature" \
  -H "Authorization: Bearer $TOKEN"
```

### 38. Unfeature Post (Moderator/Admin)
```bash
curl -X DELETE "https://aitufanpage-production.up.railway.app/api/posts/6987be909b6a7491efa27dc16e/feature" \
  -H "Authorization: Bearer $TOKEN"
```

## System Health Operations

### 39. Health Check
```bash
curl -X GET "https://aitufanpage-production.up.railway.app/health"
```