import fetchMock from 'fetch-mock';
import {HttpService} from './HttpService';

it('Returns 200 for GET request, then response is as expected', () => {

    let dummyUrl = "http://katlas.com/v1/qsl";
    let dummyParams = {qslstring: 'Cluster'};
    const response = {
      body: {foo: 'bar'},
      status: 200
    };

    fetchMock.get('http://katlas.com/v1/qsl?qslstring=Cluster', response, { overwriteRoutes: true });

    HttpService.get({
      url: dummyUrl,
      params: dummyParams
    }).then((response) => {
      expect(fetchMock.called).toBeTruthy();
      expect(response).toEqual(response);
    });

    fetchMock.reset();
});

it('Returns 204 for GET request, then response is null', () => {

    let dummyUrl = "http://katlas.com/v1/qsl";
    let dummyParams = {qslstring: 'Cluster'};
    const response = {
      body: {foo: 'bar'},
      status: 204
    };

    fetchMock.get('http://katlas.com/v1/qsl?qslstring=Cluster', response, { overwriteRoutes: true });

    HttpService.get({
      url: dummyUrl,
      params: dummyParams
    }).then((resp) => {
      expect(fetchMock.called).toBeTruthy();
      expect(resp.status).toEqual(response.status);
    });

    fetchMock.reset();
});

it('Returns 400 for GET request, then response is null', () => {

    let dummyUrl = "http://katlas.com/v1/qsl";
    let dummyParams = {qslstring: 'Cluster'};
    const response = {
      body: {},
      status: 400
    };

    fetchMock.get('http://katlas.com/v1/qsl?qslstring=Cluster', response, { overwriteRoutes: true });

    HttpService.get({
      url: dummyUrl,
      params: dummyParams
    }).then((resp) => {
      expect(fetchMock.called).toBeTruthy();
      expect(resp.status).toEqual(response.status);
    });

    fetchMock.reset();
});

