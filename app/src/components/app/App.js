import React, { Component } from "react";
import PropTypes from 'prop-types';
import { Route, Switch } from "react-router-dom";
import Paper from '@material-ui/core/Paper';
import { withStyles } from '@material-ui/core/styles';

import MenuBar from "../menuBar/MenuBar";
import Home from "../home/Home";
import Results from "../results/Results";
import GraphContainer from "../graph/GraphContainer";

import './App.css';

const styles = theme => ({
  paper: {
    padding: theme.spacing.unit * 2,
    textAlign: 'center',
    color: theme.palette.text.primary,
  },
});

class App extends Component {
  render() {
    const { classes } = this.props;
    return (
      <div className="App">
        <MenuBar/>
        <Paper className={classes.paper}>
          <Switch>
            <Route exact path="/" component={Home} />
            <Route path="/results" component={Results} />
            <Route path="/graph" component={GraphContainer} />
          </Switch>
        </Paper>
      </div>
    );
  }
}

App.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(App);