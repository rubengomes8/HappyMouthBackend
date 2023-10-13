package recipe_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/rubengomes8/HappyMouthBackend/internal/recipe"
	"github.com/rubengomes8/HappyMouthBackend/internal/recipe/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandler_CreateRecipe(t *testing.T) {

	mockSvc := mocks.NewMockService()
	handler := recipe.NewHandler(mockSvc)

	tt := map[string]struct {
		setup      func() *http.Request
		call       func(*http.Request) *httptest.ResponseRecorder
		validation func(*httptest.ResponseRecorder)
	}{
		"no include ingredients": {
			setup: func() *http.Request {

				recipe := recipe.RecipeDefinitions{
					IncludeIngredients: []string{},
					ExcludeIngredients: []string{"Tomato"},
				}
				body, err := json.Marshal(&recipe)
				if err != nil {
					t.Fatal(err)
				}

				req, err := http.NewRequest(http.MethodPost, "/api/recipe", bytes.NewBuffer(body))
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			call: func(req *http.Request) *httptest.ResponseRecorder {
				rr := httptest.NewRecorder()
				http.HandlerFunc(handler.CreateRecipe).ServeHTTP(rr, req)
				return rr
			},
			validation: func(rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
				assert.Equal(t, recipe.ErrRequiredIncludeIngredients.Error(), strings.ReplaceAll(rr.Body.String(), "\n", ""))
			},
		},
		"conflicting include and exclude ingredients": {
			setup: func() *http.Request {

				recipe := recipe.RecipeDefinitions{
					IncludeIngredients: []string{"Apple", "Tomato"},
					ExcludeIngredients: []string{"Tomato"},
				}
				body, err := json.Marshal(&recipe)
				if err != nil {
					t.Fatal(err)
				}

				req, err := http.NewRequest(http.MethodPost, "/api/recipe", bytes.NewBuffer(body))
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			call: func(req *http.Request) *httptest.ResponseRecorder {
				rr := httptest.NewRecorder()
				http.HandlerFunc(handler.CreateRecipe).ServeHTTP(rr, req)
				return rr
			},
			validation: func(rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
				assert.Equal(t, recipe.ErrConflictingIngredients.Error(), strings.ReplaceAll(rr.Body.String(), "\n", ""))
			},
		},
		"success": {
			setup: func() *http.Request {
				rec := recipe.RecipeDefinitions{
					IncludeIngredients: []string{"Apple"},
					ExcludeIngredients: []string{"Tomato"},
				}
				body, err := json.Marshal(&rec)
				if err != nil {
					t.Fatal(err)
				}

				req, err := http.NewRequest(http.MethodPost, "/api/recipe", bytes.NewBuffer(body))
				if err != nil {
					t.Fatal(err)
				}

				mockSvc.AskRecipeFunc.PushReturn(recipe.Recipe{
					ID: func() uuid.UUID {
						id, _ := uuid.FromString("A21040C9-017C-4A34-8449-F5FE26098B93")
						return id
					}(),
				}, nil)
				return req
			},
			call: func(req *http.Request) *httptest.ResponseRecorder {
				rr := httptest.NewRecorder()
				http.HandlerFunc(handler.CreateRecipe).ServeHTTP(rr, req)
				return rr
			},
			validation: func(rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rr.Code)
				var got recipe.Recipe
				err := json.Unmarshal(rr.Body.Bytes(), &got)
				if err != nil {
					t.Fatal(err)
				}
				expected := recipe.Recipe{
					ID: func() uuid.UUID {
						id, _ := uuid.FromString("A21040C9-017C-4A34-8449-F5FE26098B93")
						return id
					}(),
				}
				assert.Equal(t, expected, got)
			},
		},
	}

	for _, tc := range tt {
		req := tc.setup()
		rr := tc.call(req)
		tc.validation(rr)
	}
}
