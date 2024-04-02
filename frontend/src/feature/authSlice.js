import {createAsyncThunk, createSlice} from "@reduxjs/toolkit";
import { userLogin, registerUser } from "../actions/authActions.js";

const userToken = localStorage.getItem('userToken')
    ? localStorage.getItem('userToken')
    : null

    // const suc = localStorage.getItem('userToken')
    // ? true
    // : false   

const initialState = {
    loading: false,
    userInfo: null,
    userToken,
    error: null,
    success: false,
    message:null,
}
//const initialState=JSON.oarse(localStorage.getItem("favorite") || '[]');
const authSlice = createSlice({
    name:'auth',
    initialState,
    reducers:{
        userLogout:(state)=>{
            localStorage.removeItem('userToken')
            state.loading=false
            state.userInfo=null
            state.userToken=null
            //state.isAuth=false
            state.error=null
        },
        setCredentials:(state,{payload})=>{
            state.userInfo=payload
        },
    },
    extraReducers: (builder)=>{
        builder
            .addCase(userLogin.pending, (state) =>{
            //console.log(payload)
            state.loading = true
            state.error = null
        })
            .addCase(userLogin.fulfilled,(state,{payload}) => {
            state.loading = false
            state.userInfo = payload
            state.userToken = payload.token
            state.success =  true        })
            .addCase(userLogin.rejected,(state,{payload})=>{
            state.loading = false
            state.error = payload
        })
            .addCase(registerUser.pending, (state) =>{
            state.loading = true
            state.error = null
        }).addCase(registerUser.fulfilled, (state,{payload}) =>{
            console.log(payload)
            state.loading = false
            state.success = true
            //state.userInfo = payload
            state.userToken = payload.token
        }).addCase(registerUser.rejected,(state,{payload})=>{
            console.log(payload)
            state.loading = false
            state.error = payload
        })
    },
});

//export const {login,removeFromFavorite} = authSlice.actions;
export const {userLogout,setCredentials} = authSlice.actions
export default authSlice.reducer;