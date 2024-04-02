import { useSelector } from 'react-redux'
import { NavLink, Outlet,useLocation,useNavigate } from 'react-router-dom'
import {Container} from "../Components/Layout/Container/Container.jsx";
import { useState,useEffect } from 'react';


export const ProtectedRoute = () => {
    //const { userInfo } = useSelector((state) => state.auth)
    const {loading,userInfo,error,success} = useSelector((state)=>state.auth)

    let navigate = useNavigate();


    const location = useLocation()
    //  const params = new URLSearchParams(location.search)
  console.log(location.pathname)
      useEffect(() => {
          if(!userInfo && location.pathname!="/registration")
              navigate('/login')
          // if(userInfo)
          //     navigate('/profile')
      }, [navigate,success]);

    // show unauthorized screen if no user is found in redux store
    if (!userInfo) {
        return (
             <Container>
                 <div className='unauthorized'>
                 <h1>Unauthorized :(</h1>
                 <span>
                 
           <NavLink to='/login'>Login</NavLink> to gain access
         </span>
             </div>
                </Container>
        )
    }

    // returns child route elements
    return <Outlet />
}