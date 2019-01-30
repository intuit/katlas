//TODO:DM - looks like the validate module isn't currently being used in the app, therefore these tests shouldn't really count UNLESS we plan to use it again
import {
  validateIPaddress,
  getQSLObjTypes,
  getQSLObjTypesAndProjection
} from './validate';

describe('validation util', () => {
  it('should correctly recognize an IP addr', () => {
    expect(validateIPaddress('127.0.0.1')).toBe(true);
  });

  it('should correctly recognize a non-IP addr', () => {
    expect(validateIPaddress('localhost')).toBe(false);
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
      'deployment{@name,@availablereplicas}.replicaset{*}.pod{*}.node[@name="ip-10-83-122-52.us-west-2.compute.internal"]{*}';
    const queryProjection = getQSLObjTypesAndProjection(query);
    console.log(queryProjection);
    // expect(queryProjection.deployment).toBe('*');
    expect(queryProjection.deployment).toHaveLength(2);
    expect(queryProjection.deployment).toContain('name');
    // expect(queryProjection.pod).toContain('resourceid');
  });
});
