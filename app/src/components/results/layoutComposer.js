import React from 'react';
import { Link } from 'react-router-dom';
import _ from 'lodash';

import { getQSLObjTypesAndProjection } from '../../utils/validate';
import { CustomTableCell } from './ResultList';

// generate the layout based on the metadata and the query
// sample layout structure:
// objtype: {
//   fieldname: {
//     displayname:"",
//     representsfunc: func
//   }
// }
export const getQueryLayout = (queryStr, metadata) => {
  let layout = {};
  const queryProjection = getQSLObjTypesAndProjection(queryStr);
  const excludeFields = ['objtype', 'k8sobj', 'resourceid'];
  for (let objType in queryProjection) {
    const objProjection = queryProjection[objType];
    let fieldLayout = {};
    const objMeta = metadata[objType];
    if(objMeta && objMeta.fields && _.isArray(objMeta.fields)){
      objMeta.fields.forEach(field => {
        const [shown, projection] = getFieldProjector(objType, field);
        if (
          shown &&
          !excludeFields.includes(field.fieldname) &&
          (objProjection === '*' || objProjection.includes(field.fieldname))
        ) {
          fieldLayout[field.fieldname] = projection;
        }
      });
    }
    layout[objType] = fieldLayout;
  }

  return layout;
};

// list all the cells of a row based on the layout
export const rowCellsFromLayout = (item, layout) => {
  let cells = [];
  //copy item since we'll potentially change it below
  let obj = {...item};
  for (let objType in layout) {
    obj = navEmbeddedObject(obj, objType);
    if (obj === undefined || obj === null) {
      // it could be the depth 1 has value but once drill down(depth 2) it could be null
      // since we only show first row from the first array and not reverse back and try another row if the depth 2 has value or not
      continue;
    }
    for (let field in layout[objType]) {
      const projectionFn = layout[objType][field].representsFunc;
      const value = obj[field];

      cells.push(
        <CustomTableCell
          key={`${obj.uid}-${field}`}
          style={{
            maxWidth: '350px',
            overflow: 'hidden',
            textOverflow: 'ellipsis'
          }}
        >
          {projectionFn(obj.uid, field, value)}
        </CustomTableCell>
      );
    }
  }
  return cells;
};

// get the object using provided object type among all the attributes (depth = 1)
const navEmbeddedObject = (obj, objType) => {
  if (obj === null) {
    return obj;
  }
  if (obj.objtype === objType) {
    return obj;
  }
  for (let k in obj) {
    let value = obj[k];
    const isarray = Array.isArray(value);
    if (isarray && value.length > 0) {
      value = value[0]; // only show first item for now
    }
    const type = Object.prototype.toString.call(value);
    const isobject = type === '[object Object]';
    if (isobject) {
      if (value.objtype === objType) {
        return value;
      }
    }
  }
  return null;
};

// The input is the field of the metadata
// return value is whether it will be shown and the representer
const getFieldProjector = (objType, field) => {
  switch (field.fieldtype) {
    case 'string':
      return [
        true,
        {
          displayName: field.fieldname,
          representsFunc: stringPresenter
        }
      ];
    case 'int':
      return [
        true,
        {
          displayName: field.fieldname,
          representsFunc: intPresenter
        }
      ];
    case 'true':
      return [
        true,
        {
          displayName: field.fieldname,
          representsFunc: stringPresenter
        }
      ];
    case 'json':
      return [
        true,
        {
          displayName: field.fieldname,
          representsFunc: jsonPresenter
        }
      ];
    default:
      return [false, null];
  }
};

const stringPresenter = (uid, name, val) => {
  // we add link if the field name is 'name'
  if (name === 'name') {
    return (
      <Link
        to={{
          pathname: '/graph/' + uid
        }}
      >
        {val}
      </Link>
    );
  }
  return val;
};

const intPresenter = (uid, name, val) => {
  return val;
};

const jsonPresenter = (uid, name, val) => {
  // we only shows 3 key pairs
  let count = 0;
  let output = [];
  if (val === undefined) {
    return val;
  }
  const json = JSON.parse(val);
  for (let k in json) {
    let v = JSON.stringify(json[k]);
    output.push(
      <div key={`${uid}-${k}`}>
        {k}: {v}
      </div>
    );
    if (count === 2) {
      break;
    }
    count++;
  }
  return output;
};
