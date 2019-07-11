package store

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

func (db *DB) DeleteStaticData() error {

	// I know this is bad but I really dont care about errors here for now

	collection := db.Database.Database("eve_static").Collection("regions")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database("eve_static").Collection("constellations")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database("eve_static").Collection("solarsystems")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database("eve_static").Collection("stars")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database("eve_static").Collection("planets")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database("eve_static").Collection("moons")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database("eve_static").Collection("asteroid_belts")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database("eve_static").Collection("stargates")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database("eve_static").Collection("stations")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database("eve_static").Collection("categories")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database("eve_static").Collection("groups")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database("eve_static").Collection("types")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database("eve_static").Collection("regions")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database("eve_static").Collection("ancestries")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database("eve_static").Collection("bloodlines")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	collection = db.Database.Database("eve_static").Collection("factions")
	_, _ = collection.DeleteMany(context.Background(), bson.M{})

	return nil

}

func (db *DB) InsertRegion(region ESIRegion) error {

	collection := db.Database.Database("eve_static").Collection("regions")

	_, err := collection.InsertOne(context.TODO(), region)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve region")
	}

	return nil

}

func (db *DB) InsertConstellation(cons ESIConstellation) error {

	collection := db.Database.Database("eve_static").Collection("constellations")

	_, err := collection.InsertOne(context.TODO(), cons)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve constellation")
	}

	return nil

}

func (db *DB) InsertSystem(system ESISystem) error {

	collection := db.Database.Database("eve_static").Collection("solarsystems")

	_, err := collection.InsertOne(context.TODO(), system)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve system")
	}

	return nil

}

func (db *DB) InsertStar(star ESIStar) error {

	collection := db.Database.Database("eve_static").Collection("stars")

	_, err := collection.InsertOne(context.TODO(), star)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve star")
	}

	return nil

}

func (db *DB) InsertPlanet(planet ESIPlanet) error {

	collection := db.Database.Database("eve_static").Collection("planets")

	_, err := collection.InsertOne(context.TODO(), planet)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve planet")
	}

	return nil

}

func (db *DB) InsertMoon(moon ESIMoon) error {

	collection := db.Database.Database("eve_static").Collection("moons")

	_, err := collection.InsertOne(context.TODO(), moon)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve moon")
	}

	return nil

}

func (db *DB) InsertAsteroidBelt(belt ESIAsteroidBelt) error {

	collection := db.Database.Database("eve_static").Collection("asteroid_belts")

	_, err := collection.InsertOne(context.TODO(), belt)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve asteroid belt")
	}

	return nil

}

func (db *DB) InsertStargate(gate ESIStargate) error {

	collection := db.Database.Database("eve_static").Collection("stargates")

	_, err := collection.InsertOne(context.TODO(), gate)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve stargate")
	}

	return nil

}

func (db *DB) InsertStation(station ESIStation) error {

	collection := db.Database.Database("eve_static").Collection("stations")

	_, err := collection.InsertOne(context.TODO(), station)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve station")
	}

	return nil
}

func (db *DB) InsertType(typeESI ESIType) error {

	collection := db.Database.Database("eve_static").Collection("types")

	_, err := collection.InsertOne(context.TODO(), typeESI)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve type")
	}

	return nil
}

func (db *DB) InsertGroup(group ESIGroup) error {

	collection := db.Database.Database("eve_static").Collection("groups")

	_, err := collection.InsertOne(context.TODO(), group)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve group")
	}

	return nil
}

func (db *DB) InsertCategory(category ESICategory) error {

	collection := db.Database.Database("eve_static").Collection("categories")

	_, err := collection.InsertOne(context.TODO(), category)
	if err != nil {
		return errors.Wrap(err, "failed to insert eve category")
	}

	return nil
}

func (db *DB) GetSystems() (systems []ESISystem, err error) {
	collection := db.Database.Database("eve_static").Collection("solarsystems")

	ctx := context.Background()

	c, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return systems, errors.Wrap(err, "error retrieving existing systems")
	}

	defer c.Close(ctx)

	for c.Next(ctx) {

		var id ESISystem

		err := c.Decode(&id)
		if err != nil {
			return systems, errors.Wrap(err, "Failed to morp system into struct")
		}

		systems = append(systems, id)
	}

	return systems, nil
}
