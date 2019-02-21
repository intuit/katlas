import fetchMock from 'fetch-mock';

import ApiService from './ApiService';
import '../../public/conf.js'; // to import the configuration

describe('ApiService', () => {
  it('can get QSL result', () => {
    const response = {
      body: {},
      status: 200
    };

    // verify the url requested by the api service matches the expected url mocked.
    fetchMock.get('http://localhost/v1.1/qsl/Cluster[$$limit=25,offset=0]{*}', response);

    ApiService.getQSLResult('Cluster{*}', 0, 25);

    // verify the url requested by the api service matches the expected url mocked.
    fetchMock.get('http://localhost/v1.1/qsl/Cluster[@name="abc"$$limit=50,offset=100]{*}', response);

    ApiService.getQSLResult('Cluster[@name="abc"]{*}', 2, 50);

    fetchMock.reset();
  });
});
