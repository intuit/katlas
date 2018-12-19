import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import _ from 'lodash';
import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import CircularProgress from '@material-ui/core/CircularProgress';

import EntityDetails from '../entityDetails/EntityDetails';
import {apiService} from "../../services/apiService";
import Graph from './Graph';
import "./Graph.css";

const DATA_FETCH_PERIOD_MS = 5000;

const styles = theme => ({
  container: {
    display: 'flex',
    flexWrap: 'wrap',
  },
  root: {
    width: '100%',
    overflowX: 'auto',
  },
    progress: {
    margin: theme.spacing.unit * 2,
  },
  progressContainer: {
    textAlign: 'center',
  },
});

class GraphContainer extends Component {
  constructor(props) {
    super(props);

    this.state = {
      data: {},
      waitingOnReq: false
    };
    this._isMounted = false;
  }

  componentDidMount() {
    this._isMounted = true;
    this.setState({waitingOnReq: true});
    this.intervalHandle = setInterval(() => this.getData(), DATA_FETCH_PERIOD_MS);
    //Delay the very first call by just one scheduling cycle so that the webfont can load first
    setTimeout(() => this.getData(), 0);
  }

  componentDidUpdate(prevProps) {
    //recognize change in URL and re-issue API request as necessary
    if (this.props.location !== prevProps.location){
      this.setState({waitingOnReq: true});
      this.getData();
    }
  }

  componentWillUnmount() {
    clearInterval(this.intervalHandle);
    this._isMounted = false;
  }

  getData = () => {
    const pathComponents = this.props.location.pathname.split('/');
    const uidParam = pathComponents[pathComponents.length - 1];

    apiService.getEntity(uidParam)
      .then(json => {
        //only update state if the objects fail lodash equality check AND
        //the component is still mounted. usually, the lifecycle methods should
        //be used directly for such things that, but in testing we're getting
        //intermittent errors that setState is being called on unmounted
        //components, without this check
        if(!_.isEqual(this.state.data, json) && this._isMounted) {
          this.setState({
            data: json,
            waitingOnReq: false
          });
        }
      });
  };

  render() {
    const { classes } = this.props;
    return (
      <div>
        {//selectively show the progress indicator when we're waiting for an outstanding request
          this.state.waitingOnReq ? (
            <div className={classes.progressContainer}>
              <CircularProgress className={classes.progress} color='secondary'/>
            </div>
          ) : null //TODO:DM - ideally want to select between spinner and graph here, but loading graph only after spinner disappears leads to graph not being correctly populated; understand why this is
        }
        <Grid container>
          <Grid item sm={12} md={9} lg={8} className='Graph-scroll-container'>
            <Graph dataSet={this.state.data}/>
          </Grid>
          <Grid item sm={12} md={3} lg={4} className='Graph-scroll-container'>
            <EntityDetails selectedObj={this.state.data.objects ? this.state.data.objects[0] : {}}/>
          </Grid>
        </Grid>
      </div>
    );
  }
}

GraphContainer.propTypes = {
  classes: PropTypes.object.isRequired,
};

export default withRouter(withStyles(styles)(GraphContainer));