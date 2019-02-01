import {
  validateIpAddress,
  validateQslQuery,
  validateHexId,
  getQSLObjTypes,
  getQSLObjTypesAndProjection
} from './validate';

describe('validation util', () => {
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

  it('should correctly recognize a short hex ID', () => {
    expect(validateHexId('0x1')).toBe(true);
  });

  it('should correctly recognize a long hex ID', () => {
    expect(validateHexId('0x12345deadbeefCABFABDAB09876')).toBe(true);
  });

  it('should correctly recognize an invalid hex ID', () => {
    expect(validateHexId('0xfoo')).toBe(false);
  });

  it('should get all obj types from query', () => {
    const query = 'deployment{*}.replicaset[@count(pod)<3]{*}.pod{*}';
    const objTypes = getQSLObjTypes(query);

    expect(objTypes.length).toEqual(3);
    expect(objTypes).toContain('deployment');
    expect(objTypes).toContain('replicaset');
    expect(objTypes).toContain('pod');
  });

  it('should get all the obj types and projection', () => {
    const query =
      'deployment{@name,@availablereplicas}.replicaset{*}.pod{*}.node[@name="ip-10-83-122-52.us-west-2.compute.internal"]{@name}';
    const queryProjection = getQSLObjTypesAndProjection(query);
    expect(queryProjection.replicaset).toBe('*');
    expect(queryProjection.deployment).toHaveLength(2);
    expect(queryProjection.deployment).toContain('name');
    expect(queryProjection.deployment).toContain('availablereplicas');
    expect(queryProjection.node).toContain('name');
  });
});