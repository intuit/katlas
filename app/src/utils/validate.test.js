import { validateIpAddress, validateQslQuery, validateHexId } from './validate';

it('should correctly recognize an IP addr', () => {
  expect(validateIpAddress('127.0.0.1')).toBe(true);
});

it('should correctly recognize an invalid IP addr', () => {
  expect(validateIpAddress('localhost')).toBe(false);
});

it('should correctly recognize a QSL query', () => {
  expect(validateQslQuery(
    'Cluster[@objtype="Cluster"]{*}.Node[@objtype="Node"]{*}')).toBe(true);
});

it('should correctly recognize an invalid QSL query', () => {
  expect(validateQslQuery(
    '@Cluster[objtype="Cluster"]{*}.@Node[objtype="Node"]{*}')).toBe(false);
});

it('should correctly recognize a simple hex ID', () => {
  expect(validateHexId('0x12')).toBe(true);
});

it('should correctly recognize a long hex ID', () => {
  expect(validateHexId('0x12345deadbeefCABFABDAB09876')).toBe(true);
});

it('should correctly recognize an invalid hex ID', () => {
  expect(validateHexId('0xfoo')).toBe(false);
});