import * as apiCfg from "../config/apiConfig";

export function getQueryParam(queryUrl) {
  const params = new URLSearchParams(queryUrl);
  return {
    queryStr: params.get(apiCfg.SERVICES.queryParamName),
    page: parseInt(params.get(apiCfg.SERVICES.queryParamPage)),
    limit: parseInt(params.get(apiCfg.SERVICES.queryParamLimit))
  };
}

export function encodeQueryData(data) {
  const ret = [];
  for (let d in data)
    ret.push(encodeURIComponent(d) + '=' + encodeURIComponent(data[d]));
  return ret.join('&');
}
