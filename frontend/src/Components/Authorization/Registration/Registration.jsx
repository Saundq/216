import style from "./Registration.module.scss"
import {Navigate, NavLink, useLocation, useNavigate} from "react-router-dom";
import {useEffect, useRef, useState} from "react";
import {useDispatch, useSelector} from "react-redux";
import {userLogin,registerUser} from "../../../actions/authActions.js";
import {useForm} from 'react-hook-form'

export const Registration =()=> {

    let navigate = useNavigate();
    const {loading,userToken,userInfo,error,success} = useSelector((state)=>state.auth)
    const dispatch = useDispatch();

    const [reg,setReg]=useState(false);

    const {register, handleSubmit,reset} = useForm()
  

    const location = useLocation()
  //  const params = new URLSearchParams(location.search)

    useEffect(() => {
        if(reg)
            navigate('/login')
        // if(userInfo)
        //     navigate('/profile')
    }, [navigate,reg]);

    const submitForm = (data) => {
       // data.token=params.get("registration_token")
        console.log(data)
        if(data.password==data.password_confirmation){
            dispatch(registerUser(data))
        }

        // //console.log(data.email)
        // dispatch(userLogin(data))
        console.log(data)
         reset()
         setReg(true);
    }
    const getError = (err)=>{
        if (err){
            console.log(err)
            if (err.constructor === String) {
                console.log("строка")
                return <p className="alert alert-danger">{err}</p>
            }
            let message=[];
            for (const [key, value] of Object.entries(err)) {
                //message+=key;
                Array.from(value).forEach(item=>{
                    message.push(<p className="alert alert-danger">{item}</p>);
                    console.log(item)
                });
            }
             console.log(message)
             return message;
        } else return "";
    }

    return (
        <div className="col-md-4 m-md-auto mt-md-5">
            <div className="card card-container m-3 p-3">
                <h1 className="m-md-auto">Регистрация:</h1>
            <form onSubmit={handleSubmit(submitForm)} >
                {error && <p>{getError(error)}</p>}
                <div className="form-group">
                    <label htmlFor="name">Name</label>
                    <input
                        type="text"
                        className="form-control"
                        {...register('name')}
                        required
                    />
                </div>
                <div className="form-group">
                    <label htmlFor="email">Email</label>
                    <input
                        type="text"
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
                <div className="form-group">
                    <label htmlFor="password">Confirm Password</label>
                    <input
                        type="password"
                        className="form-control"
                        {...register('password_confirmation')}
                        required
                    />
                </div>
                <div className="form-group mt-3">
                    <button className="btn btn-primary btn-block" disabled={loading} >
                        {loading && (<div className="spinner-border spinner-border-sm" role="status">
                            <span className="visually-hidden">Загрузка...</span>
                        </div>)}
                        <span>Register</span>
                    </button>
                </div>
           </form>
            </div>
        </div>
    );};