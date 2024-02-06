package utils

import (
	"database/sql"
	"fmt"
	"math"
)

type PaginationParams struct {
	Page     int
	PageSize int
}

type TableService struct {
	DB *sql.DB
}

type PaginationResult struct {
	Rows  int
	Pages int
}

func (table TableService) GenericPagination(tableClause, conditions string, valueArgs []interface{}, paginationParams *PaginationParams) (*PaginationResult, error) {
	totalCountQuery := fmt.Sprintf("Select COUNT(*) from %s WHERE TRUE %s", tableClause, conditions)
	var totalCount int
	err := table.DB.QueryRow(totalCountQuery, valueArgs...).Scan(&totalCount)

	if paginationParams.Page == 0 {
		paginationParams.Page = 1
	}

	if paginationParams.PageSize == 0 {
		paginationParams.PageSize = 10
	}
	if err != nil {
		return nil, fmt.Errorf("Error getting total count: %w", err)
	}

	result := &PaginationResult{
		Rows:  totalCount,
		Pages: int(math.Ceil(float64(totalCount) / float64(paginationParams.PageSize))),
	}

	return result, nil
}
