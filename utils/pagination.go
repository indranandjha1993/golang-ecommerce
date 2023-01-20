package utils

import (
	"fmt"
	"math"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type Pagination struct {
	Page         int         `json:"page"`
	Limit        int         `json:"limit"`
	Total        int         `json:"total"`
	Results      interface{} `json:"results"`
	NextLink     string      `json:"next_link"`
	CurrentLink  string      `json:"current_link"`
	PreviousLink string      `json:"previous_link"`
	StartLink    string      `json:"start_link"`
	EndLink      string      `json:"end_link"`
}

func Paginate(c echo.Context, db *gorm.DB, results interface{}) (*Pagination, error) {
	// Get the query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}

	// Calculate the offset
	offset := (page - 1) * limit

	// Find the results
	if err := db.Limit(limit).Offset(offset).Find(results).Error; err != nil {
		return nil, err
	}

	// Get the total count
	var count int64
	if err := db.Model(results).Count(&count).Error; err != nil {
		return nil, err
	}

	// Calculate the next, current, previous, start and end pages
	next := page + 1
	if next*limit > int(count) {
		next = 0
	}
	current := page
	previous := page - 1
	if previous <= 0 {
		previous = 0
	}
	start := int(math.Max(1, float64(current)-2))
	end := int(math.Ceil(float64(count) / float64(limit)))

	baseURL := c.Scheme() + "://" + c.Request().Host
	path := c.Path()
	nextLink := ""
	currentLink := baseURL + path + "?page=" + strconv.Itoa(current) + "&limit=" + strconv.Itoa(limit)
	previousLink := ""

	if previous != 0 {
		previousLink = baseURL + path + "?page=" + strconv.Itoa(previous) + "&limit=" + strconv.Itoa(limit)
	}
	if next != 0 {
		nextLink = baseURL + path + "?page=" + strconv.Itoa(next) + "&limit=" + strconv.Itoa(limit)
	}

	startLink := fmt.Sprintf("%s?page=%d&limit=%d", baseURL, start, limit)
	endLink := fmt.Sprintf("%s?page=%d&limit=%d", baseURL, end, limit)

	// Create the pagination response
	pagination := &Pagination{
		Page:         page,
		Limit:        limit,
		Total:        int(count),
		Results:      results,
		NextLink:     nextLink,
		CurrentLink:  currentLink,
		PreviousLink: previousLink,
		StartLink:    startLink,
		EndLink:      endLink,
	}
	return pagination, nil
}
