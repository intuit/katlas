import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import CircularProgress from '@material-ui/core/CircularProgress';
import Typography from '@material-ui/core/Typography';
import SplitterLayout from 'react-splitter-layout';

import EntityDetails from '../entityDetails/EntityDetails';
import * as apiCfg from '../../config/apiConfig';
import * as queryActions from '../../actions/queryActions';
import ResultList from './ResultList';


const styles = theme => ({
  progress: {
    margin: theme.spacing.unit * 2,
  },
  progressContainer: {
    textAlign: 'center',
  },
  resultTitle: {
    marginTop: theme.spacing.unit * 2,
    marginLeft: theme.spacing.unit,
  },
});

function getQueryParam(locationSearchStr, queryParamStr) {
  const params = new URLSearchParams(locationSearchStr);
  return params.get(queryParamStr) || '';
}

class Results extends Component {
  constructor(props) {
    super(props);
    //TODO:DM - get rid of local component state entirely, incorporate selectedIdx into central store
    this.state = {
      selectedIdx: 0,
    };
  }

  componentDidMount() {
    const query = getQueryParam(this.props.location.search, apiCfg.SERVICES.queryParamName);
    //incase the user directly linked to this route, make sure to take query
    //from params to get store 'caught up'
    this.props.queryActions.changeQuery(query);
    this.props.queryActions.submitQuery(query);
    this.props.queryActions.fetchQuery(query);
  }

  componentDidUpdate(prevProps) {
    //recognize change in query param here and re-issue API request as necessary
    const currentQuery = getQueryParam(this.props.location.search, apiCfg.SERVICES.queryParamName);
    const prevQuery = getQueryParam(prevProps.location.search, apiCfg.SERVICES.queryParamName);
    if (prevQuery !== currentQuery) {
      //should only run if query param changes
      this.props.queryActions.submitQuery(currentQuery);
      this.props.queryActions.fetchQuery(currentQuery);
    }
  }

  handleRowClick = (event, idx) => {
    //TODO:DM - take this opportunity to distinguish the row visually?
    this.setState({ selectedIdx: idx });
  };

  render() {
    const { classes, query } = this.props;
    const { selectedIdx } = this.state;

    return (
      <div className='Results'>
        <Typography variant="h6" gutterBottom className={classes.resultTitle}>
          Search Result: {query.current}
        </Typography>
        {//selectively show progress spinner or table, once HTTP req resolves
          query.isWaiting ? (
            <div className={classes.progressContainer}>
              <CircularProgress className={classes.progress} color='secondary' />
            </div>
          ) : (
              <Grid container>
                <SplitterLayout percentage={true} secondaryInitialSize={30}>
                  <ResultList query={query} selectedIdx={selectedIdx} onRowClick={this.handleRowClick} />
                  <EntityDetails selectedObj={query.results[selectedIdx]} />
                </SplitterLayout>
              </Grid>
            )}
      </div>
    );
  }
}

Results.propTypes = {
  classes: PropTypes.object.isRequired,
  query: PropTypes.object
};

const mapStateToProps = state => ({ query: state.query });

const mapDispatchToProps = dispatch => ({
  queryActions: bindActionCreators(queryActions, dispatch)
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(
  withStyles(styles)(Results)
);