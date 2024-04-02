import {configureStore} from "@reduxjs/toolkit";
import {authApi} from "./services/auth/authService.js";
import authReducer from "./feature/authSlice.js"


const store = configureStore({
    reducer:{
        auth: authReducer,
        [authApi.reducerPath]:authApi.reducer,
    },
    middleware:(getDefaultMiddleware)=>
        getDefaultMiddleware().concat(authApi.middleware),
    devTools:import.meta.env.DEV,
})
export default store;