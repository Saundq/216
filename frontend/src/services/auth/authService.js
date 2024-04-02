import {createApi, fetchBaseQuery} from '@reduxjs/toolkit/query/react'
import {API_URL} from "../../const.js";

export const authApi = createApi({
    reducerPath:'authApi',
    baseQuery: fetchBaseQuery({
        baseUrl:API_URL,
        prepareHeaders:(headers,{getState})=> {
            const token = getState().auth.userToken
            if (token) {
                headers.set('authorization', `Bearer ${token}`)
                return headers
            }
        },
        }),
        endpoints:(builder)=>({
            getUserDetails: builder.query({
                query:()=>({
                    url: 'http://localhost:8181/api/v1/profile',
                    method:'GET',
                }),
            }),
        }),
})

export const { useGetUserDetailsQuery} = authApi