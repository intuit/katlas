export function validateIpAddress(input) {
  const ipFormat = /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
  //type coerce match array or null value to bool
  return !!input.match(ipFormat);
}

export function validateQslQuery(input) {
  //TODO:DM - confirm with @kianjones that this is a reasonable first attempt at a QSL query validation matcher
  const qslFormat = /.*\[@.*\].*/g;
  //type coerce match array or null value to bool
  return !!input.match(qslFormat);
}

export function validateHexId(input) {
  const hexIdFormat = /(0x|0X)?[a-fA-F0-9]+$/g;
  //type coerce match array or null value to bool
  return !!input.match(hexIdFormat);
}