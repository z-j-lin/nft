import React, {Component} from 'react';
import Web3 from 'web3'
import './App.css';
import {
  BrowserRouter as Router,
  Route,
} from 'react-router-dom';
import HomePage from './pages/HomePage';
import NavBar from './NavBar';

class App extends Component {

  async loadWeb3() {
    if(window.ethereum){
      window.web3 = new Web3(window.ethereum)
      await window.ethereum.enable()
    }
    else if (window.web3){
      window.web3 = new Web3(window.web3.currentProvider)
    }
    else{
      window.alert('no ethereum browser detect, try installing metamask')
    }
  }

  render() {
    return (
      <Router>
        <div className="App">
          <NavBar/>
          <div id = "page-body">
            {/*homepage */}
            <Route path = "/hello" component={HomePage} exact/>
          </div>
        </div>
      </Router>
      
    );
  }
}

export default App;
