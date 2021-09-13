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

  async componentWillMount() {
    await this.loadWeb3();
    await this.loadBlockchainData();
  };
  async loadWeb3() {
    if(window.ethereum){
      window.web3 = new Web3(window.ethereum);
      await window.ethereum.enable();
    }
    else if (window.web3){
      window.web3 = new Web3(window.web3.currentProvider)
    }
    else{
      window.alert('no ethereum browser detect, try installing metamask')
    }
  };

  async  loadBlockchainData() {
    const web3 = window.web3;
    console.log(web3)
    //load account
    const accounts = await web3.eth.getAccounts()
    this.setState({
      web3: web3,
      accounts: accounts
    })  
  };



  render(){
    return (
      <Container>
        <NavBar/>
        <Router>
          <div className="App"> 
            <Route path="/" component={HomePage} exact/>
            <Route path="/store" component={Store} exact/>
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
