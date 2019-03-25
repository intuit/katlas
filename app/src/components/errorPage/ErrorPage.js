import React from 'react';
import { withRouter } from 'react-router-dom';
import { withStyles } from '@material-ui/core/styles';

import './ErrorPage.css';
import logo from './katlas-logo-blue-300px.png';

const styles = () => ({
  title: {
    textAlign: 'center'
  }
});

class ErrorPage extends React.Component {
  render() {
    const { classes } = this.props;
    return (
      <div className='Error'>
        <div className={classes.title}>
          <h3>Something appears to have gone wrong. Please retry your operation.</h3>
        </div>
        <img src={logo} className='Error-logo-full' alt='logo' />
      </div>
    );
  }
}

export default withStyles(styles)(withRouter(ErrorPage));
