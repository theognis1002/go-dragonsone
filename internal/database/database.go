package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string
	SeedData() error
	GetAllDragons() ([]Dragon, error)
}

type service struct {
	db *mongo.Client
}

var (
	host      = os.Getenv("DB_HOST")
	port      = os.Getenv("DB_PORT")
	username  = os.Getenv("DB_USERNAME")
	password  = os.Getenv("DB_ROOT_PASSWORD")
	authDB    = os.Getenv("DB_AUTH_DB")
	dbService Service
)

func New() Service {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", username, password, host, port, authDB)))
	if err != nil {
		log.Fatal(err)
	}

	dbService = &service{
		db: client,
	}

	// Seeding data
	err = dbService.SeedData()
	if err != nil {
		log.Fatalf("Error seeding data: %v", err)
	}

	return dbService
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.Ping(ctx, nil)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Database down: %v", err))
	}

	return map[string]string{
		"message": "All calm on the Westerosi front.",
	}
}

type Dragon struct {
	Name     string `bson:"name"`
	Nickname string `bson:"email"`
	Color    string `bson:"color"`
	Age      int    `bson:"age"`
	Rider    string `bson:"rider"`
	Notes    string `bson:"notes"`
}

func (s *service) SeedData() error {
	collection := s.db.Database("main").Collection("dragons")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dragons := []interface{}{
		Dragon{
			Name:     "Balerion",
			Nickname: "The Black Dread",
			Color:    "Black",
			Age:      200,
			Rider:    "Aegon the Conqueror",
			Notes:    "The largest and oldest of the Targaryen dragons. Instrumental in the conquest of Westeros.",
		},
		Dragon{
			Name:     "Vhagar",
			Nickname: "",
			Color:    "Green and Bronze",
			Age:      181,
			Rider:    "Visenya Targaryen",
			Notes:    "One of the three dragons used in Aegon's Conquest. Later ridden by multiple Targaryens.",
		},
		Dragon{
			Name:     "Meraxes",
			Nickname: "",
			Color:    "Silver",
			Age:      100,
			Rider:    "Rhaenys Targaryen",
			Notes:    "One of the three dragons used in Aegon's Conquest. Killed in Dorne.",
		},
		Dragon{
			Name:     "Caraxes",
			Nickname: "The Blood Wyrm",
			Color:    "Red",
			Age:      50,
			Rider:    "Daemon Targaryen",
			Notes:    "Known for his ferocity and was ridden during the Dance of the Dragons.",
		},
		Dragon{
			Name:     "Vermithor",
			Nickname: "The Bronze Fury",
			Color:    "Bronze",
			Age:      100,
			Rider:    "King Jaehaerys I",
			Notes:    "The second largest dragon at the time of the Dance of the Dragons.",
		},
		Dragon{
			Name:     "Meleys",
			Nickname: "The Red Queen",
			Color:    "Red",
			Age:      80,
			Rider:    "Rhaenys Targaryen (the Queen Who Never Was)",
			Notes:    "One of the older dragons during the Dance of the Dragons, known for her speed.",
		},
		Dragon{
			Name:     "Sunfyre",
			Nickname: "Sunfyre the Golden",
			Color:    "Gold",
			Age:      40,
			Rider:    "Aegon II Targaryen",
			Notes:    "Considered the most beautiful dragon ever seen in Westeros.",
		},
		Dragon{
			Name:     "Rhaegal",
			Nickname: "",
			Color:    "Green and Bronze",
			Age:      5,
			Rider:    "Jon Snow (briefly)",
			Notes:    "One of Daenerys Targaryen's three dragons.",
		},
		Dragon{
			Name:     "Viserion",
			Nickname: "",
			Color:    "Cream and Gold",
			Age:      5,
			Rider:    "None (undead form ridden by the Night King)",
			Notes:    "One of Daenerys Targaryen's three dragons. Reanimated as an ice dragon.",
		},
		Dragon{
			Name:     "Drogon",
			Nickname: "",
			Color:    "Black with red eyes",
			Age:      7,
			Rider:    "Daenerys Targaryen",
			Notes:    "Named after Daenerys' late husband, Khal Drogo. The largest and most aggressive of Daenerys' dragons.",
		},
	}

	_, err := collection.InsertMany(ctx, dragons)
	if err != nil {
		return fmt.Errorf("could not seed data: %v", err)
	}

	fmt.Println("Data seeded successfully.")
	return nil
}

func (s *service) GetAllDragons() ([]Dragon, error) {
	collection := s.db.Database("main").Collection("dragons")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Println("Error finding dragons:", err)
		return nil, err
	}
	defer cur.Close(ctx)

	var dragons []Dragon
	for cur.Next(ctx) {
		var dragon Dragon
		if err := cur.Decode(&dragon); err != nil {
			log.Println("Error decoding dragon:", err)
			continue
		}
		dragons = append(dragons, dragon)
	}

	if err := cur.Err(); err != nil {
		log.Println("Cursor error:", err)
		return nil, err
	}

	return dragons, nil
}
