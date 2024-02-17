import React, {useEffect} from "react";
import {createBrowserRouter, createRoutesFromElements, Route, RouterProvider} from "react-router-dom";
import {ErrorPage} from "./Components/ErrorPage/ErrorPage.jsx";
import {Root} from "./routes/Root.jsx";
import {Expressions} from "./Components/Expressions/Expressions.jsx";
import {Operations} from "./Components/Operations/Operations.jsx";
import {ComputionCapabilities} from "./Components/Compution_capabilities/ComputionCapabilities.jsx";
import {AddExpression} from "./Components/AddExpression/AddExpression.jsx";

const router = createBrowserRouter(
    createRoutesFromElements(
         <Route path='/' element={<Root />}>
                     <Route path='/' element={<AddExpression />} />
                     <Route path='/expressions' element={<Expressions />} />
                     <Route path='/operations' element={<Operations />} />
                     <Route path='/computing_capabilities' element={<ComputionCapabilities />} />
         </Route>
    )
)
export const App = () =>{
    return (
        <RouterProvider router={router}>

        </RouterProvider>
    )
}

