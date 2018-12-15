import React, { Component } from 'react';
import PropTypes from 'prop-types';
import _ from 'lodash';
import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import CircularProgress from '@material-ui/core/CircularProgress';

import "./Graph.css";
import EntityDetails from '../details/EntityDetails';
import Graph from './Graph';

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
  }

  componentDidMount() {
    this.setState({waitingOnReq: true});
    //TODO:DM - unhardcode this, 5000ms -> 5 sec
    this.timer = setInterval(()=> this.getData(), 5000);
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
    clearInterval(this.timer);
  }

  getData (){
    const pathComponents = this.props.location.pathname.split('/');
    const uidParam = pathComponents[pathComponents.length - 1];
    let url = `${window.envConfig.KATLAS_API_URL}/v1/entity/uid/${uidParam}`;

    fetch(url)
      .then(this.processRawResp)
      .then(json => {
        //only update state if the objects fail lodash equality check
        if(!_.isEqual(this.state.data, json)) {
          this.setState({
            data: json,
            waitingOnReq: false
          });
        }
      });
  }

  processRawResp(resp) {
    if (!resp.ok) {
      //TODO:DM - better to just throw here and handle errors centrally across app?
      //throw Error(resp.statusText);
      console.error('Error processing API response: ' + resp.statusText);
      return {};
    }
    return resp.json();
}

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

export default withStyles(styles)(GraphContainer);