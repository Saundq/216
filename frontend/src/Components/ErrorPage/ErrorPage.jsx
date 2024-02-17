import {useRouteError} from "react-router-dom";

export const ErrorPage=()=> {
    const error = useRouteError();
    return (
        <div>
            <h2>ErrorPage 404</h2>
            <p>{error?.message || 'Неизвестная ошибка'}</p>
        </div>
    )}