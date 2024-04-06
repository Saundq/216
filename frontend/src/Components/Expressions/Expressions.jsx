import style from "./Expressions.module.scss"
import {Container} from "react-bootstrap";
import {useEffect, useState} from "react";
import { useSelector } from 'react-redux'


export const Expressions = () => {

    const [data, setData] = useState(null);
    const { token } = useSelector(state => state.auth)

    useEffect(() => {
       
        async function fetchData() {
            try {
                const response = await fetch('http://localhost:8181/api/v1/arithmetic_expressions',{
                    headers: {
                        "Authorization":`Bearer ${sessionStorage.getItem('userToken')}`,
                        "Content-Type": "application/json",
                      },
                });
                if (!response.ok) {
                    throw new Error('Network error');
                }
                const jsonData = await response.json();
                setData(jsonData);
            } catch (error) {
                console.error('Error fetching data:', error);
            }
        }

        fetchData();
        const interval = setInterval(fetchData, 5000);
        return () => clearInterval(interval);
    }, [token]);

    return (
        <Container className={style.content}>
            <div className="row row-cols-1 row-cols-md-2 g-4">
                Выражения
                {data ? (

                    <table className="table">
                        <thead>
                        <tr>
                            <th scope="col"></th>
                            <th scope="col">#</th>
                            <th scope="col">Выражение</th>
                            <th scope="col">Дата создания</th>
                            <th scope="col">Дата вычисления</th>
                            <th scope="col">Статус</th>
                            <th scope="col">Результат</th>
                        </tr>
                        </thead>
                        <tbody>
                        {data.map((item, index) => {
                            return (
                                <>
                            <tr key={item.id} className={item.result.length>0&& item.parent==null?"table-success":""} >
                                <td></td>
                                <td>{item.id}</td>
                                <td>{item.expression_string}</td>
                                <td>{new Date(item.added_at * 1000).toLocaleDateString("ru-ru")} {new Date(item.added_at * 1000).toLocaleTimeString("ru-ru")}</td>
                                <td>{new Date(item.finished_at * 1000).toLocaleDateString("ru-ru")!=="01.01.1970"?new Date(item.finished_at * 1000).toLocaleDateString("ru-ru"):""} {new Date(item.finished_at * 1000).toLocaleTimeString("ru-ru")!=="03:00:00"?new Date(item.finished_at * 1000).toLocaleTimeString("ru-ru"):""}</td>
                                <td>{item.status}</td>
                                <td>{item.result}</td>
                            </tr>
                                {item.ExpressionPart.map((v,k) =>{
                                    return (
                                        <tr key={k}>
                                            <td><svg xmlns="http://www.w3.org/2000/svg" width="50" height="50" fill="currentColor" className="bi bi-arrow-return-right" viewBox="0 0 16 16">
                                                <path fillRule="evenodd" d="M1.5 1.5A.5.5 0 0 0 1 2v4.8a2.5 2.5 0 0 0 2.5 2.5h9.793l-3.347 3.346a.5.5 0 0 0 .708.708l4.2-4.2a.5.5 0 0 0 0-.708l-4-4a.5.5 0 0 0-.708.708L13.293 8.3H3.5A1.5 1.5 0 0 1 2 6.8V2a.5.5 0 0 0-.5-.5"/>
                                            </svg></td>
                                            <td>{v.id}</td>
                                            <td>{v.expression_string } {v.next!==null ? " Вычисляется после "+v.next:""} {v.previous!=null ? "Вычисляется после "+v.previous:""}</td>
                                            <td>{new Date(v.added_at * 1000).toLocaleDateString("ru-ru")} {new Date(v.added_at * 1000).toLocaleTimeString("ru-ru")}</td>
                                            <td>{new Date(v.finished_at * 1000).toLocaleDateString("ru-ru")!=="01.01.1970"?new Date(v.finished_at * 1000).toLocaleDateString("ru-ru"):""} {new Date(v.finished_at * 1000).toLocaleTimeString("ru-ru")!=="03:00:00"?new Date(v.finished_at * 1000).toLocaleTimeString("ru-ru"):""}</td>
                                            <td>{v.status}</td>
                                            <td>{v.result}</td>
                                        </tr>
                                    )
                                })}
                                </>
                            )})

                        }
                        </tbody>
                    </table>
                ) : (
                    <p>Loading...</p>
                )}
             </div>
        </Container>
    )
}