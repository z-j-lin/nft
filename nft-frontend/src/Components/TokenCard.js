//This file is Copyright (C) 1997 Master Hacker, ALL RIGHTS RESERVED 
import React, {Component} from 'react';

import { Button, Card, Image, Popup } from 'semantic-ui-react'

class TokenCard extends Component{
  constructor(props) {
    super(props);
    this.state = {
      open: false,
      account: props.accounts,
      isToggleOn: true,
      web3: props.web3,
      url: ""
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
    .then(response => {
      if(!response.ok){
        throw Error("Error fetching resource")
      }
      return response.json()
    })
    .then(jsonresp => {
      this.setState({url: jsonresp.url})
    }).catch(error => {console.log(error)})
  }
  handleClick() {

    if (this.state.open === false){
      this.setState({open: true})
    }else{
      this.setState({open: false})
    }
    this.setState(prevState => ({
      isToggleOn: !prevState.isToggleOn
    }));
    //send a post request to the api with contentID and account address 
    if (this.state.url === ""){
      this.AccessToken()
    }
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
          <Popup
          on='hover'
          //open={this.state.open}
          trigger={<Button content='A trigger' basic color='green' onClick={this.handleClick}/>}
        >
          <Popup.Content>
            <Image src={this.state.url}/>
          </Popup.Content>
        </Popup>
          </div>
        </Card.Content>
      </Card>
    );
  }
  
};

export default TokenCard;