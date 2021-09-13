import React, {Component} from 'react';
import { Card, Grid } from 'semantic-ui-react';
import ContentCard from '../Components/ItemCard';
import { useLocation } from 'react-router-dom'
class Store extends Component {
  constructor(props) {
      super(props) //since we are extending class Table so we have to use super in order to override Component class constructor
      const web = useLocation()
      this.state = { //state is by default an object
        content: [
          { contentID: 1},
          { contentID: 2},
          { contentID: 3},
          { contentID: 4}
        ]        
      }
    }
    
    GetStore(){
      data = {"signature": signature, "account": account}
      console.log(data)
      const options = {
        method: 'POST',
        mode: 'cors',
        headers: {
          'Content-Type': 'application/json',
        },
        //credentials: 'include',
        body: JSON.stringify(data)
      };

    }
    renderItems(){
        return(
            this.state.content.map((content, index) => {
                const { contentID} = content
                return( 
                    <ContentCard contentID = {contentID} web3 ={web.web3} accounts = {web.accounts}/>
                )    
            })
        )
    }
    
    render(){
        return(
            <Card.Group itemsPerRow={2}>
                    {this.renderItems()}
            </Card.Group>
        )
    }
}

export default Store