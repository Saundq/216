import style from './Main.module.scss'
import {Container} from "react-bootstrap";
export  const Main =({children})=>(
        <div className={style.main}>
            {children}
        </div>

);