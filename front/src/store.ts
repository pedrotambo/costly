import { configureStore } from "@reduxjs/toolkit";
import { costlyAPI } from "./services/api";

export const store = configureStore({
    reducer: {
        [costlyAPI.reducerPath]: costlyAPI.reducer,
    },
    middleware: (getDefaultMiddleware) =>
        getDefaultMiddleware().concat(costlyAPI.middleware),
});

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch;