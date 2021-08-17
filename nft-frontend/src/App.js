import React, {Component} from 'react';
import ReactDOM from "react-dom";
import Web3 from 'web3'
import {
  BrowserRouter as Router,
  Route,
} from 'react-router-dom';
import NavBar from './NavBar';
//import detectEthereumProvier from "@metamask/detect-provider"
import {Container} from 'semantic-ui-react';

const App = () => (
  <Container>
    <Router>
      <div className="App">
          <NavBar/>
        <div id = "page-body">
          {/*homepage */}

        </div>
      </div>
    </Router>
  </Container>

);
  
const styleLink = document.createElement("link");
styleLink.rel = "stylesheet";
styleLink.href = "https://cdn.jsdelivr.net/npm/semantic-ui/dist/semantic.min.css";
document.head.appendChild(styleLink); 

ReactDOM.render( 
  <App>
    <NavBar />
  </App>,
  document.getElementById("root")
);
export default App;
