import React, {Component} from 'react';
import ReactDOM from "react-dom";
import Web3 from 'web3'
import {
  BrowserRouter as Router,
  Route,
} from 'react-router-dom';
import NavBar from './NavBar';
import HomePage from './pages/HomePage'
import Store from './pages/store';
//import detectEthereumProvier from "@metamask/detect-provider"
import {Container} from 'semantic-ui-react';

class App extends Component {
  constructor(props){
    super(props);
    this.state = {
      web3: {},
      accounts: []
    };
  }
  //this should take web3 client after the log in and pass it App's state
  pass_web3 = (web, account) => {
    this.setState({web3: web, accounts: account})
  }


  render(){
    return (
      <Container>
        <Router>
          <NavBar func ={this.pass_web3}/>
          <div className="App"> 
            <Route path="/" component={HomePage} exact/>
            <Route path="/store" exact><Store web3={this.state.web3} accounts={this.state.accounts}/></Route> 
          </div>
        </Router>
      </Container>
    );
  }
}
  
  
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
