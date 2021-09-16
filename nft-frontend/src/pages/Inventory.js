import React, {Component} from 'react';
import { Card, Grid } from 'semantic-ui-react';
import TokenCard from '../Components/TokenCard';

class Inventory extends Component {
  constructor(props) {

    super(props) //since we are extending class Table so we have to use super in order to override Component class constructor
      
    this.state = { //state is by default an object
      accounts: props.accounts, 
      web3: props.web3,
      content: []        
    }
  }
  componentDidMount(){
    this.GetInventory()
  }
    
  GetInventory(){
    const backendurl = 'http://127.0.0.1:8081/';
    const data = {"account": this.state.accounts.toString()}
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
    fetch(backendurl+"load")
      .then(response => response.json())
      .then(data => {
        console.log(data)
        this.setState({
          content: data
        })
      })  
  };
  renderItems(){
    return(
      this.state.content.map((content, index) => {
        console.log(this.state.accounts)
        return( 
          <TokenCard key = {index} TokenID={content} web3 ={this.state.web3} accounts={this.state.accounts}/>
          )    
      })
    )
  }

  render(){
    console.log(this.state.accounts)
    return(      
      <Card.Group itemsPerRow={2}>
        {this.renderItems()}
      </Card.Group>
    )
  }
}

export default Inventory