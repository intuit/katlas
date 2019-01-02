import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import CircularProgress from '@material-ui/core/CircularProgress';

import "./Graph.css";
import Graph from './Graph';
import EntityDetails from '../entityDetails/EntityDetails';
import * as entityActions from "../../actions/entityActions";

const FETCH_PERIOD_PER_ENTITY_MS = 2000;

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
  componentDidMount() {
    this.setRootNode();
    this.intervalHandle = setTimeout(() => this.getData(),
      FETCH_PERIOD_PER_ENTITY_MS);
    //delay the first call by no ms (therefore one scheduling cycle) which
    //allows for webfont to be loaded before data inserted to graph
    setTimeout(() => this.getData(), 0);
  }

  componentDidUpdate(prevProps) {
    //recognize change in URL and re-issue API request as necessary
    if (this.props.location !== prevProps.location){
      this.setRootNode();
      this.getData();
    }
  }

  componentWillUnmount() {
    clearInterval(this.intervalHandle);
  }

  setRootNode = () => {
    const pathComponents = this.props.location.pathname.split('/');
    //TODO:DM - simply grabbing last param after '/' feels fragile, how to more safely verify as UID?
    //could be empty string... a better default to use, if so?
    const uid = pathComponents[pathComponents.length - 1];

    this.props.entityActions.setRootEntity(uid);
    this.props.entityActions.addEntityWatch(uid);

  };

  getData = () => {
    //reschedule next automatic data request while computing time value based
    //on number of entities and a min time between fetches
    const NUM_ENTITIES = Object.keys(this.props.entity.entitiesByUid).length;
    this.intervalHandle = setTimeout(() => this.getData(), NUM_ENTITIES * FETCH_PERIOD_PER_ENTITY_MS);
    //fetch all entities currently represented as keys in the store
    this.props.entityActions.fetchEntities(Object.keys(
      this.props.entity.entitiesByUid));
  };

  render() {
    const { classes, entity } = this.props;
    return (
      <div>
        {//selectively show the progress indicator when we're waiting for an outstanding request
          entity.isWaiting ? (
            <div className={classes.progressContainer}>
              <CircularProgress className={classes.progress} color='secondary'/>
            </div>
          ) : null
        }
        <Grid container>
          <Grid item sm={12} md={9} lg={8} className='Graph-scroll-container'>
            <Graph dataSet={entity.results}/>
          </Grid>
          <Grid item sm={12} md={3} lg={4} className='Graph-scroll-container'>
            <EntityDetails selectedObj={entity.results}/>
          </Grid>
        </Grid>
      </div>
    );
  }
}

GraphContainer.propTypes = {
  classes: PropTypes.object.isRequired,
};

const mapStoreToProps = store => ({entity: store.entity});

const mapDispatchToProps = dispatch => ({
  entityActions: bindActionCreators(entityActions, dispatch)
});

export default connect(
  mapStoreToProps,
  mapDispatchToProps
)(
  withRouter(withStyles(styles)(GraphContainer))
);