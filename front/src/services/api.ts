import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'

export interface Ingredient {
  id: number
  name: string
  unit: string
  price: number
  created_at: string
  last_modified: string
}

export interface RecipeIngredient {
  ingredient: Ingredient
  units: number
}

export interface Recipe {
  id: number
  name: string
  ingredients: RecipeIngredient[]
  created_at: string
  last_modified: string
  cost: number
}

export const costlyAPI = createApi({
  reducerPath: 'costlyApi',
  baseQuery: fetchBaseQuery({
    baseUrl: 'http://localhost:1234',
    headers: {
      "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.jYyRJbb0WImFoUUdcslQQfwnXTHJzne-6tsPd8Hrw0I",
    }
  }),
  tagTypes: ['Ingredients', 'Recipes'],
  endpoints: (builder) => ({
    getIngredients: builder.query<Ingredient[], void>({
      query: () => `ingredients`,
      providesTags: ['Ingredients'],
    }),
    getRecipes: builder.query<Recipe[], void>({
      query: () => `recipes`,
      providesTags: ['Recipes'],
    }),
  }),
});

export const {
  useGetIngredientsQuery,
  useGetRecipesQuery,
} = costlyAPI;
