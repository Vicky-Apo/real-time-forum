package handlers

import (
	"net/http"

	"real-time-forum/internal/repository"
	"real-time-forum/internal/utils"
)

// GetAllCategoriesHandler retrieves all post categories
func GetAllCategoriesHandler(cr *repository.CategoryRepository, pr *repository.PostsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get all categories
		categories, err := cr.GetAllCategories()
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve categories")
			return
		}
		// Here we add the count of posts for each category
		for i := range categories {
			categories[i].Count, _ = pr.GetCountPostByCategory(categories[i].ID)
		}
		utils.RespondWithSuccess(w, http.StatusOK, categories)
	}
}
