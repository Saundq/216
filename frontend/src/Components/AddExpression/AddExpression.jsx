import style from "./AddExpression.module.scss"
import {Container} from "react-bootstrap";
import axios from "axios";
import {useState} from "react";
import CryptoJS from  "crypto-js"
import { useSelector } from 'react-redux'

export const AddExpression = () => {

    const [mathExpression, setMathExpression] = useState('');
    const [isValid, setIsValid] = useState(true); // State to track validation
    const [loading, setLoading] = useState(false); // State to track loading state
    const [responseCode, setResponseCode] = useState(null);
    const [result, setResult] = useState("");

    const handleChange = (event) => {
        const input = event.target.value;
        const regex = /^[0-9+\-\(\)\s*/]*$/;

        if (input === '' || regex.test(input)) {
            setMathExpression(input);
        }
        // Remove spaces from the input
      //  const sanitizedInput = input.replace(/\s/g, '');

        // Perform validation
        setIsValid(isValidMathExpression(input));
    };
    //const { token } = useSelector((state) => state.auth.userToken)

    const handleSubmit = async () => {
        setLoading(true);
        try {
            // Make the POST request to the API with the math expression
            const response = await axios.post('http://localhost:8181/api/v1/add/evaluation_arithmetic_expressions', { expression_string: mathExpression },{
                headers: {
                    'X-Request-ID': CryptoJS.MD5(mathExpression).toString(),
                    "Authorization":`Bearer ${localStorage.getItem('userToken')}`,
                }
            });
            console.log(response.data);
            const code = response.status;
            const result = response.data.result;
            setResult(result);
            setResponseCode(code);
            setMathExpression('');
            // Handle the response as needed
        } catch (error) {
            const code = error.response.status
            setResponseCode(code);
            setMathExpression('');
            console.error('Error submitting math expression:', error);
        }
        setLoading(false);
    };

    const isValidMathExpression = (expression) => {
        try {
            // Perform validation of the math expression
            // You can use a library like math.js to validate the expression
            // For simplicity, we'll just check for non-empty input here
            return expression !== '';
        } catch (error) {
            return false;
        }
    }

    const hash = (str) => {
            if (str == "") return 0;
            let hashString = 0;
            for (let character of str) {
                let charCode = character.charCodeAt(0);
                hashString = hashString << 5 - hashString;
                hashString += charCode;
                hashString |= hashString;
            }
            output.innerHTML += "The original string is " + str + "<br/>";
            output.innerHTML += "The hash string related to original string is " + hashString + "<br/>";
            return hashString;
    }

    return (
        <Container className={style.content}>
            <div className="row row-cols-1 row-cols-md-2 g-4">
                Добавить Выражение
                <div>
                    <div className="form-group">
                        <div className="alert alert-warning" role="alert">
                            Операнды и операторы разделяются пробелом.
                        </div>
                        <label htmlFor="expression">Выражение для вычисления:</label>
                        <input type="text" value={mathExpression} onChange={handleChange} className="form-control" id="expression"  placeholder="2 + 2 * 1"/>
                        {!isValid && <p style={{ color: 'red' }}>Invalid math expression</p>}
                    </div>
                    {responseCode===200 && <div className="alert alert-success" role="alert">
                        Ответ сервера: {responseCode} {result!="" && <p>Ответ: {result}</p>}
                    </div>}
                    {responseCode===400 && <div className="alert alert-warning" role="alert">
                        Ответ сервера: {responseCode}
                    </div>}
                    {responseCode!==null && responseCode!==200 && responseCode!==400 && <div className="alert alert-danger" role="alert">
                        Ответ сервера: {responseCode}
                    </div>}
                    <div className="form-group">
                        <button onClick={handleSubmit} disabled={!isValid || loading} className="btn btn-success">
                            {loading ? 'Loading...' : 'Submit'}
                        </button>
                    </div>
                </div>
             </div>
        </Container>
    )
}