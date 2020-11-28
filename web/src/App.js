import React from 'react';
import {BrowserRouter as Router} from 'react-router-dom';
import Routes from './Routes';
import './App.css';
import history from "./history";

const App = () => (
    <Router history={history}>
        <Routes/>
    </Router>
);

export default App;
