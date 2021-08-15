import React, {Component} from 'react';
import Web3 from 'web3'
import './App.css';
import {
  BrowserRouter as Router,
  Route,
} from 'react-router-dom';
import Login from './pages/Login';
import NavBar from './NavBar';
//import detectEthereumProvier from "@metamask/detect-provider"

class App extends Component {

  constructor(props){
    super(props)
  }

  render() {
    return (
      <Router>
        <div className="App">
          <NavBar account = {this.state.account}/>
          <div id = "page-body">
            {/*homepage */}
            {/*<Router path ="/" component={}/>*/}
            <Route path = "/login" component={Login} exact/>
          </div>
        </div>
      </Router>
      
    );
  }
}

export default App;
