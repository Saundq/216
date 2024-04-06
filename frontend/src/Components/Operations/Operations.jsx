import style from "./Operations.module.scss"
import {Container} from "react-bootstrap";
import {useEffect, useState} from "react";
import axios from "axios";

export const Operations = () => {
    const [data, setData] = useState(null);
    const [editedTime, setEditedTime] = useState("");
    const [mathExpression, setMathExpression] = useState('');
    const fetchUrl = (URL) => {
        fetch(URL,{
            headers: {
                "Authorization":`Bearer ${sessionStorage.getItem('userToken')}`,
                "Content-Type": "application/json",
              },
        })
        .then((res)=>{
            return res.json();
        })
        .then((data)=>{
            setData(data);
        })
    }

    useEffect(()=>{
        
        fetchUrl("http://localhost:8181/api/v1/arithmetic_operations");
    },[URL])
    
    const handleSubmit = (e,id) => {
        e.preventDefault();
       console.log(id)

       axios.put(`http://localhost:8181/api/v1/arithmetic_operations/${id}`, { lead_time: parseInt(editedTime.time) },{headers: {
            'Content-Type':'application/json',
            "Authorization":`Bearer ${sessionStorage.getItem('userToken')}`,
        },})
            .then(response => {
             
                setData(data.map(item => (item.id === id ? { ...item, time: editedTime } : item )));
                fetchUrl("http://localhost:8181/api/v1/arithmetic_operations");
            })
            .catch(error => {
                console.error('Ошибка отправки изменений:', error);
            });
            setEditedTime("");
            
      };

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
                                <form onSubmit={(e)=>handleSubmit(e,item.id)}>
                                    <input
                                        type="text"
                                        value={item.id === editedTime.id ?editedTime.time:item.lead_time }
                                        onChange={(e) => setEditedTime({ id: item.id, time: e.target.value })}
                                    />

                                    {item.id === editedTime.id && (
                                        <button  className="btn btn-success mx-3">Сохранить</button>
                                    )}
                                    </form>
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