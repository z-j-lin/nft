import React, {Component} from 'react';

import { Button, Card, Image } from 'semantic-ui-react'

class TokenCard extends Component{
  constructor(props) {
    super(props);
    this.state = {
      account: props.accounts,
      isToggleOn: true,
      web3: props.web3
    };

    this.handleClick = this.handleClick.bind(this);
  }
  AccessToken(){
    console.log(this.state.account.toString())
    const backendurl = 'http://127.0.0.1:8081/';
    const data = {"tokenid": this.props.TokenID, "account": this.state.account.toString()}
    console.log(data)
    const options = {
      method: 'POST',
      mode: 'cors',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      //credentials: 'include',
      body: JSON.stringify(data)
    };
    console.log(options)
    fetch(backendurl+'request', options) 
  }
  handleClick() {
    this.setState(prevState => ({
      isToggleOn: !prevState.isToggleOn
    }));
    //send a post request to the api with contentID and account address 
    this.AccessToken()
  }
  render(){
    return(
      <Card key = {this.props.TokenID}>
        <Card.Content>
          <Card.Header>{this.props.TokenID}</Card.Header>
          <Card.Description>
            Lit content
          </Card.Description>
        </Card.Content>
        <Card.Content extra>
          <div className='ui two buttons'>
            <Button basic color='green' onClick={this.handleClick}>
              Access
            </Button>
          </div>
        </Card.Content>
      </Card>
    );
  }
  
};

export default TokenCard;