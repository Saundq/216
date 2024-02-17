import style from "./ComputionCapabilities.module.scss"
import {Container} from "react-bootstrap";
import {useState, useEffect } from "react";

export const ComputionCapabilities = () => {
    const [data, setData] = useState(null);

    useEffect(() => {
        async function fetchData() {
            try {
                const response = await fetch('http://localhost:8181/api/v1/available_calculators');
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
        const interval = setInterval(fetchData, 1000);
        return () => clearInterval(interval);
    }, []);  // пустой массив зависимостей указывает, что эффект должен выполниться только один раз после монтирования компонента

    return (
        <Container className={style.content}>
            <div className="row row-cols-1 row-cols-md-2 g-4">
                Вычислительные можности

                {data ? (
                    <table className="table">
                        <thead>
                        <tr>
                            <th scope="col">#</th>
                            <th scope="col">Наименование</th>
                            <th scope="col">Состояние</th>
                            <th scope="col">Выполняется</th>
                        </tr>
                        </thead>
                        <tbody>
                        {data.map((item, index) => (
                            <tr key={index}><td>{item.Id}</td><td>{item.Name}</td><td>{item.Status}</td><td>{item.Task} ({item.TaskStr})</td></tr>
                        ))}
                        </tbody>
                    </table>
                ) : (
                    <p>Loading...</p>
                )}
             </div>
        </Container>
    )
}