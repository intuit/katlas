import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ReactJson from 'react-json-view';

import './EntityDetails.css';

//This component ended up being a pretty thin wrapper around react-json-view
//3rd party comp. But still valuable as a separate component since it will be
//used in multiple routes/views of the app.
export default class EntityDetails extends Component {
  render() {
    return (
      <div className="Details">
        <ReactJson src={this.props.selectedObj} />
      </div>
    );
  }
}

EntityDetails.propTypes = {
  selectedObj: PropTypes.object,
};