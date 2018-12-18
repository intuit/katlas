import React, { Component } from 'react';
import {connect} from 'react-redux';
import {bindActionCreators} from 'redux';
import PropTypes from 'prop-types';
import {Link, withRouter} from 'react-router-dom';
import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Paper from '@material-ui/core/Paper';
import CircularProgress from '@material-ui/core/CircularProgress';

import EntityDetails from '../entityDetails/EntityDetails';
import * as apiCfg from '../../config/apiConfig';
import * as queryActions from '../../actions/queryActions';
import './Results.css';

const styles = theme => ({
  container: {
    display: 'flex',
    flexWrap: 'wrap',
  },
  textField: {
    marginLeft: theme.spacing.unit,
    marginRight: theme.spacing.unit,
  },
  root: {
    width: '100%',
    overflowX: 'auto',
  },
  table: {
    minWidth: 700,
  },
  progress: {
    margin: theme.spacing.unit * 2,
  },
  progressContainer: {
    textAlign: 'center',
  },
});

function getQueryParam(locationSearchStr, queryParamStr){
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
    this.props.queryActions.submitQuery();
    this.props.queryActions.fetchQuery(query);
  }

  componentDidUpdate(prevProps) {
    //recognize change in query param here and re-issue API request as necessary
    const currentQuery = getQueryParam(this.props.location.search, apiCfg.SERVICES.queryParamName);
    const prevQuery = getQueryParam(prevProps.location.search, apiCfg.SERVICES.queryParamName);
    if (prevQuery !== currentQuery){
      //should only run if query param changes
      this.props.queryActions.submitQuery();
      this.props.queryActions.fetchQuery(currentQuery);
    }
  }

  handleRowClick = (event, idx) => {
    //TODO:DM - take this opportunity to distinguish the row visually?
    this.setState({ selectedIdx: idx });
  };

  //TODO:DM - should we extract table construction into ResultList comp?
  render() {
    const { classes, query } = this.props;
    return (
      <div className='Results'>
        {//selectively show progress spinner or table, once HTTP req resolves
          query.isWaiting ? (
          <div className={classes.progressContainer}>
            <CircularProgress className={classes.progress} color='secondary'/>
          </div>
        ) : (
          <Grid container>
            <Grid item sm={12} md={9} lg={8} className='Results-scroll-container'>
              <Paper className={classes.root}>
                <Table padding='dense' className={classes.table}>
                  <TableHead>
                    <TableRow>
                      <TableCell>Type</TableCell>
                      <TableCell>Name</TableCell>
                      <TableCell>Namespace</TableCell>
                      <TableCell>Creation Datetime</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {(query.results.length > 0) ? (
                      query.results.map((item, idx) => {
                        return (
                          <TableRow hover key={item.uid}
                            onClick={event => this.handleRowClick(event, idx)}>
                            <TableCell component='th' scope='row'>
                              {item.objtype}
                            </TableCell>
                            <TableCell>
                              <Link
                              to={{
                                pathname: '/graph/'+ item.uid,
                                state: {selectedObj:query.results[this.state.selectedIdx]}
                              }}>
                                {item.name}
                                </Link>
                            </TableCell>
                            <TableCell>{item.namespace ? item.namespace[0].name : ''}</TableCell>
                            <TableCell>{item.starttime}</TableCell>
                          </TableRow>
                        );
                      })
                    ) : ( //TODO:DM determine if there is a more elegant 'toggle' pattern suggested in React/jsx community
                    <TableRow>
                      <TableCell/>
                      <TableCell>No data</TableCell>
                      <TableCell/>
                    </TableRow>
                  )}
                  </TableBody>
                </Table>
              </Paper>
            </Grid>
            <Grid item sm={12} md={3} lg={4} className='Results-scroll-container'>
              <EntityDetails selectedObj={query.results[this.state.selectedIdx]}/>
            </Grid>
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

const mapStateToProps = state => ({query: state.query});

const mapDispatchToProps = dispatch => ({
    queryActions: bindActionCreators(queryActions, dispatch)
});

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(
  withStyles(styles)(
    withRouter(Results)
  )
);