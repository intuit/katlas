export default {
  query: {
    current: '',
    lastSubmitted: '',
    submitted: false,
    isWaiting: false,
    results: [],
  },
  entity: {
    rootUid: '',
    entitiesByUid: {},
    results: {},
    latestTimestamp: 0,
    isWaiting: false,
  },
  notify: {
    msg: '',
    timestamp: 0,
  }
};