export const QSLRegEx = /([a-zA-Z0-9]+)\[?(?:(@[()",@$=><!a-zA-Z0-9\-.|&:_]*|\**|\$\$[a-zA-Z0-9,=]+))\]?\{([*|[,@"=a-zA-Z0-9-]*)/;

export function validateIPaddress(ipaddress) {
  const ipformat = /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
  if (ipaddress.match(ipformat)) {
    return true;
  }
  return false;
}

// get all the objects from a query string
export function getQSLObjTypes(query) {
  const objTypes = [];
  const querySegments = query.split('.');
  querySegments.forEach(segment => {
    const matches = QSLRegEx.exec(segment);
    if (matches) {
      objTypes.push(matches[1]);
    }
  });
  return objTypes;
}
