import {createAsyncThunk} from "@reduxjs/toolkit";
import axios from "axios";
import {LOGIN_URL,REGISTRATION_URL} from "../const.js";



const backendURL = 'http://localhost:8080'

export const userLogin = createAsyncThunk(
    'auth/login',
    async ({email,password},thunkAPI)=>{
        try {
         const config = {
             headers:{
                 'Content-Type':'application/json',
             },
             rejectUnauthorized: false,
             requestCert: false,
             agent: false,
             strictSSL: false,
         }
         const {data} = await axios.post(
             LOGIN_URL,
             {email,password},
             config
         )
            sessionStorage.setItem('userToken',data.token)
            console.log(data)
            return data
        } catch (error){
            if(error.response && error.response.data.message){
                return thunkAPI.rejectWithValue(error.response.data.message);
            } else {
                return thunkAPI.rejectWithValue(error.message)
            }
        }
    }
)

export const registerUser = createAsyncThunk(
    'auth/registration',
    async ({token,email,name,password,password_confirmation},thunkAPI) => {
        try {
            const config = {
                headers: {
                    'Content-Type':'application/json',
                },
            }
            const {data} = await axios.post(
                REGISTRATION_URL,
                {token,name,email,password,password_confirmation},
                config
            )
           // localStorage.setItem('userToken',data.token)
            return data

        } catch (error){
            if(error.response && error.response.data.errors){
                return thunkAPI.rejectWithValue(error.response.data.errors);
            } else {
                return thunkAPI.rejectWithValue(error.errors)
            }
        }
    }
)