export default {
  query: {
    current: '',
    lastSubmitted: '',
    submitted: false,
    isWaiting: false,
    results: []
  },
  entity: {
    uidsObj: {},
    results: {},
    latestTimestamp: 0,
  },
  notify: {
    msg: '',
    timestamp: 0,
  }
};