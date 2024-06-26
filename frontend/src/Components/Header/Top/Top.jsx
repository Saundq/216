import {Container} from "../../Layout/Container/Container.jsx";
import style from './Top.module.scss'
import {NavLink} from "react-router-dom";
import {useDispatch, useSelector} from "react-redux";
import {useEffect} from "react";
import {setCredentials, userLogout} from "../../../feature/authSlice.js";
import {useGetUserDetailsQuery} from "../../../services/auth/authService.js";

export const Top = () => {
    const {userInfo,userToken} = useSelector((state) => state.auth)
    const dispatch = useDispatch()


    const {data,isFetching} = useGetUserDetailsQuery('userDetails',{
        pollingInterval:5000,
    })

    useEffect(()=>{

        if (data) dispatch(setCredentials(data))
        else {
            sessionStorage.removeItem('userToken')
            dispatch(setCredentials(null));
        }
    },[data,dispatch])


   return (
        <Container >
            <nav className="navbar navbar-expand-lg navbar-light bg-light">
                <NavLink className='nav-link' to='/'>
                    Добавить выражение
                </NavLink>
                <button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
                    <span className="navbar-toggler-icon"></span>
                </button>
                <div className="collapse navbar-collapse" id="navbarNav">
                    <ul className="nav navbar-nav">
                            <li className="nav-item active">
                                <NavLink className='nav-link' to='/expressions'>
                                    Выражения
                                </NavLink>
                            </li>
                            <li className="nav-item">
                                <NavLink className='nav-link' to='/operations'>
                                    Операции
                                </NavLink>
                            </li>
                            <li className="nav-item">
                                <NavLink className='nav-link' to='/computing_capabilities'>
                                    Вычислительные можности
                                </NavLink>
                            </li>
                        </ul>
                        </div>
                        <div className="d-flex align-items-center">
                        <ul className="nav navbar-nav navbar-right">
                        <li>{userInfo ?
                                (
                                    <div className="collapse navbar-collapse" id="navbarNavDarkDropdown">
      <ul className="navbar-nav">
        <li className="nav-item dropdown">
          <button className="btn btn-light dropdown-toggle" data-bs-toggle="dropdown" aria-expanded="false">
          {userInfo.email}
          </button>
          <ul className="dropdown-menu dropdown-menu-light">
            <li><button className='dropdown-item' onClick={()=>dispatch(userLogout())}>
                                         Logout
                                    </button></li>
          </ul>
        </li>
      </ul>
    </div>
                                ):(
                                    <NavLink className='nav-link' to='/login'>
                                        Login
                                    </NavLink>
                                )

                            }
                            </li>
                            

                            <li>
                            {!userInfo ? (
                                <NavLink className='nav-link' to='/registration'>Register</NavLink>):""
                            }
                            </li>
                           
             
            </ul>
                        


                    
                   
                </div>
            </nav>


        </Container>
    )
}