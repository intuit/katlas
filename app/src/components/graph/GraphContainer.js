import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import { withStyles } from '@material-ui/core/styles';
import CircularProgress from '@material-ui/core/CircularProgress';
import SplitterLayout from 'react-splitter-layout';

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
  graphContainer: {
    width: '100%',
    height: '100vh',
    overflowX: 'auto',
    minHeight: '100vh',
    textAlign: 'left',
  },
});

class GraphContainer extends Component {
  componentDidMount() {
    this.setRootNode();
    //run first data acquisition event immediately, maintain handle for
    //cancellation purpose
    this.intervalHandle = setTimeout(() => this.getDataInterval(), 0);
  }

  componentDidUpdate(prevProps) {
    //recognize change in URL and re-issue API request in that case
    if (this.props.location !== prevProps.location){
      this.setRootNode();
      this.getData();
    }
  }

  componentWillUnmount() {
    clearInterval(this.intervalHandle);
  }

  setRootNode() {
    const pathComponents = this.props.location.pathname.split('/');
    //TODO:DM - simply grabbing last param after '/' feels fragile, how to more safely verify as UID?
    //could be empty string... a better default to use, if so?
    const uid = pathComponents[pathComponents.length - 1];

    this.props.entityActions.setRootEntity(uid);
    this.props.entityActions.addEntityWatch(uid);
  }

  getDataInterval() {
    this.getData();
    //reschedule next automatic data request while computing time value based
    //on number of entities and a min time between fetches
    const NUM_ENTITIES = Object.keys(this.props.entity.entitiesByUid).length;
    this.intervalHandle = setTimeout(() => this.getDataInterval(),
      NUM_ENTITIES * FETCH_PERIOD_PER_ENTITY_MS);
  }

  getData() {
    //fetch all entities currently represented as keys in the store
    this.props.entityActions.fetchEntities(Object.keys(
      this.props.entity.entitiesByUid));
  }

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
        <SplitterLayout percentage={true} secondaryInitialSize={30}>
          <Graph dataSet={entity.results}/>
          <EntityDetails selectedObj={entity.results}/>
        </SplitterLayout>
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