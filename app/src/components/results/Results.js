import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import CircularProgress from '@material-ui/core/CircularProgress';
import TextField from '@material-ui/core/TextField';
import SearchIcon from '@material-ui/icons/Search';
import IconButton from '@material-ui/core/IconButton';
import InputAdornment from '@material-ui/core/InputAdornment';
import SplitterLayout from 'react-splitter-layout';

import { ENTER_KEYCODE } from '../../config/appConfig';
import ResultList from './ResultList';
import EntityDetails from '../entityDetails/EntityDetails';
import * as apiCfg from '../../config/apiConfig';
import * as queryActions from '../../actions/queryActions';

const styles = theme => ({
  progress: {
    margin: theme.spacing.unit * 2
  },
  progressContainer: {
    textAlign: 'center'
  },
  resultTitle: {
    marginTop: theme.spacing.unit * 2,
    marginLeft: theme.spacing.unit
  },
  searchBox: {
    marginTop: 30,
    marginLeft: 5,
    marginBottom: 20
  }
});

function getQueryParam(locationSearchStr, queryParamStr) {
  const params = new URLSearchParams(locationSearchStr);
  return params.get(queryParamStr) || '';
}

class Results extends Component {
  constructor(props) {
    super(props);
    const queryStr = getQueryParam(
      this.props.location.search,
      apiCfg.SERVICES.queryParamName
    );
    this.state = {
      selectedIdx: 0,
      queryStr: queryStr
    };
  }

  handleChange = event => {
    this.setState({ queryStr: event.target.value });
  };

  handleEnterPressCheck = event => {
    const { queryStr } = this.state;
    if (event.keyCode === ENTER_KEYCODE && queryStr !== '') {
      this.handleSubmit();
    }
  };

  handleSubmit = () => {
    const { queryStr } = this.state;
    const {
      queryActions: { submitQuery }
    } = this.props;

    submitQuery(queryStr);
  };

  componentDidMount() {
    const {
      query,
      queryActions: { fetchQuery }
    } = this.props;
    const { queryStr } = this.state;

    fetchQuery(queryStr, query.page, query.rowsPerPage);
  }

  componentDidUpdate(prevProps) {
    //recognize change in query param here and re-issue API request as necessary
    const currentQuery = getQueryParam(
      this.props.location.search,
      apiCfg.SERVICES.queryParamName
    );
    const prevQuery = getQueryParam(
      prevProps.location.search,
      apiCfg.SERVICES.queryParamName
    );
    const {
      query,
      queryActions: { fetchQuery }
    } = this.props;
    if (prevQuery !== currentQuery) {
      //should only run if query param changes
      fetchQuery(currentQuery, query.page, query.rowsPerPage);
    }
  }

  handleRowClick = (event, idx) => {
    //TODO:DM - take this opportunity to distinguish the row visually?
    this.setState({ selectedIdx: idx });
  };

  render() {
    const {
      classes,
      query,
      queryActions: { updatePagination }
    } = this.props;
    const { queryStr, selectedIdx } = this.state;

    return (
      <div className='Results'>
        <TextField
          id='outlined-full-width'
          label='Search'
          className={classes.searchBox}
          placeholder='Query String'
          fullWidth
          margin='normal'
          variant='outlined'
          InputLabelProps={{
            shrink: true
          }}
          value={queryStr}
          onChange={this.handleChange}
          onKeyUp={this.handleEnterPressCheck}
          InputProps={{
            endAdornment: (
              <InputAdornment position='end'>
                <IconButton
                  aria-label='Search'
                  onClick={() => this.handleSubmit()}
                >
                  <SearchIcon />
                </IconButton>
              </InputAdornment>
            )
          }}
        />
        {//selectively show progress spinner or table, once HTTP req resolves
        query.isWaiting ? (
          <div className={classes.progressContainer}>
            <CircularProgress className={classes.progress} color='secondary' />
          </div>
        ) : (
          <SplitterLayout percentage={true} secondaryInitialSize={30}>
            <ResultList
              query={query}
              selectedIdx={selectedIdx}
              onRowClick={this.handleRowClick}
              updatePagination={updatePagination}
            />
            <EntityDetails selectedObj={query.results[selectedIdx]} />
          </SplitterLayout>
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
)(withStyles(styles)(Results));
