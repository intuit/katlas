import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import { withStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';

import './Home.css';
import logo from './map.png';
import { ENTER_KEYCODE, ENTER_KEYSTR } from '../../config/appConfig';
import * as queryActions from '../../actions/queryActions';

const styles = theme => ({
  container: {
    display: 'flex',
    flexWrap: 'wrap',
    width: '30%'
  },
  title: {
    textAlign: 'center'
  },
  textField: {
    marginLeft: theme.spacing.unit,
    marginRight: theme.spacing.unit
  }
});

class Home extends Component {
  state = {
    queryStr: ''
  };

  handleChange = event => {
    this.setState({ queryStr: event.target.value });
  };

  handleEnterPressCheck = event => {
    //Check both for keycode (used in practice) and a key string which is all
    //the testing framework can seem to do
    if (event.keyCode === ENTER_KEYCODE || event.key === ENTER_KEYSTR) {
      this.handleSubmit();
    }
  };

  handleSubmit = () => {
    const { queryStr } = this.state;
    //Validate query in submitQuery and decide to switch to /results based on query validation.
    this.props.queryActions.submitQuery(queryStr);
  };

  render() {
    const { classes } = this.props;
    const { queryStr } = this.state;
    return (
      <div className='Home'>
        <div className={classes.title}>
          <h3>Welcome to Kubernetes Application Topology Browser</h3>
          <h1>K-Atlas Browser</h1>
        </div>
        <div className={classes.container}>
          <TextField
            label='Search...'
            className={classes.textField}
            fullWidth
            margin='normal'
            variant='filled'
            value={queryStr}
            onChange={this.handleChange}
            onKeyPress={this.handleEnterPressCheck}
          />
        </div>
        <img src={logo} className='Home-logo-full' alt='logo' />
      </div>
    );
  }
}

Home.propTypes = {
  queryActions: PropTypes.object,
  query: PropTypes.object
};

const mapStoreToProps = store => ({ query: store.query });

const mapDispatchToProps = dispatch => ({
  queryActions: bindActionCreators(queryActions, dispatch)
});

export default connect(
  mapStoreToProps,
  mapDispatchToProps
)(withStyles(styles)(withRouter(Home)));
