import React, {Component} from 'react';

import { Button, Card, Image } from 'semantic-ui-react'

class ContentCard extends Component{
  constructor(props) {
    super(props);
    this.state = {
      contentID: props.contentID,
      account: props.accounts,
      isToggleOn: true,
      web3: props.web3
    };

    this.handleClick = this.handleClick.bind(this);
  }
  buyToken(){
    console.log(this.state.account.toString())
    const backendurl = 'http://127.0.0.1:8081/';
    const data = {"resourceid": this.state.contentID, "account": this.state.account.toString()}
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
    fetch(backendurl+'buy', options) 
  }
  handleClick() {
    this.setState(prevState => ({
      isToggleOn: !prevState.isToggleOn
    }));
    //send a post request to the api with contentID and account address 
    this.buyToken()
  }
  render(){
    return(
      <Card key = {this.props.contentID}>
        <Card.Content>
          <Card.Header>{this.props.contentID}</Card.Header>
          <Card.Description>
            Lit content
          </Card.Description>
        </Card.Content>
        <Card.Content extra>
          <div className='ui two buttons'>
            <Button basic color='green' onClick={this.handleClick}>
              Purchase
            </Button>
          </div>
        </Card.Content>
      </Card>
    );
  }
  
};

export default ContentCard;