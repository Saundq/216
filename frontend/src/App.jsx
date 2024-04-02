import React, {useEffect} from "react";
import {createBrowserRouter, createRoutesFromElements, Route, RouterProvider} from "react-router-dom";
import {ErrorPage} from "./Components/ErrorPage/ErrorPage.jsx";
import {Root} from "./routes/Root.jsx";
import {Expressions} from "./Components/Expressions/Expressions.jsx";
import {Operations} from "./Components/Operations/Operations.jsx";
import {ComputionCapabilities} from "./Components/Compution_capabilities/ComputionCapabilities.jsx";
import {AddExpression} from "./Components/AddExpression/AddExpression.jsx";
import {ProtectedRoute} from "./routes/ProtectedRoute.jsx";
import {Login} from "./Components/Authorization/Login/Login.jsx";
import {Registration} from "./Components/Authorization/Registration/Registration.jsx";

const router = createBrowserRouter(
    createRoutesFromElements(
         <Route path='/' element={<Root />}>

                <Route element={<ProtectedRoute />}>
                        <Route path='/' element={<AddExpression />} /> 
                        <Route path='/expressions' element={<Expressions />} />
                        <Route path='/operations' element={<Operations />} />
                        <Route path='/computing_capabilities' element={<ComputionCapabilities />} />
                </Route>
             <Route path="/login" element={<Login />}/>
             <Route path="/registration" element={<Registration />}/>       
         </Route>
    )
)
export const App = () =>{
    return (
        <RouterProvider router={router}>

        </RouterProvider>
    )
}

