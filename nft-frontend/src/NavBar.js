import React from 'react';
import { Link } from 'react-router-dom';
import {menu} from 'semantic-ui-react';

class NavBar extends React.Component {
  constructor(props){
    super(props);
    this.handleLoginClick = this.handleLoginClick.bind(this);
    this.handleLoginClick = this.handleLogoutClick.bind(this);
    this.state = {
      isLoggedIn: false
    };

  }
  handleLoginClick(){
    this.setState({isLoggedIn: true});
  }

  handleLogoutClick(){
    this.setState({isLoggedIn:false});
  }

  render() {
    const isLoggedIn = this.state.isLoggedIn;
    /*
    let loginButton;
    if(isLoggedIn){
      mainButton = <loginButton onClick={this.handleLogoutClick}/>
    } else{
      mainButton = <logoutButton onClick={this.handleLogoutClick}/>
    }*/
    return(
      <div>
        <div class="ui menu">
          <div class="header item">
            Our Company
          </div>
          <a class="item">
            About Us
          </a>
          <a class="item">
            Jobs
          </a>
          <a class="item">
            Login
          </a>
        </div>
      </div>
    );
  }
}

function UserGreeting(props) {
  return <h1>Welcome back!</h1>;
}

function GuestGreeting(props) {
  return <h1>Please sign up.</h1>;
}

function Greeting(props) {
  const isLoggedIn = props.isLoggedIn;
  if (isLoggedIn) {
    return <UserGreeting />;
  }
  return <GuestGreeting />;
}

function LoginButton(props) {
  return (
    <button onClick={props.onClick}>
      Login
    </button>
  );
}

function LogoutButton(props) {
  return (
    <button onClick={props.onClick}>
      Logout
    </button>
  );
}

ReactDOM.render(
  <NavBar />,
  document.getElementById('root')
);
 
export default NavBar; 