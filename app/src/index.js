import React from 'react';
import { Provider } from 'react-redux';
import { render } from 'react-dom';
import { BrowserRouter } from 'react-router-dom';
import { MuiThemeProvider, createMuiTheme } from "@material-ui/core";
import WebFont from 'webfontloader';

import './index.css';
import configureStore from './store/configureStore';
import App from './components/app/App';

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

const store = configureStore();

render(
  <MuiThemeProvider theme={theme}>
    <Provider store={store}>
      <BrowserRouter basename={process.env.PUBLIC_URL}>
        <App/>
      </BrowserRouter>
    </Provider>
  </MuiThemeProvider>,
  document.getElementById('root')
);
