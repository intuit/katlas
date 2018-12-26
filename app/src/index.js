import React from 'react';
import { Provider } from 'react-redux';
import { render } from 'react-dom';
import { Router } from 'react-router-dom';
import { MuiThemeProvider, createMuiTheme } from "@material-ui/core";
import WebFont from 'webfontloader';

import App from './components/app/App';
import store from './store';
import history from './history';
import './index.css';

WebFont.load({
  custom: {
    families: ['fontawesome']
  }
});

//Construct the kernel of the material design theme colors, etc.
//(it can be expanded on, on a per component basis)
const theme = createMuiTheme({
  palette: {
    primary: {
      main: '#242321',
    },
    secondary: {
      main: '#2575E2',
    },
  },
  typography: {
    useNextVariants: true,
  },
});

render(
  <MuiThemeProvider theme={theme}>
    <Provider store={store}>
      <Router history={history}>
        <App/>
      </Router>
    </Provider>
  </MuiThemeProvider>,
  document.getElementById('root')
);
