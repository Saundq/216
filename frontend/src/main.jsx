import React from 'react'
import ReactDOM from 'react-dom/client'
import 'bootstrap/dist/js/bootstrap.bundle.min';
import "bootstrap/dist/css/bootstrap.min.css"
//import 'normalize.css'
import './index.scss'
import { App } from "./App.jsx";



ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
          <App/>
  </React.StrictMode>,
)
