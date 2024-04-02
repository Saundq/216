import style from "./Login.module.scss"
import {Navigate, NavLink, useNavigate} from "react-router-dom";
import {useEffect, useRef, useState} from "react";
import {useDispatch, useSelector} from "react-redux";
import {userLogin} from "../../../actions/authActions.js";
import {useForm} from 'react-hook-form'


//import { login } from "../../services/auth.service";

export const Login =(props)=> {

    let navigate = useNavigate();
    const {loading,userInfo,error} = useSelector((state)=>state.auth)
    
    const dispatch = useDispatch();

    const {register, handleSubmit,reset} = useForm()

    useEffect(() => {
        if(userInfo){
            navigate('/')
        }
    }, [navigate,userInfo]);
    const submitForm = (data) => {
        //console.log(data.email)
        dispatch(userLogin(data))
        reset()
    }

    return (
        <div className="col-md-4 m-md-auto mt-md-5">
            <div className="card card-container m-3 p-3">
                <h1 className="m-md-auto">Вход:</h1>
            <form onSubmit={handleSubmit(submitForm)} >
                {error && <p className="alert alert-danger">{error}</p>}
                <div className="form-group">
                    <label htmlFor="username">Email</label>
                    <input
                        type="email"
                        className="form-control"
                        {...register('email')}
                        required
                    />
                </div>
                <div className="form-group">
                    <label htmlFor="password">Password</label>
                    <input
                        type="password"
                        className="form-control"
                        {...register('password')}
                        required
                    />
                </div>
                <div className="form-group mt-3">
                    <button className="btn btn-primary btn-block" disabled={loading} >
                        {loading && (<div className="spinner-border spinner-border-sm" role="status">
                            <span className="visually-hidden">Загрузка...</span>
                        </div>)}
                        <span>Login</span>
                    </button>
                        {/* <NavLink to={"/reset"} className="m-3">Reset password</NavLink> */}
                </div>
           </form>
            </div>
        </div>
    );};