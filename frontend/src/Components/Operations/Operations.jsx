import style from "./Operations.module.scss"
import {Container} from "react-bootstrap";
import {useEffect, useState} from "react";
import axios from "axios";

export const Operations = () => {
    const [data, setData] = useState(null);
    const [editedTime, setEditedTime] = useState("");

    useEffect(() => {
        async function fetchData() {
            try {
                const response = await fetch('http://localhost:8181/api/v1/arithmetic_operations');
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
    }, []);


    const handleEditTime = (id) => {
        // Отправка изменения времени на сервер
        axios.put(`http://localhost:8181/api/v1/arithmetic_operations/${id}`, { lead_time: parseInt(editedTime.time) })
            .then(response => {
                // Обновление данных после успешной отправки изменений
                setData(data.map(item => (item.id === id ? { ...item, time: editedTime } : item )));
            })
            .catch(error => {
                console.error('Ошибка отправки изменений:', error);
            });
        window.location.reload(false);
    }

    return (
        <Container className={style.content}>
            <div className="row row-cols-1 row-cols-md-2 g-4">
                Операции
                {data ? (
                    <table className="table">
                        <thead>
                        <tr>
                            <th scope="col">Операция</th>
                            <th scope="col">Время выполнения</th>
                        </tr>
                        </thead>
                        <tbody>
                        {data.map((item, index) => (
                            <tr key={index}>
                                <td>{item.value}</td>
                                <td>
                                    <input
                                        type="text"
                                        value={item.id === editedTime.id ? editedTime.lead_time : item.lead_time }
                                        onChange={(e) => setEditedTime({ id: item.id, time: e.target.value })}
                                    />
                                    {item.id === editedTime.id && (
                                        <button onClick={() => handleEditTime(item.id)} className="btn btn-success">Сохранить</button>
                                    )}
                                </td></tr>
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