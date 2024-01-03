package helpers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Melliv/Book-Market-Server/internal/enums"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetEnvVar(key string) string {
	err := godotenv.Load(dir(".env"))
	if err != nil {
		log.Fatal("Error loading .env file. ", err)
	}

	value := os.Getenv(key)

	return value
}

func GetEnv() enums.Env {
	value := GetEnvVar("ENV")
	switch value {
	case "PROD":
		return enums.Prod
	case "DEV":
		return enums.Dev
	case "TEST":
		return enums.Test
	default:
		log.Fatal("Env not exist: ", value)
	}
	return enums.Dev
}

func BaseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func MapCursorToTarget(ctx context.Context, cursor *mongo.Cursor, target interface{}) error {
	if err := cursor.All(ctx, target); err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func GetPrimitiveObjectID(value string) primitive.ObjectID {
	objID, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		panic(err)
	}
	return objID
}

func dir(envFile string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			panic(fmt.Errorf("go.mod not found"))
		}
		currentDir = parent
	}

	return filepath.Join(currentDir, envFile)
}

func ValidateBody(obj interface{}) error {
	validate := validator.New()
	err := validate.Struct(obj)
	
	if err != nil {
		var errList []string
		for _, err := range err.(validator.ValidationErrors) {
			errList = append(errList, err.Field() + " " + err.Tag())
		}
		return fmt.Errorf(strings.Join(errList[:], ", "))
	}
	return nil
}
