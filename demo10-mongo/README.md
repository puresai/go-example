## 基于mongoDB搜索附近的人

```sql
// 查询，minDistance单位是米 https://www.docs4dev.com/docs/zh/mongodb/v3.6/reference/core-2dsphere.html
db.getCollection('persons').find({
   loc: {
     $nearSphere: {
         $geometry: {
            type : "Point",
            coordinates: [120.1856,30.300 ]
         },
         $minDistance: 1,
         $maxDistance: 10000
     }
  }
})
```