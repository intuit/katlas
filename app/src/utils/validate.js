export function validateIpAddress(input) {
  const ipFormat = /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
  //type coerce match array or null value to bool
  return !!input.match(ipFormat);
}

export const QSLRegEx = /([a-zA-Z0-9]+)\[?(?:(@[",@$=><!a-zA-Z0-9\-.|&:_]*|\**|\$\$[a-zA-Z0-9,=]+))\]?\{([*|[,@"=a-zA-Z0-9-]*)/;
export function validateQslQuery(input) {
  //type coerce match array or null value to bool
  return !!input.match(QSLRegEx);
}

export function validateHexId(input) {
  const hexIdFormat = /(0x|0X)?[a-fA-F0-9]+$/g;
  //type coerce match array or null value to bool
  return !!input.match(hexIdFormat);
}