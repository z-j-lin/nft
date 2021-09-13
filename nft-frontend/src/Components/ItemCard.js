import React, {Component} from 'react';

import { Button, Card, Image } from 'semantic-ui-react'

class ContentCard extends Component{
  constructor(props) {
    super(props);
    this.state = {isToggleOn: true};

    this.handleClick = this.handleClick.bind(this);
  }
  handleClick() {
    this.setState(prevState => ({
      isToggleOn: !prevState.isToggleOn
    }));
    
    //send a post request to the api 

  }
  render(){
    return(
      <Card>
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