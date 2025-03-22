package controllers

import (
	"context"
	"expanse-tracker/db"
	"expanse-tracker/models"
	"expanse-tracker/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var incomeExpanseCollection *mongo.Collection = db.OpenCollection(db.Client, "income-expanse")

func AddIncomeExpense(entryType models.Type) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var incomeExpanse models.IncomeExpanse

		if err := c.ShouldBindJSON(&incomeExpanse); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validate := validator.New()

		validate.RegisterValidation("account_type", utils.ValidateAccountType)
		validate.RegisterValidation("bank_name", utils.ValidateBankName)
		validate.RegisterValidation("category", utils.ValidateCategory)

		validationError := validate.Struct(incomeExpanse)

		if validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
			return
		}
		incomeExpanse.Created_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		incomeExpanse.Updated_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		incomeExpanse.Id = primitive.NewObjectID()
		incomeExpanse.Type = entryType

		_, err := incomeExpanseCollection.InsertOne(ctx, incomeExpanse)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "account added successfully", "incomeExpanse": incomeExpanse})

	}

}

func AddIncome() gin.HandlerFunc {
	return AddIncomeExpense("Income")

}
func AddExpanse() gin.HandlerFunc {
	return AddIncomeExpense("Expanse")

}

// Total

func GetTotalIncomeExpense(entryType models.Type) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		monthStr := c.Query("month")
		yearStr := c.Query("year")

		month, err := strconv.Atoi(monthStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month"})
			return
		}

		year, err := strconv.Atoi(yearStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
			return
		}

		// Calculate the first and last day of the specified month
		firstDayOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1)
		lastDayOfMonth = time.Date(lastDayOfMonth.Year(), lastDayOfMonth.Month(), lastDayOfMonth.Day(), 23, 59, 59, 999999999, lastDayOfMonth.Location())
		filter := bson.M{
			"type": entryType,
			"created_date": bson.M{
				"$gte": firstDayOfMonth,
				"$lte": lastDayOfMonth,
			},
		}

		pipeline := []bson.M{
			{
				"$match": filter,
			},
			{
				"$group": bson.M{
					"_id":   nil,
					"total": bson.M{"$sum": "$amount"},
				},
			},
		}

		cursor, err := incomeExpanseCollection.Aggregate(ctx, pipeline)
		if err != nil {
			log.Println("Aggregation error :", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate total income/expense"})
			return
		}

		var result struct {
			Total float64 `bson:"total"`
		}

		if cursor.Next(ctx) {
			if err := cursor.Decode(&result); err != nil {
				log.Println("Decode error :", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode result"})
				return
			}
		} else {
			result.Total = 0
		}

		c.JSON(http.StatusOK, gin.H{
			"month":                fmt.Sprintf("%d-%02d", year, month),
			"total_income/expense": result.Total,
		})
	}
}

func GetTotalIncomeByMonth() gin.HandlerFunc {
	return GetTotalIncomeExpense("Income")
}

func GetTotalExpanseByMonth() gin.HandlerFunc {
	return GetTotalIncomeExpense("Expanse")
}

//Transactions

func GetTransactions() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		dayStr := c.Query("day")
		monthStr := c.Query("month")

		now := time.Now()
		year := now.Year()

		var start, end time.Time

		if dayStr != "" {
			// Filter by day
			day, err := strconv.Atoi(dayStr)
			if err != nil || day < 1 || day > 31 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid day"})
				return
			}

			month := now.Month() // Use current month
			start = time.Date(year, month, day, 0, 0, 0, 0, now.Location())
			end = time.Date(year, month, day, 23, 59, 59, 999999999, now.Location())
		} else if monthStr != "" {
			// Filter by month
			month, err := strconv.Atoi(monthStr)
			if err != nil || month < 1 || month > 12 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month"})
				return
			}

			start = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, now.Location())
			end = start.AddDate(0, 1, -1)
			end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 999999999, end.Location())
		} else {
			// No parameters provided, return an error
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide 'day' or 'month' parameter"})
			return
		}

		filter := bson.M{
			"created_date": bson.M{
				"$gte": start,
				"$lte": end,
			},
		}
		cursor, err := incomeExpanseCollection.Find(ctx, filter)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch the transactions"})
			return
		}

		var transactions []models.Transaction

		if err := cursor.All(ctx, &transactions); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode transactions"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"date_range":   fmt.Sprintf("%s to %s", start.Format("2006-01-02"), end.Format("2006-01-02")),
			"transactions": transactions,
		})

	}
}
