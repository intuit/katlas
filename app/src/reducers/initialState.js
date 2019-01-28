export default {
  query: {
    current: '',
    lastSubmitted: '',
    isWaiting: false,
    results: [],
    page: 0,
    rowsPerPage: 25,
    count: 0
  },
  entity: {
    rootUid: '',
    entitiesByUid: {},
    results: {},
    latestTimestamp: 0,
    isWaiting: false
  },
  notify: {
    msg: '',
    timestamp: 0
  }
};
