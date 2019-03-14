import { HttpService } from './HttpService';

import * as notifyActions from '../actions/notifyActions';
import { addPaginationFilterQSL } from '../utils/validate';
import store from '../store.js';

const ALL_SERVICE_CONTEXT = '/v1.1';

const QUERY_KEYWORD_SERVICE_PATH = '/query';
const QUERY_KEYWORD_PARAM_NAME = 'keyword';
const QUERY_QSL_SERVICE_PATH = '/qsl/';
const QUERY_METADATA_SERVICE_PATH = '/metadata/';

const ENTITY_SERVICE_PATH = '/entity/';

export default class ApiService {
  static getQueryResult(query, page, rowsPerPage) {
    const params = {
      [QUERY_KEYWORD_PARAM_NAME]: query,
      limit: rowsPerPage,
      offset: page * rowsPerPage
    };
    //load env provided URL at query time to allow conf.js to load it in time
    //in testing
    return requestHelper(
      getServiceURL() + ALL_SERVICE_CONTEXT + QUERY_KEYWORD_SERVICE_PATH,
      params
    );
  }

  static getQSLResult(query, page = 0, rowsPerPage = 50) {
    const updatedQuery = addPaginationFilterQSL(query, page, rowsPerPage);
    return requestHelper(
      getServiceURL() +
        ALL_SERVICE_CONTEXT +
        QUERY_QSL_SERVICE_PATH +
        updatedQuery
    );
  }

  static getEntity(uid) {
    //load env provided URL at query time to allow conf.js to load it in time
    //in testing
    return requestHelper(
      getServiceURL() + ALL_SERVICE_CONTEXT + ENTITY_SERVICE_PATH + uid
    );
  }

  static getMetadata(objtype) {
    return requestHelper(
      getServiceURL() +
        ALL_SERVICE_CONTEXT +
        QUERY_METADATA_SERVICE_PATH +
        objtype
    );
  }
}

const getServiceURL = () => {
  return process.env.NODE_ENV === 'test'
    ? 'http://localhost'
    : window.envConfig.KATLAS_API_URL;
};

const requestHelper = (url, params) => {
  return HttpService.get({
    url,
    params
  })
    .then(res => makeResponse(res))
    .catch(error => {
      store.dispatch(notifyActions.showNotify(error));
      throw error;
    });
};

const makeResponse = res => {
  if (res.status === 204) {
    return null;
  }
  if (res.ok) {
    return res.json();
  } else {
    res.text().then(txt => store.dispatch(notifyActions.showNotify(txt)));
    return null;
  }
};
