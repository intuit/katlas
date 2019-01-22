import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Link, withRouter } from 'react-router-dom';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import IconButton from '@material-ui/core/IconButton';
import Typography from '@material-ui/core/Typography';
import { withStyles } from '@material-ui/core/styles';

import logo from './map.png';

const styles = theme => ({
  root: {
    width: '100%'
  },
  menuButton: {
    marginLeft: -12,
    marginRight: 20
  },
  title: {
    display: 'none',
    [theme.breakpoints.up('sm')]: {
      display: 'block'
    }
  },
  appLogoSmall: {
    height: '36px',
    width: 'auto'
  }
});

class MenuBar extends Component {
  render() {
    const { classes } = this.props;
    return (
      <AppBar className={classes.root} position='static'>
        <Toolbar>
          <IconButton className={classes.menuButton}>
            <Link to='/'>
              <img src={logo} className={classes.appLogoSmall} alt='logo' />
            </Link>
          </IconButton>
          <Typography
            className={classes.title}
            variant='h6'
            color='inherit'
            noWrap
          >
            K-Atlas Browser
          </Typography>
        </Toolbar>
      </AppBar>
    );
  }
}

MenuBar.propTypes = {
  classes: PropTypes.object.isRequired,
  queryActions: PropTypes.object,
  query: PropTypes.object
};

export default withStyles(styles)(withRouter(MenuBar));
