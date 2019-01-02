import { createBrowserHistory } from 'history';

export default createBrowserHistory({
  //this value is not needed or used during local dev or testing scenarios
  //instead, it's required in deployed scenarios where the application is
  //deployed on a subdirectory of a host (ex: 'example.com/ui'). the value is
  //provided by "homepage" attribute in package.json at webpack build time
  //further info here: https://www.npmjs.com/package/history#using-a-base-url
  basename: process.env.PUBLIC_URL
});
