package mongodb

import (
	"FirstCleanArchitecture/models"
	"FirstCleanArchitecture/services"
	"context"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/spf13/viper"
)

type App struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name"`
	Count int                `bson:"count"`
}

type MongoStorage struct {
	active *mongo.Collection
	cancel *mongo.Collection
}

func NewMongoRepository(cl *mongo.Client, dbName, active, cancel string) *MongoStorage {
	MaxCount := viper.GetInt("maxApps")
	rez := &MongoStorage{
		active: cl.Database(dbName).Collection(active),
		cancel: cl.Database(dbName).Collection(cancel),
	}

	//	delete prev documents
	filterDelete := bson.M{}
	_, err := rez.active.DeleteMany(context.TODO(), filterDelete)
	if err != nil {
		log.Fatal(err)
	}
	_, err = rez.cancel.DeleteMany(context.TODO(), filterDelete)
	if err != nil {
		log.Fatal(err)
	}

	// generate 50 active app
	for i := 0; i < MaxCount; i++ {
		tempApp := App{
			Name:  services.GetApplicationRandomNAme(),
			Count: 0,
		}
		_, err := rez.active.InsertOne(context.TODO(), tempApp)
		if err != nil {
			log.Fatal(err)
		}
	}

	//create: go func(){ every 200msec delete 1 activeApp and add newApp}
	ticker := time.Tick(200 * time.Millisecond)
	go func() {
		for range ticker {
			rez.refreshAvailableAppPool()
		}
	}()

	return rez
}

func (r *MongoStorage) refreshAvailableAppPool() {
	// get random
	randINT := services.Random(0, viper.GetInt("maxApps")-1)

	skip := int64(randINT)

	// made skip options for find one options
	opts := &options.FindOneOptions{
		Skip: &skip,
	}

	app := App{}
	// find one
	err := r.active.FindOne(context.TODO(), bson.D{}, opts).Decode(&app)
	if err != nil {
		log.Error(err)
	}

	// if app have >0 show add to cancel
	if app.Count > 0 {
		tempApp := App{
			Name:  app.Name,
			Count: app.Count,
		}

		_, err = r.cancel.InsertOne(context.TODO(), tempApp)
		if err != nil {
			log.Error(err)
		}
	}

	// made filter for update one
	filter := bson.D{{"_id", app.ID}}

	// delete one
	_, err = r.active.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Error(err)
	}

	// add new active instead deleted app
	_, err = r.active.InsertOne(context.TODO(), App{
		Name:  services.GetApplicationRandomNAme(),
		Count: 0,
	})
	if err != nil {
		log.Error(err)
	}
}

func (r *MongoStorage) GetShowedAndCancelApplications(ctx context.Context) (rezActive []models.Application, rezCancel []models.Application, err error) {
	// ADD ACTIVE
	// find many active
	filter := bson.M{
		"count": bson.M{"$gt": 0},
	}

	curA, err := r.active.Find(context.TODO(), filter, &options.FindOptions{})
	if err != nil {
		log.Error(err)
		return rezActive, rezCancel, err
	}

	defer func() {
		err = curA.Close(context.TODO())
		if err != nil {
			log.Error(err)
		}
	}()

	// parse all
	for curA.Next(context.TODO()) {
		var episode App
		if err = curA.Decode(&episode); err != nil {
			log.Error(err)
		}

		rezActive = append(rezActive, models.Application{
			Name:  episode.Name,
			Count: episode.Count,
		})
	}

	// ADD CANCEL
	// find all canceled
	curC, err := r.cancel.Find(context.TODO(), bson.M{}, &options.FindOptions{})
	if err != nil {
		log.Error(err)
		return rezActive, rezCancel, err
	}

	defer func() {
		err = curC.Close(context.TODO())
		if err != nil {
			log.Error(err)
		}
	}()

	// parse all
	for curC.Next(context.TODO()) {
		var episode App
		if err = curC.Decode(&episode); err != nil {
			log.Error(err)
		}

		rezCancel = append(rezCancel, models.Application{
			Name:  episode.Name,
			Count: episode.Count,
		})
	}

	return rezActive, rezCancel, nil
}

func (r *MongoStorage) GetRandomAliveApplication(ctx context.Context) (models.Application, error) {
	// get random
	randINT := services.Random(0, viper.GetInt("maxApps")-1)

	skip := int64(randINT)

	// made skip options for find one options
	opts := &options.FindOneOptions{
		Skip: &skip,
	}

	app := App{}
	// find one
	err := r.active.FindOne(context.TODO(), bson.D{}, opts).Decode(&app)
	if err != nil {
		log.Error(err)
		return models.Application{}, err
	}

	// made filter for update one
	filter := bson.D{{"_id", app.ID}}
	// made update data foe update one
	update := bson.D{
		{"$inc", bson.D{
			{"count", 1},
		}},
	}

	// update one
	_, err = r.active.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error(err)
		return models.Application{}, err
	}

	return models.Application{
		Name:  app.Name,
		Count: app.Count,
	}, nil
}
