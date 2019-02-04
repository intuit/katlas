import _ from 'lodash';

//TODO:DM - determine what from this file can be isolated into pure util fns and what makes most sense to incorporate directly into Graph component
//This Library Includes utilities for processing dgraph json and converting to components required by vis.js

import {
  NodeIconMap,
  NodeStatusColorMap,
  EdgeLabels,
  EdgeColorMap,
  NODE_DEFAULT_STR,
  NODE_ICON_FONT,
  NODE_DEFAULT_COLOR,
  NODE_ICON_FONT_SIZE,
  NODE_LABEL_MAX_LENGTH,
  NODE_LABEL_SPLIT_CHAR
} from "../config/appConfig";

//This attribute is added by dgraph.
const uidPropNameDgraph = "uid";

//Below are application specific attributes in the dgraph json based on Applications dgraph schema
//They are used to display Application specific data for the Node
const namePropNameDgraphApp = "name";
const objtypePropNameDgraphApp = "objtype";
const nodeStatusProp = "status";
const nodePhaseProp = "phase";

let uidMap;
let edges, nodes;
let legendTypesObj = {};
let legendStatusesObj = {};

function parseDgraphData(data, level) {
  if (data === null || data === undefined ) {
    console.error("Empty data");
    return;
  }
  _.forOwn(data, (val, key) => {
    let uid = "";
    if (key === uidPropNameDgraph) {
      uid = val;
      //Store entire block in map, so we can use to create edges later
      //Edges cannot be created here as we may not have got the uid still
      //As json from dgraph has randomized order and uid may be after Array elements
      uidMap.set(uid, {
        data,
        level
      });
    }
    //if this key is a relationship type (as defined in EdgeLabels) recurse to get children nodes
    if (EdgeLabels.indexOf(key) > -1) {
      _.forEach(val, (item) => parseDgraphData(item, level + 1));
    }
  });
}

function getNodeIcon(nodeObjtype) {
  if (nodeObjtype === undefined || nodeObjtype === null) {
    return NodeIconMap.get("default");
  }
  if (NodeIconMap.has(nodeObjtype)) {
    return NodeIconMap.get(nodeObjtype);
  } else {
    return NodeIconMap.get("default");
  }
}

function getEdgeLabelShortHand(prop) {
  let edgeLabel = "";
  if (prop !== "") {
    edgeLabel = prop;
  }
  return edgeLabel;
}

function getVisFormatEdge(fromUid, toUid, relation) {
  return {
    id: fromUid + toUid,
    from: fromUid,
    to: toUid,
    label: getEdgeLabelShortHand(relation),
    color: {
      color: EdgeColorMap.get(NODE_DEFAULT_STR),
      inherit: false
    },
    font: {
      size: 8
    },
    arrows: "to",
    smooth: {
      enabled: true,
      type: 'cubicBezier',
      forceDirection: 'vertical',
      roundness: 0.5
    }
  };
}

function getVisFormatNode(uid, nodeName, nodeObjtype, nodeStatus, level) {
  let idParam = uid;
  let titleParam = "";

  if (nodeObjtype !== undefined && nodeObjtype !== null) {
    titleParam = nodeObjtype;
  } else {
    titleParam = nodeName;
  }

  //some names are too large to render elegantly, split them by middle dash, if present, across 2 lines of text
  nodeName = nameSplitter(nodeName);

  const color = NodeStatusColorMap.get(nodeStatus || NODE_DEFAULT_STR);

  const n = {
    level,
    id: idParam,
    uid,
    label: nodeName,
    icon:{
      face: NODE_ICON_FONT,
      code: getNodeIcon(nodeObjtype),
      size: NODE_ICON_FONT_SIZE,
      color,
    },
    name: nodeName,
    title: titleParam,
    status: nodeStatus,
  };
  legendTypesObj[nodeObjtype] = {
    code: getNodeIcon(nodeObjtype),
    color: NODE_DEFAULT_COLOR,
  };
  legendStatusesObj[nodeStatus] = {
    code: getNodeIcon(NODE_DEFAULT_STR),
    color,
  };
  return n;
}

function validateJSONData(uid, nodeName, nodeObjtype) {
  if (nodeName === "") {
    console.error(`JSON Error - Attribute ${namePropNameDgraphApp} missing for uid = ${uid}`);
  }
  if (nodeObjtype === "") {
    console.error(`JSON Error - Attribute ${objtypePropNameDgraphApp} missing for uid = ${uid}`);
  }
}

export function getVisData(data) {
  let existingUids = {};
  clearVisData();
  parseDgraphData(data, 0);

  for (const [uid, val] of uidMap.entries()) {
    let nodeName = "", nodeObjtype = "", nodeStatus = "";
    const { data: block, level } = val;

    for (const prop in block) {
      if (!block.hasOwnProperty(prop)) {
        continue;
      }
      const val = block[prop];
      //determine whether we are looking at a property for this node OR a set of child nodes
      if (Array.isArray(val) && val.length > 0 &&
        typeof val[0] === "object") {
        // These are child nodes
        for (let i = 0; i < val.length; i++) {
          const fromUid = uid; //key for this map entry
          const toUid = val[i].uid;

          const e = getVisFormatEdge(fromUid, toUid, prop);
          edges.push(e);
        }
      } else {
        //get properties which we need to feed to Visjs Node
        if (prop === namePropNameDgraphApp) {
          nodeName = val;
        }
        if (prop === objtypePropNameDgraphApp) {
          nodeObjtype = val;
        }
        if (prop === nodeStatusProp ||
          prop === nodePhaseProp) {
          nodeStatus = val;
        }
      }
    }
    validateJSONData(uid, nodeName, nodeObjtype);
    let n = getVisFormatNode(uid, nodeName, nodeObjtype, nodeStatus, level);
    if(!existingUids[n.uid]){
      nodes.push(n);
      existingUids[n.uid] = true;
    }
  }
  return {nodes, edges};
}

export function getLegends(){
  return {
    types: legendTypesObj,
    statuses: legendStatusesObj,
  };
}

export function clearVisData(){
  uidMap = new Map();
  edges = [];
  nodes = [];
  legendTypesObj = {};
  legendStatusesObj = {};
}

//colorChannelA and colorChannelB are ints ranging from 0 to 255
function colorChannelMixer(colorChannelA, colorChannelB, amountToMix){
  let channelA = colorChannelA*amountToMix;
  let channelB = colorChannelB*(1-amountToMix);
  return parseInt(channelA+channelB);
}
//rgbA and rgbB are arrays, amountToMix ranges from 0.0 to 1.0
//example (red): rgbA = [255,0,0]
export function colorMixer(rgbA, rgbB, amountToMix){
  let r = colorChannelMixer(rgbA[0],rgbB[0],amountToMix);
  let g = colorChannelMixer(rgbA[1],rgbB[1],amountToMix);
  let b = colorChannelMixer(rgbA[2],rgbB[2],amountToMix);
  return "rgb("+r+","+g+","+b+")";
}

//split long label names; if too long, split by its middle char across 2 lines
function nameSplitter(name){
  let splitName = name;
  if (typeof name !== 'string') return name;
  if (name.length > NODE_LABEL_MAX_LENGTH){
    splitName = name.slice(0, name.length/2) + '\n' +
      name.slice(name.length/2);
  }
  return splitName;
}
