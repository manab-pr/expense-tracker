package controllers

import (
	"context"
	"expanse-tracker/db"
	"expanse-tracker/models"
	"expanse-tracker/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var budgetCollection *mongo.Collection = db.OpenCollection(db.Client, "budget")

func CreateBudget() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var budget models.Budget

		if err := c.BindJSON(&budget); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validate := validator.New()
		validate.RegisterValidation("category", utils.ValidateCategory)

		validationErr := validate.Struct(budget)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		currentTime := time.Now()

		budget.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		budget.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		budget.Month = currentTime.Format("January")
		budget.Id = primitive.NewObjectID()
		if budget.Recieve_Alert == nil {
			defaultRecieveAlert := false
			budget.Recieve_Alert = &defaultRecieveAlert
		}
		if budget.Amount == nil {
			defaultBudgetAmount := int64(0)
			budget.Amount = &defaultBudgetAmount
		}

		if budget.Alert_Percentage == nil {
			defaultAlertPercentage := 50
			budget.Alert_Percentage = &defaultAlertPercentage
		}

		_, insertErr := budgetCollection.InsertOne(ctx, budget)

		if insertErr != nil {
			msg := fmt.Sprintf("Budget is not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": budget})

	}
}

func UpdateBudget() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		category := c.Query("category")
		if category == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "category query parameter is required"})
			return
		}
		var budget models.Budget

		if err := c.BindJSON(&budget); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		filter := bson.M{
			"category": category,
		}
		update := bson.M{
			"$set": bson.M{
				"amount":           budget.Amount,
				"month":            budget.Month,
				"receive_alert":    budget.Recieve_Alert,
				"alert_percentage": budget.Alert_Percentage,
				"updated_at":       time.Now(), // Update the timestamp
			},
		}

		result, err := budgetCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to update document: %v", err)})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusOK, gin.H{"error": "no budget matched this id "})
			return
		}
		if result.ModifiedCount == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "document matched but was not modified"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "budget updated successfully"})

	}
}

func DeleteBudget() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		budgetId := c.Param("_id")
		objID, err := primitive.ObjectIDFromHex(budgetId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget ID"})
			return
		}
		filter := bson.M{"_id": objID}

		result, err := budgetCollection.DeleteOne(ctx, filter)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while deleting the budget"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Budget deleted successfully"})

	}
}

func GetBudgetAmount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		category := c.Query("category")
		month := c.Query("month")

		if category == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "category query parameter is required"})
			return
		}

		var budget models.Budget

		if err := c.BindJSON(&budget); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{
			"category": category,
			"month":    month,
		}

		err := budgetCollection.FindOne(ctx, filter).Decode(&budget)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while fetching the budget"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "budget for this month fetched successfully", "Budget": budget})

	}
}

func GetRemainingAmount() gin.HandlerFunc {
	return func(c *gin.Context) {

		category := c.Query("category")
		month := c.Query("month")
		year := c.Query("year")

		if category == "" || month == "" || year == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "category, month, and year query parameters are required"})
			return
		}

		budgetCtx := c.Copy()
		GetBudgetAmount()(budgetCtx)

		budgetResponse, exists := budgetCtx.Get("Budget")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch budget"})
			return
		}

		budget := budgetResponse.(models.Budget)

		expenseCtx := c.Copy() // Create a copy of the context to avoid conflicts
		GetTotalExpanseByMonth()(expenseCtx)

		// Extract the total expenses from the response
		expenseResponse, exists := expenseCtx.Get("total_income/expense")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses"})
			return
		}
		totalExpenses := expenseResponse.(float64)

		remainingBudget := float64(*budget.Amount) - totalExpenses

		c.JSON(http.StatusOK, gin.H{
			"category":         category,
			"month":            fmt.Sprintf("%s-%s", year, month),
			"budget_amount":    budget.Amount,
			"total_expenses":   totalExpenses,
			"remaining_budget": remainingBudget,
		})
	}
}
