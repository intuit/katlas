import React, { Component } from 'react';
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import { withRouter } from 'react-router-dom';
import _ from 'lodash';
import vis from 'vis';
import { Grid } from '@material-ui/core';

import * as entityActions from "../../actions/entityActions";
import { options } from "../../config/visjsConfig";
import { getVisData, clearVisData, colorMixer, getLegends } from "../../utils/graph";
import { NodeStatusPulseColors } from "../../config/appConfig";
import './Graph.css';

const PULSE_TIME_STEP_MS = 100;

class Graph extends Component {
  constructor(props) {
    super(props);

    this.state = {
      data: {},
      errors: [],
      detailsTab: [],
    };

    this._nodes = [];
    this._edges = [];
    this._data = {};
    this._network = null;
    this.nodesDataset = new vis.DataSet();
    this.edgesDataset = new vis.DataSet();
    this.pulseIntervalHandle = null;
    this._legends = {
      types: {},
      statuses: {}
    };

    //TODO:DM - can I mitigate need for these hard binds with => fns or something else?
    this.validateInputs = this.validateInputs.bind(this);
    this.renderGraph = this.renderGraph.bind(this);
    this.renderVisGraph = this.renderVisGraph.bind(this);
    this.clearNetwork = this.clearNetwork.bind(this);
  }

  componentDidMount() {
    this.clearNetwork();
  }

  componentWillReceiveProps(nextProps){
    if(!_.isEqual(nextProps.dataSet, this.props.dataSet)){
      this.renderGraph(nextProps.dataSet)
    }
  }

  componentWillUnmount(){
    this.clearNetwork();
  }

  validateInputs() {
    const errorList = [];
    if (this.state.data.length <= 0) {
      errorList.push("emptyDataField");
    }
    this.setState({errors: errorList});
    return errorList;
  }

  render() {
      return (
        <div className="Graph">
          {/*Graph Visualization*/}
          <div className="Graph-container" align="center" id="graph"/>
          {/*Graph Legend*/}
          <div className='Graph-legend-container'>
            Legend
            <Grid
              container
              spacing={16}
              className={''}
              alignItems="center"
              direction="row"
              justify="center"
            >
              {/*First Render Types*/}
              {Object.keys(this._legends.types).map(typeKey => {
                return (
                  <Grid key={typeKey} item className="Graph-legend-cell">
                    <span className="Graph-legend-icon"
                      style={{color: '#000000'}}>
                      {this._legends.types[typeKey].code}
                    </span>
                    <p>{typeKey}</p>
                  </Grid>
                )
              })}
              {/*Then a divider and Status Indications*/}
              <Grid key={0} item>
                <span className="Graph-legend-vr">|</span>
              </Grid>
              {Object.keys(this._legends.statuses).map(statusKey => {
                return (
                  <Grid key={statusKey} item className="Graph-legend-cell">
                    <span className="Graph-legend-icon"
                      style={{color: this._legends.statuses[statusKey].color}}>
                      {this._legends.statuses[statusKey].code}
                    </span>
                    <p>{statusKey || 'No Status'}</p>
                  </Grid>
                )
              })}
            </Grid>
          </div>
        </div>
      );
  }

  renderGraph(jsonData) {
    if (_.isEmpty(jsonData)) return;

    const {nodes, edges} = getVisData(jsonData);
    this._edges = edges;
    this._legends = getLegends();

    for (let i = 0; i < nodes.length; i++) {
      this._nodes.push(nodes[i]);
    }
    this.renderVisGraph();
  }

  renderVisGraph() {
    //determine if we have any newnodes to add to the graph dataset
    for (let i = 0; i < this._nodes.length; i++) {
      if(!this.nodesDataset.get(this._nodes[i].uid)){
        this.nodesDataset.add(this._nodes[i]);
      }
    }
    //determine if we have any new edges to add to the graph dataset
    for (let i = 0; i < this._edges.length; i++) {
      //edge IDs are a concatenation of from and to IDs
      if(!this.edgesDataset.get(this._edges[i].from + this._edges[i].to)){
        this.edgesDataset.add(this._edges[i]);
      }
    }
    this._data = {nodes: this.nodesDataset, edges: this.edgesDataset};

    //id of div container for graph
    const container = document.getElementById("graph");
    this._network = new vis.Network(container, this._data, options);

    this.configNetwork(this._network);
    //must bind call to handleColorPulse since it'll be called by browser
    //otherwise with window as "this" context
    this.pulseIntervalHandle = setInterval(this.handleColorPulse.bind(this),
      PULSE_TIME_STEP_MS);
  }

  handleColorPulse() {
    let mixRatio = ((new Date()).getTime() % 1000) / 1000;
    this.nodesDataset.forEach((node) => {
      if (NodeStatusPulseColors.has(node.status)) {
        let extremes = NodeStatusPulseColors.get(node.status);
        this.nodesDataset.update([{
          id: node.id,
          icon: {
            face: node.icon.face,
            code: node.icon.code,
            size: node.icon.size,
            color: colorMixer(extremes[0], extremes[1], mixRatio)
          }
        }]);
      }
    });
  }

  clearNetwork() {
    //cancel node pulse timers, before this._data is cleared
    clearInterval(this.pulseIntervalHandle);

    this._nodes = [];
    this._edges = [];
    this._data = {};
    this._network = {};
    this.nodesDataset.clear();
    this.edgesDataset.clear();

    clearVisData();
  }

  //configure custom behaviors for network object
  configNetwork(network) {
    network.on("doubleClick", element => {
      //ensure the double click was on a graph node
      if (element.nodes.length > 0) {
        const targetNodeUid = element.nodes[0];
        const pathComponents = this.props.location.pathname.split('/');
        const currentNodeUid = pathComponents[pathComponents.length - 1];
        //only add node data if target node is not current node
        //TODO:DM - rather than just current node, could skip any already watched UIDs
        if (targetNodeUid !== currentNodeUid) {
          this.props.entityActions.addEntityWatch(targetNodeUid);
          //and immediately attempt to fetch for the indicated UID
          this.props.entityActions.fetchEntity(targetNodeUid);
        }
      }
    });
  }
}

const mapStoreToProps = store => ({entity: store.entity});

const mapDispatchToProps = dispatch => ({
  entityActions: bindActionCreators(entityActions, dispatch)
});

export default connect(
  mapStoreToProps,
  mapDispatchToProps
)(withRouter(Graph));