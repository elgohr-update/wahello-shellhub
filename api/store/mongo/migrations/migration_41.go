package migrations

import (
	"context"

	"github.com/sirupsen/logrus"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var migration41 = migrate.Migration{
	Version:     41,
	Description: "update online index from devices collection",
	Up: func(db *mongo.Database) error {
		logrus.WithFields(logrus.Fields{
			"component": "migration",
			"version":   41,
			"action":    "Up",
		}).Info("Applying migration")

		if _, err := db.Collection("connected_devices").Indexes().DropOne(context.TODO(), "last_seen"); err != nil {
			return err
		}

		mod := mongo.IndexModel{
			Keys:    bson.D{{"last_seen", 1}},
			Options: options.Index().SetName("last_seen").SetExpireAfterSeconds(120),
		}
		if _, err := db.Collection("connected_devices").Indexes().CreateOne(context.TODO(), mod); err != nil {
			return err
		}

		return nil
	},
	Down: func(db *mongo.Database) error {
		logrus.WithFields(logrus.Fields{
			"component": "migration",
			"version":   41,
			"action":    "Down",
		}).Info("Applying migration")

		if _, err := db.Collection("connected_devices").Indexes().DropOne(context.TODO(), "last_seen"); err != nil {
			return err
		}

		mod := mongo.IndexModel{
			Keys:    bson.D{{"last_seen", 1}},
			Options: options.Index().SetName("last_seen").SetExpireAfterSeconds(60),
		}
		if _, err := db.Collection("connected_devices").Indexes().CreateOne(context.TODO(), mod); err != nil {
			return err
		}

		return nil
	},
}