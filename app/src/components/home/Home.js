import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import { withStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';

import './Home.css';
import logo from './katlas-logo-blue-300px.png';
import { ENTER_KEYCODE, ENTER_KEYSTR } from '../../config/appConfig';
import * as queryActions from '../../actions/queryActions';
import { validateQslQuery } from '../../utils/validate';

const styles = theme => ({
  container: {
    display: 'flex',
    flexWrap: 'wrap',
    width: '80%'
  },
  title: {
    textAlign: 'center'
  },
  textField: {
    marginLeft: 0,
    marginRight: 0,
  },
  leftButton: {
    width: '55%',
    margin: 0,
    marginTop: -theme.spacing.unit,
    borderRadius: 0,
    borderBottomLeftRadius: theme.spacing.unit,
  },
  rightButton: {
    width: '45%',
    margin: 0,
    marginTop: -theme.spacing.unit,
    borderRadius: 0,
    borderBottomRightRadius: theme.spacing.unit,
  },
  marginLeft: theme.spacing.unit,
  marginRight: theme.spacing.unit
});

class Home extends React.Component {
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

  handleQslSubmit = () => {
    const { queryStr } = this.state;
    this.props.queryActions.submitQslQuery(queryStr);
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
            label='Search string or QSL...'
            className={classes.textField}
            fullWidth
            margin='normal'
            variant='filled'
            value={queryStr}
            onChange={this.handleChange}
            onKeyPress={this.handleEnterPressCheck}
          />
        </div>
        <div className={classes.container}>
          <Button variant='contained' color='secondary' className={classes.leftButton}
            disabled={!validateQslQuery(queryStr)} onClick={this.handleQslSubmit}>
            I'm Feeling Graphy!
          </Button>
          <Button variant='contained' color='primary' className={classes.rightButton}
            onClick={this.handleSubmit}>
            Search
          </Button>
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
