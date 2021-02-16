## 基于mongoDB搜索附近的人

```sql
// https://docs.mongodb.com/manual/reference/operator/aggregation/geoNear/#example

db.persons.aggregate([
   {
     $geoNear: {
        near: { type: "Point", coordinates: [ 120.110893,30.2078490] },
        distanceField: "dist.calculated",
        minDistance: 2,
        includeLocs: "dist.location",
        spherical: true
     }
   }
])
```