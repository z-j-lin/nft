import React from 'react';
import { Link } from 'react-router-dom';


const NavBar = (props) => (
  <nav>
    <ul>
      <li><Link to = "/">Home</Link></li>
      <li><Link to="/hello">hello</Link></li>
      <li><span id = "account">{props.account}</span> </li>
    </ul>
  </nav>
);

export default NavBar; 