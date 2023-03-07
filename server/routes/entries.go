package routes

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/menyasosali/go-react-calorie-tracker/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

var validate = validator.New()
var entryCollection *mongo.Collection = openCollection(Client, "calories")

func AddEntry(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var entry models.Entry

	if err := c.Bind(&entry); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())

	}

	entry.ID = primitive.NewObjectID()

	validateErr := validate.Struct(entry)
	if validateErr != nil {
		fmt.Println(validateErr)
		return c.JSON(http.StatusInternalServerError, validateErr.Error())

	}
	result, err := entryCollection.InsertOne(ctx, entry)
	if err != nil {
		msg := fmt.Sprintf("entry item was not created")
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, msg)

	}

	return c.JSON(http.StatusOK, result)
}

func GetEntries(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := entryCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())

	}

	var entries []bson.M

	if err = result.All(ctx, &entries); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fmt.Println(entries)
	return c.JSON(http.StatusOK, entries)
}

func GetEntryById(c echo.Context) error {
	entryId := c.Param("id")
	docId, _ := primitive.ObjectIDFromHex(entryId)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.M{"_id": docId}

	var entry bson.M

	if err := entryCollection.FindOne(ctx, filter).Decode(&entry); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())

	}

	fmt.Println(entry)
	return c.JSON(http.StatusOK, entry)
}

func UpdateEntry(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	entryId := c.Param("id")
	docId, _ := primitive.ObjectIDFromHex(entryId)

	filter := bson.M{"_id": docId}

	var entry models.Entry

	if err := c.Bind(&entry); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())

	}

	validateErr := validate.Struct(entry)
	if validateErr != nil {
		fmt.Println(validateErr)
		return c.JSON(http.StatusInternalServerError, validateErr.Error())

	}

	result, err := entryCollection.UpdateOne(
		ctx,
		filter,
		bson.M{
			"dish":          entry.Dish,
			"fat":           entry.Fat,
			"carbohydrates": entry.Carbohydrates,
			"protein":       entry.Protein,
			"calories":      entry.Calories,
			"ingredients":   entry.Ingredients,
		},
	)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())

	}

	return c.JSON(http.StatusOK, result.ModifiedCount)
}

func DeleteEntry(c echo.Context) error {
	entryId := c.Param("id")
	docID, _ := primitive.ObjectIDFromHex(entryId)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := entryCollection.DeleteOne(ctx, bson.M{"_id": docID})
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())

	}

	return c.JSON(http.StatusOK, result.DeletedCount)
}

func GetEntriesByIngredient(c echo.Context) error {
	ingredient := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var entries []bson.M

	filter := bson.M{"ingredients": ingredient}

	result, err := entryCollection.Find(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())

	}

	if err = result.All(ctx, &entries); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())

	}
	fmt.Println(entries)
	return c.JSON(http.StatusOK, entries)
}

func UpdateIngredient(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	entryId := c.Param("id")
	docId, _ := primitive.ObjectIDFromHex(entryId)

	type Ingredient struct {
		Ingredients *string `json:"ingredients"`
	}

	var ingredient Ingredient

	if err := c.Bind(&ingredient); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	filter := bson.M{"_id": docId}

	result, err := entryCollection.UpdateOne(ctx, filter, bson.D{{"$set", ingredient.Ingredients}})
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, err.Error())

	}
	return c.JSON(http.StatusOK, result.ModifiedCount)
}
