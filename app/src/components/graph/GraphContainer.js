import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { withRouter } from 'react-router-dom';
import { withStyles } from '@material-ui/core/styles';
import CircularProgress from '@material-ui/core/CircularProgress';
import SplitterLayout from 'react-splitter-layout';
import '../../shared/reactSplitterLayoutWithOverrides.css';

import "./Graph.css";
import Graph from './Graph';
import EntityDetails from '../entityDetails/EntityDetails';
import * as entityActions from "../../actions/entityActions";
import { validateQslQuery, validateHexId } from "../../utils/validate";

const FETCH_PERIOD_PER_ENTITY_MS = 5000;

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
    this.determineInitialWatch();
    //run first data acquisition event immediately, maintain handle for
    //cancellation purpose
    this.intervalHandle = setTimeout(() => this.getDataInterval(), 0);
  }

  componentDidUpdate(prevProps) {
    //recognize change in URL and re-issue API request in that case
    if (this.props.location !== prevProps.location){
      this.determineInitialWatch();
      this.getData();
    }
  }

  componentWillUnmount() {
    const { clearWatches } = this.props.entityActions;
    clearWatches();
    clearInterval(this.intervalHandle);
  }

  determineInitialWatch() {
    const pathParam = this.props.match.params.uidOrQsl;
    if (pathParam && validateQslQuery(pathParam)) {
      this.props.entityActions.fetchQslQuery(pathParam);
      this.props.entityActions.addWatchQslQuery(pathParam);
    } else if (pathParam && validateHexId(pathParam)) {
      this.props.entityActions.setRootUid(pathParam);
      this.props.entityActions.addWatchUid(pathParam);
    }
  }

  getDataInterval() {
    let numEntities;
    this.getData();
    //reschedule next automatic data request while computing time value based
    //on number of entities and a min time between fetches
    if (this.props.entity.qslQuery) {
      numEntities = 10;
    } else {
      numEntities = Object.keys(this.props.entity.entitiesByUid).length;
    }
    this.intervalHandle = setTimeout(() => this.getDataInterval(),
      numEntities * FETCH_PERIOD_PER_ENTITY_MS);
  }

  getData() {
    //fetch QSL query, if one is registered
    if (this.props.entity.qslQuery) {
     this.props.entityActions.fetchQslQuery(this.props.entity.qslQuery);
    }
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
        <SplitterLayout percentage={true} secondaryInitialSize={0}>
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