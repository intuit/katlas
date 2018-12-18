import React, { Component } from 'react';
import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import { withStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';

import { ENTER_KEYCODE } from "../config/appConfig";
import * as queryActions from '../actions/queryActions';
import logo from './map.png';
import './Home.css';

const styles = theme => ({
  container: {
    display: 'flex',
    flexWrap: 'wrap',
    width: '30%'
  },
  textField: {
    marginLeft: theme.spacing.unit,
    marginRight: theme.spacing.unit,
  },
});

class Home extends Component {
  handleChange = event => {
    this.props.queryActions.changeQuery(event.target.value);
  };

  handleEnterPressCheck = event => {
    if(event.keyCode === ENTER_KEYCODE) {
      this.handleSubmit();
    }
  };

  handleSubmit = () => {
    //Only carryout submission if string is present
    if(this.props.query.current !== ''){
      this.props.queryActions.submitQuery();
      //TODO:DM - should we also do a fetch here? we do in menu bar for cases where the history push doesn't change route handler
      //no need to do xhr here, will do that upon a route change to /results
      this.props.history.push('/results?query=' + encodeURIComponent(this.props.query.current));
    }
  };

  render() {
    return (
      <div className="Home">
        <div>
          <h3>Welcome to Kubernetes Application Topology Browser</h3>
          <h1>K-Atlas Browser</h1>
        </div>
        <div className={this.props.classes.container}>
          <TextField
            label="Search..."
            className={this.props.classes.textField}
            fullWidth
            margin="normal"
            variant="filled"
            value={this.props.query.current}
            onChange={this.handleChange}
            onKeyUp={this.handleEnterPressCheck}
          />
        </div>
        <img src={logo} className="Home-logo-full" alt="logo"/>
      </div>
    );
  }
}

Home.propTypes = {
  queryActions: PropTypes.object,
  query: PropTypes.object
};

const mapStateToProps = state => ({query: state.query});

const mapDispatchToProps = dispatch => ({
  queryActions: bindActionCreators(queryActions, dispatch)
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(
  withStyles(styles)(
    withRouter(Home)
  )
);
