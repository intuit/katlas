import * as actions from './queryActions'
import * as types from './actionTypes'

describe('query actions', () => {
  it('should create an action to change query text', () => {
    const query = 'foo';
    const expectedAction = {
      type: types.CHANGE_QUERY,
      query
    };
    expect(actions.changeQuery(query)).toEqual(expectedAction);
  });

  it('should create an action to submit query', () => {
    const expectedAction = {
      type: types.SUBMIT_QUERY
    };
    expect(actions.submitQuery()).toEqual(expectedAction);
  });

  it('should create an action to receieve query', () => {
    const results = {'foo':'bar'};
    const expectedAction = {
      type: types.RECEIVE_QUERY,
      results
    };
    expect(actions.receiveQuery(results)).toEqual(expectedAction);
  });
});