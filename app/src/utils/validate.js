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

//TODO:DM - figure out how to extract shared code between following two functions into another helper fn which both can use
export const addResourceIdFilterQSL = (query = '', resourceId = '') => {
  if (!validateQslQuery(query)) return query;
  const resourceIdFilter = `@resourceid="${resourceId}"`;
  const querySegments = query.split('.');
  const rootEntityQuery = querySegments[0];
  const matches = QSLRegEx.exec(rootEntityQuery);
  if (matches) {
    const objType = matches[1];
    const filter = matches[2];
    const fields = matches[3];
    querySegments[0] = `${objType}[${filter}${resourceIdFilter}]{${fields}}`;
  }
  return querySegments.join('.');
};

export const addPaginationFilterQSL = (
  query = '',
  page = 0,
  rowsPerPage = 50
) => {
  if (!validateQslQuery(query)) return query;
  //TODO:DM - splitting based on such a simple aspect of the QSL query pattern feels fragile; could regex solidify this?
  const splitChars = '}.';
  const querySegments = query.split(splitChars);
  //decode segment incase it was URI encoded thru a route navigation
  const rootEntityQuery = decodeURIComponent(querySegments[0]);
  //inject the pagination if not already present
  if (!rootEntityQuery.includes('$$')) {
    const matches = QSLRegEx.exec(rootEntityQuery);
    if (matches) {
      const objType = matches[1];
      const filter = matches[2];
      const fields = matches[3];
      const pagination = `$$limit=${rowsPerPage},offset=${page * rowsPerPage}`;
      querySegments[0] = `${objType}[${filter}${pagination}]{${fields}`;
    }
  }
  return querySegments.length === 1
    ? querySegments[0] + '}'
    : querySegments.join(splitChars);
};
