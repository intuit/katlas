export default {
  query: {
    current: '',
    lastSubmitted: '',
    isWaiting: false,
    isQSL: false,
    results: [],
    page: 0,
    rowsPerPage: 25,
    count: 0,
    metadata: {}
  },
  entity: {
    rootUid: '',
    entitiesByUid: {},
    qslQuery: '',
    results: {},
    latestTimestamp: 0,
    isWaiting: false
  },
  notify: {
    msg: '',
    timestamp: 0
  }
};
