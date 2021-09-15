import React, {Component} from 'react';
import { Card, Grid } from 'semantic-ui-react';
import ContentCard from '../Components/ItemCard';
import { useLocation } from 'react-router-dom'
class Store extends Component {
  constructor(props) {

    super(props) //since we are extending class Table so we have to use super in order to override Component class constructor
      
    this.state = { //state is by default an object
      accounts: props.accounts, 
      web3: props.web3,
      content: []        
    }
  }
  componentDidMount(){
    this.GetStore()
  }
    
  GetStore(){
    const backendurl = 'http://127.0.0.1:8081/';
    fetch(backendurl+"getstore")
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
          <ContentCard key = {index} contentID={content} web3 ={this.state.web3} accounts={this.state.accounts}/>
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

export default Store