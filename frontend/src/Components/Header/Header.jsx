import {Top} from "./Top/Top.jsx";
import style from "./Header.module.scss";

export const Header = () => (
    <header className={style.header}>
        <Top/>
    </header>
)