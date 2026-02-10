# CHECK ALL INDEXES IN DATABASE

## USERS collection indexes
```bash
printjson(db.users.getIndexes());
```

## POSTS collection indexes
```bash
printjson(db.posts.getIndexes());
```

## COMMENTS collection indexes
```bash
printjson(db.comments.getIndexes());
```

## CHECK INDEX USAGE STATISTICS
```bash
printjson(db.users.aggregate([{ $indexStats: {} }]).toArray());
printjson(db.posts.aggregate([{ $indexStats: {} }]).toArray());
printjson(db.comments.aggregate([{ $indexStats: {} }]).toArray());
```

# EXPLAIN QUERIES TO SEE INDEX USAGE

## Get user by email
```bash
printjson(db.users.find({ email: "yerasylhello@gmail.com" }).explain("executionStats").executionStats);
```

## Get posts by category with pagination
```bash
printjson(db.posts.find({ category: "meme" })
.sort({ created_at: -1 })
.limit(10)
.explain("executionStats").executionStats);
```

## Get posts by author
```bash
printjson(db.posts.find({ author_id: ObjectId("6988f32b0e3fbc33b8f0e459") })
.sort({ created_at: -1 })
.limit(10)
.explain("executionStats").executionStats);
```
## Get comments for a post
```bash
printjson(db.comments.find({ post_id: ObjectId("6988f32b0e3fbc33b8f0e459") })
.sort({ created_at: -1 })
.limit(50)
.explain("executionStats").executionStats);
```

## Quick check all in one line
```bash
db.getCollectionNames().forEach(coll => {
print(`${coll}: ${db[coll].getIndexes().length} indexes`);
});
```