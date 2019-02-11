export function validateIpAddress(input) {
  const ipFormat = /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
  //type coerce match array or null value to bool
  return !!input.match(ipFormat);
}

export const QSLRegEx = /([a-zA-Z0-9]+)\[?(?:(@[()",@$~=><!a-zA-Z0-9\-.|&:_^*]*|\**|\$\$[a-zA-Z0-9,=]+))\]?\{([*|[,@"=a-zA-Z0-9-]*)/;

export function validateQslQuery(input) {
  if (input === null) {
    return false;
  }
  //type coerce match array or null value to bool
  return !!input.match(QSLRegEx);
}

export function validateHexId(input) {
  if (input === null) {
    return false;
  }
  const hexIdFormat = /(0x|0X)?[a-fA-F0-9]+$/g;
  //type coerce match array or null value to bool
  return !!input.match(hexIdFormat);
}

//TODO:DM+FZ - determine if the following util fns belong here, not strictly validations
// get all the objects from a query string
export function getQSLObjTypes(query) {
  const objTypes = [];
  const querySegments = query.split('}.');
  querySegments.forEach(segment => {
    const matches = QSLRegEx.exec(segment);
    if (matches) {
      objTypes.push(matches[1]);
    }
  });
  return objTypes;
}

// get all the object and the projection requested
export function getQSLObjTypesAndProjection(query) {
  const queryProjection = {};
  const querySegments = query.split('}.');
  querySegments.forEach(segment => {
    const matches = QSLRegEx.exec(segment);
    if (matches) {
      const objType = matches[1];
      const projection = matches[3];
      const fields =
        projection === '*'
          ? '*'
          : projection.split(',').map(p => p.substring(1));

      queryProjection[objType] = fields;
    }
  });
  return queryProjection;
}