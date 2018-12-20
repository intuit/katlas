import React, { Component } from 'react';
import { withRouter } from 'react-router-dom';
import _ from 'lodash';
import vis from 'vis';
import { Grid } from '@material-ui/core';

import {options} from "../../config/visjsConfig";
import {getVisData, clearVisData, colorMixer, getLegends} from "../../utils/graph";
import {NodeStatusPulseColors} from "../../config/appConfig";
import './Graph.css';

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


    this.setDetailsTab = this.setDetailsTab.bind(this);
    this.validateInputs = this.validateInputs.bind(this);
    this.renderGraph = this.renderGraph.bind(this);
    this.renderGraphExpanded = this.renderGraphExpanded.bind(this);
    this.renderVisGraph = this.renderVisGraph.bind(this);
    this.clearNetwork = this.clearNetwork.bind(this);
  }

  componentDidMount() {
    this.clearNetwork();
  }

  componentWillReceiveProps(nextProps){
    if(!_.isEqual(nextProps, this.props)){
      this.clearNetwork();
      this.setState({data: nextProps.dataSet});
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
    this.renderGraphExpanded(this.state.data);
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
    console.debug(`Got data ${JSON.stringify(jsonData)}`);

    const {nodes, edges} = getVisData(jsonData);
    this._nodes = nodes;
    this._edges = edges;

    this.renderVisGraph();
  }

  renderGraphExpanded(jsonData) {
    if (_.isEmpty(jsonData)) return;

    this.nodesDataset.clear();
    this.edgesDataset.clear();
    this._nodes = [];
    this._edges = [];

    const {nodes, edges} = getVisData(jsonData);
    const nodesNew = nodes;
    const edgesNew = edges;
    this._legends = getLegends();

    //Avoid adding duplicate selected node, as vis.js does not allow duplicates.
    for (let i = 0; i < nodesNew.length; i++) {
      this._nodes.push(nodesNew[i]);
    }
    this._edges.push(edgesNew);

    this.renderVisGraph();
  }

  renderVisGraph() {
      this.nodesDataset = new vis.DataSet();
      for (let i = 0; i < this._nodes.length; i++) {
        this.nodesDataset.add(this._nodes[i]);
      }
      this.edgesDataset = new vis.DataSet();
      for (let i = 0; i < this._edges.length; i++) {
        this.edgesDataset.add(this._edges[i]);
      }
      this._data = {nodes: this.nodesDataset, edges: this.edgesDataset};

      const container = document.getElementById("graph"); //id of div container for graph.
      this._network = new vis.Network(container, this._data, options);

      this.configNetwork(this._network);
      //must bind call to handleColorPulse since it'll be called by browser otherwise with Window as "this" context
      this.pulseIntervalHandle = setInterval(this.handleColorPulse.bind(this), 100);//TODO:DM-constantize time!
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

  setDetailsTab(attributesMap) {
      const attributes = [];
      for (const [k, v] of attributesMap.entries()) {
          if (v !== undefined && v !== "") {
              const entry = `  ${k} : ${v}`;
              attributes.push(entry);
          }
      }
      this.setState({detailsTab: attributes});
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
    network.on("doubleClick", params => {
      //ensure the double click was on a graph node
      if (params.nodes.length > 0) {
        const targetNodeUid = params.nodes[0];
        const pathComponents = this.props.location.pathname.split('/');
        const currentNodeUid = pathComponents[pathComponents.length - 1];
        //only update props if target node is not current node
        if (targetNodeUid !== currentNodeUid) {
          this.clearNetwork();
          this.props.history.push('/graph/'+targetNodeUid);
        }
      }
    });
  }
}

export default withRouter(Graph);