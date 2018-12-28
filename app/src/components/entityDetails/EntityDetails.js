import React, { Component } from 'react';
import { withStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import PropTypes from 'prop-types';
import ReactJson from 'react-json-view';

import './EntityDetails.css';

const styles = theme => ({
  root: {
    width: '100%',
    overflowX: 'auto',
  },
});

//This component ended up being a pretty thin wrapper around react-json-view
//3rd party comp. But still valuable as a separate component since it will be
//used in multiple routes/views of the app.
class EntityDetails extends Component {
  render() {
    const { classes } = this.props;

    return (
      <Paper className={classes.root}>
        <ReactJson src={this.props.selectedObj} theme="ocean"/>
      </Paper>
    );
  }
}

EntityDetails.propTypes = {
  selectedObj: PropTypes.object,
};

export default withStyles(styles)(EntityDetails);