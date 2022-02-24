package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

func StringsToBsonId(ids []string) (objectIds []primitive.ObjectID) {
	for _, id := range ids {
		objId, _ := primitive.ObjectIDFromHex(id)
		objectIds = append(objectIds, objId)
	}
	return
}
