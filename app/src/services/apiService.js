import {HttpService} from './httpService';

const ALL_SERVICE_CONTEXT = '/v1';
const QUERY_SERVICE_PATH = '/query';
const ENTITY_SERVICE_PATH = '/entity/uid/';

export class apiService {
  static getKeyword(input) {
    const params = {
      keyword: input
    };
    const ALL_SERVICE_URL = window.envConfig.KATLAS_API_URL;
    return request_helper(ALL_SERVICE_URL + ALL_SERVICE_CONTEXT +
      QUERY_SERVICE_PATH, params);
  }

  static getEntity(uid) {
    const ALL_SERVICE_URL = window.envConfig.KATLAS_API_URL;
    return request_helper(ALL_SERVICE_URL + ALL_SERVICE_CONTEXT +
      ENTITY_SERVICE_PATH + uid);
  }
}

const request_helper = (url, params) => {
  return HttpService.get({
    url,
    params
  }).then((response) => {
    return response;
  }).catch((error) => {
    throw error;
  });
};
