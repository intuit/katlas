import React, { Component } from 'react';
import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';
import PropTypes from 'prop-types';
import { Link, withRouter } from 'react-router-dom';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import IconButton from '@material-ui/core/IconButton';
import Typography from '@material-ui/core/Typography';
import InputBase from '@material-ui/core/InputBase';
import { fade } from '@material-ui/core/styles/colorManipulator';
import { withStyles } from '@material-ui/core/styles';
import SearchIcon from '@material-ui/icons/Search';

import { ENTER_KEYCODE } from '../../config/appConfig';
import * as queryActions from '../../actions/queryActions';
import logo from './map.png';

const styles = theme => ({
  root: {
    width: '100%',
  },
  grow: {
    flexGrow: 1,
  },
  menuButton: {
    marginLeft: -12,
    marginRight: 20,
  },
  title: {
    display: 'none',
    [theme.breakpoints.up('sm')]: {
      display: 'block',
    },
  },
  search: {
    position: 'relative',
    borderRadius: theme.shape.borderRadius,
    backgroundColor: fade(theme.palette.common.white, 0.15),
    '&:hover': {
      backgroundColor: fade(theme.palette.common.white, 0.25),
    },
    marginLeft: 0,
    width: '100%',
    [theme.breakpoints.up('sm')]: {
      marginLeft: theme.spacing.unit,
      width: 'auto',
    },
  },
  searchIcon: {
    width: theme.spacing.unit * 5,
    height: '100%',
    position: 'absolute',
    pointerEvents: 'auto',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
  inputRoot: {
    color: 'inherit',
    width: '100%',
  },
  inputInput: {
    paddingTop: theme.spacing.unit,
    paddingRight: theme.spacing.unit,
    paddingBottom: theme.spacing.unit,
    paddingLeft: theme.spacing.unit * 5,
    transition: theme.transitions.create('width'),
    width: '100%',
    [theme.breakpoints.up('sm')]: {
      width: 120,
      '&:focus': {
        width: 200,
      },
    },
  },
  appLogoSmall: {
    height: '36px',
    width: 'auto',
  },
});

class MenuBar extends Component {
  handleChange = event => {
    //TODO:DM - value in doing search as you type here?
    this.props.queryActions.changeQuery(event.target.value);
  };

  handleEnterPressCheck = event => {
    if(event.keyCode === ENTER_KEYCODE && this.props.query.current !== '') {
      this.handleSubmit();
    }
  };

  //TODO:DM-why isn't this working for icon clicks directly?
  handleSubmit = () => {
    this.props.queryActions.submitQuery(this.props.query.current);
  };

  render() {
    const { classes } = this.props;
    return (
      <AppBar className={classes.root} position='static'>
        <Toolbar>
          <IconButton className={classes.menuButton}>
            <Link to="/">
              <img src={logo} className={classes.appLogoSmall} alt="logo"/>
            </Link>
          </IconButton>
          <Typography className={classes.title} variant="h6" color="inherit" noWrap>
            K-Atlas Browser
          </Typography>
          <div className={classes.grow} />
          <div className={classes.search}>
            <div className={classes.searchIcon} onClick={this.handleSubmit}>
              <SearchIcon />
            </div>
            <InputBase
              placeholder="Search string…"
              classes={{
                root: classes.inputRoot,
                input: classes.inputInput,
              }}
              value={this.props.query.current}
              onChange={this.handleChange}
              onKeyUp={this.handleEnterPressCheck}
            />
          </div>
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

const mapStateToProps = state => ({query: state.query});

const mapDispatchToProps = dispatch => ({
  queryActions: bindActionCreators(queryActions, dispatch)
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(
  withStyles(styles)(
    withRouter(MenuBar)
  )
);
